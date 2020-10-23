package parse

import (
	"testing"
)

func TestKeyVaultPrivateEndpointConnectionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *KeyVaultPrivateEndpointConnectionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Vaults Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/",
			Expected: nil,
		},
		{
			Name:     "Missing Private Endpoint Connection segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1",
			Expected: nil,
		},
		{
			Name:     "Missing Private Endpoint Connection Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/privateEndpointConnections",
			Expected: nil,
		},
		{
			Name:  "Key Vault Private Endpoint Connection ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/privateEndpointConnections/conn1",
			Expected: &KeyVaultPrivateEndpointConnectionId{
				Name:          "conn1",
				VaultName:     "vault1",
				ResourceGroup: "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/PrivateEndpointConnections/conn1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := KeyVaultPrivateEndpointConnectionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.VaultName != v.Expected.VaultName {
			t.Fatalf("Expected %q but got %q for VaultName", v.Expected.VaultName, actual.VaultName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
