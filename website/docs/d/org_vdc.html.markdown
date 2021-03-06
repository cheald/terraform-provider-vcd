---
layout: "vcd"
page_title: "vCloudDirector: vcd_org_vdc"
sidebar_current: "docs-vcd-data-source-org-vdc"
description: |-
  Provides an organization VDC data source.
---

# vcd\_org\_vcd

Provides a vCloud Director Organization VDC data source. An Organization VDC can be used to reference a VCD and use its 
data within other resources or data sources.

Supported in provider *v2.5+*

## Example Usage

```hcl
data "vcd_org_vdc" "my-org-vdc" {
  org     = "my-org"
  name    = "my-vdc"
}

output "provider_vdc" {
 value   = data.vcd_org_vdc.my-org-vdc.provider_vdc_name
}

```

## Argument Reference

The following arguments are supported:

* `org` - (Optional, but required if not set at provider level) Org name 
* `name` - (Required) Organization VDC name

## Attribute Reference

* `description` - VDC friendly description
* `provider_vdc_name` -A name of the Provider VDC from which this organization VDC is provisioned.
* `allocation_model` - The allocation model used by this VDC; must be one of {AllocationVApp ("Pay as you go"), AllocationPool ("Allocation pool"), ReservationPool ("Reservation pool")}
* `compute_capacity` - The compute capacity allocated to this VDC.  See [Compute Capacity](#computecapacity) below for details.
* `nic_quota` - Maximum number of virtual NICs allowed in this VDC. Defaults to 0, which specifies an unlimited number.
* `network_quota` - Maximum number of network objects that can be deployed in this VDC. Defaults to 0, which means no networks can be deployed.
* `vm_quota` - The maximum number of VMs that can be created in this VDC. Includes deployed and undeployed VMs in vApps and vApp templates. Defaults to 0, which specifies an unlimited number.
* `enabled` - True if this VDC is enabled for use by the organization VDCs. Default is true.
* `storage_profile` - Storage profiles supported by this VDC.  See [Storage Profile](#storageprofile) below for details.
* `memory_guaranteed` - Percentage of allocated memory resources guaranteed to vApps deployed in this VDC. For example, if this value is 0.75, then 75% of allocated resources are guaranteed. Required when AllocationModel is AllocationVApp or AllocationPool. When Allocation model is AllocationPool minimum value is 0.2. If left empty, vCD sets a value.
* `cpu_guaranteed` - Percentage of allocated CPU resources guaranteed to vApps deployed in this VDC. For example, if this value is 0.75, then 75% of allocated resources are guaranteed. Required when AllocationModel is AllocationVApp or AllocationPool. If left empty, vCD sets a value.
* `cpu_speed` - Specifies the clock frequency, in Megahertz, for any virtual CPU that is allocated to a VM. A VM with 2 vCPUs will consume twice as much of this value. Ignored for ReservationPool. Required when AllocationModel is AllocationVApp or AllocationPool, and may not be less than 256 MHz. Defaults to 1000 MHz if value isn't provided.
* `metadata` - Key value map of metadata to assign to this VDC
* `enable_thin_provisioning` - Boolean to request thin provisioning. Request will be honored only if the underlying data store supports it. Thin provisioning saves storage space by committing it on demand. This allows over-allocation of storage.
* `enable_fast_provisioning` - (Request fast provisioning. Request will be honored only if the underlying datastore supports it. Fast provisioning can reduce the time it takes to create virtual machines by using vSphere linked clones. If you disable fast provisioning, all provisioning operations will result in full clones.
* `network_pool_name` - Reference to a network pool in the Provider VDC. Required if this VDC will contain routed or isolated networks.
* `allow_over_commit` - Set to false to disallow creation of the VDC if the AllocationModel is AllocationPool or ReservationPool and the ComputeCapacity you specified is greater than what the backing Provider VDC can supply. Default is true.
* `enable_vm_discovery` - If true, discovery of vCenter VMs is enabled for resource pools backing this VDC. If false, discovery is disabled. If left unspecified, the actual behaviour depends on enablement at the organization level and at the system level.

<a id="storageprofile"></a>
## Storage Profile

* `name` - Name of Provider VDC storage profile.
* `enabled` - True if this storage profile is enabled for use in the VDC. Default is true.
* `limit` - Maximum number of MB allocated for this storage profile. A value of 0 specifies unlimited MB.
* `default` - True if this is default storage profile for this VDC. The default storage profile is used when an object that can specify a storage profile is created with no storage profile specified.

<a id="computecapacity"></a>
## Compute Capacity

Capacity must be specified twice, once for `memory` and another for `cpu`.  Each has the same structure:

* `allocated` - Capacity that is committed to be available. Value in MB or MHz. Used with AllocationPool ("Allocation pool") and ReservationPool ("Reservation pool").
* `limit` - Capacity limit relative to the value specified for Allocation. It must not be less than that value. If it is greater than that value, it implies over provisioning. A value of 0 specifies unlimited units. Value in MB or MHz. Used with AllocationVApp ("Pay as you go").