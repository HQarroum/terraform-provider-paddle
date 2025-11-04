package resources

import (
	"context"
	"fmt"

	"github.com/HQarroum/terraform-provider-paddle/internal/helpers"
	"github.com/HQarroum/terraform-provider-paddle/internal/validators"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PriceResource{}
var _ resource.ResourceWithImportState = &PriceResource{}

// Creates a new Paddle price resource.
func NewPriceResource() resource.Resource {
	return &PriceResource{}
}

// Price resource manages Paddle prices.
type PriceResource struct {
	client *paddle.SDK
}

// Price resource model describes the resource data model.
type priceResourceModel struct {
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

type unitPriceModel struct {
	Amount       types.String `tfsdk:"amount"`
	CurrencyCode types.String `tfsdk:"currency_code"`
}

type billingCycleModel struct {
	Frequency types.Int64  `tfsdk:"frequency"`
	Interval  types.String `tfsdk:"interval"`
}

type trialPeriodModel struct {
	Frequency types.Int64  `tfsdk:"frequency"`
	Interval  types.String `tfsdk:"interval"`
}

type quantityModel struct {
	Minimum types.Int64 `tfsdk:"minimum"`
	Maximum types.Int64 `tfsdk:"maximum"`
}

type unitPriceOverrideModel struct {
	CountryCodes types.List   `tfsdk:"country_codes"`
	UnitPrice    types.Object `tfsdk:"unit_price"`
}

// Metadata returns the resource type name.
func (r *PriceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_price"
}

// Schema returns the resource schema definition.
func (r *PriceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Paddle price resource. Prices determine how much and how often you charge for a product.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Paddle price ID (format: pri_...)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Paddle product ID that this price is for. Cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Internal description for the price, not shown to customers.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Name of this price, shown to customers. Typically describes the billing frequency.",
			},
			"tax_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "How tax is calculated for this price. One of: `account_setting` (use account default), `internal` (tax included in price), `external` (tax added on top of price). Defaults to `account_setting`.",
				Validators: []validator.String{
					validators.TaxModeValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"unit_price": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Base price. Prices are denominated in cents. Cannot be changed after creation - changing this will force a new resource.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"amount": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Amount in cents as a string (e.g., '2900' for $29.00).",
					},
					"currency_code": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Three-letter ISO 4217 currency code in uppercase (e.g., 'USD', 'EUR').",
						Validators: []validator.String{
							validators.CurrencyCodeValidator{},
						},
					},
				},
			},
			"unit_price_overrides": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Country-specific price overrides for regional pricing strategies. Useful for purchasing power parity or local market pricing. Cannot be changed after creation - changing this will force a new resource.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"country_codes": schema.ListAttribute{
							Required:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of two-letter ISO 3166-1 alpha-2 country codes (e.g., ['GB', 'FR', 'DE']).",
						},
						"unit_price": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "Overridden price for these countries.",
							Attributes: map[string]schema.Attribute{
								"amount": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Amount in cents as a string.",
								},
								"currency_code": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Three-letter ISO 4217 currency code.",
								},
							},
						},
					},
				},
			},
			"billing_cycle": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "How often this price should be charged. null for one-time prices. Cannot be changed after creation - changing this will force a new resource.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"frequency": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Billing frequency (e.g., 1 for monthly, 3 for quarterly).",
					},
					"interval": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Billing interval unit. One of: day, week, month, year.",
					},
				},
			},
			"trial_period": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Trial period for the price. Customers will not be charged during trial. Cannot be changed after creation - changing this will force a new resource.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"frequency": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Trial period length (e.g., 7 for 7 days, 1 for 1 month).",
					},
					"interval": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Trial period unit. One of: day, week, month, year.",
					},
				},
			},
			"quantity": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Limits on purchase quantity. Defaults to minimum=1, maximum=100 if omitted. Useful for discount campaigns and limiting quantities.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"minimum": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Minimum quantity that can be purchased. Must be at least 1.",
					},
					"maximum": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Maximum quantity that can be purchased. Must be greater than or equal to minimum.",
					},
				},
			},
			"custom_data": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom data for this price. Max 10 keys, 1KB total.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the price. Either 'active' or 'archived'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the price was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the price was last updated.",
			},
		},
	}
}

