package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestLiffAppResource_Metadata(t *testing.T) {
	r := &LiffAppResource{}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "liff"}, resp)

	if resp.TypeName != "liff_app" {
		t.Errorf("expected TypeName 'liff_app', got %q", resp.TypeName)
	}
}

func TestLiffAppResource_Schema(t *testing.T) {
	r := &LiffAppResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	s := resp.Schema
	if s.Attributes == nil {
		t.Fatal("schema attributes should not be nil")
	}

	// Check liff_id
	liffID, ok := s.Attributes["liff_id"]
	if !ok {
		t.Fatal("missing liff_id attribute")
	}
	liffIDAttr := liffID.(rschema.StringAttribute)
	if !liffIDAttr.Computed {
		t.Error("liff_id should be computed")
	}

	// Check description
	desc, ok := s.Attributes["description"]
	if !ok {
		t.Fatal("missing description attribute")
	}
	descAttr := desc.(rschema.StringAttribute)
	if !descAttr.Optional {
		t.Error("description should be optional")
	}

	// Check view (required)
	view, ok := s.Attributes["view"]
	if !ok {
		t.Fatal("missing view attribute")
	}
	viewAttr := view.(rschema.SingleNestedAttribute)
	if !viewAttr.Required {
		t.Error("view should be required")
	}
	if _, ok := viewAttr.Attributes["type"]; !ok {
		t.Error("view should have 'type' attribute")
	}
	if _, ok := viewAttr.Attributes["url"]; !ok {
		t.Error("view should have 'url' attribute")
	}
	if _, ok := viewAttr.Attributes["module_mode"]; !ok {
		t.Error("view should have 'module_mode' attribute")
	}

	// Check features (optional + computed)
	features, ok := s.Attributes["features"]
	if !ok {
		t.Fatal("missing features attribute")
	}
	featuresAttr := features.(rschema.SingleNestedAttribute)
	if !featuresAttr.Optional {
		t.Error("features should be optional")
	}
	if !featuresAttr.Computed {
		t.Error("features should be computed")
	}

	// Check ble is computed only
	ble := featuresAttr.Attributes["ble"].(rschema.BoolAttribute)
	if !ble.Computed {
		t.Error("features.ble should be computed")
	}
	if ble.Optional {
		t.Error("features.ble should NOT be optional (read-only)")
	}

	// Check bot_prompt default
	botPrompt, ok := s.Attributes["bot_prompt"]
	if !ok {
		t.Fatal("missing bot_prompt attribute")
	}
	botPromptAttr := botPrompt.(rschema.StringAttribute)
	if !botPromptAttr.Optional {
		t.Error("bot_prompt should be optional")
	}
	if !botPromptAttr.Computed {
		t.Error("bot_prompt should be computed")
	}

	// Check scope
	scope, ok := s.Attributes["scope"]
	if !ok {
		t.Fatal("missing scope attribute")
	}
	scopeAttr := scope.(rschema.ListAttribute)
	if !scopeAttr.Optional {
		t.Error("scope should be optional")
	}
	if !scopeAttr.Computed {
		t.Error("scope should be computed")
	}

	// Check permanent_link_pattern
	plp, ok := s.Attributes["permanent_link_pattern"]
	if !ok {
		t.Fatal("missing permanent_link_pattern attribute")
	}
	plpAttr := plp.(rschema.StringAttribute)
	if !plpAttr.Optional {
		t.Error("permanent_link_pattern should be optional")
	}
}

func TestLiffAppResource_Interfaces(t *testing.T) {
	r := NewLiffAppResource()
	if _, ok := r.(resource.Resource); !ok {
		t.Error("should implement resource.Resource")
	}
	if _, ok := r.(resource.ResourceWithImportState); !ok {
		t.Error("should implement resource.ResourceWithImportState")
	}
	if _, ok := r.(resource.ResourceWithConfigure); !ok {
		t.Error("should implement resource.ResourceWithConfigure")
	}
}
