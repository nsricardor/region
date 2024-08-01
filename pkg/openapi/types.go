// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package openapi

import (
	externalRef0 "github.com/unikorn-cloud/core/pkg/openapi"
)

const (
	Oauth2AuthenticationScopes = "oauth2Authentication.Scopes"
)

// Defines values for GpuVendor.
const (
	AMD    GpuVendor = "AMD"
	NVIDIA GpuVendor = "NVIDIA"
)

// Defines values for ImageVirtualization.
const (
	Any         ImageVirtualization = "any"
	Baremetal   ImageVirtualization = "baremetal"
	Virtualized ImageVirtualization = "virtualized"
)

// Defines values for RegionType.
const (
	Openstack RegionType = "openstack"
)

// ExternalNetwork An Openstack external network.
type ExternalNetwork struct {
	// Id The resource ID.
	Id string `json:"id"`

	// Name The resource name.
	Name string `json:"name"`
}

// ExternalNetworks A list of openstack external networks.
type ExternalNetworks = []ExternalNetwork

// Flavor A flavor.
type Flavor struct {
	// Metadata This metadata is for resources that just exist, and don't require
	// any provisioning and health status, but benefit from a standarized
	// metadata format.
	Metadata externalRef0.StaticResourceMetadata `json:"metadata"`

	// Spec A flavor.
	Spec FlavorSpec `json:"spec"`
}

// FlavorSpec A flavor.
type FlavorSpec struct {
	// Baremetal Whether the flavor is for a dedicated machine.
	Baremetal *bool `json:"baremetal,omitempty"`

	// CpuFamily A free form CPU family description e.g. model number, architecture.
	CpuFamily *string `json:"cpuFamily,omitempty"`

	// Cpus The number of CPUs.
	Cpus int `json:"cpus"`

	// Disk The amount of ephemeral disk in GB.
	Disk int `json:"disk"`

	// Gpu GPU specification.
	Gpu *GpuSpec `json:"gpu,omitempty"`

	// Memory The amount of memory in GiB.
	Memory int `json:"memory"`
}

// Flavors A list of flavors.
type Flavors = []Flavor

// GpuModel A GPU model number.
type GpuModel = string

// GpuModelList A list of GPU model numbers.
type GpuModelList = []GpuModel

// GpuSpec GPU specification.
type GpuSpec struct {
	// Count The number of GPUs available.
	Count int `json:"count"`

	// Memory GPU memory in GiB.
	Memory int `json:"memory"`

	// Model A GPU model.
	Model string `json:"model"`

	// Vendor The GPU vendor.
	Vendor GpuVendor `json:"vendor"`
}

// GpuVendor The GPU vendor.
type GpuVendor string

// IdentitiesRead A list of provider specific identities.
type IdentitiesRead = []IdentityRead

// IdentityRead A provider specific identity.
type IdentityRead struct {
	Metadata externalRef0.ProjectScopedResourceReadMetadata `json:"metadata"`

	// Spec A provider specific identity, while the client can list regions to infer the
	// type, we don't requires this and return it with the response.  That can then
	// be used in turn to determine which provider specification to examine.
	Spec IdentitySpec `json:"spec"`
}

// IdentitySpec A provider specific identity, while the client can list regions to infer the
// type, we don't requires this and return it with the response.  That can then
// be used in turn to determine which provider specification to examine.
type IdentitySpec struct {
	// Openstack Everything an OpenStack client needs to function.
	Openstack *IdentitySpecOpenStack `json:"openstack,omitempty"`

	// RegionId The region an identity is provisioned in.
	RegionId string `json:"regionId"`

	// Tags A list of tags.
	Tags *TagList `json:"tags,omitempty"`

	// Type The region's provider type.
	Type RegionType `json:"type"`
}

// IdentitySpecOpenStack Everything an OpenStack client needs to function.
type IdentitySpecOpenStack struct {
	// Cloud The name of the cloud in the cloud config.
	Cloud string `json:"cloud"`

	// CloudConfig A base64 encoded cloud config file.
	CloudConfig string `json:"cloudConfig"`

	// ProjectId Project identifier allocated for the infrastructure.
	ProjectId string `json:"projectId"`

	// ServerGroupId Server group identifier allocated for the intrastructure.
	ServerGroupId *string `json:"serverGroupId,omitempty"`

	// UserId User identitifer allocated for the infrastructure.
	UserId string `json:"userId"`
}

// IdentityWrite An identity request.
type IdentityWrite struct {
	// Metadata Resource metadata valid for all API resource reads and writes.
	Metadata externalRef0.ResourceWriteMetadata `json:"metadata"`

	// Spec Request parameters for creating an identity.
	Spec IdentityWriteSpec `json:"spec"`
}

// IdentityWriteSpec Request parameters for creating an identity.
type IdentityWriteSpec struct {
	// RegionId The region an identity is provisioned in.
	RegionId string `json:"regionId"`

	// Tags A list of tags.
	Tags *TagList `json:"tags,omitempty"`
}

// Image An image.
type Image struct {
	// Metadata This metadata is for resources that just exist, and don't require
	// any provisioning and health status, but benefit from a standarized
	// metadata format.
	Metadata externalRef0.StaticResourceMetadata `json:"metadata"`

	// Spec An image.
	Spec ImageSpec `json:"spec"`
}

// ImageGpu The GPU driver if installed.
type ImageGpu struct {
	// Driver The GPU driver version, this is vendor specific.
	Driver string `json:"driver"`

	// Models A list of GPU model numbers.
	Models *GpuModelList `json:"models,omitempty"`

	// Vendor The GPU vendor.
	Vendor GpuVendor `json:"vendor"`
}

