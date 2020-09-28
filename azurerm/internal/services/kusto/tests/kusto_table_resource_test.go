package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMKustoTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// CheckDestroy: testCheckAzureRMKustoClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoTable_basic(data),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckAzureRMKustoClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKustoTable_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_kusto_table" "test" {
  name                = "table1"
  kusto_database_id   = "/subscriptions/85b3dbca-5974-4067-9669-67a141095a76/resourceGroups/cwz/providers/Microsoft.Kusto/Clusters/czcz/Databases/a"
  schema              = "Level:string, Timestamp:datetime, UserId:string, TraceId:string, Message:string, ProcessId:int32"
}
`)
}
