package parse

import "testing"

func TestParseAppServiceVirtualNetworkConnectionID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *AppServiceVirtualNetworkConnectionId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Error: true,
		},
		{
			Name:  "Missing App Service Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web/sites/",
			Error: true,
		},
		{
			Name:  "Missing Virtual Network Connection Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web/sites/as1/virtualNetworkConnections/",
			Error: true,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Web/sites/as1/virtualNetworkConnections/vnet1",
			Error: false,
			Expect: &AppServiceVirtualNetworkConnectionId{
				ResourceGroup:  "resGroup1",
				AppServiceName: "as1",
				Name:           "vnet1",
			},
		},
		{
			Name:  "Wrong Case",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Web/sites/as1/virtualnetworkconnections/vnet1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := AppServiceVirtualNetworkConnectionID(v.Input)
		if err != nil {
			if v.Expect == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
