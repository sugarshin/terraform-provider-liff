// Copyright sugarshin 2026
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLiffProvider_Metadata(t *testing.T) {
	p := &LiffProvider{version: "1.0.0"}
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), provider.MetadataRequest{}, resp)

	if resp.TypeName != "liff" {
		t.Errorf("expected TypeName 'liff', got %q", resp.TypeName)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected Version '1.0.0', got %q", resp.Version)
	}
}

func TestLiffProvider_Schema(t *testing.T) {
	p := &LiffProvider{}
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, resp)

	s := resp.Schema
	if s.Attributes == nil {
		t.Fatal("schema attributes should not be nil")
	}

	for _, attrName := range []string{"channel_id", "channel_secret", "channel_access_token"} {
		attr, ok := s.Attributes[attrName]
		if !ok {
			t.Errorf("missing attribute %q", attrName)
			continue
		}
		sa, ok := attr.(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q should be StringAttribute", attrName)
			continue
		}
		if !sa.Optional {
			t.Errorf("attribute %q should be optional", attrName)
		}
	}

	secretAttr, ok := s.Attributes["channel_secret"].(schema.StringAttribute)
	if !ok {
		t.Fatal("channel_secret should be StringAttribute")
	}
	if !secretAttr.Sensitive {
		t.Error("channel_secret should be sensitive")
	}

	tokenAttr, ok := s.Attributes["channel_access_token"].(schema.StringAttribute)
	if !ok {
		t.Fatal("channel_access_token should be StringAttribute")
	}
	if !tokenAttr.Sensitive {
		t.Error("channel_access_token should be sensitive")
	}
}

func TestLiffProvider_Resources(t *testing.T) {
	p := &LiffProvider{}
	resources := p.Resources(context.Background())
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
}

func TestLiffProvider_DataSources(t *testing.T) {
	p := &LiffProvider{}
	dataSources := p.DataSources(context.Background())
	if len(dataSources) != 2 {
		t.Fatalf("expected 2 data sources, got %d", len(dataSources))
	}
}

func TestNew(t *testing.T) {
	factory := New("test-version")
	p := factory()
	if p == nil {
		t.Fatal("provider factory returned nil")
	}

	lp, ok := p.(*LiffProvider)
	if !ok {
		t.Fatalf("expected *LiffProvider, got %T", p)
	}
	if lp.version != "test-version" {
		t.Errorf("expected version 'test-version', got %q", lp.version)
	}
}

func TestStringValueOrEnv(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		isNull   bool
		envKey   string
		envValue string
		want     string
	}{
		{
			name:  "value set",
			value: "direct-value",
			want:  "direct-value",
		},
		{
			name:     "null falls back to env",
			isNull:   true,
			envKey:   "TEST_LIFF_VALUE",
			envValue: "env-value",
			want:     "env-value",
		},
		{
			name:   "null no env returns empty",
			isNull: true,
			envKey: "TEST_LIFF_UNSET_VALUE",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}

			var v types.String
			if tt.isNull {
				v = types.StringNull()
			} else {
				v = types.StringValue(tt.value)
			}

			got := stringValueOrEnv(v, tt.envKey)
			if got != tt.want {
				t.Errorf("stringValueOrEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}
