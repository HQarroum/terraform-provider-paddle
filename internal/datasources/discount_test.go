package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDiscountDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDiscountDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.paddle_discount.test", "description", "Test discount for data source"),
					resource.TestCheckResourceAttr("data.paddle_discount.test", "type", "percentage"),
					resource.TestCheckResourceAttr("data.paddle_discount.test", "amount", "15"),
					resource.TestCheckResourceAttrSet("data.paddle_discount.test", "id"),
				),
			},
		},
	})
}

func testAccDiscountDataSourceConfig() string {
	return `
resource "paddle_discount" "test" {
  description = "Test discount for data source"
  type        = "percentage"
  amount      = "15"
}

data "paddle_discount" "test" {
  id = paddle_discount.test.id
}
`
}
