package convert

import (
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func Test_StringArrayToInterfaceArray(t *testing.T) {
	s := []string{"111", "ddd", "aaa"}
	t.Log(StringArrayToInterfaceArray(s))
}

func Test_Float64ToString(t *testing.T) {
	p0 := Float64ToString(123456789.12)
	assertEqual(t, p0, "123456789")
	p1 := Float64ToString(123456789)
	assertEqual(t, p1, "123456789")
}

func Test_Float64ToStringPrecision(t *testing.T) {
	p0 := Float64ToStringPrecision(123456789.123456789, 0)
	assertEqual(t, p0, "123456789")
	p2 := Float64ToStringPrecision(123456789.123456789, 2)
	assertEqual(t, p2, "123456789.12")
	//4舍5入
	p5 := Float64ToStringPrecision(123456789.123456789, 5)
	assertEqual(t, p5, "123456789.12346")
}
