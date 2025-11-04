package resources

import (
	"context"
	"fmt"

	"github.com/HQarroum/terraform-provider-paddle/internal/helpers"
	"github.com/HQarroum/terraform-provider-paddle/internal/validators"
	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProductResource{}
var _ resource.ResourceWithImportState = &ProductResource{}

// Creates a new Paddle product resource.
func NewProductResource() resource.Resource {
	return &ProductResource{}
}

// Product resource manages Paddle products.
type ProductResource struct {
	client *paddle.SDK
}

// Product resource model describes the resource data model.
type productResourceModel struct {
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

// Metadata returns the resource type name.
func (r *ProductResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

// Schema returns the resource schema definition.
func (r *ProductResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Paddle product resource. Products are the items that you sell. Prices determine how much you charge for them.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Paddle product ID (format: pro_...)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the product. Displayed in customer-facing areas.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Short description of the product. Displayed in customer-facing areas.",
			},
			"tax_category": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Tax category for the product. One of: standard, digital-goods, ebooks, implementation-services, professional-services, saas, software-programming-services, training-services.",
				Validators: []validator.String{
					validators.TaxCategoryValidator{},
				},
			},
			"image_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "URL of the product image. Displayed in customer-facing areas.",
			},
			"custom_data": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom data for this product. Max 10 keys, 1KB total.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the product. Either 'active' or 'archived'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the product was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC 3339 timestamp when the product was last updated.",
			},
		},
	}
}

// Configure initializes the resource with the Paddle SDK client.
func (r *ProductResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

// Create creates a new Paddle product.
func (r *ProductResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data productResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert custom_data from types.Map to map[string]string
	var customData map[string]string
	if !data.CustomData.IsNull() && !data.CustomData.IsUnknown() {
		customData = make(map[string]string)
		diags := data.CustomData.ElementsAs(ctx, &customData, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build create request
	createReq := &paddle.CreateProductRequest{
		Name:        data.Name.ValueString(),
		TaxCategory: paddle.TaxCategory(data.TaxCategory.ValueString()),
	}

	if !data.Description.IsNull() {
		desc := data.Description.ValueString()
		createReq.Description = &desc
	}

	if !data.ImageURL.IsNull() {
		imageURL := data.ImageURL.ValueString()
		createReq.ImageURL = &imageURL
	}

	if customData != nil {
		createReq.CustomData = helpers.MapToCustomData(customData)
	}

	// Create product via Paddle API
	product, err := r.client.CreateProduct(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating product",
			fmt.Sprintf("Could not create product with name '%s': %s", data.Name.ValueString(), err.Error()),
		)
		return
	}

	// Map response to model
	data.ID = types.StringValue(product.ID)
	data.Status = types.StringValue(string(product.Status))
	data.CreatedAt = types.StringValue(product.CreatedAt)
	data.UpdatedAt = types.StringValue(product.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read retrieves the current state of a Paddle product.
func (r *ProductResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data productResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read product from Paddle API
	product, err := r.client.GetProduct(ctx, &paddle.GetProductRequest{
		ProductID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading product",
			fmt.Sprintf("Could not read product ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update model with response data
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update modifies an existing Paddle product.
func (r *ProductResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data productResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert custom_data from types.Map to map[string]string
	var customData map[string]string
	if !data.CustomData.IsNull() && !data.CustomData.IsUnknown() {
		customData = make(map[string]string)
		diags := data.CustomData.ElementsAs(ctx, &customData, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build update request
	updateReq := &paddle.UpdateProductRequest{
		ProductID:   data.ID.ValueString(),
		Name:        paddle.NewPatchField(data.Name.ValueString()),
		TaxCategory: paddle.NewPatchField(paddle.TaxCategory(data.TaxCategory.ValueString())),
	}

	if !data.Description.IsNull() {
		desc := data.Description.ValueString()
		updateReq.Description = paddle.NewPatchField(&desc)
	} else {
		updateReq.Description = paddle.NewPatchField[*string](nil)
	}

	if !data.ImageURL.IsNull() {
		imageURL := data.ImageURL.ValueString()
		updateReq.ImageURL = paddle.NewPatchField(&imageURL)
	} else {
		updateReq.ImageURL = paddle.NewPatchField[*string](nil)
	}

	if customData != nil {
		updateReq.CustomData = paddle.NewPatchField(helpers.MapToCustomData(customData))
	} else if data.CustomData.IsNull() {
		updateReq.CustomData = paddle.NewPatchField[paddle.CustomData](nil)
	}

	// Update product via Paddle API
	product, err := r.client.UpdateProduct(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating product",
			fmt.Sprintf("Could not update product ID %s (name: %s): %s", data.ID.ValueString(), data.Name.ValueString(), err.Error()),
		)
		return
	}

	// Update model with response data
	data.Status = types.StringValue(string(product.Status))
	data.UpdatedAt = types.StringValue(product.UpdatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete archives a Paddle product by setting its status to archived.
// Products cannot be hard-deleted in Paddle, only archived.
func (r *ProductResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data productResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Archive the product by setting its status to "archived"
	updateReq := &paddle.UpdateProductRequest{
		ProductID: data.ID.ValueString(),
		Status:    paddle.NewPatchField(paddle.StatusArchived),
	}

	_, err := r.client.UpdateProduct(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error archiving product",
			fmt.Sprintf("Could not archive product ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
}

// ImportState imports an existing Paddle product by its ID.
func (r *ProductResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
