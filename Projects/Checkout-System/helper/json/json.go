package json

import (
	"encoding/json"

	"Checkout-System/helper/error"
)

func Parse(v any) []byte {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		error.Error(err, "error marshalling config to JSON")
		return nil
	} else {
		return jsonData
	}
}
