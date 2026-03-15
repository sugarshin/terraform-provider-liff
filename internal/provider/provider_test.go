// Copyright (c) sugarshin
// SPDX-License-Identifier: MIT

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"liff": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("LIFF_CHANNEL_ACCESS_TOKEN"); v != "" {
		return
	}

	if os.Getenv("LIFF_CHANNEL_ID") == "" {
		t.Fatal("LIFF_CHANNEL_ID or LIFF_CHANNEL_ACCESS_TOKEN must be set for acceptance tests")
	}
	if os.Getenv("LIFF_CHANNEL_SECRET") == "" {
		t.Fatal("LIFF_CHANNEL_SECRET must be set for acceptance tests")
	}
}
