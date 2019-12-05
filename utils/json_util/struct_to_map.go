package json_util

import (
	"encoding/json"
	"net/url"
)

// StructToStringMap 结构体转对应结构的字符串Map
// 只支持一层结构
func StructToStringMap(v interface{}, m *map[string]string) error {
	jsonBytes, _ := json.Marshal(v)
	return json.Unmarshal(jsonBytes, m)
}

// StructToFormData 结构体转对应结构的url.Values
// 只支持一层结构
func StructToFormData(v interface{}, formData *url.Values) error {
	m := make(map[string]string)
	err := StructToStringMap(v, &m)
	if err != nil {
		return err
	}
	for key, value := range m {
		formData.Set(key, value)
	}

	return nil
}
