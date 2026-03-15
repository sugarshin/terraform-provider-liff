package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/line/line-bot-sdk-go/v8/linebot/liff"
)

func TestBuildAddLiffAppRequest(t *testing.T) {
	ctx := context.Background()

	viewObj, diags := types.ObjectValue(viewAttrTypes, map[string]attr.Value{
		"type":        types.StringValue("full"),
		"url":         types.StringValue("https://example.com"),
		"module_mode": types.BoolValue(false),
	})
	if diags.HasError() {
		t.Fatalf("creating view object: %v", diags.Errors())
	}

	featuresObj, diags := types.ObjectValue(featuresAttrTypes, map[string]attr.Value{
		"ble":     types.BoolValue(false),
		"qr_code": types.BoolValue(true),
	})
	if diags.HasError() {
		t.Fatalf("creating features object: %v", diags.Errors())
	}

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, []string{"openid", "profile"})
	if diags.HasError() {
		t.Fatalf("creating scope list: %v", diags.Errors())
	}

	data := &LiffAppResourceModel{
		Description:          types.StringValue("Test App"),
		PermanentLinkPattern: types.StringValue("concat"),
		BotPrompt:            types.StringValue("normal"),
		Scope:                scopeList,
		View:                 viewObj,
		Features:             featuresObj,
	}

	var d diag.Diagnostics
	req := buildAddLiffAppRequest(ctx, data, &d)
	if d.HasError() {
		t.Fatalf("buildAddLiffAppRequest: %v", d.Errors())
	}

	if req.View == nil {
		t.Fatal("view should not be nil")
	}
	if req.View.Type != liff.LiffViewTYPE_FULL {
		t.Errorf("expected view type FULL, got %q", req.View.Type)
	}
	if req.View.Url != "https://example.com" {
		t.Errorf("expected view URL 'https://example.com', got %q", req.View.Url)
	}
	if req.View.ModuleMode {
		t.Error("expected module_mode false")
	}
	if req.Description != "Test App" {
		t.Errorf("expected description 'Test App', got %q", req.Description)
	}
	if req.PermanentLinkPattern != "concat" {
		t.Errorf("expected permanent_link_pattern 'concat', got %q", req.PermanentLinkPattern)
	}
	if req.BotPrompt != liff.LiffBotPrompt_NORMAL {
		t.Errorf("expected bot_prompt 'normal', got %q", req.BotPrompt)
	}
	if len(req.Scope) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(req.Scope))
	}
	if req.Scope[0] != liff.LiffScope_OPENID {
		t.Errorf("expected first scope 'openid', got %q", req.Scope[0])
	}
	if req.Features == nil {
		t.Fatal("features should not be nil")
	}
	if !req.Features.QrCode {
		t.Error("expected qr_code true")
	}
}

func TestBuildAddLiffAppRequest_MinimalFields(t *testing.T) {
	ctx := context.Background()

	viewObj, diags := types.ObjectValue(viewAttrTypes, map[string]attr.Value{
		"type":        types.StringValue("compact"),
		"url":         types.StringValue("https://minimal.example.com"),
		"module_mode": types.BoolValue(false),
	})
	if diags.HasError() {
		t.Fatalf("creating view object: %v", diags.Errors())
	}

	data := &LiffAppResourceModel{
		Description:          types.StringNull(),
		PermanentLinkPattern: types.StringNull(),
		BotPrompt:            types.StringNull(),
		Scope:                types.ListNull(types.StringType),
		View:                 viewObj,
		Features:             types.ObjectNull(featuresAttrTypes),
	}

	var d diag.Diagnostics
	req := buildAddLiffAppRequest(ctx, data, &d)
	if d.HasError() {
		t.Fatalf("buildAddLiffAppRequest: %v", d.Errors())
	}

	if req.Description != "" {
		t.Errorf("expected empty description, got %q", req.Description)
	}
	if req.PermanentLinkPattern != "" {
		t.Errorf("expected empty permanent_link_pattern, got %q", req.PermanentLinkPattern)
	}
	if req.BotPrompt != "" {
		t.Errorf("expected empty bot_prompt, got %q", req.BotPrompt)
	}
	if req.Scope != nil {
		t.Errorf("expected nil scope, got %v", req.Scope)
	}
	if req.Features != nil {
		t.Errorf("expected nil features, got %v", req.Features)
	}
}

func TestBuildUpdateLiffAppRequest(t *testing.T) {
	ctx := context.Background()

	viewObj, diags := types.ObjectValue(viewAttrTypes, map[string]attr.Value{
		"type":        types.StringValue("tall"),
		"url":         types.StringValue("https://updated.example.com"),
		"module_mode": types.BoolValue(true),
	})
	if diags.HasError() {
		t.Fatalf("creating view object: %v", diags.Errors())
	}

	data := &LiffAppResourceModel{
		Description:          types.StringValue("Updated"),
		PermanentLinkPattern: types.StringNull(),
		BotPrompt:            types.StringValue("aggressive"),
		Scope:                types.ListNull(types.StringType),
		View:                 viewObj,
		Features:             types.ObjectNull(featuresAttrTypes),
	}

	var d diag.Diagnostics
	req := buildUpdateLiffAppRequest(ctx, data, &d)
	if d.HasError() {
		t.Fatalf("buildUpdateLiffAppRequest: %v", d.Errors())
	}

	if req.View == nil {
		t.Fatal("view should not be nil")
	}
	if req.View.Type != liff.UpdateLiffViewTYPE_TALL {
		t.Errorf("expected view type TALL, got %q", req.View.Type)
	}
	if !req.View.ModuleMode {
		t.Error("expected module_mode true")
	}
	if req.Description != "Updated" {
		t.Errorf("expected description 'Updated', got %q", req.Description)
	}
	if req.BotPrompt != liff.LiffBotPrompt_AGGRESSIVE {
		t.Errorf("expected bot_prompt 'aggressive', got %q", req.BotPrompt)
	}
}

