package datasources

import (
	"context"
	"fmt"

	"github.com/HQarroum/terraform-provider-paddle/internal/helpers"
	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &PriceDataSource{}

// Creates a new price data source.
func NewPriceDataSource() datasource.DataSource {
	return &PriceDataSource{}
}

// A single Paddle price.
type PriceDataSource struct {
	client *paddle.SDK
}

type priceDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ProductID          types.String `tfsdk:"product_id"`
	Description        types.String `tfsdk:"description"`
	Name               types.String `tfsdk:"name"`
	TaxMode            types.String `tfsdk:"tax_mode"`
	UnitPrice          types.Object `tfsdk:"unit_price"`
	UnitPriceOverrides types.List   `tfsdk:"unit_price_overrides"`
	BillingCycle       types.Object `tfsdk:"billing_cycle"`
	TrialPeriod        types.Object `tfsdk:"trial_period"`
	Quantity           types.Object `tfsdk:"quantity"`
	CustomData         types.Map    `tfsdk:"custom_data"`
	Status             types.String `tfsdk:"status"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *PriceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_price"
}

// Schema returns the data source schema.
func (d *PriceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves details of a specific Paddle price by ID.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Paddle price ID (format: pri_...).",
			},
			"product_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Paddle product ID that this price is for.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal description for the price.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of this price shown to customers.",
			},
			"tax_mode": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "How tax is calculated for this price.",
			},
			"unit_price": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Base price.",
				Attributes: map[string]schema.Attribute{
					"amount": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Amount in cents as a string.",
					},
					"currency_code": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Three-letter ISO 4217 currency code.",
					},
				},
			},
			"unit_price_overrides": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Country-specific price overrides.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"country_codes": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of country codes.",
						},
						"unit_price": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Overridden price.",
							Attributes: map[string]schema.Attribute{
								"amount": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Amount in cents.",
								},
								"currency_code": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Currency code.",
								},
							},
						},
					},
				},
			},
			"quantity": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Purchase quantity limits.",
				Attributes: map[string]schema.Attribute{
					"minimum": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Minimum quantity.",
					},
					"maximum": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Maximum quantity.",
					},
				},
			},
			"billing_cycle": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "How often this price should be charged. null for one-time prices.",
				Attributes: map[string]schema.Attribute{
					"frequency": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Billing frequency.",
					},
					"interval": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Billing interval unit.",
					},
				},
			},
			"trial_period": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Trial period for the price.",
				Attributes: map[string]schema.Attribute{
					"frequency": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Trial period length.",
					},
					"interval": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Trial period unit.",
					},
				},
			},
			"custom_data": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom data for this price.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the price (active or archived).",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the price was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the price was last updated.",
			},
		},
	}
}

// Configure initializes the data source with the Paddle SDK client.
func (d *PriceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves the price data from Paddle API.
func (d *PriceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data priceDataSourceModel

	// Read configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get price from Paddle API
	price, err := d.client.GetPrice(ctx, &paddle.GetPriceRequest{
		PriceID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading price",
			fmt.Sprintf("Could not read price ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map response to model
	data.ProductID = types.StringValue(price.ProductID)
	data.Description = types.StringValue(price.Description)
	data.Status = types.StringValue(string(price.Status))
	data.CreatedAt = types.StringValue(price.CreatedAt)
	data.UpdatedAt = types.StringValue(price.UpdatedAt)

	if price.Name != nil {
		data.Name = types.StringValue(*price.Name)
	} else {
		data.Name = types.StringNull()
	}

	// Map tax_mode
	data.TaxMode = types.StringValue(string(price.TaxMode))

	// Map unit_price
	unitPriceAttrTypes := map[string]attr.Type{
		"amount":        types.StringType,
		"currency_code": types.StringType,
	}
	unitPriceObj, diags := types.ObjectValue(unitPriceAttrTypes, map[string]attr.Value{
		"amount":        types.StringValue(price.UnitPrice.Amount),
		"currency_code": types.StringValue(string(price.UnitPrice.CurrencyCode)),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.UnitPrice = unitPriceObj

	// Map unit_price_overrides if present
	if len(price.UnitPriceOverrides) > 0 {
		overrideElements := []attr.Value{}
		for _, override := range price.UnitPriceOverrides {
			// Convert country codes to string list
			countryCodeValues := make([]attr.Value, len(override.CountryCodes))
			for i, code := range override.CountryCodes {
				countryCodeValues[i] = types.StringValue(string(code))
			}
			countryCodesList, diags := types.ListValue(types.StringType, countryCodeValues)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// Create unit price object
			overrideUnitPriceObj, diags := types.ObjectValue(unitPriceAttrTypes, map[string]attr.Value{
				"amount":        types.StringValue(override.UnitPrice.Amount),
				"currency_code": types.StringValue(string(override.UnitPrice.CurrencyCode)),
			})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// Create override object
			overrideAttrTypes := map[string]attr.Type{
				"country_codes": types.ListType{ElemType: types.StringType},
				"unit_price":    types.ObjectType{AttrTypes: unitPriceAttrTypes},
			}
			overrideObj, diags := types.ObjectValue(overrideAttrTypes, map[string]attr.Value{
				"country_codes": countryCodesList,
				"unit_price":    overrideUnitPriceObj,
			})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			overrideElements = append(overrideElements, overrideObj)
		}

		overrideAttrTypes := map[string]attr.Type{
			"country_codes": types.ListType{ElemType: types.StringType},
			"unit_price":    types.ObjectType{AttrTypes: unitPriceAttrTypes},
		}
		overridesList, diags := types.ListValue(types.ObjectType{AttrTypes: overrideAttrTypes}, overrideElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.UnitPriceOverrides = overridesList
	} else {
		overrideAttrTypes := map[string]attr.Type{
			"country_codes": types.ListType{ElemType: types.StringType},
			"unit_price":    types.ObjectType{AttrTypes: unitPriceAttrTypes},
		}
		data.UnitPriceOverrides = types.ListNull(types.ObjectType{AttrTypes: overrideAttrTypes})
	}

	// Map quantity if present
	quantityAttrTypes := map[string]attr.Type{
		"minimum": types.Int64Type,
		"maximum": types.Int64Type,
	}
	if price.Quantity.Minimum > 0 || price.Quantity.Maximum > 0 {
		quantityObj, diags := types.ObjectValue(quantityAttrTypes, map[string]attr.Value{
			"minimum": types.Int64Value(int64(price.Quantity.Minimum)),
			"maximum": types.Int64Value(int64(price.Quantity.Maximum)),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Quantity = quantityObj
	} else {
		data.Quantity = types.ObjectNull(quantityAttrTypes)
	}

	// Map billing_cycle if present
	billingCycleAttrTypes := map[string]attr.Type{
		"frequency": types.Int64Type,
		"interval":  types.StringType,
	}
	if price.BillingCycle != nil {
		billingCycleObj, diags := types.ObjectValue(billingCycleAttrTypes, map[string]attr.Value{
			"frequency": types.Int64Value(int64(price.BillingCycle.Frequency)),
			"interval":  types.StringValue(string(price.BillingCycle.Interval)),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.BillingCycle = billingCycleObj
	} else {
		data.BillingCycle = types.ObjectNull(billingCycleAttrTypes)
	}

	// Map trial_period if present
	trialPeriodAttrTypes := map[string]attr.Type{
		"frequency": types.Int64Type,
		"interval":  types.StringType,
	}
	if price.TrialPeriod != nil {
		trialPeriodObj, diags := types.ObjectValue(trialPeriodAttrTypes, map[string]attr.Value{
			"frequency": types.Int64Value(int64(price.TrialPeriod.Frequency)),
			"interval":  types.StringValue(string(price.TrialPeriod.Interval)),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.TrialPeriod = trialPeriodObj
	} else {
		data.TrialPeriod = types.ObjectNull(trialPeriodAttrTypes)
	}

	// Map custom_data if present
	if price.CustomData != nil {
		customDataMap, err := helpers.CustomDataToMap(price.CustomData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error processing custom data",
				fmt.Sprintf("Could not convert custom_data for price %s: %s", data.ID.ValueString(), err.Error()),
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
