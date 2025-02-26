package abrha

import (
	"context"
	"strings"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/config"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestURLOverride(t *testing.T) {
	customEndpoint := "https://mock-api.internal.example.com/"

	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":        "12345",
		"api_endpoint": customEndpoint,
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil")
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	if client.BaseURL.String() != customEndpoint {
		t.Fatalf("Expected %s, got %s", customEndpoint, client.BaseURL.String())
	}
}

func TestURLDefault(t *testing.T) {
	rawProvider := Provider()
	raw := map[string]interface{}{
		"token": "12345",
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatal("Expected metadata, got nil")
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()
	if client.BaseURL.String() != "https://my.abrha.net/cserver/api" {
		t.Fatalf("Expected %s, got %s", "https://my.abrha.net/cserver/api", client.BaseURL.String())
	}
}

func diagnosticsToString(diags diag.Diagnostics) string {
	diagsAsStrings := make([]string, len(diags))
	for i, diag := range diags {
		diagsAsStrings[i] = diag.Summary
	}

	return strings.Join(diagsAsStrings, "; ")
}

func TestSpaceAPIDefaultEndpoint(t *testing.T) {
	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":             "12345",
		"spaces_access_id":  "abcdef",
		"spaces_secret_key": "xyzzy",
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil")
	}

	var client *session.Session
	var err error
	client, err = meta.(*config.CombinedConfig).SpacesClient("sfo2")
	if err != nil {
		t.Fatalf("Failed to create Spaces client: %s", err)
	}

	expectedEndpoint := "https://sfo2.my.abrha.com"
	if *client.Config.Endpoint != expectedEndpoint {
		t.Fatalf("Expected %s, got %s", expectedEndpoint, *client.Config.Endpoint)
	}
}

func TestSpaceAPIEndpointOverride(t *testing.T) {
	customSpacesEndpoint := "https://{{.Region}}.not-parspack-domain.com"

	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":             "12345",
		"spaces_endpoint":   customSpacesEndpoint,
		"spaces_access_id":  "abcdef",
		"spaces_secret_key": "xyzzy",
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatal("Expected metadata, got nil")
	}

	var client *session.Session
	var err error
	client, err = meta.(*config.CombinedConfig).SpacesClient("sfo2")
	if err != nil {
		t.Fatalf("Failed to create Spaces client: %s", err)
	}

	expectedEndpoint := "https://sfo2.not-parspack-domain.com"
	if *client.Config.Endpoint != expectedEndpoint {
		t.Fatalf("Expected %s, got %s", expectedEndpoint, *client.Config.Endpoint)
	}
}