func TestMapLiffAppToModel(t *testing.T) {
	ctx := context.Background()
	app := &liff.LiffApp{
		LiffId:      "1234567890-AbCdEfGh",
		Description: "Test App",
		View: &liff.LiffView{
			Type:       liff.LiffViewTYPE_FULL,
			Url:        "https://example.com",
			ModuleMode: true,
		},
		Features: &liff.LiffFeatures{
			Ble:    false,
			QrCode: true,
		},
		PermanentLinkPattern: "concat",
		BotPrompt:            liff.LiffBotPrompt_NORMAL,
		Scope:                []liff.LiffScope{liff.LiffScope_OPENID, liff.LiffScope_EMAIL},
	}

	var data LiffAppResourceModel
	var diags diag.Diagnostics
	mapLiffAppToModel(ctx, app, &data, &diags)

	if diags.HasError() {
		t.Fatalf("mapLiffAppToModel: %v", diags.Errors())
	}

	if data.LiffID.ValueString() != "1234567890-AbCdEfGh" {
		t.Errorf("expected liff_id '1234567890-AbCdEfGh', got %q", data.LiffID.ValueString())
	}
	if data.Description.ValueString() != "Test App" {
		t.Errorf("expected description 'Test App', got %q", data.Description.ValueString())
	}
	if data.PermanentLinkPattern.ValueString() != "concat" {
		t.Errorf("expected permanent_link_pattern 'concat', got %q", data.PermanentLinkPattern.ValueString())
	}
	if data.BotPrompt.ValueString() != "normal" {
		t.Errorf("expected bot_prompt 'normal', got %q", data.BotPrompt.ValueString())
	}

	// Check view
	if data.View.IsNull() {
		t.Fatal("view should not be null")
	}
	viewAttrs := data.View.Attributes()
	if viewAttrs["type"].(types.String).ValueString() != "full" {
		t.Errorf("expected view type 'full', got %q", viewAttrs["type"])
	}
	if viewAttrs["url"].(types.String).ValueString() != "https://example.com" {
		t.Errorf("expected view url 'https://example.com', got %q", viewAttrs["url"])
	}
	if !viewAttrs["module_mode"].(types.Bool).ValueBool() {
		t.Error("expected module_mode true")
	}

	// Check features
	if data.Features.IsNull() {
		t.Fatal("features should not be null")
	}
	featAttrs := data.Features.Attributes()
	if featAttrs["ble"].(types.Bool).ValueBool() {
		t.Error("expected ble false")
	}
	if !featAttrs["qr_code"].(types.Bool).ValueBool() {
		t.Error("expected qr_code true")
	}

	// Check scope
	if data.Scope.IsNull() {
		t.Fatal("scope should not be null")
	}
	scopeElems := data.Scope.Elements()
	if len(scopeElems) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(scopeElems))
	}
	if scopeElems[0].(types.String).ValueString() != "openid" {
		t.Errorf("expected first scope 'openid', got %q", scopeElems[0])
	}
}

func TestMapLiffAppToModel_EmptyFields(t *testing.T) {
	ctx := context.Background()
	app := &liff.LiffApp{
		LiffId: "1234567890-Empty",
		View: &liff.LiffView{
			Type: liff.LiffViewTYPE_COMPACT,
			Url:  "https://empty.example.com",
		},
		BotPrompt: liff.LiffBotPrompt_NONE,
	}

	var data LiffAppResourceModel
	var diags diag.Diagnostics
	mapLiffAppToModel(ctx, app, &data, &diags)

	if diags.HasError() {
		t.Fatalf("mapLiffAppToModel: %v", diags.Errors())
	}

	if !data.Description.IsNull() {
		t.Error("expected description to be null")
	}
	if !data.PermanentLinkPattern.IsNull() {
		t.Error("expected permanent_link_pattern to be null")
	}
	if !data.Scope.IsNull() {
		t.Error("expected scope to be null")
	}
	// Features should default to false/false when nil
	if data.Features.IsNull() {
		t.Fatal("features should not be null (should default)")
	}
}

func TestMapLiffAppToDataSourceModel(t *testing.T) {
	ctx := context.Background()
	app := &liff.LiffApp{
		LiffId:      "1234567890-DS",
		Description: "DS App",
		View: &liff.LiffView{
			Type:       liff.LiffViewTYPE_TALL,
			Url:        "https://ds.example.com",
			ModuleMode: false,
		},
		Features: &liff.LiffFeatures{
			Ble:    true,
			QrCode: false,
		},
		BotPrompt: liff.LiffBotPrompt_AGGRESSIVE,
		Scope:     []liff.LiffScope{liff.LiffScope_PROFILE},
	}

	var diags diag.Diagnostics
	m := mapLiffAppToDataSourceModel(ctx, app, &diags)

	if diags.HasError() {
		t.Fatalf("mapLiffAppToDataSourceModel: %v", diags.Errors())
	}

	if m.LiffID.ValueString() != "1234567890-DS" {
		t.Errorf("expected liff_id '1234567890-DS', got %q", m.LiffID.ValueString())
	}
	if m.Description.ValueString() != "DS App" {
		t.Errorf("expected description 'DS App', got %q", m.Description.ValueString())
	}
	if m.BotPrompt.ValueString() != "aggressive" {
		t.Errorf("expected bot_prompt 'aggressive', got %q", m.BotPrompt.ValueString())
	}

	featAttrs := m.Features.Attributes()
	if !featAttrs["ble"].(types.Bool).ValueBool() {
		t.Error("expected ble true")
	}
}
