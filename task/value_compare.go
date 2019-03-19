package task

import (
	"strconv"
	"strings"
)

func ValueCompare(field string, operate string, value string) bool {
	rs := false
	switch operate {
	case ">":
		rs = numberCompare(field, operate, value)
	case ">=":
		rs = numberCompare(field, operate, value)
	case "<":
		rs = numberCompare(field, operate, value)
	case "<=":
		rs = numberCompare(field, operate, value)
	case "<>":
		rs = numberCompare(field, operate, value)
	case "=":
		rs = numberCompare(field, operate, value)
	case "between": //1-5
		rs = numberBetween(field, value)
	case "in":
		rs = stringHasField(value, ",", field)
	case "not in":
		rs = !stringHasField(value, ",", field)
	}

	return rs
}

//filed与value转成float64后，基于op进行比较
func numberCompare(field string, operate string, value string) bool {
	rs := false

	//输入数据转float64
	fpf, _ := strconv.ParseFloat(field, 64)
	vpf, _ := strconv.ParseFloat(value, 64)

	switch operate {
	case ">":
		rs = fpf > vpf
	case ">=":
		rs = fpf >= vpf
	case "<":
		rs = fpf < vpf
	case "<=":
		rs = fpf <= vpf
	case "<>":
		rs = fpf != vpf
	case "=":
		rs = fpf == vpf
	}
	return rs
}

//field转成float64后，是否在value之间
func numberBetween(field string, value string) bool {
	rs := false
	ss := strings.Split(value, "-")
	fpf, _ := strconv.ParseFloat(field, 64)
	min, _ := strconv.ParseFloat(ss[0], 64)
	if len(ss) > 1 {
		max, _ := strconv.ParseFloat(ss[1], 64)
		rs = fpf >= min && fpf <= max
	} else {
		rs = fpf >= min
	}

	return rs
}

//字符串s基于sep割裂后，是否含有field项字符串
func stringHasField(s string, sep string, field string) bool {
	rs := false
	ss := strings.Split(s, sep)
	for _, v := range ss {
		if v == field {
			rs = true
			break
		}
	}
	return rs
}
