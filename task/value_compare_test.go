package task_test

import (
	"audit-center/task"
	"testing"
)

func TestValueCompareTest(t *testing.T) {
	type fov struct {
		field   string
		operate string
		value   string
	}

	type uCase struct {
		in   fov
		want bool
	}

	var cases = []uCase{
		//<=
		{in: fov{"1.0", "<=", "2.0"}, want: true},
		{in: fov{"1.0", "<", "2.0"}, want: true},
		//>
		{in: fov{"1.0", ">", "2.0"}, want: false},
		{in: fov{"1.2", ">", "2.0"}, want: false},
		//>=
		{in: fov{"1.0", ">=", "2.0"}, want: false},
		{in: fov{"47.00", ">=", "12"}, want: true},
		//<>
		{in: fov{"1.10", "<>", "2.0"}, want: true},
		{in: fov{"1.0", "<>", "2.0"}, want: true},
		{in: fov{"1.0", "<>", "1.0"}, want: false},
		//between
		{in: fov{"1.0", "between", "1-2.0"}, want: true},
		{in: fov{"5.000000000001", "between", "3.1-5.1"}, want: true},
		//in
		{in: fov{"1.0", "in", "1.0,11"}, want: true},
		{in: fov{"GB", "in", "GBUK,GBFR"}, want: false},
		{in: fov{"GB", "in", "GB,"}, want: true},
		{in: fov{"GB", "in", "GB"}, want: true},
		{in: fov{"GB", "in", "GBUK,GB,GBFR"}, want: true},
	}

	for _, c := range cases {
		got := task.ValueCompare(c.in.field, c.in.operate, c.in.value)
		if got != c.want {
			t.Errorf("valueCompare(%q, %q, %q) == %v, want: %v", c.in.field, c.in.operate, c.in.value, got, c.want)
		}
	}

}
