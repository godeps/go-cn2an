package gocn2an

import (
	"testing"
)

func TestTransformStrictPairs(t *testing.T) {
	pairs := map[string]string{
		"小王捡了100块钱":         "小王捡了一百块钱",
		"用户增长最快的3个城市":       "用户增长最快的三个城市",
		"小王的生日是2001年3月4日":   "小王的生日是二零零一年三月四日",
		"小王的生日是2012年12月12日": "小王的生日是二零一二年十二月十二日",
		"今天股价上涨了8%":         "今天股价上涨了百分之八",
		"第2天股价下降了-3.8%":     "第二天股价下降了百分之负三点八",
		"抛出去的硬币为正面的概率是1/2":  "抛出去的硬币为正面的概率是二分之一",
		"现在室内温度为39℃，很热啊！":   "现在室内温度为三十九摄氏度，很热啊！",
		"创业板指9月9日早盘低开1.57%": "创业板指九月九日早盘低开百分之一点五七",
		"今年盈利增长率为12.34%":    "今年盈利增长率为百分之十二点三四",
		"实验成功率是0.5%":        "实验成功率是百分之零点五",
		"股票价格下跌了-7.25%":     "股票价格下跌了百分之负七点二五",
		"预计需要3/8的时间完成":      "预计需要八分之三的时间完成",
		"室外温度是-5℃":          "室外温度是-五摄氏度",
		"我们有2500个用户":        "我们有二千五百个用户",
		"连续发布3天":            "连续发布三天",
		"第10期节目":            "第十期节目",
	}

	tr := NewTransform()
	for arabic, chinese := range pairs {
		gotChinese, err := tr.Transform(arabic, "an2cn")
		if err != nil {
			t.Errorf("Transform(%q, an2cn) error: %v", arabic, err)
		} else if gotChinese != chinese {
			t.Errorf("Transform(%q, an2cn) = %s, want %s", arabic, gotChinese, chinese)
		}

		gotArabic, err := tr.Transform(chinese, "cn2an")
		if err != nil {
			t.Errorf("Transform(%q, cn2an) error: %v", chinese, err)
		} else if gotArabic != arabic {
			t.Errorf("Transform(%q, cn2an) = %s, want %s", chinese, gotArabic, arabic)
		}
	}
}

func TestTransformSmartCn2an(t *testing.T) {
	testData := map[string]string{
		"约2.5亿年~6500万年": "约250000000年~65000000年",
		"廿二日，日出东方":      "22日，日出东方",
		"大陆":            "大陆",
		"半斤":            "0.5斤",
		"两个":            "2个",
	}

	tr := NewTransform()
	for input, expected := range testData {
		result, err := tr.Transform(input, "cn2an")
		if err != nil {
			t.Errorf("Transform(%q, cn2an) error: %v", input, err)
			continue
		}
		if result != expected {
			t.Errorf("Transform(%q, cn2an) = %s, want %s", input, result, expected)
		}
	}
}

func TestTransformMathSymbols(t *testing.T) {
	tr := NewTransform()
	cases := map[string]string{
		"1+2-3":              "一加二减三",
		"x-y":                "x减y",
		"-y":                 "负y",
		"a/b":                "a除以b",
		"x>=y":               "x大于等于y",
		"x<=y":               "x小于等于y",
		"a!=b":               "a不等于b",
		"∑x_i":               "求和x_i",
		"∵A>0, ∴B>0":         "因为A大于零, 所以B大于零",
		"∫f(x)dx":            "积分f(x)dx",
		"A∪B":                "A并集B",
		"A∩B":                "A交集B",
		"sigma=0":            "求和等于零",
		"πr^2":               "派r的二次方",
		"∞>0":                "无穷大大于零",
		"go-cn2an":           "go-cn二an",
		"http://example.com": "http://example.com",
		"x=y":                "x等于y",
		"x^2":                "x的二次方",
		"√9":                 "根号九",
		"|x|":                "x的绝对值",
	}

	for input, expected := range cases {
		result, err := tr.Transform(input, "an2cn")
		if err != nil {
			t.Fatalf("Transform(%q, an2cn) error: %v", input, err)
		}
		if result != expected {
			t.Errorf("Transform(%q, an2cn) = %q, want %q", input, result, expected)
		}
	}
}

func TestTransformAn2cnDateWithSpaces(t *testing.T) {
	tr := NewTransform()
	input := "现在是2025 年 10 月 30 日"
	expected := "现在是二零二五年十月三十日"
	got, err := tr.Transform(input, "an2cn")
	if err != nil {
		t.Fatalf("Transform(%q, an2cn) error: %v", input, err)
	}
	if got != expected {
		t.Errorf("Transform(%q, an2cn) = %q, want %q", input, got, expected)
	}
}
