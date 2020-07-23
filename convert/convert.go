package convert

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

var structTag = "json"

//UintToByte 无符号整数转[]byte
func UintToByte(number uint64) []byte {
	numberStr := strconv.FormatUint(number, 10)
	return []byte(numberStr)
}

//ByteToUint []byte转无符号整数
func ByteToUint(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

func StringToInt64(str string) int64 {
	i64, _ := strconv.ParseInt(str, 10, 64)
	return i64
}

func StringToInt(str string) int {
	i64, _ := strconv.ParseInt(str, 10, 64)
	return int(i64)
}

func StringToFloat64(str string) float64 {
	f64, _ := strconv.ParseFloat(str, 64)
	return f64
}

func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

//注意，这个方法会删除小数。如需保留小数，使用Float64ToStringPrecision
func Float64ToString(i float64) string {
	return fmt.Sprintf("%.0f", i)
}

func Float64ToStringPrecision(i float64, precision int) string {
	return fmt.Sprintf("%."+IntToString(precision)+"f", i)
}

func AbsInt64(i int64) int64 {
	return int64(math.Abs(float64(i)))
}

//todo:这个方法是不是有点多余
func InterfaceToInt(i interface{}) int {
	number := InterfaceToInt64(i)
	return int(number)
}

func InterfaceToMaps(arr interface{}) (data []map[string]string, err error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		err = errors.New("toslice arr not slice")
		return
	}
	l := v.Len()
	var m map[string]string
	for i := 0; i < l; i++ {
		m, err = InterfaceToMap(v.Index(i).Interface())
		if err != nil {
			break
		}
		data = append(data, m)
	}
	return
}

func InterfaceToMap(in interface{}) (map[string]string, error) {
	v := reflect.ValueOf(in)

	if v.Kind() != reflect.Map {
		return nil, fmt.Errorf("ToMap only accepts map; got %T", v)
	}

	i := v.Interface()
	interMap := i.(map[string]interface{})
	out := map[string]string{}
	for key, value := range interMap {
		valueV := reflect.ValueOf(value)
		switch valueV.Kind() {
		case reflect.String:
			out[key] = value.(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out[key] = strconv.Itoa(value.(int))
		case reflect.Float64:
			out[key] = Float64ToString(value.(float64))
		default:
			out[key] = ""
		}

	}

	return out, nil
}

func MutipleInterfaceToMap(in interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(in)

	if v.Kind() != reflect.Map {
		return nil, fmt.Errorf("ToMap only accepts map; got %T", v)
	}

	i := v.Interface()
	interMap := i.(map[string]interface{})
	out := make(map[string]interface{})
	for key, value := range interMap {
		valueV := reflect.ValueOf(value)
		switch valueV.Kind() {
		case reflect.Map:
			out[key], _ = MutipleInterfaceToMap(value)
		case reflect.String:
			out[key] = value.(string)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out[key] = InterfaceToInt64(value)
		case reflect.Float32, reflect.Float64:
			out[key] = InterfaceToFloat64(value)
		case reflect.Slice:
			out[key] = value
		default:
			out[key] = ""
		}
	}

	return out, nil
}

// interface to map data
func StructToLowerMap(in interface{}) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		fieldName := strings.ToLower(fi.Name)
		switch v.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out[fieldName] = strconv.FormatInt(v.Field(i).Int(), 10)
		default:
			out[fieldName] = v.Field(i).String()
		}

	}
	return out, nil
}

// interface to map data
func StructToLowerInterface(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		fieldName := strings.ToLower(fi.Name)
		out[fieldName] = v.Field(i).Interface()
	}
	return out, nil
}

// StructToMapByTag //根据struct tag 转为map
func StructToMapByTag(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	kind := v.Kind()
	if kind == reflect.Ptr {
		v = v.Elem()
	}
	if kind != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		tag := fi.Tag.Get(structTag)
		if tag == "-" || tag == "" {
			continue
		}
		out[tag] = v.Field(i).Interface()
	}
	return out, nil
}

// interface to map datas
func StructsToLowerMaps(arr interface{}) (data []map[string]string, err error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		err = errors.New("toslice arr not slice")
		return
	}
	l := v.Len()
	var m map[string]string
	for i := 0; i < l; i++ {
		m, err = StructToLowerMap(v.Index(i).Interface())
		if err != nil {
			break
		}
		data = append(data, m)
	}
	return
}

// interface to map datas
func StructsToLowerInterfaces(arr interface{}) (data []map[string]interface{}, err error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		err = errors.New("toslice arr not slice")
		return
	}
	l := v.Len()
	var m map[string]interface{}
	for i := 0; i < l; i++ {
		m, err = StructToLowerInterface(v.Index(i).Interface())
		if err != nil {
			break
		}
		data = append(data, m)
	}
	return
}

func StringArrayToInterfaceArray(s []string) (interfaceArray []interface{}) {
	interfaceArray = make([]interface{}, len(s))
	for i, v := range s {
		interfaceArray[i] = v
	}
	return
}

func InterfaceToString(in interface{}) (out string) {
	v := reflect.ValueOf(in)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float64, reflect.Float32:
		out = Float64ToString(v.Float())
	case reflect.String:
		out = in.(string)
	default:
		out = ""
	}
	return
}

func InterfaceToInt64(i interface{}) (out int64) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = int64(v.Uint())
	case reflect.String:
		out, _ = strconv.ParseInt(v.String(), 10, 64)
	case reflect.Float64:
		out = int64(i.(float64))
	case reflect.Float32:
		out = int64(i.(float32))
	default:
		out = 0
	}
	return
}

func InterfaceToFloat64(i interface{}) (out float64) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Float64:
		out = i.(float64)
	case reflect.Float32:
		out = float64(i.(float32))
	case reflect.String:
		out, _ = strconv.ParseFloat(i.(string), 64)
	default:
		out = 0
	}
	return
}

//InterfaceToJSON 转json
func InterfaceToJSON(param interface{}) (map[string]interface{}, error) {
	paramString, ok := param.(string)
	if ok == false {
		return nil, fmt.Errorf("%s\t%v", "parse param into string faild", param)
	}

	listJSON, err := simplejson.NewJson([]byte(paramString))
	if err != nil {
		return nil, err
	}

	return listJSON.Map()
}

func StringToUint64(in string) (out uint64) {
	out, _ = strconv.ParseUint(in, 10, 64)
	return
}
