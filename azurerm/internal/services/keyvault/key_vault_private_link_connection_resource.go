package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultPrivateEndpointConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultPrivateEndpointConnectionCreate,
		Read:   resourceArmKeyVaultPrivateEndpointConnectionRead,
		Update: resourceArmKeyVaultPrivateEndpointConnectionUpdate,
		Delete: resourceArmKeyVaultPrivateEndpointConnectionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.KeyVaultPrivateEndpointConnectionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KeyVaultName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KeyVaultID,
			},

			"private_link_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmKeyVaultPrivateEndpointConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	keyVaultId, err := parse.KeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}
	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, keyVaultId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for presence of existing Kusto Attached Database Configuration %q (Resource Group %q, Cluster %q): %s", name, resourceGroup, keyVaultId.Name, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_private_link_connection", *resp.ID)
		}
	}

	conn := keyvault.PrivateEndpointConnection{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		PrivateEndpointConnectionProperties: &keyvault.PrivateEndpointConnectionProperties{
			PrivateEndpoint: &keyvault.PrivateEndpoint{
				ID: utils.String(d.Get("private_endpoint_id").(string)),
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Put(ctx, resourceGroup, keyVaultId.Name, name, conn); err != nil {
		return fmt.Errorf("creating or updating Kusto Attached Database Configuration %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, keyVaultId.Name, err)
	}

	read, err := client.Get(ctx, resourceGroup, keyVaultId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Kusto Attached Database Configuration %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, keyVaultId.Name, err)
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("empty or nil ID for Automation Connection '%s' (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmKeyVaultRead(d, meta)
}

func resourceArmKeyVaultPrivateEndpointConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*clients.Client).KeyVault.PrivateEndpointConnectionsClient
	//ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	//defer cancel()

	//id, err := azure.ParseAzureResourceID(d.Id())
	//if err != nil {
	//	return err
	//}
	//
	//if d.HasChange("tags") {
	//	t := d.Get("tags").(map[string]interface{})
	//	update.Tags = tags.Expand(t)
	//}
	//
	//if _, err := client.Update(ctx, resourceGroup, name, update); err != nil {
	//	return fmt.Errorf("Error updating Key Vault %q (Resource Group %q): %+v", name, resourceGroup, err)
	//}
	//

	return resourceArmKeyVaultPrivateEndpointConnectionRead(d, meta)
}

func resourceArmKeyVaultPrivateEndpointConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyVaultPrivateEndpointConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VaultName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Key Vault %q private Endpoint Connection %q was not found in Resource Group %q - removing from state!", id.VaultName, id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on KeyVault %q private Endpoint Connection %q (Resource Group %q): %+v", id.VaultName, id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if resp.PrivateEndpointConnectionProperties != nil && resp.PrivateEndpointConnectionProperties.PrivateEndpoint != nil {
		d.Set("private_link_id", resp.PrivateEndpointConnectionProperties.PrivateEndpoint.ID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmKeyVaultPrivateEndpointConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyVaultPrivateEndpointConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.VaultName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Private Endpoint Connection %q (Resource Group %q / Key Vault %q): %+v", id.Name, id.ResourceGroup, id.VaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Private Endpoint Connection %q (Resource Group %q / Key Vault %q): %+v", id.Name, id.ResourceGroup, id.VaultName, err)
	}

	return nil
}
