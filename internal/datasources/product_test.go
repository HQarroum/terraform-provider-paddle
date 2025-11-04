package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProductDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProductDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.paddle_product.test", "name", "Test Product for Data Source"),
					resource.TestCheckResourceAttr("data.paddle_product.test", "tax_category", "saas"),
					resource.TestCheckResourceAttr("data.paddle_product.test", "status", "active"),
					resource.TestCheckResourceAttrSet("data.paddle_product.test", "id"),
					resource.TestCheckResourceAttrSet("data.paddle_product.test", "created_at"),
				),
			},
		},
	})
}

func testAccProductDataSourceConfig() string {
	return `
resource "paddle_product" "test" {
  name         = "Test Product for Data Source"
  tax_category = "saas"
  description  = "A test product for data source testing"
}

data "paddle_product" "test" {
  id = paddle_product.test.id
}
`
}
