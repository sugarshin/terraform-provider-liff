// Copyright sugarshin 2026
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestLiffAppsDataSource_Metadata(t *testing.T) {
	d := &LiffAppsDataSource{}
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "liff"}, resp)

	if resp.TypeName != "liff_apps" {
		t.Errorf("expected TypeName 'liff_apps', got %q", resp.TypeName)
	}
}

func TestLiffAppsDataSource_Schema(t *testing.T) {
	d := &LiffAppsDataSource{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	s := resp.Schema
	apps, ok := s.Attributes["apps"]
	if !ok {
		t.Fatal("missing apps attribute")
	}

	appsAttr, ok := apps.(dschema.ListNestedAttribute)
	if !ok {
		t.Fatal("apps should be ListNestedAttribute")
	}
	if !appsAttr.Computed {
		t.Error("apps should be computed")
	}

	nested := appsAttr.NestedObject.Attributes
	for _, name := range []string{"liff_id", "description", "permanent_link_pattern", "bot_prompt", "scope", "view", "features"} {
		if _, ok := nested[name]; !ok {
			t.Errorf("missing nested attribute %q", name)
		}
	}
}

func TestLiffAppsDataSource_Interfaces(t *testing.T) {
	d := NewLiffAppsDataSource()
	if _, ok := d.(datasource.DataSourceWithConfigure); !ok {
		t.Error("should implement datasource.DataSourceWithConfigure")
	}
}
