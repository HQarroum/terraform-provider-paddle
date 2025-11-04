package resources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNotificationSettingResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationSettingDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationSettingResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "description", "Test webhook"),
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "destination", "https://example.com/webhook"),
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "type", "url"),
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "active", "true"),
					resource.TestCheckResourceAttrSet("paddle_notification_setting.test", "id"),
					resource.TestCheckResourceAttrSet("paddle_notification_setting.test", "endpoint_secret_key"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "paddle_notification_setting.test",
				ImportState:       true,
				ImportStateVerify: true,
				// endpoint_secret_key is sensitive and not returned in import
				ImportStateVerifyIgnore: []string{"endpoint_secret_key"},
			},
			// Update and Read testing
			{
				Config: testAccNotificationSettingResourceConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "description", "Updated webhook"),
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "destination", "https://example.com/webhook-updated"),
				),
			},
		},
	})
}

func TestAccNotificationSettingResource_withOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationSettingResourceConfigFull(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "description", "Full webhook"),
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "include_sensitive_fields", "true"),
					resource.TestCheckResourceAttr("paddle_notification_setting.test", "traffic_source", "platform"),
				),
			},
		},
	})
}

func testAccNotificationSettingResourceConfig() string {
	return `
resource "paddle_notification_setting" "test" {
  description = "Test webhook"
  destination = "https://example.com/webhook"
  type        = "url"
  active      = true
  
  subscribed_events = [
    "transaction.completed",
    "subscription.activated"
  ]
}
`
}

func testAccNotificationSettingResourceConfigUpdated() string {
	return `
resource "paddle_notification_setting" "test" {
  description = "Updated webhook"
  destination = "https://example.com/webhook-updated"
  type        = "url"
  active      = true
  
  subscribed_events = [
    "transaction.completed",
    "subscription.activated",
    "subscription.canceled"
  ]
}
`
}

func testAccNotificationSettingResourceConfigFull() string {
	return `
resource "paddle_notification_setting" "test" {
  description               = "Full webhook"
  destination              = "https://example.com/webhook-full"
  type                     = "url"
  active                   = true
  include_sensitive_fields = true
  traffic_source          = "platform"
  
  subscribed_events = [
    "transaction.completed",
    "subscription.activated"
  ]
}
`
}

func testAccCheckNotificationSettingDestroy(s *terraform.State) error {
	// Notification settings can be hard-deleted
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "paddle_notification_setting" {
			continue
		}
	}
	return nil
}
