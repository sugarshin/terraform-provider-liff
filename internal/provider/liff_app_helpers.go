// Copyright (c) sugarshin
// SPDX-License-Identifier: MIT

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/line/line-bot-sdk-go/v8/linebot/liff"
)

var viewAttrTypes = map[string]attr.Type{
	"type":        types.StringType,
	"url":         types.StringType,
	"module_mode": types.BoolType,
}

var featuresAttrTypes = map[string]attr.Type{
	"ble":     types.BoolType,
	"qr_code": types.BoolType,
}

func buildAddLiffAppRequest(ctx context.Context, data *LiffAppResourceModel, diags *diag.Diagnostics) *liff.AddLiffAppRequest {
	var view LiffAppViewModel
	diags.Append(data.View.As(ctx, &view, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	req := &liff.AddLiffAppRequest{
		View: &liff.LiffView{
			Type:       liff.LiffViewTYPE(view.Type.ValueString()),
			Url:        view.URL.ValueString(),
			ModuleMode: view.ModuleMode.ValueBool(),
		},
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		req.Description = data.Description.ValueString()
	}

	if !data.PermanentLinkPattern.IsNull() && !data.PermanentLinkPattern.IsUnknown() {
		req.PermanentLinkPattern = data.PermanentLinkPattern.ValueString()
	}

	if !data.BotPrompt.IsNull() && !data.BotPrompt.IsUnknown() {
		req.BotPrompt = liff.LiffBotPrompt(data.BotPrompt.ValueString())
	}

	if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
		var scopes []string
		diags.Append(data.Scope.ElementsAs(ctx, &scopes, false)...)
		if diags.HasError() {
			return nil
		}
		for _, s := range scopes {
			req.Scope = append(req.Scope, liff.LiffScope(s))
		}
	}

	if !data.Features.IsNull() && !data.Features.IsUnknown() {
		var features LiffAppFeaturesModel
		diags.Append(data.Features.As(ctx, &features, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
		req.Features = &liff.LiffFeatures{
			QrCode: features.QRCode.ValueBool(),
		}
	}

	return req
}

func buildUpdateLiffAppRequest(ctx context.Context, data *LiffAppResourceModel, diags *diag.Diagnostics) *liff.UpdateLiffAppRequest {
	var view LiffAppViewModel
	diags.Append(data.View.As(ctx, &view, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}

	req := &liff.UpdateLiffAppRequest{
		View: &liff.UpdateLiffView{
			Type:       liff.UpdateLiffViewTYPE(view.Type.ValueString()),
			Url:        view.URL.ValueString(),
			ModuleMode: view.ModuleMode.ValueBool(),
		},
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		req.Description = data.Description.ValueString()
	}

	if !data.PermanentLinkPattern.IsNull() && !data.PermanentLinkPattern.IsUnknown() {
		req.PermanentLinkPattern = data.PermanentLinkPattern.ValueString()
	}

	if !data.BotPrompt.IsNull() && !data.BotPrompt.IsUnknown() {
		req.BotPrompt = liff.LiffBotPrompt(data.BotPrompt.ValueString())
	}

	if !data.Scope.IsNull() && !data.Scope.IsUnknown() {
		var scopes []string
		diags.Append(data.Scope.ElementsAs(ctx, &scopes, false)...)
		if diags.HasError() {
			return nil
		}
		for _, s := range scopes {
			req.Scope = append(req.Scope, liff.LiffScope(s))
		}
	}

	if !data.Features.IsNull() && !data.Features.IsUnknown() {
		var features LiffAppFeaturesModel
		diags.Append(data.Features.As(ctx, &features, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil
		}
		req.Features = &liff.LiffFeatures{
			QrCode: features.QRCode.ValueBool(),
		}
	}

	return req
}

func mapLiffAppToModel(ctx context.Context, app *liff.LiffApp, data *LiffAppResourceModel, diags *diag.Diagnostics) {
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
		scopeValues := make([]attr.Value, len(app.Scope))
		for i, s := range app.Scope {
			scopeValues[i] = types.StringValue(string(s))
		}
		scopeList, d := types.ListValue(types.StringType, scopeValues)
		diags.Append(d...)
		data.Scope = scopeList
	} else {
		data.Scope = types.ListNull(types.StringType)
	}

	if app.View != nil {
		viewObj, d := types.ObjectValue(viewAttrTypes, map[string]attr.Value{
			"type":        types.StringValue(string(app.View.Type)),
			"url":         types.StringValue(app.View.Url),
			"module_mode": types.BoolValue(app.View.ModuleMode),
		})
		diags.Append(d...)
		data.View = viewObj
	}

	if app.Features != nil {
		featuresObj, d := types.ObjectValue(featuresAttrTypes, map[string]attr.Value{
			"ble":     types.BoolValue(app.Features.Ble),
			"qr_code": types.BoolValue(app.Features.QrCode),
		})
		diags.Append(d...)
		data.Features = featuresObj
	} else {
		featuresObj, d := types.ObjectValue(featuresAttrTypes, map[string]attr.Value{
			"ble":     types.BoolValue(false),
			"qr_code": types.BoolValue(false),
		})
		diags.Append(d...)
		data.Features = featuresObj
	}
}

func mapLiffAppToDataSourceModel(ctx context.Context, app *liff.LiffApp, diags *diag.Diagnostics) LiffAppDataModel {
	m := LiffAppDataModel{
		LiffID: types.StringValue(app.LiffId),
	}

	if app.Description != "" {
		m.Description = types.StringValue(app.Description)
	} else {
		m.Description = types.StringNull()
	}

	if app.PermanentLinkPattern != "" {
		m.PermanentLinkPattern = types.StringValue(app.PermanentLinkPattern)
	} else {
		m.PermanentLinkPattern = types.StringNull()
	}

	m.BotPrompt = types.StringValue(string(app.BotPrompt))

	if len(app.Scope) > 0 {
		scopeValues := make([]attr.Value, len(app.Scope))
		for i, s := range app.Scope {
			scopeValues[i] = types.StringValue(string(s))
		}
		scopeList, d := types.ListValue(types.StringType, scopeValues)
		diags.Append(d...)
		m.Scope = scopeList
	} else {
		m.Scope = types.ListNull(types.StringType)
	}

	if app.View != nil {
		viewObj, d := types.ObjectValue(viewAttrTypes, map[string]attr.Value{
			"type":        types.StringValue(string(app.View.Type)),
			"url":         types.StringValue(app.View.Url),
			"module_mode": types.BoolValue(app.View.ModuleMode),
		})
		diags.Append(d...)
		m.View = viewObj
	} else {
		m.View = types.ObjectNull(viewAttrTypes)
	}

	if app.Features != nil {
		featuresObj, d := types.ObjectValue(featuresAttrTypes, map[string]attr.Value{
			"ble":     types.BoolValue(app.Features.Ble),
			"qr_code": types.BoolValue(app.Features.QrCode),
		})
		diags.Append(d...)
		m.Features = featuresObj
	} else {
		featuresObj, d := types.ObjectValue(featuresAttrTypes, map[string]attr.Value{
			"ble":     types.BoolValue(false),
			"qr_code": types.BoolValue(false),
		})
		diags.Append(d...)
		m.Features = featuresObj
	}

	return m
}
