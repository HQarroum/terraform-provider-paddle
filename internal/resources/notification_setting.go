package resources

import (
	"context"
	"fmt"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationSettingResource{}
var _ resource.ResourceWithImportState = &NotificationSettingResource{}

// Creates a new Paddle notification setting resource.
func NewNotificationSettingResource() resource.Resource {
	return &NotificationSettingResource{}
}

// Manages Paddle notification settings (webhooks).
type NotificationSettingResource struct {
	client *paddle.SDK
}

// notificationSettingResourceModel describes the resource data model.
type notificationSettingResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Description            types.String `tfsdk:"description"`
	Type                   types.String `tfsdk:"type"`
	Destination            types.String `tfsdk:"destination"`
	Active                 types.Bool   `tfsdk:"active"`
	SubscribedEvents       types.List   `tfsdk:"subscribed_events"`
	EndpointSecretKey      types.String `tfsdk:"endpoint_secret_key"`
	APIVersion             types.Int64  `tfsdk:"api_version"`
	IncludeSensitiveFields types.Bool   `tfsdk:"include_sensitive_fields"`
	TrafficSource          types.String `tfsdk:"traffic_source"`
}

// Metadata returns the resource type name.
func (r *NotificationSettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_setting"
}

// Schema returns the resource schema definition.
func (r *NotificationSettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Paddle notification setting resource. Notification settings allow you to subscribe to events and receive webhooks.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Paddle notification setting ID (format: ntfset_...)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Short description for this notification destination.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Type of notification destination. One of: `url` (webhook), `email`. Defaults to `url`. Cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destination": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Webhook endpoint URL or email address for notifications. URLs must be HTTPS.",
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the notification destination is active. Defaults to true.",
			},
			"subscribed_events": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of event types to subscribe to (e.g., 'transaction.completed', 'subscription.activated').",
			},
			"endpoint_secret_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Webhook secret key for signature verification. Keep this secure!",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "API version for event payloads.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"include_sensitive_fields": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to include sensitive fields in webhook payloads.",
			},
			"traffic_source": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Filter events by source. One of: `platform` (events from dashboard), `api` (events from API). Omit to receive all events.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure initializes the resource with the Paddle SDK client.
