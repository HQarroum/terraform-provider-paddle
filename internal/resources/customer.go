package resources

import (
	"context"
	"fmt"

	"github.com/HQarroum/terraform-provider-paddle/internal/helpers"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &CustomerResource{}
var _ resource.ResourceWithImportState = &CustomerResource{}

// Creates a new Paddle customer resource.
func NewCustomerResource() resource.Resource {
	return &CustomerResource{}
}

type CustomerResource struct {
	client *paddle.SDK
}

type customerResourceModel struct {
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

func (r *CustomerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customer"
}

func (r *CustomerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Paddle customer resource. Customers represent people and businesses that make purchases.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Paddle customer ID (format: ctm_...)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Full name of this customer. Required when creating transactions where `collection_mode` is `manual` (invoices).",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email address for this customer.",
			},
			"marketing_consent": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this customer opted into marketing. Set automatically by Paddle when customers check the marketing consent box at checkout.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this customer can be used. Either `active` or `archived`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_data": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom key-value data. Max 10 keys, 1KB total.",
			},
			"locale": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Valid IETF BCP 47 short form locale tag (e.g., 'en', 'en-US'). Defaults to 'en'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the customer was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the customer was last updated.",
			},
		},
	}
}

func (r *CustomerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data customerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &paddle.CreateCustomerRequest{
		Email: data.Email.ValueString(),
	}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		createReq.Name = &name
	}

	if !data.Locale.IsNull() {
		locale := data.Locale.ValueString()
		createReq.Locale = &locale
	}

	if !data.CustomData.IsNull() && !data.CustomData.IsUnknown() {
		customDataStr := make(map[string]string)
		resp.Diagnostics.Append(data.CustomData.ElementsAs(ctx, &customDataStr, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.CustomData = helpers.MapToCustomData(customDataStr)
	}

	customer, err := r.client.CreateCustomer(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating customer",
			fmt.Sprintf("Could not create customer with email '%s': %s", data.Email.ValueString(), err.Error()),
		)
		return
	}

	data.ID = types.StringValue(customer.ID)
	data.MarketingConsent = types.BoolValue(customer.MarketingConsent)
	data.Status = types.StringValue(string(customer.Status))
	data.CreatedAt = types.StringValue(customer.CreatedAt)
	data.UpdatedAt = types.StringValue(customer.UpdatedAt)

	if customer.Locale != "" {
		data.Locale = types.StringValue(customer.Locale)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data customerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	customer, err := r.client.GetCustomer(ctx, &paddle.GetCustomerRequest{
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
		customDataValue, diags := types.MapValueFrom(ctx, types.StringType, customDataMap)
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

func (r *CustomerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data customerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &paddle.UpdateCustomerRequest{
		CustomerID: data.ID.ValueString(),
		Email:      paddle.NewPatchField(data.Email.ValueString()),
	}

	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		updateReq.Name = paddle.NewPatchField(&name)
	} else {
		updateReq.Name = paddle.NewPatchField[*string](nil)
	}

	if !data.Locale.IsNull() {
		locale := data.Locale.ValueString()
		updateReq.Locale = paddle.NewPatchField(locale)
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

	customer, err := r.client.UpdateCustomer(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating customer",
			fmt.Sprintf("Could not update customer ID %s (email: %s): %s", data.ID.ValueString(), data.Email.ValueString(), err.Error()),
		)
		return
	}

	data.MarketingConsent = types.BoolValue(customer.MarketingConsent)
	data.Status = types.StringValue(string(customer.Status))
	data.UpdatedAt = types.StringValue(customer.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data customerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &paddle.UpdateCustomerRequest{
		CustomerID: data.ID.ValueString(),
		Status:     paddle.NewPatchField(paddle.StatusArchived),
	}

	_, err := r.client.UpdateCustomer(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error archiving customer",
			fmt.Sprintf("Could not archive customer ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *CustomerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
