package convert

import (
	"encoding/json"
	"bytes"
)

func JsonByteToMap(data []byte) (map[string]interface{}, error) {
	var m map[string]interface{}

	err := json.Unmarshal(data, &m)

	return m, err
}

func JsonByteToArray(data []byte) ([]interface{}, error) {
	var m []interface{}

	err := json.Unmarshal(data, &m)

	return m, err
}

func MapToJson(m map[string]interface{}) (string, error) {
	rs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}

func MapStringToJson(m map[string]interface{}) (string, error) {
	rs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}

func StructToJson(m map[string]interface{}) (string, error) {
	rs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}

func StructToString(in interface{}) (string, error) {
	m, err := StructToLowerMap(in)
	if err != nil {
		return "", err
	}

	rs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}


/**
* @StringJoin		高效字符串拼接
* @separator		string		拼接分隔符
* @args...			slice		参数
* @s[return]		string
*/
func StringJoin(separator string, args ...string) (s string) {
	b := bytes.Buffer{}
	for k, v := range args {
		b.WriteString(v)
		if k < len(args) - 1 {
			b.WriteString(separator)
		}
    }
    s = b.String()
	return
}
