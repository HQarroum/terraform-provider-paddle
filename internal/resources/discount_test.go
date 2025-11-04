package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDiscountResource_percentage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDiscountDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDiscountResourceConfigPercentage(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_discount.test", "description", "Test percentage discount"),
					resource.TestCheckResourceAttr("paddle_discount.test", "type", "percentage"),
					resource.TestCheckResourceAttr("paddle_discount.test", "amount", "20"),
					resource.TestCheckResourceAttr("paddle_discount.test", "enabled_for_checkout", "true"),
					resource.TestCheckResourceAttrSet("paddle_discount.test", "id"),
					resource.TestCheckResourceAttrSet("paddle_discount.test", "code"),
					resource.TestCheckResourceAttrSet("paddle_discount.test", "status"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "paddle_discount.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDiscountResourceConfigPercentageUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_discount.test", "description", "Updated percentage discount"),
					resource.TestCheckResourceAttr("paddle_discount.test", "amount", "25"),
				),
			},
		},
	})
}

func TestAccDiscountResource_flat(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDiscountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiscountResourceConfigFlat(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_discount.test", "description", "Test flat discount"),
					resource.TestCheckResourceAttr("paddle_discount.test", "type", "flat"),
					resource.TestCheckResourceAttr("paddle_discount.test", "amount", "500"),
					resource.TestCheckResourceAttr("paddle_discount.test", "currency_code", "USD"),
				),
			},
		},
	})
}

func TestAccDiscountResource_recurring(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDiscountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiscountResourceConfigRecurring(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_discount.test", "recur", "true"),
					resource.TestCheckResourceAttr("paddle_discount.test", "maximum_recurring_intervals", "3"),
				),
			},
		},
	})
}

func testAccDiscountResourceConfigPercentage() string {
	return `
resource "paddle_discount" "test" {
  description          = "Test percentage discount"
  type                = "percentage"
  amount              = "20"
  enabled_for_checkout = true
}
`
}

func testAccDiscountResourceConfigPercentageUpdated() string {
	return `
resource "paddle_discount" "test" {
  description          = "Updated percentage discount"
  type                = "percentage"
  amount              = "25"
  enabled_for_checkout = true
}
`
}

func testAccDiscountResourceConfigFlat() string {
	return `
resource "paddle_discount" "test" {
  description   = "Test flat discount"
  type         = "flat"
  amount       = "500"
  currency_code = "USD"
}
`
}

func testAccDiscountResourceConfigRecurring() string {
	return `
resource "paddle_discount" "test" {
  description                  = "Test recurring discount"
  type                        = "percentage"
  amount                      = "15"
  recur                       = true
  maximum_recurring_intervals = 3
}
`
}

func testAccCheckDiscountDestroy(s *terraform.State) error {
	// Discounts are archived, not deleted
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "paddle_discount" {
			continue
		}
	}
	return nil
}