// Configure initializes the resource with the Paddle SDK client.
func (r *PriceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*paddle.SDK)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *paddle.SDK, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates a new Paddle price.
func (r *PriceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data priceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse unit_price
	var unitPrice unitPriceModel
	diags := data.UnitPrice.As(ctx, &unitPrice, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build create request
	createReq := &paddle.CreatePriceRequest{
		ProductID:   data.ProductID.ValueString(),
		Description: data.Description.ValueString(),
		UnitPrice: paddle.Money{
			Amount:       unitPrice.Amount.ValueString(),
			CurrencyCode: paddle.CurrencyCode(unitPrice.CurrencyCode.ValueString()),
		},
	}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		createReq.Name = &name
	}

	// Set tax_mode if present
	if !data.TaxMode.IsNull() && !data.TaxMode.IsUnknown() {
		taxMode := paddle.TaxMode(data.TaxMode.ValueString())
		createReq.TaxMode = &taxMode
	}

	// Parse unit_price_overrides if present
	if !data.UnitPriceOverrides.IsNull() && !data.UnitPriceOverrides.IsUnknown() {
		var overrides []unitPriceOverrideModel
		diags := data.UnitPriceOverrides.ElementsAs(ctx, &overrides, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var paddleOverrides []paddle.UnitPriceOverride
		for _, override := range overrides {
			var countryCodes []string
			diags := override.CountryCodes.ElementsAs(ctx, &countryCodes, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// Convert string country codes to paddle.CountryCode
			paddleCountryCodes := make([]paddle.CountryCode, len(countryCodes))
			for i, code := range countryCodes {
				paddleCountryCodes[i] = paddle.CountryCode(code)
			}

			var unitPrice unitPriceModel
			diags = override.UnitPrice.As(ctx, &unitPrice, basetypes.ObjectAsOptions{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			paddleOverrides = append(paddleOverrides, paddle.UnitPriceOverride{
				CountryCodes: paddleCountryCodes,
				UnitPrice: paddle.Money{
					Amount:       unitPrice.Amount.ValueString(),
					CurrencyCode: paddle.CurrencyCode(unitPrice.CurrencyCode.ValueString()),
				},
			})
		}
		createReq.UnitPriceOverrides = paddleOverrides
	}

	// Parse billing_cycle if present
	if !data.BillingCycle.IsNull() {
		var billingCycle billingCycleModel
		diags := data.BillingCycle.As(ctx, &billingCycle, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.BillingCycle = &paddle.Duration{
			Frequency: int(billingCycle.Frequency.ValueInt64()),
			Interval:  paddle.Interval(billingCycle.Interval.ValueString()),
		}
	}

	// Parse trial_period if present
	if !data.TrialPeriod.IsNull() {
		var trialPeriod trialPeriodModel
		diags := data.TrialPeriod.As(ctx, &trialPeriod, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.TrialPeriod = &paddle.Duration{
			Frequency: int(trialPeriod.Frequency.ValueInt64()),
			Interval:  paddle.Interval(trialPeriod.Interval.ValueString()),
		}
	}

	// Parse quantity if present
	if !data.Quantity.IsNull() && !data.Quantity.IsUnknown() {
		var quantity quantityModel
		diags := data.Quantity.As(ctx, &quantity, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Quantity = &paddle.PriceQuantity{
			Minimum: int(quantity.Minimum.ValueInt64()),
			Maximum: int(quantity.Maximum.ValueInt64()),
		}
	}

	// Parse custom_data if present
	if !data.CustomData.IsNull() && !data.CustomData.IsUnknown() {
		customDataStr := make(map[string]string)
		diags := data.CustomData.ElementsAs(ctx, &customDataStr, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.CustomData = helpers.MapToCustomData(customDataStr)
	}

	// Create price via Paddle API
	price, err := r.client.CreatePrice(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating price",
			fmt.Sprintf("Could not create price for product %s (description: %s): %s", data.ProductID.ValueString(), data.Description.ValueString(), err.Error()),
		)
		return
	}

	// Map response to model
	data.ID = types.StringValue(price.ID)
	data.Status = types.StringValue(string(price.Status))
	data.TaxMode = types.StringValue(string(price.TaxMode))
	data.CreatedAt = types.StringValue(price.CreatedAt)
	data.UpdatedAt = types.StringValue(price.UpdatedAt)

	// Read quantity from API response (Paddle may return defaults)
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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read retrieves the current state of a Paddle price.
func (r *PriceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data priceResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read price from Paddle API
	price, err := r.client.GetPrice(ctx, &paddle.GetPriceRequest{
		PriceID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading price",
			fmt.Sprintf("Could not read price ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update model with response data
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

	// Map quantity if present (check if not zero value)
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update modifies an existing Paddle price.
// Note: unit_price, billing_cycle, and trial_period are immutable and will trigger replacement.
func (r *PriceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan priceResourceModel
	var state priceResourceModel

	// Read Terraform plan and state data into the models
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if any mutable fields actually changed
	hasChanges := !plan.Description.Equal(state.Description) ||
		!plan.Name.Equal(state.Name) ||
		!plan.TaxMode.Equal(state.TaxMode) ||
		!plan.Quantity.Equal(state.Quantity) ||
		!plan.CustomData.Equal(state.CustomData)

	// If nothing changed, preserve state computed values and return
	if !hasChanges {
		plan.UpdatedAt = state.UpdatedAt
		plan.Status = state.Status
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	// Parse unit_price
	var unitPrice unitPriceModel
	diags := plan.UnitPrice.As(ctx, &unitPrice, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update request
	updateReq := &paddle.UpdatePriceRequest{
		PriceID:     plan.ID.ValueString(),
		Description: paddle.NewPatchField(plan.Description.ValueString()),
	}

	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		namePtr := &name
		updateReq.Name = paddle.NewPatchField(namePtr)
	} else {
		updateReq.Name = paddle.NewPatchField[*string](nil)
	}

	if !plan.TaxMode.IsNull() && !plan.TaxMode.IsUnknown() {
		taxMode := paddle.TaxMode(plan.TaxMode.ValueString())
		updateReq.TaxMode = paddle.NewPatchField(taxMode)
	} else if plan.TaxMode.IsNull() {
		updateReq.TaxMode = paddle.NewPatchField(paddle.TaxModeAccountSetting)
	}

	if !plan.Quantity.IsNull() && !plan.Quantity.IsUnknown() {
		var quantity quantityModel
		diags := plan.Quantity.As(ctx, &quantity, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		quantityValue := paddle.PriceQuantity{
			Minimum: int(quantity.Minimum.ValueInt64()),
			Maximum: int(quantity.Maximum.ValueInt64()),
		}
		updateReq.Quantity = paddle.NewPatchField(quantityValue)
	}

	if !plan.CustomData.IsNull() && !plan.CustomData.IsUnknown() {
		customDataStr := make(map[string]string)
		diags := plan.CustomData.ElementsAs(ctx, &customDataStr, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.CustomData = paddle.NewPatchField(helpers.MapToCustomData(customDataStr))
	} else {
		updateReq.CustomData = paddle.NewPatchField[paddle.CustomData](nil)
	}

	// Update price via Paddle API
	price, err := r.client.UpdatePrice(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating price",
			fmt.Sprintf("Could not update price ID %s (product: %s, description: %s): %s", plan.ID.ValueString(), plan.ProductID.ValueString(), plan.Description.ValueString(), err.Error()),
		)
		return
	}

	// Update model with response data
	plan.Status = types.StringValue(string(price.Status))
	plan.TaxMode = types.StringValue(string(price.TaxMode))
	plan.UpdatedAt = types.StringValue(price.UpdatedAt)

	// Update quantity from API response (Paddle may return defaults)
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
		plan.Quantity = quantityObj
	} else {
		plan.Quantity = types.ObjectNull(quantityAttrTypes)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete archives a Paddle price by setting its status to archived.
// Prices cannot be hard-deleted in Paddle, only archived.
func (r *PriceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data priceResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Archive the price by setting status to "archived"
	updateReq := &paddle.UpdatePriceRequest{
		PriceID: data.ID.ValueString(),
		Status:  paddle.NewPatchField(paddle.StatusArchived),
	}

	_, err := r.client.UpdatePrice(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error archiving price",
			fmt.Sprintf("Could not archive price ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports an existing Paddle price by its ID.
func (r *PriceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
