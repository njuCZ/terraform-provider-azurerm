package kusto

import (
	"fmt"
	"log"
	"time"

	dataplaneKusto "github.com/Azure/azure-kusto-go/kusto"
	"github.com/Azure/azure-kusto-go/kusto/data/table"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceArmKustoTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoTableCreateUpdate,
		Read:   resourceArmKustoTableRead,
		Update: resourceArmKustoTableCreateUpdate,
		Delete: resourceArmKustoTableDelete,

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: validateAzureRMKustoTableName,
			},

			"kusto_database_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KustoDatabaseID,
			},

			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},

			"folder": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"doc": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArmKustoTableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Database creation.")

	name := d.Get("name").(string)
	kustoDatabaseId, err := parse.KustoDatabaseID(d.Get("kusto_database_id").(string))
	if err != nil {
		return err
	}

	cluster, err := client.Get(ctx, kustoDatabaseId.ResourceGroup, kustoDatabaseId.Cluster)
	if err != nil {
		return fmt.Errorf("retrieving kusto cluster %q (Resource Group %q): %+v", kustoDatabaseId.Cluster, kustoDatabaseId.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil || cluster.ClusterProperties.URI == nil {
		return fmt.Errorf("kusto cluster %q (Resource Group %q) uri property is nil or empty", kustoDatabaseId.Cluster, kustoDatabaseId.ResourceGroup)
	}

	dataplaneClient, err := meta.(*clients.Client).Kusto.NewDataPlaneClient(*cluster.URI)
	if err != nil {
		return fmt.Errorf("init kusto data plane client: %+v", err)
	}

	if d.IsNewResource() {
		rowIter, err := dataplaneClient.Mgmt(ctx, kustoDatabaseId.Name, dataplaneKusto.NewStmt(".show tables"))
		if err != nil {
			return fmt.Errorf("Error listing tables in Kusto Database %q (Resource Group %q, Cluster %q): %s", kustoDatabaseId.Name, kustoDatabaseId.ResourceGroup, kustoDatabaseId.Cluster, err)
		}
		defer rowIter.Stop()

		found := false
		err = rowIter.Do(func(row *table.Row) error {
			if len(row.Values) > 0 && row.Values[0].String() == name {
				found = true
				return nil
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error iterating tables in Kusto Database %q (Resource Group %q, Cluster %q): %s", kustoDatabaseId.Name, kustoDatabaseId.ResourceGroup, kustoDatabaseId.Cluster, err)
		}

		if found {
			id := fmt.Sprintf("%s|%s|%s", *cluster.URI, kustoDatabaseId.Name, name)
			return tf.ImportAsExistsError("azurerm_kusto_table", id)
		}
	}

	sql := fmt.Sprintf(".create-merge table %s (%s)", name, d.Get("schema").(string))
	folder := d.Get("folder").(string)
	doc := d.Get("doc").(string)
	if folder != "" || doc != "" {
		sql = sql + " with ("
		innerSql := ""
		if doc != "" {
			innerSql = "docstring=" + doc
		}
		if folder != "" {
			if innerSql != "" {
				innerSql = innerSql + ","
			}
			innerSql = "folder=" + folder
		}
		sql = sql + innerSql + ")"
	}
	rowIter, err := dataplaneClient.Mgmt(ctx, kustoDatabaseId.Name, dataplaneKusto.NewStmt(".show tables"))
	if err != nil {
		return fmt.Errorf("Error listing tables in Kusto Database %q (Resource Group %q, Cluster %q): %s", kustoDatabaseId.Name, kustoDatabaseId.ResourceGroup, kustoDatabaseId.Cluster, err)
	}
	defer rowIter.Stop()

	// future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, name, readWriteDatabase)
	// if err != nil {
	// 	return fmt.Errorf("Error creating or updating Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	// }

	// if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
	// 	return fmt.Errorf("Error waiting for completion of Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	// }

	// resp, err := client.Get(ctx, resourceGroup, clusterName, name)
	// if err != nil {
	// 	return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	// }
	// if resp.Value == nil {
	// 	return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): Invalid resource response", name, resourceGroup, clusterName)
	// }

	// database, ok := resp.Value.AsReadWriteDatabase()
	// if !ok {
	// 	return fmt.Errorf("Resource is not a Read/Write Database %q (Resource Group %q, Cluster %q)", name, resourceGroup, clusterName)
	// }
	// if database.ID == nil || *database.ID == "" {
	// 	return fmt.Errorf("Cannot read ID for Kusto Database %q (Resource Group %q, Cluster %q)", name, resourceGroup, clusterName)
	// }

	// d.SetId(*database.ID)

	return resourceArmKustoTableRead(d, meta)
}

func resourceArmKustoTableRead(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*clients.Client).Kusto.DatabasesClient
	// ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	// defer cancel()

	// id, err := parse.KustoTableID(d.Id())
	// if err != nil {
	// 	return err
	// }

	// resp, err := client.Get(ctx, id.ResourceGroup, id.Cluster, id.Name)
	// if err != nil {
	// 	if utils.ResponseWasNotFound(resp.Response) {
	// 		d.SetId("")
	// 		return nil
	// 	}
	// 	return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", id.Name, id.ResourceGroup, id.Cluster, err)
	// }

	// if resp.Value == nil {
	// 	return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): Invalid resource response", id.Name, id.ResourceGroup, id.Cluster)
	// }

	// database, ok := resp.Value.AsReadWriteDatabase()
	// if !ok {
	// 	return fmt.Errorf("Existing resource is not a Read/Write Database (Resource Group %q, Cluster %q): %q", id.ResourceGroup, id.Cluster, id.Name)
	// }

	// d.Set("name", id.Name)
	// d.Set("resource_group_name", id.ResourceGroup)
	// d.Set("cluster_name", id.Cluster)

	// if location := database.Location; location != nil {
	// 	d.Set("location", azure.NormalizeLocation(*location))
	// }

	// if props := database.ReadWriteDatabaseProperties; props != nil {
	// 	d.Set("hot_cache_period", props.HotCachePeriod)
	// 	d.Set("soft_delete_period", props.SoftDeletePeriod)

	// 	if statistics := props.Statistics; statistics != nil {
	// 		d.Set("size", statistics.Size)
	// 	}
	// }

	return nil
}

func resourceArmKustoTableDelete(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*clients.Client).Kusto.DatabasesClient
	// ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	// defer cancel()

	// id, err := azure.ParseAzureResourceID(d.Id())
	// if err != nil {
	// 	return err
	// }

	// resGroup := id.ResourceGroup
	// clusterName := id.Path["Clusters"]
	// name := id.Path["Databases"]

	// future, err := client.Delete(ctx, resGroup, clusterName, name)
	// if err != nil {
	// 	return fmt.Errorf("Error deleting Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resGroup, clusterName, err)
	// }

	// if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
	// 	return fmt.Errorf("Error waiting for deletion of Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resGroup, clusterName, err)
	// }

	return nil

}
