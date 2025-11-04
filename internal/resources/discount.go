package resources

import (
	"context"
	"fmt"

	"github.com/HQarroum/terraform-provider-paddle/internal/helpers"
	"github.com/HQarroum/terraform-provider-paddle/internal/validators"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DiscountResource{}
var _ resource.ResourceWithImportState = &DiscountResource{}

// Creates a new Paddle discount resource.
func NewDiscountResource() resource.Resource {
	return &DiscountResource{}
}

type DiscountResource struct {
	client *paddle.SDK
}

type discountResourceModel struct {
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

func (r *DiscountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_discount"
}

func (r *DiscountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Paddle discount resource. Discounts reduce transaction totals by percentage or amount. Sometimes called coupons or promo codes.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Paddle discount ID (format: dsc_...)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this discount can be used. Either `active` or `archived`. Set to `archived` to archive.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Short description for this discount. Not shown to customers.",
			},
			"enabled_for_checkout": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether customers can redeem this discount at checkout using a code. Defaults to `true`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"code": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Unique code that customers use to redeem this discount. Not case-sensitive. Omit for discounts applied by the seller.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of discount. One of: `percentage`, `flat` (flat amount per transaction), `flat_per_seat` (flat amount per unit).",
				Validators: []validator.String{
					validators.DiscountTypeValidator{},
				},
			},
			"mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Discount mode. One of: `standard` (shown in dashboard), `custom` (not shown). Defaults to `standard`. Cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"amount": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Amount to discount by. For `percentage`: value between 0.01 and 100. For `flat`/`flat_per_seat`: amount in lowest denomination (e.g., cents).",
			},
			"currency_code": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Three-letter ISO 4217 currency code. Required for `flat` and `flat_per_seat` discount types. Cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"recur": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this discount applies for multiple subscription billing periods. Defaults to `false`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_recurring_intervals": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Number of billing periods this discount recurs for. Requires `recur` to be `true`. Omit for unlimited recurring discount.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"usage_limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum number of times this discount can be redeemed overall (not per customer). Omit for unlimited usage.",
			},
			"restrict_to": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Product or price IDs this discount applies to (e.g., ['pro_...', 'pri_...']). Omit to apply to all products and prices.",
			},
			"expires_at": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "RFC 3339 datetime when this discount expires. Omit for discount that never expires.",
			},
			"custom_data": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom key-value data. Max 10 keys, 1KB total.",
			},
			"times_used": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "How many times this discount has been redeemed. Set automatically by Paddle.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"discount_group_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Paddle ID of the discount group this discount belongs to (format: dsg_...). Omit if not in a group.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the discount was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the discount was last updated.",
			},
		},
	}
}

