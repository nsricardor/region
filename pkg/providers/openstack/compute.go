/*
Copyright 2022-2024 EscherCloud.
Copyright 2024 the Unikorn Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openstack

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/availabilityzones"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/quotasets"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/unikorn-cloud/core/pkg/util/cache"
	unikornv1 "github.com/unikorn-cloud/region/pkg/apis/unikorn/v1alpha1"
	"github.com/unikorn-cloud/region/pkg/constants"

	"k8s.io/utils/ptr"
)

// ComputeClient wraps the generic client because gophercloud is unsafe.
type ComputeClient struct {
	options *unikornv1.RegionOpenstackComputeSpec
	client  *gophercloud.ServiceClient

	flavorCache *cache.TimeoutCache[[]flavors.Flavor]
}

// NewComputeClient provides a simple one-liner to start computing.
func NewComputeClient(ctx context.Context, provider CredentialProvider, options *unikornv1.RegionOpenstackComputeSpec) (*ComputeClient, error) {
	providerClient, err := provider.Client(ctx)
	if err != nil {
		return nil, err
	}

	client, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	// Need at least 2.15 for soft-anti-affinity policy.
	// Need at least 2.64 for new server group interface.
	client.Microversion = "2.90"

	c := &ComputeClient{
		options:     options,
		client:      client,
		flavorCache: cache.New[[]flavors.Flavor](time.Hour),
	}

	return c, nil
}

// CreateKeypair creates a new keypair.
// NOTE: while OpenStack can generate one for us, we have far more control doing it ourselves
// thus allowing us to impose stricter security, and it's more provider agnostic that way.
func (c *ComputeClient) CreateKeypair(ctx context.Context, name, publicKey string) error {
	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, "POST /compute/v2/os-keypairs", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	opts := &keypairs.CreateOpts{
		Name:      name,
		Type:      "ssh",
		PublicKey: publicKey,
	}

	return keypairs.Create(ctx, c.client, opts).Err
}

func (c *ComputeClient) DeleteKeypair(ctx context.Context, name string) error {
	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, fmt.Sprintf("DELETE /compute/v2/os-keypairs/%s", name), trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	return keypairs.Delete(ctx, c.client, name, nil).Err
}

// KeyPairs returns a list of key pairs.
func (c *ComputeClient) KeyPairs(ctx context.Context) ([]keypairs.KeyPair, error) {
	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, "GET /compute/v2/os-keypairs", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	page, err := keypairs.List(c.client, &keypairs.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	return keypairs.ExtractKeyPairs(page)
}

// mutateFlavors allows nova's view of fact to be altered...
func (c *ComputeClient) mutateFlavors(f []flavors.Flavor) {
	if c.options == nil || c.options.Flavors == nil {
		return
	}

	for _, metadata := range c.options.Flavors.Metadata {
		index := slices.IndexFunc(f, func(flavor flavors.Flavor) bool {
			return flavor.ID == metadata.ID
		})

		if index < 0 {
			continue
		}

		if metadata.CPU != nil && metadata.CPU.Count != nil {
			f[index].VCPUs = *metadata.CPU.Count
		}

		if metadata.Memory != nil {
			// Convert from bytes to MiB
			f[index].RAM = int(metadata.Memory.Value() >> 20)
		}
	}
}

// Flavors returns a list of flavors.
//
//nolint:cyclop
func (c *ComputeClient) Flavors(ctx context.Context) ([]flavors.Flavor, error) {
	if result, ok := c.flavorCache.Get(); ok {
		return result, nil
	}

	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, "GET /compute/v2/flavors", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	page, err := flavors.ListDetail(c.client, &flavors.ListOpts{SortKey: "name"}).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	result, err := flavors.ExtractFlavors(page)
	if err != nil {
		return nil, err
	}

	// Mutate any flavors first, as this may alter their selection criteria.
	c.mutateFlavors(result)

	result = slices.DeleteFunc(result, func(flavor flavors.Flavor) bool {
		// We are admin, so see all the things, throw out private flavors.
		// TODO: we _could_ allow if our project is in the allowed IDs.
		if !flavor.IsPublic {
			return true
		}

		// Kubeadm requires 2 VCPU, 2 "GB" of RAM (I'll pretend it's GiB) and no swap:
		// https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
		if flavor.VCPUs < 2 || flavor.RAM < 2048 || flavor.Swap != 0 {
			return true
		}

		// Don't remove the flavor if it's implicitly selected by a lack of configuration.
		if c.options == nil || c.options.Flavors == nil || c.options.Flavors.Selector == nil {
			return false
		}

		if len(c.options.Flavors.Selector.IDs) > 0 {
			if !slices.Contains(c.options.Flavors.Selector.IDs, flavor.ID) {
				return true
			}
		}

		return false
	})

	c.flavorCache.Set(result)

	return result, nil
}

// AvailabilityZones returns a list of availability zones.
func (c *ComputeClient) AvailabilityZones(ctx context.Context) ([]availabilityzones.AvailabilityZone, error) {
	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, "GET /compute/v2/os-availability-zones", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	page, err := availabilityzones.List(c.client).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	result, err := availabilityzones.ExtractAvailabilityZones(page)
	if err != nil {
		return nil, err
	}

	filtered := []availabilityzones.AvailabilityZone{}

	for _, az := range result {
		if !az.ZoneState.Available {
			continue
		}

		filtered = append(filtered, az)
	}

	return filtered, nil
}

// CreateServerGroup creates the named server group with the given policy and returns
// the result.
func (c *ComputeClient) CreateServerGroup(ctx context.Context, name string) (*servergroups.ServerGroup, error) {
	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, "POST /compute/v2/os-server-groups", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	opts := &servergroups.CreateOpts{
		Name:   name,
		Policy: "soft-anti-affinity",
	}

	if c.options != nil && c.options.ServerGroupPolicy != nil {
		opts.Policy = *c.options.ServerGroupPolicy
	}

	return servergroups.Create(ctx, c.client, opts).Extract()
}

// DeleteServerGroup removes a server group, this exists because nova does do any cleanup
// on project deletion and just orphans the resource.
func (c *ComputeClient) DeleteServerGroup(ctx context.Context, id string) error {
	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, "DELETE /compute/v2/os-server-groups/"+id)
	defer span.End()

	return servergroups.Delete(ctx, c.client, id).ExtractErr()
}

func (c *ComputeClient) UpdateQuotas(ctx context.Context, projectID string) error {
	tracer := otel.GetTracerProvider().Tracer(constants.Application)

	_, span := tracer.Start(ctx, "PUT /compute/v2/os-quota-sets")
	defer span.End()

	opts := &quotasets.UpdateOpts{
		// TODO: instances, cores and ram need to be driven by client input.
		Instances: ptr.To(-1),
		Cores:     ptr.To(-1),
		RAM:       ptr.To(-1),
	}

	return quotasets.Update(ctx, c.client, projectID, opts).Err
}
