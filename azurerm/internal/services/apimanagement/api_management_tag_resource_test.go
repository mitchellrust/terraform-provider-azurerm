package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementTagResource struct {
}

func TestAccApiManagementTag_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_tag", "test")
	r := ApiManagementTagResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Tag"),
				check.That(data.ResourceName).Key("tag_id").HasValue("test-tag"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementTag_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_tag", "test")
	r := ApiManagementTagResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementTag_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_tag", "test")
	r := ApiManagementTagResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Tag"),
				check.That(data.ResourceName).Key("tag_id").HasValue("test-tag"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Updated Tag"),
				check.That(data.ResourceName).Key("tag_id").HasValue("test-tag"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Tag"),
				check.That(data.ResourceName).Key("tag_id").HasValue("test-tag"),
			),
		},
	})
}

func TestAccApiManagementTag_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")
	r := ApiManagementTagResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Tag"),
				check.That(data.ResourceName).Key("tag_id").HasValue("test-tag"),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementTagResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	tagID := id.Path["tags"]

	resp, err := clients.ApiManagement.TagClient.Get(ctx, resourceGroup, serviceName, tagID)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Tag (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementTagResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_tag" "test" {
  tag_id	            = "test-tag"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Tag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementTagResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_tag" "import" {
  tag_id	            = azurerm_api_management_product.test.tag_id
  api_management_name   = azurerm_api_management_product.test.api_management_name
  resource_group_name   = azurerm_api_management_product.test.resource_group_name
  display_name          = azurerm_api_management_product.test.display_name
}
`, r.basic(data))
}

func (ApiManagementTagResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_tag" "test" {
  product_id            = "test-tag"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Updated Tag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementTagResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_tag" "test" {
  tag_id	            = "test-tag"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Tag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
