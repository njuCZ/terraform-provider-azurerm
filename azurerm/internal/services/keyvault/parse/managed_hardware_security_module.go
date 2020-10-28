package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagedHardwareSecurityModuleId struct {
	ResourceGroup string
	Name          string
}

func ManagedHardwareSecurityModuleID(input string) (*ManagedHardwareSecurityModuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing ManagedHardwareSecurityModule ID %q: %+v", input, err)
	}

	managedHardwareSecurityModule := ManagedHardwareSecurityModuleId{
		ResourceGroup: id.ResourceGroup,
	}

	if managedHardwareSecurityModule.Name, err = id.PopSegment("managedHSMs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &managedHardwareSecurityModule, nil
}
