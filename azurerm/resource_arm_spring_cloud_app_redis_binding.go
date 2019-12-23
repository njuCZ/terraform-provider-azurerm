package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudAppRedisBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudAppRedisBindingCreate,
		Read:   resourceArmSpringCloudAppRedisBindingRead,
		Update: resourceArmSpringCloudAppRedisBindingUpdate,
		Delete: resourceArmSpringCloudAppRedisBindingDelete,

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

			"name": {
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

			"spring_cloud_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azappplatform.ValidateSpringCloudName,
			},

			// the return key is like: abc******
			"key": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: azure.SuppressKeyDiff,
			},

			"redis_cache_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"ssl_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"generated_properties": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmSpringCloudAppRedisBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Spring Cloud App Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app_redis_binding", *existing.ID)
		}
	}

	key := d.Get("key").(string)
	redisCacheID := d.Get("redis_cache_id").(string)
	enableSSL := d.Get("ssl_enabled").(bool)
	bindingParameters := map[string]interface{}{
		"useSsl": enableSSL,
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: bindingParameters,
			Key:               utils.String(key),
			ResourceID:        utils.String(redisCacheID),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, springCloudName, appName, name, &bindingResource); err != nil {
		return fmt.Errorf("Error creating Spring Cloud App Redis Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Spring Cloud App Redis Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Spring Cloud App Redis Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q) ID", name, appName, springCloudName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmSpringCloudAppRedisBindingRead(d, meta)
}

func resourceArmSpringCloudAppRedisBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]
	name := id.Path["bindings"]

	resp, err := client.Get(ctx, resourceGroup, springCloudName, appName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App Binding %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud App Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("spring_cloud_name", springCloudName)
	d.Set("spring_cloud_app_name", appName)
	if bindingResourceProperties := resp.Properties; bindingResourceProperties != nil {
		d.Set("generated_properties", bindingResourceProperties.GeneratedProperties)
		if bindingResourceProperties.GeneratedProperties != nil {
			d.Set("key", azure.FindValueInGeneratedProperties(*bindingResourceProperties.GeneratedProperties, "spring.redis.password", "\n"))
		}
		d.Set("redis_cache_id", bindingResourceProperties.ResourceID)
		d.Set("resource_name", bindingResourceProperties.ResourceName)
		d.Set("resource_type", bindingResourceProperties.ResourceType)
		if enableSSL, ok := bindingResourceProperties.BindingParameters["useSsl"]; ok {
			d.Set("ssl_enabled", strings.EqualFold(enableSSL.(string), "true"))
		}
	}
	return nil
}

func resourceArmSpringCloudAppRedisBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	appName := d.Get("spring_cloud_app_name").(string)
	springCloudName := d.Get("spring_cloud_name").(string)

	key := d.Get("key").(string)
	redisCacheID := d.Get("redis_cache_id").(string)
	enableSSL := d.Get("ssl_enabled").(bool)
	bindingParameters := map[string]interface{}{
		"useSsl": enableSSL,
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: bindingParameters,
			Key:               utils.String(key),
			ResourceID:        utils.String(redisCacheID),
		},
	}

	if _, err := client.Update(ctx, resourceGroup, springCloudName, appName, name, &bindingResource); err != nil {
		return fmt.Errorf("Error updating Spring Cloud App Redis Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	return resourceArmSpringCloudAppRedisBindingRead(d, meta)
}

func resourceArmSpringCloudAppRedisBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]
	appName := id.Path["apps"]
	name := id.Path["bindings"]

	if _, err := client.Delete(ctx, resourceGroup, springCloudName, appName, name); err != nil {
		return fmt.Errorf("Error deleting Spring Cloud App Redis Binding %q (App Name %q / Spring Cloud Name %q / Resource Group %q): %+v", name, appName, springCloudName, resourceGroup, err)
	}

	return nil
}
