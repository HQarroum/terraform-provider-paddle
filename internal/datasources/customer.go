package datasources

import (
	"context"
	"fmt"

	"github.com/HQarroum/terraform-provider-paddle/internal/helpers"
	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &CustomerDataSource{}

// Creates a new customer data source.
func NewCustomerDataSource() datasource.DataSource {
	return &CustomerDataSource{}
}

type CustomerDataSource struct {
	client *paddle.SDK
}

type customerDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Email            types.String `tfsdk:"email"`
	MarketingConsent types.Bool   `tfsdk:"marketing_consent"`
	Status           types.String `tfsdk:"status"`
	CustomData       types.Map    `tfsdk:"custom_data"`
	Locale           types.String `tfsdk:"locale"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

func (d *CustomerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customer"
}

func (d *CustomerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves details of a specific Paddle customer by ID.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Paddle customer ID (format: ctm_...).",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Full name of this customer.",
			},
			"email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Email address for this customer.",
			},
			"marketing_consent": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this customer opted into marketing.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this customer can be used.",
			},
			"custom_data": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom key-value data.",
			},
			"locale": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "IETF BCP 47 locale tag.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the customer was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the customer was last updated.",
			},
		},
	}
}

// Configure sets the Paddle client for the data source.
func (d *CustomerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*paddle.SDK)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *paddle.SDK, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read retrieves the customer details from Paddle and
// sets the data source state.
func (d *CustomerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data customerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	customer, err := d.client.GetCustomer(ctx, &paddle.GetCustomerRequest{
		CustomerID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading customer",
			fmt.Sprintf("Could not read customer ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	data.Email = types.StringValue(customer.Email)
	data.MarketingConsent = types.BoolValue(customer.MarketingConsent)
	data.Status = types.StringValue(string(customer.Status))
	data.Locale = types.StringValue(customer.Locale)
	data.CreatedAt = types.StringValue(customer.CreatedAt)
	data.UpdatedAt = types.StringValue(customer.UpdatedAt)

	if customer.Name != nil {
		data.Name = types.StringValue(*customer.Name)
	} else {
		data.Name = types.StringNull()
	}

	if customer.CustomData != nil {
		customDataMap, err := helpers.CustomDataToMap(customer.CustomData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error processing custom data",
				fmt.Sprintf("Could not convert custom_data for customer %s: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}
		customDataValue, diags := types.MapValue(types.StringType, customDataMap)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomData = customDataValue
	} else {
		data.CustomData = types.MapNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
