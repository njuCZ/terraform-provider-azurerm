subcategory: "appplatform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_deployment"
sidebar_current: "docs-azurerm-resource-spring-cloud-deployment"
description: |-
  Manage Azure Spring Cloud Deployment instance.
---

# azurerm_spring_cloud_deployment

Manage Azure Spring Cloud Deployment instance. An Azure Spring Cloud Application can contain multiple spring cloud deployments

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

resource "azurerm_spring_cloud_deployment" "example" {
  spring_cloud_name             = azurerm_spring_cloud.example.name
  resource_group_name           = azurerm_resource_group.example.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.example.name
  name                          = "scd-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Spring Cloud Application resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group that contains the resource. Changing this forces a new resource to be created.

* `spring_cloud_name` - (Required) The name of the Spring Cloud Service resource. Changing this forces a new resource to be created. 

* `spring_cloud_app_name` - (Required) The name of the Spring Cloud Application resource. Changing this forces a new resource to be created. 

* `cpu` - (Optional) The number of CPU assigned to this deployment. The value Should be between 1 and 4, Default 1.

* `memory_in_gb` - (Optional) The size (GB) of memory assigned to this deployment. The value should be between 1 amd 8, Default 1.

* `runtime_version` - (Optional) The runtime version for the deployment. It should be either `Java_8` or `Java_11`, Default `Java_8`.

* `instance_count` - (Optional) The number of instance of the deployment. It should be between 1 and 20, Default 1.

* `jvm_options` - (Optional) The jvm option of the deployment.

* `env` - (Optional) the environment variables for the deployment. It is a map of string key value pair

* `jar_file` - (Optional) the jar file running in the deployment. If it's empty, a default jar will be used for the deployment

## Attributes Reference

The following attributes are exported:

* `id` - Fully qualified resource Id for the resource.

## Import

Spring CLoud Deployment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/Spring/myservice/apps/myapp/deployments/mydeployment
```
