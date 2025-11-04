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

var _ datasource.DataSource = &DiscountDataSource{}

// Creates a new discount data source.
func NewDiscountDataSource() datasource.DataSource {
	return &DiscountDataSource{}
}

type DiscountDataSource struct {
	client *paddle.SDK
}

type discountDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Status                    types.String `tfsdk:"status"`
	Description               types.String `tfsdk:"description"`
	EnabledForCheckout        types.Bool   `tfsdk:"enabled_for_checkout"`
	Code                      types.String `tfsdk:"code"`
	Type                      types.String `tfsdk:"type"`
	Mode                      types.String `tfsdk:"mode"`
	Amount                    types.String `tfsdk:"amount"`
	CurrencyCode              types.String `tfsdk:"currency_code"`
	Recur                     types.Bool   `tfsdk:"recur"`
	MaximumRecurringIntervals types.Int64  `tfsdk:"maximum_recurring_intervals"`
	UsageLimit                types.Int64  `tfsdk:"usage_limit"`
	RestrictTo                types.List   `tfsdk:"restrict_to"`
	ExpiresAt                 types.String `tfsdk:"expires_at"`
	CustomData                types.Map    `tfsdk:"custom_data"`
	TimesUsed                 types.Int64  `tfsdk:"times_used"`
	DiscountGroupID           types.String `tfsdk:"discount_group_id"`
	CreatedAt                 types.String `tfsdk:"created_at"`
	UpdatedAt                 types.String `tfsdk:"updated_at"`
}

func (d *DiscountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_discount"
}

func (d *DiscountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves details of a specific Paddle discount by ID.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Paddle discount ID (format: dsc_...).",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this discount can be used.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Short description for this discount.",
			},
			"enabled_for_checkout": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether customers can redeem this discount at checkout.",
			},
			"code": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique code that customers use to redeem this discount.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of discount.",
			},
			"mode": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Discount mode.",
			},
			"amount": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Amount to discount by.",
			},
			"currency_code": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Three-letter ISO 4217 currency code.",
			},
			"recur": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this discount applies for multiple billing periods.",
			},
			"maximum_recurring_intervals": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of billing periods this discount recurs for.",
			},
			"usage_limit": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Maximum number of times this discount can be redeemed.",
			},
			"restrict_to": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Product or price IDs this discount applies to.",
			},
			"expires_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 datetime when this discount expires.",
			},
			"custom_data": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom key-value data.",
			},
			"times_used": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "How many times this discount has been redeemed.",
			},
			"discount_group_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Paddle ID of the discount group this discount belongs to.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the discount was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the discount was last updated.",
			},
		},
	}
}

func (d *DiscountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DiscountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data discountDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	discount, err := d.client.GetDiscount(ctx, &paddle.GetDiscountRequest{
		DiscountID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading discount",
			fmt.Sprintf("Could not read discount ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	data.Status = types.StringValue(string(discount.Status))
	data.Description = types.StringValue(discount.Description)
	data.EnabledForCheckout = types.BoolValue(discount.EnabledForCheckout)
	data.Type = types.StringValue(string(discount.Type))
	data.Mode = types.StringValue(string(discount.Mode))
	data.Amount = types.StringValue(discount.Amount)
	data.Recur = types.BoolValue(discount.Recur)
	data.TimesUsed = types.Int64Value(int64(discount.TimesUsed))
	data.CreatedAt = types.StringValue(discount.CreatedAt)
	data.UpdatedAt = types.StringValue(discount.UpdatedAt)

	if discount.Code != nil {
		data.Code = types.StringValue(*discount.Code)
	} else {
		data.Code = types.StringNull()
	}

	if discount.CurrencyCode != nil {
		data.CurrencyCode = types.StringValue(string(*discount.CurrencyCode))
	} else {
		data.CurrencyCode = types.StringNull()
	}

	if discount.MaximumRecurringIntervals != nil {
		data.MaximumRecurringIntervals = types.Int64Value(int64(*discount.MaximumRecurringIntervals))
	} else {
		data.MaximumRecurringIntervals = types.Int64Null()
	}

	if discount.UsageLimit != nil {
		data.UsageLimit = types.Int64Value(int64(*discount.UsageLimit))
	} else {
		data.UsageLimit = types.Int64Null()
	}

	if len(discount.RestrictTo) > 0 {
		restrictToList, diags := types.ListValueFrom(ctx, types.StringType, discount.RestrictTo)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.RestrictTo = restrictToList
	} else {
		data.RestrictTo = types.ListNull(types.StringType)
	}

	if discount.ExpiresAt != nil {
		data.ExpiresAt = types.StringValue(*discount.ExpiresAt)
	} else {
		data.ExpiresAt = types.StringNull()
	}

	if discount.CustomData != nil {
		customDataMap, err := helpers.CustomDataToMap(discount.CustomData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error processing custom data",
				fmt.Sprintf("Could not convert custom_data for discount %s: %s", data.ID.ValueString(), err.Error()),
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

	if discount.DiscountGroupID != nil {
		data.DiscountGroupID = types.StringValue(*discount.DiscountGroupID)
	} else {
		data.DiscountGroupID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
