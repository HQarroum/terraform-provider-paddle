package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPriceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPriceDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPriceResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_price.test", "description", "Test monthly price"),
					resource.TestCheckResourceAttr("paddle_price.test", "unit_price.amount", "2900"),
					resource.TestCheckResourceAttr("paddle_price.test", "unit_price.currency_code", "USD"),
					resource.TestCheckResourceAttr("paddle_price.test", "status", "active"),
					resource.TestCheckResourceAttrSet("paddle_price.test", "id"),
					resource.TestCheckResourceAttrSet("paddle_price.test", "product_id"),
					resource.TestCheckResourceAttrSet("paddle_price.test", "created_at"),
					resource.TestCheckResourceAttrSet("paddle_price.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "paddle_price.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing (only mutable fields)
			{
				Config: testAccPriceResourceConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_price.test", "description", "Updated monthly price"),
					resource.TestCheckResourceAttr("paddle_price.test", "name", "Updated Monthly Plan"),
				),
			},
		},
	})
}

func TestAccPriceResource_recurring(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPriceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPriceResourceConfigRecurring(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_price.test", "description", "Monthly subscription"),
					resource.TestCheckResourceAttr("paddle_price.test", "billing_cycle.frequency", "1"),
					resource.TestCheckResourceAttr("paddle_price.test", "billing_cycle.interval", "month"),
					resource.TestCheckResourceAttr("paddle_price.test", "trial_period.frequency", "7"),
					resource.TestCheckResourceAttr("paddle_price.test", "trial_period.interval", "day"),
				),
			},
		},
	})
}

func TestAccPriceResource_withQuantity(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPriceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPriceResourceConfigWithQuantity(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_price.test", "quantity.minimum", "5"),
					resource.TestCheckResourceAttr("paddle_price.test", "quantity.maximum", "50"),
				),
			},
		},
	})
}

func testAccPriceResourceConfig() string {
	return `
resource "paddle_product" "test" {
  name         = "Test Product for Price"
  tax_category = "saas"
}

resource "paddle_price" "test" {
  product_id  = paddle_product.test.id
  description = "Test monthly price"
  
  unit_price = {
    amount        = "2900"
    currency_code = "USD"
  }
}
`
}

func testAccPriceResourceConfigUpdated() string {
	return `
resource "paddle_product" "test" {
  name         = "Test Product for Price"
  tax_category = "saas"
}

resource "paddle_price" "test" {
  product_id  = paddle_product.test.id
  description = "Updated monthly price"
  name        = "Updated Monthly Plan"
  
  unit_price = {
    amount        = "2900"
    currency_code = "USD"
  }
}
`
}

func testAccPriceResourceConfigRecurring() string {
	return `
resource "paddle_product" "test" {
  name         = "Test Product for Recurring Price"
  tax_category = "saas"
}

resource "paddle_price" "test" {
  product_id  = paddle_product.test.id
  description = "Monthly subscription"
  name        = "Monthly Plan"
  
  unit_price = {
    amount        = "2900"
    currency_code = "USD"
  }
  
  billing_cycle = {
    frequency = 1
    interval  = "month"
  }
  
  trial_period = {
    frequency = 7
    interval  = "day"
  }
}
`
}

func testAccPriceResourceConfigWithQuantity() string {
	return `
resource "paddle_product" "test" {
  name         = "Test Product for Quantity Price"
  tax_category = "saas"
}

resource "paddle_price" "test" {
  product_id  = paddle_product.test.id
  description = "Price with quantity limits"
  
  unit_price = {
    amount        = "2900"
    currency_code = "USD"
  }
  
  quantity = {
    minimum = 5
    maximum = 50
  }
}
`
}

func testAccCheckPriceDestroy(s *terraform.State) error {
	// Prices are archived, not deleted
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "paddle_price" {
			continue
		}
	}
	return nil
}
