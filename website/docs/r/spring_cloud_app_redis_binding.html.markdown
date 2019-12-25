subcategory: "App Platform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_redis_binding"
sidebar_current: "docs-azurerm-resource-spring-cloud-app-redis-binding"
description: |-
  Manage Azure Spring Cloud App Redis Binding instance.
---

# azurerm_spring_cloud_app_redis_binding

Manage Azure Spring Cloud App Redis Binding instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resGroup"
  location = "Southeast Asia"
}

resource "azurerm_spring_cloud" "example" {
  name                     = "sc-example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
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

resource "azurerm_spring_cloud_active_deployment" "example" {
  resource_group_name           = azurerm_resource_group.example.name
  spring_cloud_name             = azurerm_spring_cloud.example.name
  spring_cloud_app_name         = azurerm_spring_cloud_app.example.name
  spring_cloud_deployment_name  = azurerm_spring_cloud_deployment.example.name
}

resource "azurerm_redis_cache" "example" {
  name                     = "redis"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  capacity           	   = 0
  family                   = "C"
  sku_name                 = "Basic"
  enable_non_ssl_port      = true
}

resource "azurerm_spring_cloud_app_redis_binding" "example" {
	name 					= "app_redis_binding"
    resource_group_name    	= azurerm_resource_group.example.name
    spring_cloud_name       = azurerm_spring_cloud.example.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.example.name
    redis_cache_id          = azurerm_redis_cache.example.id
    key                     = azurerm_redis_cache.example.primary_access_key
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Binding resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `spring_cloud_name` - (Required) The name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `spring_cloud_app_name` - (Required) The name of the Spring Cloud App resource. Changing this forces a new resource to be created.

* `redis_cache_id` - (Required) The resource id of redis that Spring Cloud App wants to bind

* `key` - (Required) The access key of redis

## Attributes Reference

The following attributes are exported:

* `generated_properties` - The generated Spring Boot property file for this binding. The secret will be deducted.

* `resource_name` - The Azure resource name of the bounded resource

* `resource_type` - The resource type of the bounded resource

* `id` - Fully qualified resource Id for the resource.
