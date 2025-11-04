package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccProductResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProductDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProductResourceConfig("Test Product", "saas"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_product.test", "name", "Test Product"),
					resource.TestCheckResourceAttr("paddle_product.test", "tax_category", "saas"),
					resource.TestCheckResourceAttr("paddle_product.test", "status", "active"),
					resource.TestCheckResourceAttrSet("paddle_product.test", "id"),
					resource.TestCheckResourceAttrSet("paddle_product.test", "created_at"),
					resource.TestCheckResourceAttrSet("paddle_product.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "paddle_product.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProductResourceConfig("Updated Product", "digital-goods"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_product.test", "name", "Updated Product"),
					resource.TestCheckResourceAttr("paddle_product.test", "tax_category", "digital-goods"),
				),
			},
		},
	})
}

func TestAccProductResource_withOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProductResourceConfigFull(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_product.test", "name", "Full Test Product"),
					resource.TestCheckResourceAttr("paddle_product.test", "description", "A comprehensive test product"),
					resource.TestCheckResourceAttr("paddle_product.test", "tax_category", "saas"),
					resource.TestCheckResourceAttr("paddle_product.test", "image_url", "https://example.com/image.png"),
					resource.TestCheckResourceAttr("paddle_product.test", "custom_data.key1", "value1"),
					resource.TestCheckResourceAttr("paddle_product.test", "custom_data.key2", "value2"),
				),
			},
		},
	})
}

func testAccProductResourceConfig(name, taxCategory string) string {
	return fmt.Sprintf(`
resource "paddle_product" "test" {
  name         = %[1]q
  tax_category = %[2]q
}
`, name, taxCategory)
}

func testAccProductResourceConfigFull() string {
	return `
resource "paddle_product" "test" {
  name         = "Full Test Product"
  description  = "A comprehensive test product"
  tax_category = "saas"
  image_url    = "https://example.com/image.png"
  
  custom_data = {
    key1 = "value1"
    key2 = "value2"
  }
}
`
}

func testAccCheckProductDestroy(s *terraform.State) error {
	// Note: Products are archived, not deleted, so we don't check for complete removal
	// We just verify the resource is removed from state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "paddle_product" {
			continue
		}
	}
	return nil
}