func (r *NotificationSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new Paddle notification setting (webhook).
func (r *NotificationSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data notificationSettingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert subscribed_events from types.List to []paddle.EventTypeName
	var subscribedEventsStr []string
	resp.Diagnostics.Append(data.SubscribedEvents.ElementsAs(ctx, &subscribedEventsStr, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscribedEvents := make([]paddle.EventTypeName, len(subscribedEventsStr))
	for i, event := range subscribedEventsStr {
		subscribedEvents[i] = paddle.EventTypeName(event)
	}

	// Build create request
	createReq := &paddle.CreateNotificationSettingRequest{
		Description:      data.Description.ValueString(),
		Destination:      data.Destination.ValueString(),
		SubscribedEvents: subscribedEvents,
	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		createReq.Type = paddle.NotificationSettingType(data.Type.ValueString())
	} else {
		createReq.Type = paddle.NotificationSettingTypeURL
	}

	if !data.IncludeSensitiveFields.IsNull() {
		includeSensitive := data.IncludeSensitiveFields.ValueBool()
		createReq.IncludeSensitiveFields = &includeSensitive
	}

	if !data.TrafficSource.IsNull() && !data.TrafficSource.IsUnknown() {
		trafficSource := paddle.TrafficSource(data.TrafficSource.ValueString())
		createReq.TrafficSource = &trafficSource
	}

	// Create notification setting via Paddle API
	notifSetting, err := r.client.CreateNotificationSetting(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification setting",
			fmt.Sprintf("Could not create notification setting: %s", err.Error()),
		)
		return
	}

	// Map response to model
	data.ID = types.StringValue(notifSetting.ID)
	data.Type = types.StringValue(string(notifSetting.Type))
	data.Active = types.BoolValue(notifSetting.Active)
	data.EndpointSecretKey = types.StringValue(notifSetting.EndpointSecretKey)
	data.APIVersion = types.Int64Value(int64(notifSetting.APIVersion))
	data.IncludeSensitiveFields = types.BoolValue(notifSetting.IncludeSensitiveFields)

	if notifSetting.TrafficSource != "" {
		data.TrafficSource = types.StringValue(string(notifSetting.TrafficSource))
	} else {
		data.TrafficSource = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read retrieves the current state of a Paddle notification setting.
func (r *NotificationSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data notificationSettingResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read notification setting from Paddle API
	notifSetting, err := r.client.GetNotificationSetting(ctx, &paddle.GetNotificationSettingRequest{
		NotificationSettingID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading notification setting",
			fmt.Sprintf("Could not read notification setting ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Update model with response data
	data.Description = types.StringValue(notifSetting.Description)
	data.Type = types.StringValue(string(notifSetting.Type))
	data.Destination = types.StringValue(notifSetting.Destination)
	data.Active = types.BoolValue(notifSetting.Active)
	data.EndpointSecretKey = types.StringValue(notifSetting.EndpointSecretKey)
	data.APIVersion = types.Int64Value(int64(notifSetting.APIVersion))
	data.IncludeSensitiveFields = types.BoolValue(notifSetting.IncludeSensitiveFields)

	// Map traffic_source if present
	if notifSetting.TrafficSource != "" {
		data.TrafficSource = types.StringValue(string(notifSetting.TrafficSource))
	} else {
		data.TrafficSource = types.StringNull()
	}

	// For subscribed_events, preserve the order from state to avoid drift
	// Paddle may return events in a different order
	var stateEvents []string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("subscribed_events"), &stateEvents)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert API response to a map for lookup
	apiEventsMap := make(map[string]bool)
	for _, event := range notifSetting.SubscribedEvents {
		apiEventsMap[string(event.Name)] = true
	}

	// Verify all state events are still in API response, use state order
	orderedEvents := make([]string, 0, len(stateEvents))
	for _, event := range stateEvents {
		if apiEventsMap[event] {
			orderedEvents = append(orderedEvents, event)
		}
	}

	// Add any new events from API that aren't in state (shouldn't happen normally)
	for _, event := range notifSetting.SubscribedEvents {
		found := false
		for _, stateEvent := range stateEvents {
			if string(event.Name) == stateEvent {
				found = true
				break
			}
		}
		if !found {
			orderedEvents = append(orderedEvents, string(event.Name))
		}
	}

	eventList, diags := types.ListValueFrom(ctx, types.StringType, orderedEvents)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.SubscribedEvents = eventList

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update modifies an existing Paddle notification setting.
func (r *NotificationSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data notificationSettingResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert subscribed_events from types.List to []paddle.EventTypeName
	var subscribedEventsStr []string
	resp.Diagnostics.Append(data.SubscribedEvents.ElementsAs(ctx, &subscribedEventsStr, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscribedEvents := make([]paddle.EventTypeName, len(subscribedEventsStr))
	for i, event := range subscribedEventsStr {
		subscribedEvents[i] = paddle.EventTypeName(event)
	}

	// Build update request
	updateReq := &paddle.UpdateNotificationSettingRequest{
		NotificationSettingID: data.ID.ValueString(),
		Description:           paddle.NewPatchField(data.Description.ValueString()),
		Destination:           paddle.NewPatchField(data.Destination.ValueString()),
		SubscribedEvents:      paddle.NewPatchField(subscribedEvents),
	}

	if !data.Active.IsNull() && !data.Active.IsUnknown() {
		updateReq.Active = paddle.NewPatchField(data.Active.ValueBool())
	}

	if !data.IncludeSensitiveFields.IsNull() && !data.IncludeSensitiveFields.IsUnknown() {
		updateReq.IncludeSensitiveFields = paddle.NewPatchField(data.IncludeSensitiveFields.ValueBool())
	}

	if !data.TrafficSource.IsNull() && !data.TrafficSource.IsUnknown() {
		trafficSource := paddle.TrafficSource(data.TrafficSource.ValueString())
		updateReq.TrafficSource = paddle.NewPatchField(trafficSource)
	} else if data.TrafficSource.IsNull() {
		updateReq.TrafficSource = paddle.NewPatchField(paddle.TrafficSource(""))
	}

	// Update notification setting via Paddle API
	notifSetting, err := r.client.UpdateNotificationSetting(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating notification setting",
			fmt.Sprintf("Could not update notification setting ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	data.Type = types.StringValue(string(notifSetting.Type))
	data.Active = types.BoolValue(notifSetting.Active)
	data.APIVersion = types.Int64Value(int64(notifSetting.APIVersion))
	data.IncludeSensitiveFields = types.BoolValue(notifSetting.IncludeSensitiveFields)

	if notifSetting.TrafficSource != "" {
		data.TrafficSource = types.StringValue(string(notifSetting.TrafficSource))
	} else {
		data.TrafficSource = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes a Paddle notification setting.
func (r *NotificationSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data notificationSettingResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete notification setting via Paddle API
	err := r.client.DeleteNotificationSetting(ctx, &paddle.DeleteNotificationSettingRequest{
		NotificationSettingID: data.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting notification setting",
			fmt.Sprintf("Could not delete notification setting ID %s: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
}

// Imports an existing Paddle notification setting by its ID.
func (r *NotificationSettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
