package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sugarshin/terraform-provider-liff/internal/client"
)

var _ datasource.DataSource = &LiffAppsDataSource{}
var _ datasource.DataSourceWithConfigure = &LiffAppsDataSource{}

type LiffAppsDataSource struct {
	client *client.Client
}

type LiffAppsDataSourceModel struct {
	Apps types.List `tfsdk:"apps"`
}

type LiffAppDataModel struct {
	LiffID               types.String `tfsdk:"liff_id"`
	Description          types.String `tfsdk:"description"`
	PermanentLinkPattern types.String `tfsdk:"permanent_link_pattern"`
	BotPrompt            types.String `tfsdk:"bot_prompt"`
	Scope                types.List   `tfsdk:"scope"`
	View                 types.Object `tfsdk:"view"`
	Features             types.Object `tfsdk:"features"`
}

var liffAppDataModelAttrTypes = map[string]attr.Type{
	"liff_id":                types.StringType,
	"description":            types.StringType,
	"permanent_link_pattern": types.StringType,
	"bot_prompt":             types.StringType,
	"scope":                  types.ListType{ElemType: types.StringType},
	"view":                   types.ObjectType{AttrTypes: viewAttrTypes},
	"features":               types.ObjectType{AttrTypes: featuresAttrTypes},
}

func NewLiffAppsDataSource() datasource.DataSource {
	return &LiffAppsDataSource{}
}

func (d *LiffAppsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apps"
}

func (d *LiffAppsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches all LIFF apps in the channel.",
		Attributes: map[string]schema.Attribute{
			"apps": schema.ListNestedAttribute{
				MarkdownDescription: "List of LIFF apps.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: liffAppDataSourceAttributes(),
				},
			},
		},
	}
}

func (d *LiffAppsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LiffAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	apps, err := d.client.GetAllLiffApps(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list LIFF apps", err.Error())
		return
	}

	appModels := make([]LiffAppDataModel, len(apps))
	for i, app := range apps {
		appModels[i] = mapLiffAppToDataSourceModel(ctx, &app, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	appObjValues := make([]attr.Value, len(appModels))
	for i, m := range appModels {
		obj, diags := types.ObjectValueFrom(ctx, liffAppDataModelAttrTypes, m)
		resp.Diagnostics.Append(diags...)
		appObjValues[i] = obj
	}
	if resp.Diagnostics.HasError() {
		return
	}

	appsList, diags := types.ListValue(types.ObjectType{AttrTypes: liffAppDataModelAttrTypes}, appObjValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := LiffAppsDataSourceModel{
		Apps: appsList,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func liffAppDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"liff_id": schema.StringAttribute{
			MarkdownDescription: "The LIFF app ID.",
			Computed:            true,
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
	}
}
