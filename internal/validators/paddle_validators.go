package validators

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// TaxCategoryValidator validates Paddle tax categories
type TaxCategoryValidator struct{}

func (v TaxCategoryValidator) Description(ctx context.Context) string {
	return "must be one of the valid Paddle tax categories"
}

func (v TaxCategoryValidator) MarkdownDescription(ctx context.Context) string {
	return "must be one of: `standard`, `digital-goods`, `ebooks`, `implementation-services`, `professional-services`, `saas`, `software-programming-services`, `training-services`"
}

func (v TaxCategoryValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	validCategories := map[string]bool{
		"standard":                      true,
		"digital-goods":                 true,
		"ebooks":                        true,
		"implementation-services":       true,
		"professional-services":         true,
		"saas":                          true,
		"software-programming-services": true,
		"training-services":             true,
	}

	if !validCategories[value] {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Tax Category",
			fmt.Sprintf("Tax category must be one of: standard, digital-goods, ebooks, implementation-services, professional-services, saas, software-programming-services, training-services. Got: %s", value),
		)
	}
}

// DiscountTypeValidator validates Paddle discount types
type DiscountTypeValidator struct{}

func (v DiscountTypeValidator) Description(ctx context.Context) string {
	return "must be one of the valid Paddle discount types"
}

func (v DiscountTypeValidator) MarkdownDescription(ctx context.Context) string {
	return "must be one of: `percentage`, `flat`, `flat_per_seat`"
}

func (v DiscountTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	validTypes := map[string]bool{
		"percentage":    true,
		"flat":          true,
		"flat_per_seat": true,
	}

	if !validTypes[value] {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Discount Type",
			fmt.Sprintf("Discount type must be one of: percentage, flat, flat_per_seat. Got: %s", value),
		)
	}
}

// PercentageAmountValidator validates percentage amounts (0.01 to 100)
type PercentageAmountValidator struct{}

func (v PercentageAmountValidator) Description(ctx context.Context) string {
	return "percentage amount must be between 0.01 and 100"
}

func (v PercentageAmountValidator) MarkdownDescription(ctx context.Context) string {
	return "percentage amount must be between `0.01` and `100`"
}

func (v PercentageAmountValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Percentage Amount",
			fmt.Sprintf("Percentage amount must be a valid number. Got: %s", value),
		)
		return
	}

	if amount < 0.01 || amount > 100 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Percentage Amount",
			fmt.Sprintf("Percentage amount must be between 0.01 and 100. Got: %s", value),
		)
	}
}

// CurrencyCodeValidator validates ISO 4217 currency codes
type CurrencyCodeValidator struct{}

func (v CurrencyCodeValidator) Description(ctx context.Context) string {
	return "must be a valid three-letter ISO 4217 currency code"
}

func (v CurrencyCodeValidator) MarkdownDescription(ctx context.Context) string {
	return "must be a valid three-letter ISO 4217 currency code in uppercase (e.g., `USD`, `EUR`, `GBP`)"
}

func (v CurrencyCodeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Basic validation: must be exactly 3 uppercase letters
	matched, _ := regexp.MatchString("^[A-Z]{3}$", value)
	if !matched {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Currency Code",
			fmt.Sprintf("Currency code must be a three-letter uppercase ISO 4217 code (e.g., USD, EUR, GBP). Got: %s", value),
		)
	}
}

// CountryCodeValidator validates ISO 3166-1 alpha-2 country codes
type CountryCodeValidator struct{}

func (v CountryCodeValidator) Description(ctx context.Context) string {
	return "must be a valid two-letter ISO 3166-1 alpha-2 country code"
}

func (v CountryCodeValidator) MarkdownDescription(ctx context.Context) string {
	return "must be a valid two-letter ISO 3166-1 alpha-2 country code in uppercase (e.g., `US`, `GB`, `FR`)"
}

func (v CountryCodeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Basic validation: must be exactly 2 uppercase letters
	matched, _ := regexp.MatchString("^[A-Z]{2}$", value)
	if !matched {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Country Code",
			fmt.Sprintf("Country code must be a two-letter uppercase ISO 3166-1 alpha-2 code (e.g., US, GB, FR). Got: %s", value),
		)
	}
}

// TaxModeValidator validates Paddle tax modes
type TaxModeValidator struct{}

func (v TaxModeValidator) Description(ctx context.Context) string {
	return "must be one of the valid Paddle tax modes"
}

func (v TaxModeValidator) MarkdownDescription(ctx context.Context) string {
	return "must be one of: `account_setting`, `internal`, `external`"
}

func (v TaxModeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	validModes := map[string]bool{
		"account_setting": true,
		"internal":        true,
		"external":        true,
	}

	if !validModes[value] {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Tax Mode",
			fmt.Sprintf("Tax mode must be one of: account_setting, internal, external. Got: %s", value),
		)
	}
}

// IntervalValidator validates Paddle duration intervals
type IntervalValidator struct{}

func (v IntervalValidator) Description(ctx context.Context) string {
	return "must be one of the valid Paddle interval types"
}

func (v IntervalValidator) MarkdownDescription(ctx context.Context) string {
	return "must be one of: `day`, `week`, `month`, `year`"
}

func (v IntervalValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	validIntervals := map[string]bool{
		"day":   true,
		"week":  true,
		"month": true,
		"year":  true,
	}

	if !validIntervals[value] {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Interval",
			fmt.Sprintf("Interval must be one of: day, week, month, year. Got: %s", value),
		)
	}
}
