package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sugarshin/terraform-provider-liff/internal/client"
)

var _ provider.Provider = &LiffProvider{}

type LiffProvider struct {
	version string
}

type LiffProviderModel struct {
	ChannelID          types.String `tfsdk:"channel_id"`
	ChannelSecret      types.String `tfsdk:"channel_secret"`
	ChannelAccessToken types.String `tfsdk:"channel_access_token"`
}

func (p *LiffProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "liff"
	resp.Version = p.version
}

func (p *LiffProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The LIFF provider manages LINE LIFF (LINE Front-end Framework) applications.",
		Attributes: map[string]schema.Attribute{
			"channel_id": schema.StringAttribute{
				MarkdownDescription: "LINE Login Channel ID. Can also be set via `LIFF_CHANNEL_ID` environment variable.",
				Optional:            true,
			},
			"channel_secret": schema.StringAttribute{
				MarkdownDescription: "LINE Login Channel Secret. Can also be set via `LIFF_CHANNEL_SECRET` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"channel_access_token": schema.StringAttribute{
				MarkdownDescription: "Channel Access Token. Takes precedence over channel_id/channel_secret. Can also be set via `LIFF_CHANNEL_ACCESS_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *LiffProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LiffProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelAccessToken := stringValueOrEnv(data.ChannelAccessToken, "LIFF_CHANNEL_ACCESS_TOKEN")
	channelID := stringValueOrEnv(data.ChannelID, "LIFF_CHANNEL_ID")
	channelSecret := stringValueOrEnv(data.ChannelSecret, "LIFF_CHANNEL_SECRET")

	var c *client.Client
	var err error

	if channelAccessToken != "" {
		c, err = client.NewWithToken(channelAccessToken)
	} else if channelID != "" && channelSecret != "" {
		c, err = client.NewWithCredentials(channelID, channelSecret)
	} else {
		resp.Diagnostics.AddError(
			"Missing Authentication Configuration",
			"Provide either channel_access_token or both channel_id and channel_secret. "+
				"These can be set in the provider block or via LIFF_CHANNEL_ACCESS_TOKEN, "+
				"LIFF_CHANNEL_ID, and LIFF_CHANNEL_SECRET environment variables.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to create LIFF API client", err.Error())
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *LiffProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLiffAppResource,
	}
}

func (p *LiffProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLiffAppsDataSource,
		NewLiffAppDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LiffProvider{
			version: version,
		}
	}
}

func stringValueOrEnv(v types.String, envKey string) string {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueString()
	}
	return os.Getenv(envKey)
}
