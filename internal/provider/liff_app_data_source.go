// Copyright (c) sugarshin
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sugarshin/terraform-provider-liff/internal/client"
)

var _ datasource.DataSource = &LiffAppDataSource{}
var _ datasource.DataSourceWithConfigure = &LiffAppDataSource{}

type LiffAppDataSource struct {
	client *client.Client
}

type LiffAppDataSourceModel struct {
	LiffID               types.String `tfsdk:"liff_id"`
	Description          types.String `tfsdk:"description"`
	PermanentLinkPattern types.String `tfsdk:"permanent_link_pattern"`
	BotPrompt            types.String `tfsdk:"bot_prompt"`
	Scope                types.List   `tfsdk:"scope"`
	View                 types.Object `tfsdk:"view"`
	Features             types.Object `tfsdk:"features"`
}

func NewLiffAppDataSource() datasource.DataSource {
	return &LiffAppDataSource{}
}

func (d *LiffAppDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (d *LiffAppDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a specific LIFF app by its ID.",
		Attributes: map[string]schema.Attribute{
			"liff_id": schema.StringAttribute{
				MarkdownDescription: "The LIFF app ID to look up.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Name of the LIFF app.",
				Computed:            true,
			},
			"permanent_link_pattern": schema.StringAttribute{
				MarkdownDescription: "How additional information in LIFF URLs is handled.",
				Computed:            true,
			},
			"bot_prompt": schema.StringAttribute{
				MarkdownDescription: "Bot link feature setting.",
				Computed:            true,
			},
			"scope": schema.ListAttribute{
				MarkdownDescription: "Array of scopes.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"view": schema.SingleNestedAttribute{
				MarkdownDescription: "LIFF app view settings.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Size of the LIFF app view.",
						Computed:            true,
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "Endpoint URL.",
						Computed:            true,
					},
					"module_mode": schema.BoolAttribute{
						MarkdownDescription: "Modular mode.",
						Computed:            true,
					},
				},
			},
			"features": schema.SingleNestedAttribute{
				MarkdownDescription: "LIFF app feature settings.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"ble": schema.BoolAttribute{
						MarkdownDescription: "BLE support.",
						Computed:            true,
					},
					"qr_code": schema.BoolAttribute{
						MarkdownDescription: "2D code reader.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *LiffAppDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *LiffAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LiffAppDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := d.client.GetLiffApp(ctx, data.LiffID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read LIFF app", err.Error())
		return
	}

	data.LiffID = types.StringValue(app.LiffId)

	if app.Description != "" {
		data.Description = types.StringValue(app.Description)
	} else {
		data.Description = types.StringNull()
	}

	if app.PermanentLinkPattern != "" {
		data.PermanentLinkPattern = types.StringValue(app.PermanentLinkPattern)
	} else {
		data.PermanentLinkPattern = types.StringNull()
	}

	data.BotPrompt = types.StringValue(string(app.BotPrompt))

	if len(app.Scope) > 0 {
		scopeValues := make([]types.String, len(app.Scope))
		for i, s := range app.Scope {
			scopeValues[i] = types.StringValue(string(s))
		}
		scopeList, diags := types.ListValueFrom(ctx, types.StringType, scopeValues)
		resp.Diagnostics.Append(diags...)
		data.Scope = scopeList
	} else {
		data.Scope = types.ListNull(types.StringType)
	}

	if app.View != nil {
		viewObj, diags := types.ObjectValueFrom(ctx, viewAttrTypes, LiffAppViewModel{
			Type:       types.StringValue(string(app.View.Type)),
			URL:        types.StringValue(app.View.Url),
			ModuleMode: types.BoolValue(app.View.ModuleMode),
		})
		resp.Diagnostics.Append(diags...)
		data.View = viewObj
	} else {
		data.View = types.ObjectNull(viewAttrTypes)
	}

	if app.Features != nil {
		featuresObj, diags := types.ObjectValueFrom(ctx, featuresAttrTypes, LiffAppFeaturesModel{
			BLE:    types.BoolValue(app.Features.Ble),
			QRCode: types.BoolValue(app.Features.QrCode),
		})
		resp.Diagnostics.Append(diags...)
		data.Features = featuresObj
	} else {
		featuresObj, diags := types.ObjectValueFrom(ctx, featuresAttrTypes, LiffAppFeaturesModel{
			BLE:    types.BoolValue(false),
			QRCode: types.BoolValue(false),
		})
		resp.Diagnostics.Append(diags...)
		data.Features = featuresObj
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
