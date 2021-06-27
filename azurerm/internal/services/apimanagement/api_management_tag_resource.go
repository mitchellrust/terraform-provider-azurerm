package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementTag() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementTagCreateUpdate,
		Read:   resourceApiManagementTagRead,
		Update: resourceApiManagementTagCreateUpdate,
		Delete: resourceApiManagementTagDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"tag_id": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementTagCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tagID := d.Get("tag_id").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	displayName := d.Get("display_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, tagID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Tag (API Management Service %q / Resource Group %q): %s", serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_tag", *existing.ID)
		}
	}

	parameters := apimanagement.TagCreateUpdateParameters{
		TagContractProperties: &apimanagement.TagContractProperties{
			DisplayName: utils.String(displayName),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, tagID, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Tag %q (Resource Group %q / API Management Service %q): %+v", tagID, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, tagID)
	if err != nil {
		return fmt.Errorf("retrieving Tag %q (Resource Group %q / API Management Service %q): %+v", tagID, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Tag %q (Resource Group %q / API Management Service %q): %+v", tagID, resourceGroup, serviceName, err)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementAPIPolicyRead(d, meta)
}

func resourceApiManagementTagRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TagID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	tagID := id.Name

	resp, err := client.Get(ctx, resourceGroup, serviceName, tagID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Tag %q was not found in API Management Service %q / Resource Group %q - removing from state!", tagID, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Tag %q (API Management Service %q / Resource Group %q): %+v", tagID, serviceName, resourceGroup, err)
	}

	d.Set("tag_id", tagID)
	d.Set("api_management_name", serviceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.TagContractProperties; props != nil {
		d.Set("display_name", props.DisplayName)
	}

	return nil
}

func resourceApiManagementTagDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TagID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	tagID := id.Name

	log.Printf("[DEBUG] Deleting Tag %q (API Management Service %q / Resource Grouo %q)", tagID, serviceName, resourceGroup)
	resp, err := client.Delete(ctx, resourceGroup, serviceName, tagID, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Tag %q (API Management Service %q / Resource Group %q): %+v", tagID, serviceName, resourceGroup, err)
		}
	}

	return nil
}
