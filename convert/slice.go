package convert

func SliceStringToSliceInterface(s []string) (i []interface{}) {
	for _, v := range s {
		i = append(i, v)
	}
	return
}

func InterfaceToSliceString(i interface{}) (s []string) {
	for _, v := range i.([]interface{}) {
		s = append(s, v.(string))
	}
	return
}
