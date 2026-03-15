// Copyright (c) sugarshin
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestLiffAppDataSource_Metadata(t *testing.T) {
	d := &LiffAppDataSource{}
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "liff"}, resp)

	if resp.TypeName != "liff_app" {
		t.Errorf("expected TypeName 'liff_app', got %q", resp.TypeName)
	}
}

func TestLiffAppDataSource_Schema(t *testing.T) {
	d := &LiffAppDataSource{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema

	liffID, ok := s.Attributes["liff_id"]
	if !ok {
		t.Fatal("missing liff_id attribute")
	}
	liffIDAttr, ok := liffID.(dschema.StringAttribute)
	if !ok {
		t.Fatal("liff_id should be StringAttribute")
	}
	if !liffIDAttr.Required {
		t.Error("liff_id should be required")
	}

	for _, name := range []string{"description", "permanent_link_pattern", "bot_prompt"} {
		attr, ok := s.Attributes[name]
		if !ok {
			t.Errorf("missing attribute %q", name)
			continue
		}
		sa, ok := attr.(dschema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q should be StringAttribute", name)
			continue
		}
		if !sa.Computed {
			t.Errorf("attribute %q should be computed", name)
		}
	}

	scope, ok := s.Attributes["scope"]
	if !ok {
		t.Fatal("missing scope attribute")
	}
	scopeAttr, ok := scope.(dschema.ListAttribute)
	if !ok {
		t.Fatal("scope should be ListAttribute")
	}
	if !scopeAttr.Computed {
		t.Error("scope should be computed")
	}

	view, ok := s.Attributes["view"]
	if !ok {
		t.Fatal("missing view attribute")
	}
	viewAttr, ok := view.(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("view should be SingleNestedAttribute")
	}
	if !viewAttr.Computed {
		t.Error("view should be computed")
	}

	features, ok := s.Attributes["features"]
	if !ok {
		t.Fatal("missing features attribute")
	}
	featuresAttr, ok := features.(dschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("features should be SingleNestedAttribute")
	}
	if !featuresAttr.Computed {
		t.Error("features should be computed")
	}
}

func TestLiffAppDataSource_Interfaces(t *testing.T) {
	d := NewLiffAppDataSource()
	if _, ok := d.(datasource.DataSourceWithConfigure); !ok {
		t.Error("should implement datasource.DataSourceWithConfigure")
	}
}
