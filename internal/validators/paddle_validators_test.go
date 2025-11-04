package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestTaxCategoryValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid standard", "standard", false},
		{"valid digital-goods", "digital-goods", false},
		{"valid ebooks", "ebooks", false},
		{"valid implementation-services", "implementation-services", false},
		{"valid professional-services", "professional-services", false},
		{"valid saas", "saas", false},
		{"valid software-programming-services", "software-programming-services", false},
		{"valid training-services", "training-services", false},
		{"invalid category", "invalid", true},
		{"empty string", "", true},
	}

	v := TaxCategoryValidator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

func TestDiscountTypeValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid percentage", "percentage", false},
		{"valid flat", "flat", false},
		{"valid flat_per_seat", "flat_per_seat", false},
		{"invalid type", "invalid", true},
		{"empty string", "", true},
	}

	v := DiscountTypeValidator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

func TestCurrencyCodeValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid USD", "USD", false},
		{"valid EUR", "EUR", false},
		{"valid GBP", "GBP", false},
		{"lowercase usd", "usd", true},
		{"too short", "US", true},
		{"too long", "USDD", true},
		{"with numbers", "US1", true},
		{"empty string", "", true},
	}

	v := CurrencyCodeValidator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

func TestTaxModeValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid account_setting", "account_setting", false},
		{"valid internal", "internal", false},
		{"valid external", "external", false},
		{"invalid mode", "invalid", true},
		{"empty string", "", true},
	}

	v := TaxModeValidator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

func TestIntervalValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid day", "day", false},
		{"valid week", "week", false},
		{"valid month", "month", false},
		{"valid year", "year", false},
		{"invalid interval", "quarter", true},
		{"empty string", "", true},
	}

	v := IntervalValidator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}

func TestCountryCodeValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{"valid US", "US", false},
		{"valid GB", "GB", false},
		{"valid FR", "FR", false},
		{"lowercase us", "us", true},
		{"too short", "U", true},
		{"too long", "USA", true},
		{"with numbers", "U1", true},
		{"empty string", "", true},
	}

	v := CountryCodeValidator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}
