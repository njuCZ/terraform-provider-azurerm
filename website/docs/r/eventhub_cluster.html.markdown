---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_cluster"
description: |-
  Manages an EventHub Cluster.
---

# azurerm_eventhub_cluster

Manages an EventHub Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location			  = azurerm_resource_group.example.location

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Cluster resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the EventHub Cluster. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of EventHub Cluster resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Cluster.

## Import

EventHub Cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1
```
