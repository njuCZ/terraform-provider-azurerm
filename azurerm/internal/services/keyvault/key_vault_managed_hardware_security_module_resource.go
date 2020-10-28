package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultManagedHardwareSecurityModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultManagedHardwareSecurityModuleCreateUpdate,
		Read:   resourceArmKeyVaultManagedHardwareSecurityModuleRead,
		Update: resourceArmKeyVaultManagedHardwareSecurityModuleCreateUpdate,
		Delete: resourceArmKeyVaultManagedHardwareSecurityModuleDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagedHardwareSecurityModuleID(id)
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
				ValidateFunc: validate.KeyVaultManagedHardwareSecurityModuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			// currently supported location are: eastus2,southcentralus,northeurope,westeurope
			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.StandardB1),
				}, false),
			},

			"admin_object_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			// https://github.com/Azure/azure-rest-api-specs/issues/11419
			// field `purge_protection_enabled`, `soft_delete_retention_days` could not be updated
			"purge_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"soft_delete_retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      90,
				ValidateFunc: validation.IntBetween(7, 90),
			},

			"hsm_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmKeyVaultManagedHardwareSecurityModuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] Preparing arguments for Key Vault Managed Hardware Security Module")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Key Vault Managed Hardware Security Module %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_managed_hardware_security_module", *existing.ID)
		}
	}

	tenantId := uuid.FromStringOrNil(d.Get("tenant_id").(string))
	parameters := keyvault.ManagedHsm{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &keyvault.ManagedHsmProperties{
			TenantID:              &tenantId,
			InitialAdminObjectIds: utils.ExpandStringSlice(d.Get("admin_object_ids").(*schema.Set).List()),
			CreateMode:            keyvault.CreateModeDefault,
			// the rest api will report error when setting to false. Fixed it to true.
			EnableSoftDelete:          utils.Bool(true),
			SoftDeleteRetentionInDays: utils.Int32(int32(d.Get("soft_delete_retention_days").(int))),
			EnablePurgeProtection:     utils.Bool(d.Get("purge_protection_enabled").(bool)),
		},
		Sku: &keyvault.ManagedHsmSku{
			Family: utils.String("B"),
			Name:   keyvault.ManagedHsmSkuName(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Key Vault Managed Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Key Vault Managed Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Key Vault Managed Hardware Security Module %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Key Vault Managed Hardware Security Module %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	return resourceArmKeyVaultManagedHardwareSecurityModuleRead(d, meta)
}

func resourceArmKeyVaultManagedHardwareSecurityModuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[ERROR] Key Vault Managed Hardware Security Module %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Key Vault Managed Hardware Security Module %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		if tid := props.TenantID; tid != nil {
			d.Set("tenant_id", tid.String())
		}
		d.Set("admin_object_ids", utils.FlattenStringSlice(props.InitialAdminObjectIds))
		d.Set("hsm_uri", props.HsmURI)
		d.Set("soft_delete_retention_days", props.SoftDeleteRetentionInDays)
		d.Set("purge_protection_enabled", props.EnablePurgeProtection)
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", string(sku.Name))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmKeyVaultManagedHardwareSecurityModuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedHardwareSecurityModuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Key Vault Managed Hardware Security Module %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("waiting on deleting future for Key Vault Managed Hardware Security Module %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
