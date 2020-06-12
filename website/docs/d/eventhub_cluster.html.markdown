---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_cluster"
description: |-
  Gets information about an EventHub Cluster.
---

# azurerm_eventhub_cluster

Use this data source to access information about an existing EventHub Cluster.

## Example Usage

```hcl
data "azurerm_eventhub_cluster" "example" {
  name                = azurerm_eventhub_cluster.example.name
  resource_group_name = azurerm_eventhub_cluster.example.resource_group_name
}

output "eventhub_cluster_id" {
  value = data.azurerm_eventhub_cluster.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Cluster resource.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the EventHub Cluster exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of EventHub Cluster resource.

* `location` - The location of EventHub Cluster resource.

* `tags` - A mapping of tags assigned to EventHub Cluster resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Cluster.
