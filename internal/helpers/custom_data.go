package helpers

import (
	"encoding/json"
	"fmt"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Safely converts Paddle CustomData to a Terraform map.
// Paddle `CustomData` can contain any JSON-serializable values
// We convert everything to strings for simplicity in Terraform.
func CustomDataToMap(customData paddle.CustomData) (map[string]attr.Value, error) {
	if customData == nil {
		return nil, nil
	}

	result := make(map[string]attr.Value)
	for key, value := range customData {
		// Convert value to string representation
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case float64, int, int64, bool:
			// Convert numbers and booleans to string
			strValue = fmt.Sprintf("%v", v)
		default:
			// For complex types (maps, arrays), marshal to JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal custom_data value for key %s: %w", key, err)
			}
			strValue = string(jsonBytes)
		}
		result[key] = types.StringValue(strValue)
	}

	return result, nil
}

// Converts a Terraform map to Paddle `CustomData`.
func MapToCustomData(tfMap map[string]string) paddle.CustomData {
	if tfMap == nil {
		return nil
	}

	customData := make(paddle.CustomData)
	for key, value := range tfMap {
		customData[key] = value
	}
	return customData
}
