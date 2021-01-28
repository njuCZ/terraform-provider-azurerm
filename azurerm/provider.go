package azurerm

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func Provider() terraform.ResourceProvider {
	log.Printf("[INFO] *************************************************************************************\n")
	log.Printf("[INFO] ***** Using Terraform Azure provider with Spring Cloud Service Vnet Integration *****\n")
	log.Printf("[INFO] *************************************************************************************\n\n")
	return provider.AzureProvider()
}
