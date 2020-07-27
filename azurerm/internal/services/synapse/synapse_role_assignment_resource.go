package synapse

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/2020-02-01-preview/accesscontrol"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSynapseRoleAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSynapseRoleAssignmentCreate,
		Read:   resourceArmSynapseRoleAssignmentRead,
		Delete: resourceArmSynapseRoleAssignmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SynapseRoleAssignmentID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"workspace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SynapseWorkspaceName,
			},

			"principal_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmSynapseRoleAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceName := d.Get("workspace_name").(string)
	principalID := d.Get("principal_id").(string)
	roleName := d.Get("role_name").(string)

	client := synapseClient.AccessControlClient(workspaceName)
	roleId, err := getRoleIdByName(ctx, client, roleName)
	if err != nil {
		return err
	}

	// check exist
	listResp, err := client.GetRoleAssignments(ctx, roleId, principalID, "")
	if err != nil {
		if !utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("checking for present of existing Synapse Role Assignment (workspace %q): %+v", workspaceName, err)
		}
	}
	if listResp.Value != nil && len(*listResp.Value) != 0 {
		existing := (*listResp.Value)[0]
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_synapse_role_assignment", *existing.ID)
		}
	}

	// create
	roleAssignment := accesscontrol.RoleAssignmentOptions{
		RoleID:      utils.String(roleId),
		PrincipalID: utils.String(principalID),
	}
	resp, err := client.CreateRoleAssignment(ctx, roleAssignment)
	if err != nil {
		return fmt.Errorf("creating Synapse RoleAssignment %q: %+v", roleName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse RoleAssignment %q", roleName)
	}

	id := fmt.Sprintf("%s|%s", workspaceName, *resp.ID)
	d.SetId(id)
	return resourceArmSynapseRoleAssignmentRead(d, meta)
}

func resourceArmSynapseRoleAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SynapseRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	client := synapseClient.AccessControlClient(id.WorkspaceName)

	resp, err := client.GetRoleAssignmentByID(ctx, id.Id)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] synapse role assignment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse RoleAssignment (Resource Group %q): %+v", id.WorkspaceName, err)
	}

	if resp.RoleID != nil {
		role, err := client.GetRoleDefinitionByID(ctx, *resp.RoleID)
		if err != nil {
			return fmt.Errorf("retrieving role definition by ID %q: %+v", *resp.RoleID, err)
		}
		d.Set("role_name", role.Name)
	}

	d.Set("workspace_name", id.WorkspaceName)
	d.Set("principal_id", resp.PrincipalID)
	return nil
}

func resourceArmSynapseRoleAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SynapseRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	client := synapseClient.AccessControlClient(id.WorkspaceName)

	if _, err := client.DeleteRoleAssignmentByID(ctx, id.Id); err != nil {
		return fmt.Errorf("deleting Synapse RoleAssignment %q (workspace %q): %+v", id, id.WorkspaceName, err)
	}

	return nil
}

func getRoleIdByName(ctx context.Context, client *accesscontrol.BaseClient, roleName string) (string, error) {
	resp, err := client.GetRoleDefinitions(ctx)
	if err != nil {
		return "", fmt.Errorf("listing synapse role definitions %+v", err)
	}

	var availableRoleName []string
	for resp.NotDone() {
		for _, role := range resp.Values() {
			if role.Name != nil {
				if *role.Name == roleName && role.ID != nil {
					return *role.ID, nil
				}
				availableRoleName = append(availableRoleName, *role.Name)
			}
		}
		if err := resp.NextWithContext(ctx); err != nil {
			return "", fmt.Errorf("retrieving next page of storage accounts: %+v", err)
		}
	}
	return "", fmt.Errorf("role name %q invalid. Available role names are %q", roleName, strings.Join(availableRoleName, ","))
}
