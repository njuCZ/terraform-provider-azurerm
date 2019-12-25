subcategory: "App Platform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_binding"
sidebar_current: "docs-azurerm-resource-spring-cloud-app-binding"
description: |-
  Gets information about an existing Spring Cloud App Binding Instance
---

# azurerm_spring_cloud_app_binding

Use this data source to access information about an existing Spring Cloud App Binding Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resGroup"
  location = "Southeast Asia"
}

data "azurerm_spring_cloud_app_binding" "example" {
    name = "abce"
    resource_group_name = "ms-cloud-env-test-000"
    spring_cloud_name = "testabc"
    spring_cloud_app_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Binding resource. Changing this forces a new resource to be created.

* `resource_group` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `app_name` - (Required) The name of the App resource. Changing this forces a new resource to be created.

* `service_name` - (Required) The name of the Service resource. Changing this forces a new resource to be created.

* `binding_parameters` - (Optional) Binding parameters of the Binding resource

* `key` - (Optional) The key of the bound resource

* `resource_id` - (Optional) The Azure resource id of the bound resource

* `resource_name` - (Optional) The name of the bound resource

* `resource_type` - (Optional) The standard Azure resource type of the bound resource

## Attributes Reference

The following attributes are exported:

* `generated_properties` - The generated Spring Boot property file for this binding. The secret will be deducted.

* `created_at` - Creation time of the Binding resource

* `updated_at` - Update time of the Binding resource

* `id` - Fully qualified resource Id for the resource.

* `name` - The name of the resource.

* `type` - The type of the resource.
