package parse

import (
	"testing"
)

func TestManagedHardwareSecurityModuleID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ManagedHardwareSecurityModuleId
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
			Name:     "Missing ManagedHardwareSecurityModule Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.KeyVault/managedHSMs",
			Expected: nil,
		},
		{
			Name:  "ManagedHardwareSecurityModuleID ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.KeyVault/managedHSMs/hsm1",
			Expected: &ManagedHardwareSecurityModuleId{
				ResourceGroup: "group1",
				Name:          "hsm1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.KeyVault/ManagedHSMs/hsm1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := ManagedHardwareSecurityModuleID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
