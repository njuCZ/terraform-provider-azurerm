subcategory: "App Platform"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_cosmosdb_binding"
sidebar_current: "docs-azurerm-resource-spring-cloud-app-cosmosdb-binding"
description: |-
  Manage Azure Spring Cloud App cosmosdb Binding instance.
---

# azurerm_spring_cloud_app_cosmosdb_binding

Manage Azure Spring Cloud App cosmosdb Binding instance.

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

resource "azurerm_cosmosdb_account" "example" {
  name                = "cosmosdb-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  capabilities {
    name = "EnableCassandra"
  }

  consistency_policy {
    consistency_level = "BoundedStaleness"
  }

  geo_location {
    location            = azurerm_resource_group.example.location
    failover_priority   = 0
  }
}

resource "azurerm_cosmosdb_cassandra_keyspace" "example" {
  name                = "cassandra-keyspace"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
}

resource "azurerm_spring_cloud_app_cosmosdb_binding" "example" {
    name                    = "app_cosmosdb_binding"
    resource_group_name    	= azurerm_resource_group.example.name
    spring_cloud_name       = azurerm_spring_cloud.example.name
    spring_cloud_app_name 	= azurerm_spring_cloud_app.example.name
    cosmosdb_account_id     = azurerm_cosmosdb_account.example.id
    api_type                = "cassandra"
    cassandra_key_space     = azurerm_cosmosdb_cassandra_keyspace.example.name
    key                     = azurerm_cosmosdb_account.example.primary_master_key
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Binding resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `spring_cloud_name` - (Required) The name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `spring_cloud_app_name` - (Required) The name of the Spring Cloud App resource. Changing this forces a new resource to be created.

* `cosmosdb_account_id` - (Required) The resource id of cosmosdb account that Spring Cloud App wants to bind.

* `api_type` - (Required) The type of resource within cosmosdb account that Spring Cloud App wants to bind. The value should be one of `cassandra`, `mongo`, `sql`, `table`.

* `cassandra_key_space` - (Optional) The key space of cassandra. If `api_type` is `cassandra`, the field is Required. Conflicts with `mongo_database_name`, `sql_database_name`.

* `mongo_database_name` - (Optional) The database name of mongo. If `api_type` is `mongo`, the field is Required. Conflicts with `cassandra_key_space`, `sql_database_name`.

* `sql_database_name` - (Optional) The database name of sql. If `api_type` is `sql`, the field is Required. Conflicts with `cassandra_key_space`, `mongo_database_name`.

* `key` - (Required) The access key of cosmosdb account.

## Attributes Reference

The following attributes are exported:

* `generated_properties` - The generated Spring Boot property file for this binding. The secret will be deducted.

* `resource_name` - The Azure resource name of the bounded resource

* `resource_type` - The resource type of the bounded resource

* `id` - Fully qualified resource Id for the resource.
