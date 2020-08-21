provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-rg"
  location = var.location
}

resource "azurerm_app_service_plan" "test" {
  name                = "${var.prefix}-asp"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "${var.prefix}-as"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}