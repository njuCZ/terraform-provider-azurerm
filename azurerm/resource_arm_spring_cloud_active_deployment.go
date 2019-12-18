package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudActiveDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudActiveDeploymentCreateUpdate,
		Read:   resourceArmSpringCloudActiveDeploymentRead,
		Update: resourceArmSpringCloudActiveDeploymentCreateUpdate,
		Delete: resourceArmSpringCloudActiveDeploymentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"spring_cloud_app_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"spring_cloud_deployment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			"persistent_disk": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "/persistent",
							ValidateFunc: azappplatform.ValidateMountPath,
						},
						"size_in_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntInSlice([]int{0, 50}),
						},
					},
				},
			},

			"public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"temporary_disk": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "/tmp",
							ValidateFunc: azappplatform.ValidateMountPath,
						},
						"size_in_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(0, 5),
						},
					},
				},
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSpringCloudActiveDeploymentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	springCloudName := d.Get("spring_cloud_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	deploymentName := d.Get("spring_cloud_deployment_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("[DEBUG] Spring Cloud App %q (Spring Cloud service %q / resource group %q) was not found.", appName, springCloudName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Spring Cloud App %q (Spring Cloud service %q / resource group %q): %+v", appName, springCloudName, resourceGroup, err)
	}

	persistentDisk := d.Get("persistent_disk").([]interface{})
	public := d.Get("public").(bool)
	temporaryDisk := d.Get("temporary_disk").([]interface{})

	resp.Properties = &appplatform.AppResourceProperties{
		ActiveDeploymentName: &deploymentName,
		PersistentDisk:       expandArmSpringCloudAppPersistentDisk(persistentDisk),
		Public:               utils.Bool(public),
		TemporaryDisk:        expandArmSpringCloudAppTemporaryDisk(temporaryDisk),
	}

	future, err := client.Update(ctx, resourceGroup, springCloudName, appName, &resp)
	if err != nil {
		return fmt.Errorf("Error swapping active deployment %q (Spring Cloud Service %q / Spring Cloud App %q / Resource Group %q): %+v", deploymentName, springCloudName, appName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error swapping active deployment %q (Spring Cloud Service %q / Spring Cloud App %q / Resource Group %q): %+v", deploymentName, springCloudName, appName, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudActiveDeploymentRead(d, meta)
}

func resourceArmSpringCloudActiveDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", appName, springCloudName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", resp.Name)

	if appResourceProperties := resp.Properties; appResourceProperties != nil {
		if err := d.Set("persistent_disk", flattenArmSpringCloudAppPersistentDisk(appResourceProperties.PersistentDisk)); err != nil {
			return fmt.Errorf("Error setting `persistent_disk`: %+v", err)
		}
		if err := d.Set("temporary_disk", flattenArmSpringCloudAppTemporaryDisk(appResourceProperties.TemporaryDisk)); err != nil {
			return fmt.Errorf("Error setting `temporary_disk`: %+v", err)
		}
		d.Set("public", appResourceProperties.Public)
		d.Set("url", appResourceProperties.URL)
		d.Set("spring_cloud_deployment_name", appResourceProperties.ActiveDeploymentName)
	}

	return nil
}

func resourceArmSpringCloudActiveDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	// There is nothing to delete and the server side can not update app active_deployment_name to empty
	// so return nil
	return nil
}

func expandArmSpringCloudAppPersistentDisk(input []interface{}) *appplatform.PersistentDisk {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	sizeInGb := v["size_in_gb"].(int)
	mountPath := v["mount_path"].(string)

	result := appplatform.PersistentDisk{
		MountPath: utils.String(mountPath),
		SizeInGB:  utils.Int32(int32(sizeInGb)),
	}
	return &result
}

func expandArmSpringCloudAppTemporaryDisk(input []interface{}) *appplatform.TemporaryDisk {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	sizeInGb := v["size_in_gb"].(int)
	mountPath := v["mount_path"].(string)

	result := appplatform.TemporaryDisk{
		MountPath: utils.String(mountPath),
		SizeInGB:  utils.Int32(int32(sizeInGb)),
	}
	return &result
}

func flattenArmSpringCloudAppPersistentDisk(input *appplatform.PersistentDisk) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if mountPath := input.MountPath; mountPath != nil {
		result["mount_path"] = *mountPath
	}
	if sizeInGb := input.SizeInGB; sizeInGb != nil {
		result["size_in_gb"] = int(*sizeInGb)
	}

	return []interface{}{result}
}

func flattenArmSpringCloudAppTemporaryDisk(input *appplatform.TemporaryDisk) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if mountPath := input.MountPath; mountPath != nil {
		result["mount_path"] = *mountPath
	}
	if sizeInGb := input.SizeInGB; sizeInGb != nil {
		result["size_in_gb"] = int(*sizeInGb)
	}

	return []interface{}{result}
}
