package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultCertificateIssuer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultCertificateIssuerCreateOrUpdate,
		Update: resourceArmKeyVaultCertificateIssuerCreateOrUpdate,
		Read:   resourceArmKeyVaultCertificateIssuerRead,
		Delete: resourceArmKeyVaultCertificateIssuerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceArmKeyVaultChildResourceImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KeyVaultID,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KeyVaultCertificateIssuerName,
			},

			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DigiCert",
					"GlobalSign",
					"OneCertV2-PrivateCA",
					"OneCertV2-PublicCA",
					"SslAdminV2",
				}, false),
			},

			"org_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArmKeyVaultCertificateIssuerCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("failed to look up Certificate Issuer %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	if d.IsNewResource() {
		existing, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("failed to check for presence of existing Certificate Issuer %q (Key Vault %q): %s", name, keyVaultBaseUri, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_certificate_issuer", *existing.ID)
		}
	}

	parameter := keyvault.CertificateIssuerSetParameters{
		Provider:            utils.String(d.Get("provider_name").(string)),
		OrganizationDetails: &keyvault.OrganizationDetails{},
	}

	if orgIdRaw, ok := d.GetOk("org_id"); ok {
		parameter.OrganizationDetails.ID = utils.String(orgIdRaw.(string))
	}

	if adminsRaw, ok := d.GetOk("admin"); ok {
		parameter.OrganizationDetails.AdminDetails = expandKeyVaultCertificateIssuerOrganizationDetailsAdminDetails(adminsRaw.([]interface{}))
	}

	accountId, gotAccountId := d.GetOk("account_id")
	password, gotPassword := d.GetOk("password")

	if gotAccountId && gotPassword {
		parameter.Credentials = &keyvault.IssuerCredentials{
			AccountID: utils.String(accountId.(string)),
			Password:  utils.String(password.(string)),
		}
	}

	_, err = client.SetCertificateIssuer(ctx, keyVaultBaseUri, name, parameter)
	if err != nil {
		return fmt.Errorf("failed to set Certificate Issuer %q (Key Vault %q): %s", name, keyVaultId, err)
	}

	resp, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		return err
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("failure reading Key Vault Certificate Issuer ID for %q", name)
	}
	d.SetId(*resp.ID)

	return resourceArmKeyVaultCertificateIssuerRead(d, meta)
}

func resourceArmKeyVaultCertificateIssuerRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyVaultCertificateIssuerId(d.Id())
	if err != nil {
		return err
	}

	vault, err := keyVaultClient.FindKeyVault(ctx, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if vault == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	resp, err := client.GetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failed to make Read request on Azure KeyVault Certificate Issuer %s: %+v", id.Name, err)
	}

	d.Set("name", id.Name)

	if resp.Provider != nil {
		d.Set("provider_name", resp.Provider)
	}
	if resp.OrganizationDetails != nil {
		if resp.OrganizationDetails.ID != nil {
			d.Set("org_id", resp.OrganizationDetails.ID)
		}
		d.Set("admin", flattenKeyVaultCertificateIssuerAdmins(resp.OrganizationDetails.AdminDetails))
	}
	if resp.Credentials != nil {
		if resp.Credentials.AccountID != nil {
			d.Set("account_id", resp.Credentials.AccountID)
		}
	}

	return nil
}

func resourceArmKeyVaultCertificateIssuerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	keyVaultClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.KeyVaultCertificateIssuerId(d.Id())
	if err != nil {
		return err
	}

	vault, err := keyVaultClient.FindKeyVault(ctx, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if vault == nil {
		log.Printf("[DEBUG] Issuer %q (Key Vault %q) was not found in Key Vault at URI %q - removing from state", id.Name, vault.ID, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	_, err = client.DeleteCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)

	return err
}

func expandKeyVaultCertificateIssuerOrganizationDetailsAdminDetails(vs []interface{}) *[]keyvault.AdministratorDetails {
	results := make([]keyvault.AdministratorDetails, 0, len(vs))

	for _, v := range vs {
		administratorDetails := keyvault.AdministratorDetails{}
		args := v.(map[string]interface{})
		if firstName, ok := args["first_name"]; ok {
			administratorDetails.FirstName = utils.String(firstName.(string))
		}
		if lastName, ok := args["last_name"]; ok {
			administratorDetails.LastName = utils.String(lastName.(string))
		}
		if emailAddress, ok := args["email_address"]; ok {
			administratorDetails.EmailAddress = utils.String(emailAddress.(string))
		}
		if phone, ok := args["phone"]; ok {
			administratorDetails.Phone = utils.String(phone.(string))
		}
		results = append(results, administratorDetails)
	}

	return &results
}

func flattenKeyVaultCertificateIssuerAdmins(input *[]keyvault.AdministratorDetails) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, admin := range *input {
		emailAddress := ""
		if admin.EmailAddress != nil {
			emailAddress = *admin.EmailAddress
		}

		firstName := ""
		if admin.FirstName != nil {
			firstName = *admin.FirstName
		}

		lastName := ""
		if admin.LastName != nil {
			lastName = *admin.LastName
		}

		phone := ""
		if admin.Phone != nil {
			phone = *admin.Phone
		}

		results = append(results, map[string]interface{}{
			"email_address": emailAddress,
			"first_name":    firstName,
			"last_name":     lastName,
			"phone":         phone,
		})
	}

	return results
}