// ImageSpec An image.
type ImageSpec struct {
	// Gpu The GPU driver if installed.
	Gpu *ImageGpu `json:"gpu,omitempty"`

	// SoftwareVersions Image preinstalled version version metadata.
	SoftwareVersions *SoftwareVersions `json:"softwareVersions,omitempty"`

	// Virtualization What type of machine the image is for.
	Virtualization ImageVirtualization `json:"virtualization"`
}

// ImageVirtualization What type of machine the image is for.
type ImageVirtualization string

// Images A list of images that are compatible with this platform.
type Images = []Image

// KubernetesNameParameter A Kubernetes name. Must be a valid DNS containing only lower case characters, numbers or hyphens, start and end with a character or number, and be at most 63 characters in length.
type KubernetesNameParameter = string

// PhysicalNetworkRead A physical network.
type PhysicalNetworkRead struct {
	Metadata externalRef0.ProjectScopedResourceReadMetadata `json:"metadata"`

	// Spec A phyical network's specification.
	Spec *PhysicalNetworkSpec `json:"spec,omitempty"`
}

// PhysicalNetworkSpec A phyical network's specification.
type PhysicalNetworkSpec struct {
	// Tags A list of tags.
	Tags *TagList `json:"tags,omitempty"`
}

// PhysicalNetworkWrite A physical network request.
type PhysicalNetworkWrite struct {
	// Metadata Resource metadata valid for all API resource reads and writes.
	Metadata externalRef0.ResourceWriteMetadata `json:"metadata"`

	// Spec A phyical network's specification.
	Spec *PhysicalNetworkSpec `json:"spec,omitempty"`
}

// RegionFeatures A set of features the region may provide to clients.
type RegionFeatures struct {
	// PhysicalNetworks If set, this indicates that the region supports physical networks and
	// one should be provisioned for clusters to use.  The impliciation here is
	// the region supports base-metal machines, and these must be provisioned
	// on a physical VLAN etc.
	PhysicalNetworks bool `json:"physicalNetworks"`
}

// RegionRead A region.
type RegionRead struct {
	// Metadata Resource metadata valid for all reads.
	Metadata externalRef0.ResourceReadMetadata `json:"metadata"`

	// Spec Information about the region.
	Spec RegionSpec `json:"spec"`
}

// RegionSpec Information about the region.
type RegionSpec struct {
	// Features A set of features the region may provide to clients.
	Features RegionFeatures `json:"features"`

	// Type The region's provider type.
	Type RegionType `json:"type"`
}

// RegionType The region's provider type.
type RegionType string

// Regions A list of regions.
type Regions = []RegionRead

// SoftwareVersions Image preinstalled version version metadata.
type SoftwareVersions struct {
	// Kubernetes A semantic version.
	Kubernetes *externalRef0.Semver `json:"kubernetes,omitempty"`
}

// Tag An arbitrary tag name and value.
type Tag struct {
	// Name A unique tag name.
	Name string `json:"name"`

	// Value The value of the tag.
	Value string `json:"value"`
}

// TagList A list of tags.
type TagList = []Tag

// IdentityIDParameter A Kubernetes name. Must be a valid DNS containing only lower case characters, numbers or hyphens, start and end with a character or number, and be at most 63 characters in length.
type IdentityIDParameter = KubernetesNameParameter

// OrganizationIDParameter defines model for organizationIDParameter.
type OrganizationIDParameter = string

// ProjectIDParameter A Kubernetes name. Must be a valid DNS containing only lower case characters, numbers or hyphens, start and end with a character or number, and be at most 63 characters in length.
type ProjectIDParameter = KubernetesNameParameter

// RegionIDParameter A Kubernetes name. Must be a valid DNS containing only lower case characters, numbers or hyphens, start and end with a character or number, and be at most 63 characters in length.
type RegionIDParameter = KubernetesNameParameter

// ExternalNetworksResponse A list of openstack external networks.
type ExternalNetworksResponse = ExternalNetworks

// FlavorsResponse A list of flavors.
type FlavorsResponse = Flavors

// IdentitiesResponse A list of provider specific identities.
type IdentitiesResponse = IdentitiesRead

// IdentityResponse A provider specific identity.
type IdentityResponse = IdentityRead

// ImagesResponse A list of images that are compatible with this platform.
type ImagesResponse = Images

// PhysicalNetworkResponse A physical network.
type PhysicalNetworkResponse = PhysicalNetworkRead

// RegionsResponse A list of regions.
type RegionsResponse = Regions

// IdentityRequest An identity request.
type IdentityRequest = IdentityWrite

// PhysicalNetworkRequest A physical network request.
type PhysicalNetworkRequest = PhysicalNetworkWrite

// PostApiV1OrganizationsOrganizationIDProjectsProjectIDIdentitiesJSONRequestBody defines body for PostApiV1OrganizationsOrganizationIDProjectsProjectIDIdentities for application/json ContentType.
type PostApiV1OrganizationsOrganizationIDProjectsProjectIDIdentitiesJSONRequestBody = IdentityWrite

// PostApiV1OrganizationsOrganizationIDProjectsProjectIDIdentitiesIdentityIDPhysicalNetworksJSONRequestBody defines body for PostApiV1OrganizationsOrganizationIDProjectsProjectIDIdentitiesIdentityIDPhysicalNetworks for application/json ContentType.
type PostApiV1OrganizationsOrganizationIDProjectsProjectIDIdentitiesIdentityIDPhysicalNetworksJSONRequestBody = PhysicalNetworkWrite