func (r *DiscountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DiscountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data discountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &paddle.CreateDiscountRequest{
		Description: data.Description.ValueString(),
		Type:        paddle.DiscountType(data.Type.ValueString()),
		Amount:      data.Amount.ValueString(),
	}

	if !data.EnabledForCheckout.IsNull() {
		enabled := data.EnabledForCheckout.ValueBool()
		createReq.EnabledForCheckout = &enabled
	}

	if !data.Code.IsNull() {
		code := data.Code.ValueString()
		createReq.Code = &code
	}

	if !data.Mode.IsNull() {
		mode := paddle.DiscountMode(data.Mode.ValueString())
		createReq.Mode = &mode
	}

	if !data.CurrencyCode.IsNull() {
		currencyCode := paddle.CurrencyCode(data.CurrencyCode.ValueString())
		createReq.CurrencyCode = &currencyCode
	}

	if !data.Recur.IsNull() {
		recur := data.Recur.ValueBool()
		createReq.Recur = &recur
	}

	if !data.MaximumRecurringIntervals.IsNull() {
		maxIntervals := int(data.MaximumRecurringIntervals.ValueInt64())
		createReq.MaximumRecurringIntervals = &maxIntervals
	}

	if !data.UsageLimit.IsNull() {
		usageLimit := int(data.UsageLimit.ValueInt64())
		createReq.UsageLimit = &usageLimit
	}

	if !data.RestrictTo.IsNull() && !data.RestrictTo.IsUnknown() {
		var restrictTo []string
		resp.Diagnostics.Append(data.RestrictTo.ElementsAs(ctx, &restrictTo, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.RestrictTo = restrictTo
	}

	if !data.ExpiresAt.IsNull() {
		expiresAt := data.ExpiresAt.ValueString()
		createReq.ExpiresAt = &expiresAt
	}

	if !data.CustomData.IsNull() && !data.CustomData.IsUnknown() {
		customDataStr := make(map[string]string)
		resp.Diagnostics.Append(data.CustomData.ElementsAs(ctx, &customDataStr, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.CustomData = helpers.MapToCustomData(customDataStr)
	}

	if !data.DiscountGroupID.IsNull() {
		groupID := data.DiscountGroupID.ValueString()
		createReq.DiscountGroupID = &groupID
	}

	discount, err := r.client.CreateDiscount(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating discount",
			fmt.Sprintf("Could not create discount with code '%s' (type: %s): %s", data.Code.ValueString(), data.Type.ValueString(), err.Error()),
		)
		return
	}

	data.ID = types.StringValue(discount.ID)
	data.Status = types.StringValue(string(discount.Status))
	data.EnabledForCheckout = types.BoolValue(discount.EnabledForCheckout)
	data.Recur = types.BoolValue(discount.Recur)
	data.TimesUsed = types.Int64Value(int64(discount.TimesUsed))
	data.CreatedAt = types.StringValue(discount.CreatedAt)
	data.UpdatedAt = types.StringValue(discount.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiscountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data discountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	discount, err := r.client.GetDiscount(ctx, &paddle.GetDiscountRequest{
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
		customDataValue, diags := types.MapValueFrom(ctx, types.StringType, customDataMap)
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

func (r *DiscountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data discountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &paddle.UpdateDiscountRequest{
		DiscountID:  data.ID.ValueString(),
		Description: paddle.NewPatchField(data.Description.ValueString()),
	}

	if !data.EnabledForCheckout.IsNull() {
		updateReq.EnabledForCheckout = paddle.NewPatchField(data.EnabledForCheckout.ValueBool())
	}

	if !data.Code.IsNull() {
		code := data.Code.ValueString()
		updateReq.Code = paddle.NewPatchField(&code)
	} else {
		updateReq.Code = paddle.NewPatchField[*string](nil)
	}

	updateReq.Amount = paddle.NewPatchField(data.Amount.ValueString())
	updateReq.Type = paddle.NewPatchField(paddle.DiscountType(data.Type.ValueString()))

	if !data.Recur.IsNull() {
		updateReq.Recur = paddle.NewPatchField(data.Recur.ValueBool())
	}

	if !data.UsageLimit.IsNull() {
		usageLimit := int(data.UsageLimit.ValueInt64())
		updateReq.UsageLimit = paddle.NewPatchField(&usageLimit)
	} else {
		updateReq.UsageLimit = paddle.NewPatchField[*int](nil)
	}

	if !data.RestrictTo.IsNull() && !data.RestrictTo.IsUnknown() {
		var restrictTo []string
		resp.Diagnostics.Append(data.RestrictTo.ElementsAs(ctx, &restrictTo, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.RestrictTo = paddle.NewPatchField(restrictTo)
	} else if data.RestrictTo.IsNull() {
		updateReq.RestrictTo = paddle.NewPatchField[[]string](nil)
	}

	if !data.ExpiresAt.IsNull() {
		expiresAt := data.ExpiresAt.ValueString()
		updateReq.ExpiresAt = paddle.NewPatchField(&expiresAt)
	} else {
		updateReq.ExpiresAt = paddle.NewPatchField[*string](nil)
	}

	if !data.CustomData.IsNull() && !data.CustomData.IsUnknown() {
		customDataStr := make(map[string]string)
		resp.Diagnostics.Append(data.CustomData.ElementsAs(ctx, &customDataStr, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.CustomData = paddle.NewPatchField(helpers.MapToCustomData(customDataStr))
	} else {
		updateReq.CustomData = paddle.NewPatchField[paddle.CustomData](nil)
	}

	if !data.DiscountGroupID.IsNull() {
		groupID := data.DiscountGroupID.ValueString()
		updateReq.DiscountGroupID = paddle.NewPatchField(&groupID)
	} else {
		updateReq.DiscountGroupID = paddle.NewPatchField[*string](nil)
	}

	discount, err := r.client.UpdateDiscount(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating discount",
			fmt.Sprintf("Could not update discount ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	data.Status = types.StringValue(string(discount.Status))
	data.EnabledForCheckout = types.BoolValue(discount.EnabledForCheckout)
	data.Recur = types.BoolValue(discount.Recur)
	data.TimesUsed = types.Int64Value(int64(discount.TimesUsed))
	data.UpdatedAt = types.StringValue(discount.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiscountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data discountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &paddle.UpdateDiscountRequest{
		DiscountID: data.ID.ValueString(),
		Status:     paddle.NewPatchField(paddle.DiscountStatusArchived),
	}

	_, err := r.client.UpdateDiscount(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error archiving discount",
			fmt.Sprintf("Could not archive discount ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *DiscountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
