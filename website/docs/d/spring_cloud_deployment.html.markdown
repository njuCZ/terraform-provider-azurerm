subcategory: "appplatform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_deployment"
sidebar_current: "docs-azurerm-datasource-spring-cloud-deployment"
description: |-
  Gets information about an existing Azure Spring Cloud Deployment instance.
---

# Data Source: azurerm_spring_cloud_deployment

Use this data source to access information about an existing Azure Spring Cloud Deployment instance.

## Example Usage

```hcl
data "azurerm_spring_cloud_deployment" "example" {
  spring_cloud_name             = azurerm_spring_cloud.example.name
  resource_group_name           = azurerm_resource_group.example.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.example.name
  name                          = "scd-example"
}

output "spring_cloud_deployment_id" {
  value = "${data.azurerm_spring_cloud_deployment.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Spring Cloud Application resource.

* `resource_group_name` - (Required) The name of the resource group that contains the resource.

* `spring_cloud_name` - (Required) The name of the Spring Cloud Service resource.

* `spring_cloud_app_name` - (Required) The name of the Spring Cloud Application resource.

## Attributes Reference

The following attributes are exported:

* `id` - Fully qualified resource Id for the resource.

* `cpu` - The number of CPU assigned to this deployment. The value Should be between 1 and 4, Default 1.

* `memory_in_gb` - The size (GB) of memory assigned to this deployment. The value should be between 1 amd 8, Default 1.

* `runtime_version` - The runtime version for the deployment. It should be either `Java_8` or `Java_11`, Default `Java_8`.

* `instance_count` - The number of instance of the deployment. It should be between 1 and 20, Default 1.

* `jvm_options` - The jvm option of the deployment.

* `env` - the environment variables for the deployment. It is a map of string key value pair