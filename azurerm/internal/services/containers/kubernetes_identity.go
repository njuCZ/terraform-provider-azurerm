package containers

import (
	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msivalidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-12-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func schemaKubernetesIdentity() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(containerservice.ResourceIdentityTypeSystemAssigned),
						string(containerservice.ResourceIdentityTypeUserAssigned),
					}, false),
				},
				"user_assigned_identity_id": {
					Type:         schema.TypeString,
					ValidateFunc: msivalidate.UserAssignedIdentityID,
					Optional:     true,
				},
				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandKubernetesClusterManagedClusterIdentity(input []interface{}) *containerservice.ManagedClusterIdentity {
	if len(input) == 0 || input[0] == nil {
		return &containerservice.ManagedClusterIdentity{
			Type: containerservice.ResourceIdentityTypeNone,
		}
	}

	values := input[0].(map[string]interface{})

	if containerservice.ResourceIdentityType(values["type"].(string)) == containerservice.ResourceIdentityTypeUserAssigned {
		userAssignedIdentities := map[string]*containerservice.ManagedClusterIdentityUserAssignedIdentitiesValue{
			values["user_assigned_identity_id"].(string): {},
		}

		return &containerservice.ManagedClusterIdentity{
			Type:                   containerservice.ResourceIdentityType(values["type"].(string)),
			UserAssignedIdentities: userAssignedIdentities,
		}
	}

	return &containerservice.ManagedClusterIdentity{
		Type: containerservice.ResourceIdentityType(values["type"].(string)),
	}
}

func flattenKubernetesClusterManagedClusterIdentity(input *containerservice.ManagedClusterIdentity) ([]interface{}, error) {
	// if it's none, omit the block
	if input == nil || input.Type == containerservice.ResourceIdentityTypeNone {
		return []interface{}{}, nil
	}

	identity := make(map[string]interface{})

	identity["principal_id"] = ""
	if input.PrincipalID != nil {
		identity["principal_id"] = *input.PrincipalID
	}

	identity["tenant_id"] = ""
	if input.TenantID != nil {
		identity["tenant_id"] = *input.TenantID
	}

	identity["user_assigned_identity_id"] = ""
	if input.UserAssignedIdentities != nil {
		keys := []string{}
		for key := range input.UserAssignedIdentities {
			keys = append(keys, key)
		}
		if len(keys) > 0 {
			parsedId, err := msiparse.UserAssignedIdentityID(keys[0])
			if err != nil {
				return nil, err
			}
			identity["user_assigned_identity_id"] = parsedId.ID()
		}
	}

	identity["type"] = string(input.Type)

	return []interface{}{identity}, nil
}

// when update, we should set the value of `Identity.UserAssignedIdentities` empty
// otherwise the rest api will report error
func normalizeClusterUserAssignedIdentities(instance *containerservice.ManagedCluster) {
	if instance == nil || instance.Identity == nil || instance.Identity.UserAssignedIdentities == nil {
		return
	}

	for k := range instance.Identity.UserAssignedIdentities {
		instance.Identity.UserAssignedIdentities[k] = &containerservice.ManagedClusterIdentityUserAssignedIdentitiesValue{}
	}
}
