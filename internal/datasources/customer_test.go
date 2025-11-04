package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomerDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.paddle_customer.test", "email", "datasource@example.com"),
					resource.TestCheckResourceAttr("data.paddle_customer.test", "name", "Data Source Test"),
					resource.TestCheckResourceAttrSet("data.paddle_customer.test", "id"),
					resource.TestCheckResourceAttrSet("data.paddle_customer.test", "status"),
				),
			},
		},
	})
}

func testAccCustomerDataSourceConfig() string {
	return `
resource "paddle_customer" "test" {
  email = "datasource@example.com"
  name  = "Data Source Test"
}

data "paddle_customer" "test" {
  id = paddle_customer.test.id
}
`
}
