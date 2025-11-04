package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPriceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPriceDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.paddle_price.test", "description", "Test price for data source"),
					resource.TestCheckResourceAttr("data.paddle_price.test", "unit_price.amount", "1900"),
					resource.TestCheckResourceAttr("data.paddle_price.test", "unit_price.currency_code", "USD"),
					resource.TestCheckResourceAttrSet("data.paddle_price.test", "id"),
					resource.TestCheckResourceAttrSet("data.paddle_price.test", "product_id"),
				),
			},
		},
	})
}

func testAccPriceDataSourceConfig() string {
	return `
resource "paddle_product" "test" {
  name         = "Test Product for Price Data Source"
  tax_category = "saas"
}

resource "paddle_price" "test" {
  product_id  = paddle_product.test.id
  description = "Test price for data source"
  
  unit_price = {
    amount        = "1900"
    currency_code = "USD"
  }
}

data "paddle_price" "test" {
  id = paddle_price.test.id
}
`
}
