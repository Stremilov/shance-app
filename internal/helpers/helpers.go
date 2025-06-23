package helpers

import (
	"encoding/json"
)

func GetFirst(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func MarshalPhotos(photos []string) string {
	b, _ := json.Marshal(photos)
	return string(b)
}