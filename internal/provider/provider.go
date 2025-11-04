package provider

import (
	"context"
	"fmt"
	"os"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HQarroum/terraform-provider-paddle/internal/datasources"
	"github.com/HQarroum/terraform-provider-paddle/internal/resources"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &paddleProvider{}
)

// New creates a new instance of the Paddle provider with the specified version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &paddleProvider{
			version: version,
		}
	}
}

// paddleProvider is the provider implementation.
type paddleProvider struct {
	version string
}

// paddleProviderModel describes the provider data model.
type paddleProviderModel struct {
	ApiKey      types.String `tfsdk:"api_key"`
	Environment types.String `tfsdk:"environment"`
}

// Metadata returns the provider type name and version.
func (p *paddleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "paddle"
	resp.Version = p.version
}

// Schema returns the provider configuration schema.
func (p *paddleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Paddle Billing API to manage products, prices, and notification settings.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Paddle API key for authentication. May also be provided via PADDLE_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"environment": schema.StringAttribute{
				Description: "Paddle environment: 'sandbox' or 'production'. Defaults to 'sandbox'. May also be provided via PADDLE_ENVIRONMENT environment variable.",
				Optional:    true,
			},
		},
	}
}

// Configure initializes the Paddle SDK client using credentials from the provider configuration or environment variables.
func (p *paddleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Paddle client")

	// Retrieve provider data from configuration
	var config paddleProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Paddle API Key",
			"The provider cannot create the Paddle API client as there is an unknown configuration value for the Paddle API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PADDLE_API_KEY environment variable.",
		)
	}

	if config.Environment.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Unknown Paddle Environment",
			"The provider cannot create the Paddle API client as there is an unknown configuration value for the Paddle environment. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PADDLE_ENVIRONMENT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiKey := os.Getenv("PADDLE_API_KEY")
	environment := os.Getenv("PADDLE_ENVIRONMENT")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if !config.Environment.IsNull() {
		environment = config.Environment.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Paddle API Key",
			"The provider cannot create the Paddle API client as there is a missing or empty value for the Paddle API key. "+
				"Set the api_key value in the configuration or use the PADDLE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if environment == "" {
		environment = "sandbox" // Default to sandbox
	}

	if environment != "sandbox" && environment != "production" {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Invalid Paddle Environment",
			fmt.Sprintf("The environment must be either 'sandbox' or 'production', got: %s", environment),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "paddle_environment", environment)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "paddle_api_key")

	// Create Paddle client
	var client *paddle.SDK
	var err error

	if environment == "sandbox" {
		client, err = paddle.NewSandbox(apiKey)
	} else {
		client, err = paddle.New(apiKey)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Paddle API Client",
			"An unexpected error occurred when creating the Paddle API client. "+
				"Paddle Client Error: "+err.Error(),
		)
		return
	}

	// Make the Paddle client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Paddle client", map[string]any{"success": true})
}

// DataSources returns the list of data sources supported by this provider.
func (p *paddleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewProductDataSource,
		datasources.NewPriceDataSource,
		datasources.NewDiscountDataSource,
		datasources.NewCustomerDataSource,
	}
}

// Resources returns the list of resources supported by this provider.
func (p *paddleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewProductResource,
		resources.NewPriceResource,
		resources.NewNotificationSettingResource,
		resources.NewDiscountResource,
		resources.NewCustomerResource,
	}
}
