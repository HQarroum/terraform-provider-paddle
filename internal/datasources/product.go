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

var _ datasource.DataSource = &ProductDataSource{}

// Creates a new product data source.
func NewProductDataSource() datasource.DataSource {
	return &ProductDataSource{}
}

// A single Paddle product.
type ProductDataSource struct {
	client *paddle.SDK
}

type productDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	TaxCategory types.String `tfsdk:"tax_category"`
	ImageURL    types.String `tfsdk:"image_url"`
	CustomData  types.Map    `tfsdk:"custom_data"`
	Status      types.String `tfsdk:"status"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *ProductDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

// Schema returns the data source schema.
func (d *ProductDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves details of a specific Paddle product by ID.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Paddle product ID (format: pro_...).",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the product.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the product.",
			},
			"tax_category": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Tax category for the product.",
			},
			"image_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "URL of the product image.",
			},
			"custom_data": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom data for this product.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the product (active or archived).",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the product was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the product was last updated.",
			},
		},
	}
}

// Configure initializes the data source with the Paddle SDK client.
func (d *ProductDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves the product data from Paddle API.
func (d *ProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data productDataSourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get product from Paddle API
	product, err := d.client.GetProduct(ctx, &paddle.GetProductRequest{
		ProductID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading product",
			fmt.Sprintf("Could not read product ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to model
	data.Name = types.StringValue(product.Name)
	data.TaxCategory = types.StringValue(string(product.TaxCategory))
	data.Status = types.StringValue(string(product.Status))
	data.CreatedAt = types.StringValue(product.CreatedAt)
	data.UpdatedAt = types.StringValue(product.UpdatedAt)

	if product.Description != nil {
		data.Description = types.StringValue(*product.Description)
	} else {
		data.Description = types.StringNull()
	}

	if product.ImageURL != nil {
		data.ImageURL = types.StringValue(*product.ImageURL)
	} else {
		data.ImageURL = types.StringNull()
	}

	if product.CustomData != nil {
		customDataMap, err := helpers.CustomDataToMap(product.CustomData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error processing custom data",
				fmt.Sprintf("Could not convert custom_data for product %s: %s", data.ID.ValueString(), err.Error()),
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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
