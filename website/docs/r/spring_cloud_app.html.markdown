subcategory: "appplatform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app"
sidebar_current: "docs-azurerm-resource-spring-cloud-app"
description: |-
  Manage Azure Spring Cloud Application instance.
---

# azurerm_spring_cloud_app

Manage Azure Spring Cloud Application instance. An Azure Spring Cloud Service can contain multiple spring cloud apps

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg"
  location = "Southeast Asia"
}

resource "azurerm_spring_cloud" "example" {
  name                     = "sc-example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location

  tags = {
    environment = "staging"
  }
}

resource "azurerm_spring_cloud_app" "example" {
  spring_cloud_name             = azurerm_spring_cloud.example.name
  resource_group_name           = azurerm_resource_group.example.name
  name                          = "sca-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Spring Cloud Application resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group that contains the resource. Changing this forces a new resource to be created.

* `spring_cloud_name` - (Required) The name of the Spring Cloud Service resource. Changing this forces a new resource to be created. 

## Attributes Reference

The following attributes are exported:

* `id` - Fully qualified resource Id for the resource.

* `created_time` - The creation time of the resource.

## Import

Spring CLoud App can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/Spring/myservice/apps/myapp
```
