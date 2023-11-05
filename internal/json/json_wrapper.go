package json

import "encoding/json"

type RawJsonArray []json.RawMessage

func IsValidJson(str string) bool {
	return json.Valid([]byte(str))
}

func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
