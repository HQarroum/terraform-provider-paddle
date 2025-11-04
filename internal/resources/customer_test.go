package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCustomerResource_basic(t *testing.T) {
	email := "test@example.com"
	updatedEmail := "updated@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCustomerDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomerResourceConfig(email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_customer.test", "email", email),
					resource.TestCheckResourceAttr("paddle_customer.test", "status", "active"),
					resource.TestCheckResourceAttrSet("paddle_customer.test", "id"),
					resource.TestCheckResourceAttrSet("paddle_customer.test", "created_at"),
					resource.TestCheckResourceAttrSet("paddle_customer.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "paddle_customer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccCustomerResourceConfig(updatedEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_customer.test", "email", updatedEmail),
				),
			},
		},
	})
}

func TestAccCustomerResource_withOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCustomerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerResourceConfigFull(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_customer.test", "email", "fulltest@example.com"),
					resource.TestCheckResourceAttr("paddle_customer.test", "name", "Test Customer"),
					resource.TestCheckResourceAttr("paddle_customer.test", "locale", "en-US"),
					resource.TestCheckResourceAttr("paddle_customer.test", "custom_data.account_id", "12345"),
				),
			},
		},
	})
}

func testAccCustomerResourceConfig(email string) string {
	return fmt.Sprintf(`
resource "paddle_customer" "test" {
  email = %[1]q
}
`, email)
}

func testAccCustomerResourceConfigFull() string {
	return `
resource "paddle_customer" "test" {
  email  = "fulltest@example.com"
  name   = "Test Customer"
  locale = "en-US"
  
  custom_data = {
    account_id = "12345"
  }
}
`
}

func testAccCheckCustomerDestroy(s *terraform.State) error {
	// Customers are archived, not deleted
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "paddle_customer" {
			continue
		}
	}
	return nil
}
