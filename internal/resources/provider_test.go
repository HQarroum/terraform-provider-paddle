package resources_test

import (
	"os"
	"testing"

	"github.com/HQarroum/terraform-provider-paddle/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Used to instantiate a provider during acceptance testing.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"paddle": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// Ensures required environment variables are set for acceptance tests.
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PADDLE_API_KEY"); v == "" {
		t.Fatal("PADDLE_API_KEY must be set for acceptance tests")
	}
}
