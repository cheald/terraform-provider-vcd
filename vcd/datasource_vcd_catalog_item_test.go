// +build catalog ALL functional

package vcd

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Test catalog and catalog item data sources
// Using a catalog data source we reference a catalog item data source
// Using a catalog item data source we create another catalog item
// where the description is the first data source ID
func TestAccVcdCatalogAndItemDatasource(t *testing.T) {
	var TestCatalogItemDS = "TestCatalogItemDS"

	var params = StringMap{
		"Org":             testConfig.VCD.Org,
		"Catalog":         testConfig.VCD.Catalog.Name,
		"CatalogItem":     testConfig.VCD.Catalog.CatalogItem,
		"NewCatalogItem":  TestCatalogItemDS,
		"OvaPath":         testConfig.Ova.OvaPath,
		"UploadPieceSize": testConfig.Ova.UploadPieceSize,
		"UploadProgress":  testConfig.Ova.UploadProgress,
		"Tags":            "catalog",
	}

	configText := templateFill(testAccCheckVcdCatalogItemDS, params)
	if vcdShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}
	debugPrintf("#[DEBUG] CONFIGURATION: %s", configText)

	datasourceCatalog := "data.vcd_catalog." + testConfig.VCD.Catalog.Name
	datasourceCatalogItem := "data.vcd_catalog_item." + testConfig.VCD.Catalog.CatalogItem
	resourceCatalogItem := "vcd_catalog_item." + TestCatalogItemDS
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preRunChecks(t) },
		Providers:    testAccProviders,
		CheckDestroy: catalogItemDestroyed(testConfig.VCD.Catalog.Name, TestCatalogItemDS),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcdCatalogItemExists("vcd_catalog_item."+TestCatalogItemDS),
					resource.TestCheckResourceAttr(
						resourceCatalogItem, "name", TestCatalogItemDS),
					resource.TestCheckResourceAttrPair(
						datasourceCatalog, "name",
						resourceCatalogItem, "catalog"),
					// The description of the new catalog item was created using
					// the ID of the catalog item data source
					resource.TestCheckResourceAttrPair(
						datasourceCatalogItem, "id",
						resourceCatalogItem, "description"),
					resource.TestCheckResourceAttr(
						resourceCatalogItem, "metadata.catalogItem_metadata", "catalogItem Metadata"),
					resource.TestCheckResourceAttr(
						resourceCatalogItem, "metadata.catalogItem_metadata2", "catalogItem Metadata2"),
				),
			},
			resource.TestStep{
				ResourceName:      "vcd_catalog_item." + TestCatalogItemDS + "-import",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdByCatalogItem(TestCatalogItemDS),
				// These fields can't be retrieved from catalog item data
				ImportStateVerifyIgnore: []string{"ova_path", "upload_piece_size", "show_upload_progress"},
			},
		},
	})
}

func catalogItemDestroyed(catalog, itemName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*VCDClient)
		org, err := conn.GetOrgByName(testConfig.VCD.Org)
		if err != nil {
			return err
		}
		cat, err := org.GetCatalogByName(catalog, false)
		if err != nil {
			return err
		}
		_, err = cat.GetCatalogItemByName(itemName, false)
		if err == nil {
			return fmt.Errorf("catalog item %s not deleted", itemName)
		}
		return nil
	}
}

func importStateIdByCatalogItem(objectName string) resource.ImportStateIdFunc {
	return func(*terraform.State) (string, error) {
		importId := testConfig.VCD.Org + "." + testConfig.VCD.Catalog.Name + "." + objectName
		if testConfig.VCD.Org == "" || testConfig.VCD.Catalog.Name == "" || objectName == "" {
			return "", fmt.Errorf("missing information to generate import path: %s", importId)
		}
		return importId, nil
	}
}

const testAccCheckVcdCatalogItemDS = `
data "vcd_catalog" "{{.Catalog}}" {
  org  = "{{.Org}}"
  name = "{{.Catalog}}"
}

data "vcd_catalog_item" "{{.CatalogItem}}" {
  org     = "{{.Org}}"
  catalog = "${data.vcd_catalog.{{.Catalog}}.name}"
  name    = "{{.CatalogItem}}"
}

resource "vcd_catalog_item" "{{.NewCatalogItem}}" {
  org     = "{{.Org}}"
  catalog = "${data.vcd_catalog.{{.Catalog}}.name}"

  name                 = "{{.NewCatalogItem}}"
  description          = "${data.vcd_catalog_item.{{.CatalogItem}}.id}"
  ova_path             = "{{.OvaPath}}"
  upload_piece_size    = {{.UploadPieceSize}}
  show_upload_progress = "{{.UploadProgress}}"

  metadata = {
    catalogItem_metadata = "catalogItem Metadata"
    catalogItem_metadata2 = "catalogItem Metadata2"
  }
}
`
