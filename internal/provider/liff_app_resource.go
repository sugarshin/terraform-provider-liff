// Copyright (c) sugarshin
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sugarshin/terraform-provider-liff/internal/client"
)

var (
	_ resource.Resource                = &LiffAppResource{}
	_ resource.ResourceWithImportState = &LiffAppResource{}
	_ resource.ResourceWithConfigure   = &LiffAppResource{}
)

type LiffAppResource struct {
	client *client.Client
}

type LiffAppResourceModel struct {
	LiffID               types.String `tfsdk:"liff_id"`
	Description          types.String `tfsdk:"description"`
	PermanentLinkPattern types.String `tfsdk:"permanent_link_pattern"`
	BotPrompt            types.String `tfsdk:"bot_prompt"`
	Scope                types.List   `tfsdk:"scope"`
	View                 types.Object `tfsdk:"view"`
	Features             types.Object `tfsdk:"features"`
}

type LiffAppViewModel struct {
	Type       types.String `tfsdk:"type"`
	URL        types.String `tfsdk:"url"`
	ModuleMode types.Bool   `tfsdk:"module_mode"`
}

type LiffAppFeaturesModel struct {
	BLE    types.Bool `tfsdk:"ble"`
	QRCode types.Bool `tfsdk:"qr_code"`
}

func NewLiffAppResource() resource.Resource {
	return &LiffAppResource{}
}

func (r *LiffAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *LiffAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a LINE LIFF application.",
		Attributes: map[string]schema.Attribute{
			"liff_id": schema.StringAttribute{
				MarkdownDescription: "The LIFF app ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Name of the LIFF app.",
				Optional:            true,
			},
			"permanent_link_pattern": schema.StringAttribute{
				MarkdownDescription: "How additional information in LIFF URLs is handled. Specify `concat`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("concat"),
				},
			},
			"bot_prompt": schema.StringAttribute{
				MarkdownDescription: "Specify the setting for bot link feature. One of `normal`, `aggressive`, `none`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf("normal", "aggressive", "none"),
				},
			},
			"scope": schema.ListAttribute{
				MarkdownDescription: "Array of scopes.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"view": schema.SingleNestedAttribute{
				MarkdownDescription: "LIFF app view settings.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Size of the LIFF app view. One of `compact`, `tall`, `full`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("compact", "tall", "full"),
						},
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "Endpoint URL (HTTPS).",
						Required:            true,
					},
					"module_mode": schema.BoolAttribute{
						MarkdownDescription: "Use the LIFF app in modular mode.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"features": schema.SingleNestedAttribute{
				MarkdownDescription: "LIFF app feature settings.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"ble": schema.BoolAttribute{
						MarkdownDescription: "Whether the LIFF app supports BLE (read-only).",
						Computed:            true,
					},
					"qr_code": schema.BoolAttribute{
						MarkdownDescription: "Whether to use the 2D code reader.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

func (r *LiffAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *LiffAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LiffAppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := buildAddLiffAppRequest(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.AddLiffApp(ctx, addReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create LIFF app", err.Error())
		return
	}

	data.LiffID = types.StringValue(result.LiffId)

	// Read back to get full state (including computed fields)
	app, err := r.client.GetLiffApp(ctx, result.LiffId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read LIFF app after creation", err.Error())
		return
	}

	mapLiffAppToModel(app, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LiffAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LiffAppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.GetLiffApp(ctx, data.LiffID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read LIFF app", err.Error())
		return
	}

	mapLiffAppToModel(app, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LiffAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LiffAppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state LiffAppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	liffID := state.LiffID.ValueString()
	updateReq := buildUpdateLiffAppRequest(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateLiffApp(ctx, liffID, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update LIFF app", err.Error())
		return
	}

	app, err := r.client.GetLiffApp(ctx, liffID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read LIFF app after update", err.Error())
		return
	}

	data.LiffID = types.StringValue(liffID)
	mapLiffAppToModel(app, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LiffAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LiffAppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLiffApp(ctx, data.LiffID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete LIFF app", err.Error())
		return
	}
}

func (r *LiffAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("liff_id"), req, resp)
}
