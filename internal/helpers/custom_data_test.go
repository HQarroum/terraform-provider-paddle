package helpers

import (
	"testing"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCustomDataToMap(t *testing.T) {
	tests := []struct {
		name        string
		input       paddle.CustomData
		expected    map[string]attr.Value
		expectError bool
	}{
		{
			name:     "nil custom data",
			input:    nil,
			expected: nil,
		},
		{
			name: "string values",
			input: paddle.CustomData{
				"key1": "value1",
				"key2": "value2",
			},
			expected: map[string]attr.Value{
				"key1": types.StringValue("value1"),
				"key2": types.StringValue("value2"),
			},
		},
		{
			name: "numeric values",
			input: paddle.CustomData{
				"count":   float64(42),
				"enabled": true,
			},
			expected: map[string]attr.Value{
				"count":   types.StringValue("42"),
				"enabled": types.StringValue("true"),
			},
		},
		{
			name: "mixed types",
			input: paddle.CustomData{
				"string": "test",
				"number": float64(123),
				"bool":   false,
			},
			expected: map[string]attr.Value{
				"string": types.StringValue("test"),
				"number": types.StringValue("123"),
				"bool":   types.StringValue("false"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CustomDataToMap(tt.input)

			if tt.expectError && err == nil {
				t.Fatal("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expected == nil && result != nil {
				t.Fatalf("expected nil but got %v", result)
			}
			if tt.expected != nil && result == nil {
				t.Fatal("expected result but got nil")
			}

			if tt.expected != nil {
				if len(result) != len(tt.expected) {
					t.Fatalf("expected %d items but got %d", len(tt.expected), len(result))
				}

				for key, expectedVal := range tt.expected {
					actualVal, ok := result[key]
					if !ok {
						t.Fatalf("missing key %s in result", key)
					}
					if !expectedVal.Equal(actualVal) {
						t.Fatalf("for key %s: expected %v but got %v", key, expectedVal, actualVal)
					}
				}
			}
		})
	}
}

func TestMapToCustomData(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected paddle.CustomData
	}{
		{
			name:     "nil map",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: paddle.CustomData{},
		},
		{
			name: "string values",
			input: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			expected: paddle.CustomData{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "numeric string values",
			input: map[string]string{
				"count":   "42",
				"enabled": "true",
			},
			expected: paddle.CustomData{
				"count":   "42",
				"enabled": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapToCustomData(tt.input)

			if tt.expected == nil && result != nil {
				t.Fatalf("expected nil but got %v", result)
			}
			if tt.expected != nil && result == nil {
				t.Fatal("expected result but got nil")
			}

			if tt.expected != nil {
				if len(result) != len(tt.expected) {
					t.Fatalf("expected %d items but got %d", len(tt.expected), len(result))
				}

				for key, expectedVal := range tt.expected {
					actualVal, ok := result[key]
					if !ok {
						t.Fatalf("missing key %s in result", key)
					}
					if expectedVal != actualVal {
						t.Fatalf("for key %s: expected %v but got %v", key, expectedVal, actualVal)
					}
				}
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Test that converting map -> CustomData -> map gives the same result
	original := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	customData := MapToCustomData(original)
	result, err := CustomDataToMap(customData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(original) {
		t.Fatalf("expected %d items but got %d", len(original), len(result))
	}

	for key, originalVal := range original {
		resultVal, ok := result[key]
		if !ok {
			t.Fatalf("missing key %s in result", key)
		}
		expectedAttr := types.StringValue(originalVal)
		if !expectedAttr.Equal(resultVal) {
			t.Fatalf("for key %s: expected %v but got %v", key, expectedAttr, resultVal)
		}
	}
}
