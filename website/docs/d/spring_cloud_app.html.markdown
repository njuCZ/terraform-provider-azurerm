subcategory: "appplatform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app"
sidebar_current: "docs-azurerm-datasource-spring-cloud-app"
description: |-
  Gets information about an existing Azure Spring Cloud Application instance.
---

# Data Source: azurerm_spring_cloud_app

Use this data source to access information about an existing Azure Spring Cloud Application instance.

## Example Usage

```hcl
data "azurerm_spring_cloud_app" "example" {
  spring_cloud_name             = azurerm_spring_cloud.example.name
  resource_group_name           = azurerm_resource_group.example.name
  name                          = "sca-example"
}

output "spring_cloud_app_id" {
  value = "${data.azurerm_spring_cloud_app.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Spring Cloud Application resource.

* `resource_group_name` - (Required) The name of the resource group that contains the resource.

* `spring_cloud_name` - (Required) The name of the Spring Cloud Service resource.

## Attributes Reference

The following attributes are exported:

* `id` - Fully qualified resource Id for the resource.

* `created_time` - The creation time of the resource.

* `active_deployment_name` - The active Deployment of the resource

* `public` - whether This application has public access endpoint

* `url` - The public access url of this application

* `persistent_disk` - An `persistent_disk` block as defined below.

* `temporary_disk` - An `temporary_disk` block as defined below.

---

An `persistent_disk` block exports the following:

* `mount_path` - The path that persistent disk is mounted.

* `size_in_gb` - The size of persistent disk.

---

An `temporary_disk` block exports the following:

* `mount_path` - The path that temporary disk is mounted.

* `size_in_gb` - The size of temporary disk.
