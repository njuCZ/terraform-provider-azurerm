subcategory: "App Platform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_mysql_binding"
sidebar_current: "docs-azurerm-resource-spring-cloud-app-mysql-binding"
description: |-
  Manage Azure Spring Cloud App mysql Binding instance.
---

# azurerm_spring_cloud_app_mysql_binding

Manage Azure Spring Cloud App mysql Binding instance.

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

resource "azurerm_mysql_server" "example" {
  name                = "mysql-server"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 5120
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mysql_database" "example" {
  name                = "mysql-db"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_server.example.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}

resource "azurerm_spring_cloud_app_mysql_binding" "example" {
    name                    = "app-mysql-binding"
    resource_group_name    	= azurerm_resource_group.example.name
    spring_cloud_name       = azurerm_spring_cloud.example.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.example.name
    mysql_server_id         = azurerm_mysql_server.example.id
    database_name           = azurerm_mysql_database.example.name
    username                = azurerm_mysql_server.example.administrator_login
    password                = azurerm_mysql_server.example.administrator_login_password
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Binding resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `spring_cloud_name` - (Required) The name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `spring_cloud_app_name` - (Required) The name of the Spring Cloud App resource. Changing this forces a new resource to be created.

* `mysql_server_id` - (Required) The resource id of `mysql_server` that Spring Cloud App wants to bind

* `database_name` - (Required) The name of database within `mysql_server_id` that Spring Cloud App wants to bind

* `username` - (Required) The username of database to be accessed by the spring cloud app

* `password` - (Required) The password of database to be accessed by the spring cloud app

## Attributes Reference

The following attributes are exported:

* `generated_properties` - The generated Spring Boot property file for this binding. The secret will be deducted.

* `resource_name` - The Azure resource name of the bounded resource

* `resource_type` - The resource type of the bounded resource

* `id` - Fully qualified resource Id for the resource.
