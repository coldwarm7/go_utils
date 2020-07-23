package convert

import (
	"encoding/json"
	"errors"
	"reflect"
)

// JsonStructsToMaps 可以用于缓存互转类型
/* 例如
cacheData, err := cache.Get(cacheKey)
if err != nil {
  return
}
if cacheData != nil {
  data, _ = convert.JsonStructsToMaps(cacheData)
  return
}
*/
func JsonStructsToMaps(in interface{}) (data []map[string]interface{}, err error) {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Slice {
		err = errors.New("param is not slice")
		return
	}
	l := v.Len()
	var m map[string]interface{}
	for i := 0; i < l; i++ {
		m, err = JsonStructToMap(v.Index(i).Interface())
		if err != nil {
			break
		}
		data = append(data, m)
	}

	return
}

func JsonStructToMap(s interface{}) (data map[string]interface{}, err error) {
	b, err := json.Marshal(s)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return
	}
	return
}
