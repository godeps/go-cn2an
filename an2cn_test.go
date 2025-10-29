package gocn2an

import (
	"testing"
)

func TestAn2cnLow(t *testing.T) {
	testData := map[interface{}]string{
		0:                    "零",
		1:                    "一",
		11:                   "十一",
		1000000:              "一百万",
		1000054:              "一百万零五十四",
		31000054:             "三千一百万零五十四",
		9876543298765432:     "九千八百七十六万五千四百三十二亿九千八百七十六万五千四百三十二",
		10000000000000:       "十万亿",
		-1:                   "负一",
		-11:                  "负十一",
		0.000500050005005:    "零点零零零五零零零五零零零五零零五",
		0.00005:              "零点零零零零五",
		0.4321:               "零点四三二一",
		1000054.4321:         "一百万零五十四点四三二一",
		1.01:                 "一点零一",
		1.2:                  "一点二",
		0.01:                 "零点零一",
		-0.1:                 "负零点一",
		1.10:                 "一点一",
		12.0:                 "十二点零",
		2.0:                  "二点零",
		0.10:                 "零点一",
	}

	a := NewAn2Cn()
	for input, expected := range testData {
		result, err := a.An2cn(input, "low")
		if err != nil {
			t.Errorf("An2cn(%v, low) error: %v", input, err)
			continue
		}
		if result != expected {
			t.Errorf("An2cn(%v, low) = %s, want %s", input, result, expected)
		}
	}
}

func TestAn2cnUp(t *testing.T) {
	testData := map[interface{}]string{
		0:        "零",
		1:        "壹",
		11:       "壹拾壹",
		1000000:  "壹佰万",
		1000054:  "壹佰万零伍拾肆",
		31000054: "叁仟壹佰万零伍拾肆",
		-1:       "负壹",
		-11:      "负壹拾壹",
		0.00005:  "零点零零零零伍",
		0.4321:   "零点肆叁贰壹",
		1.01:     "壹点零壹",
		1.2:      "壹点贰",
		0.01:     "零点零壹",
		-0.1:     "负零点壹",
		1.10:     "壹点壹",
		12.0:     "壹拾贰点零",
		2.0:      "贰点零",
		0.10:     "零点壹",
	}

	a := NewAn2Cn()
	for input, expected := range testData {
		result, err := a.An2cn(input, "up")
		if err != nil {
			t.Errorf("An2cn(%v, up) error: %v", input, err)
			continue
		}
		if result != expected {
			t.Errorf("An2cn(%v, up) = %s, want %s", input, result, expected)
		}
	}
}

func TestAn2cnRmb(t *testing.T) {
	testData := map[interface{}]string{
		0:            "零元整",
		1:            "壹元整",
		11:           "壹拾壹元整",
		1000000:      "壹佰万元整",
		1000054:      "壹佰万零伍拾肆元整",
		31000054:     "叁仟壹佰万零伍拾肆元整",
		10000000000000: "壹拾万亿元整",
		-1:           "负壹元整",
		-11:          "负壹拾壹元整",
		0.00005:      "零元整",
		0.4321:       "肆角叁分",
		1000054.4321: "壹佰万零伍拾肆元肆角叁分",
		1.01:         "壹元零壹分",
		1.2:          "壹元贰角",
		0.01:         "壹分",
		-0.1:         "负壹角",
		1.10:         "壹元壹角",
		12.0:         "壹拾贰元整",
		2.0:          "贰元整",
		0.10:         "壹角",
	}

	a := NewAn2Cn()
	for input, expected := range testData {
		result, err := a.An2cn(input, "rmb")
		if err != nil {
			t.Errorf("An2cn(%v, rmb) error: %v", input, err)
			continue
		}
		if result != expected {
			t.Errorf("An2cn(%v, rmb) = %s, want %s", input, result, expected)
		}
	}
}

func TestAn2cnDirect(t *testing.T) {
	testData := map[interface{}]string{
		0:            "零",
		1:            "一",
		11:           "一一",
		1000000:      "一零零零零零零",
		1000054:      "一零零零零五四",
		31000054:     "三一零零零零五四",
		-1:           "负一",
		-11:          "负一一",
		0.00005:      "零点零零零零五",
		0.4321:       "零点四三二一",
		1000054.4321: "一零零零零五四点四三二一",
		1.01:         "一点零一",
		1.2:          "一点二",
		0.01:         "零点零一",
		1.10:         "一点一",
		12.0:         "一二点零",
		2.0:          "二点零",
		0.10:         "零点一",
	}

	a := NewAn2Cn()
	for input, expected := range testData {
		result, err := a.An2cn(input, "direct")
		if err != nil {
			t.Errorf("An2cn(%v, direct) error: %v", input, err)
			continue
		}
		if result != expected {
			t.Errorf("An2cn(%v, direct) = %s, want %s", input, result, expected)
		}
	}
}

func TestAn2cnError(t *testing.T) {
	errorData := []string{
		"123.1.1",
		"0.1零",
	}

	a := NewAn2Cn()
	for _, input := range errorData {
		_, err := a.An2cn(input, "low")
		if err == nil {
			t.Errorf("An2cn(%q, low) should return error but got nil", input)
		}
	}
}
