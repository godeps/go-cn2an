package gocn2an

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Transform 句子转换器
type Transform struct {
	allNum           string
	allUnit          string
	cn2an            *Cn2An
	an2cn            *An2Cn
	cnPattern        string
	smartCnPattern   string
	cnPatternRe      *regexp.Regexp
	smartCnPatternRe *regexp.Regexp
}

// NewTransform 创建新的句子转换器
func NewTransform() *Transform {
	t := &Transform{
		allNum: "零一二三四五六七八九",
		cn2an:  NewCn2An(),
		an2cn:  NewAn2Cn(),
	}

	for r := range UnitCN2AN {
		t.allUnit += string(r)
	}

	t.cnPattern = fmt.Sprintf(`负?([%s%s]+点)?[%s%s]+`, t.allNum, t.allUnit, t.allNum, t.allUnit)
	t.smartCnPattern = fmt.Sprintf(`-?([0-9]+.)?[0-9]+[%s]+`, t.allUnit)
	t.cnPatternRe = regexp.MustCompile(t.cnPattern)
	t.smartCnPatternRe = regexp.MustCompile(t.smartCnPattern)

	return t
}

// Transform 转换句子中的数字
// inputs: 输入句子
// method: cn2an(中文转阿拉伯) 或 an2cn(阿拉伯转中文)
func (t *Transform) Transform(inputs, method string) (string, error) {
	if method == "cn2an" {
		inputs = strings.ReplaceAll(inputs, "廿", "二十")
		inputs = strings.ReplaceAll(inputs, "半", "0.5")
		inputs = strings.ReplaceAll(inputs, "两", "2")

		// 日期
		datePattern := fmt.Sprintf(`(((%s)|(%s))年)?([%s十]+月)?([%s十]+日)?`, t.smartCnPattern, t.cnPattern, t.allNum, t.allNum)
		dateRe := regexp.MustCompile(datePattern)
		inputs = dateRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "cn2an", "date")
		})

		// 分数
		fractionPattern := fmt.Sprintf(`%s分之%s`, t.cnPattern, t.cnPattern)
		fractionRe := regexp.MustCompile(fractionPattern)
		inputs = fractionRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "cn2an", "fraction")
		})

		// 百分比
		percentPattern := fmt.Sprintf(`百分之%s`, t.cnPattern)
		percentRe := regexp.MustCompile(percentPattern)
		inputs = percentRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "cn2an", "percent")
		})

		// 摄氏度
		celsiusPattern := fmt.Sprintf(`%s摄氏度`, t.cnPattern)
		celsiusRe := regexp.MustCompile(celsiusPattern)
		inputs = celsiusRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "cn2an", "celsius")
		})

		// 数字
		output := t.cnPatternRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "cn2an", "number")
		})

		return output, nil
	} else if method == "an2cn" {
		// 日期
		dateRe := regexp.MustCompile(`(\d{2,4}年)?(\d{1,2}月)?(\d{1,2}日)?`)
		inputs = dateRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "an2cn", "date")
		})

		// 分数
		fractionRe := regexp.MustCompile(`\d+/\d+`)
		inputs = fractionRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "an2cn", "fraction")
		})

		// 百分比
		percentRe := regexp.MustCompile(`-?(\d+\.)?\d+%`)
		inputs = percentRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "an2cn", "percent")
		})

		// 摄氏度
		celsiusRe := regexp.MustCompile(`\d+℃`)
		inputs = celsiusRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "an2cn", "celsius")
		})

		// 数字
		numberRe := regexp.MustCompile(`-?(\d+\.)?\d+`)
		output := numberRe.ReplaceAllStringFunc(inputs, func(match string) string {
			return t.subUtil(match, "an2cn", "number")
		})

		return output, nil
	}

	return "", fmt.Errorf("error method: %s, only support 'cn2an' and 'an2cn'", method)
}

// subUtil 替换辅助函数
func (t *Transform) subUtil(inputs, method, subMode string) string {
	if inputs == "" {
		return inputs
	}

	defer func() {
		if r := recover(); r != nil {
			// 发生错误时返回原始输入
		}
	}()

	if method == "cn2an" {
		switch subMode {
		case "date":
			// 匹配日期中的中文数字
			combinedPattern := fmt.Sprintf(`((%s)|(%s))`, t.smartCnPattern, t.cnPattern)
			re := regexp.MustCompile(combinedPattern)
			return re.ReplaceAllStringFunc(inputs, func(match string) string {
				result, err := t.cn2an.Cn2an(match, "smart")
				if err != nil {
					return match
				}
				// 转换为整数字符串
				return fmt.Sprintf("%.0f", result)
			})

		case "fraction":
			if strings.HasPrefix(inputs, "百") {
				return inputs
			}
			result := t.cnPatternRe.ReplaceAllStringFunc(inputs, func(match string) string {
				val, err := t.cn2an.Cn2an(match, "smart")
				if err != nil {
					return match
				}
				return fmt.Sprintf("%.0f", val)
			})
			parts := strings.Split(result, "分之")
			if len(parts) == 2 {
				return fmt.Sprintf("%s/%s", parts[1], parts[0])
			}
			return inputs

		case "percent":
			if !strings.HasPrefix(inputs, "百分之") {
				return inputs
			}
			target := strings.TrimPrefix(inputs, "百分之")
			val, err := t.cn2an.Cn2an(target, "smart")
			if err != nil {
				return inputs
			}
			return strconv.FormatFloat(val, 'f', -1, 64) + "%"

		case "celsius":
			if !strings.HasSuffix(inputs, "摄氏度") {
				return inputs
			}
			target := strings.TrimSuffix(inputs, "摄氏度")
			val, err := t.cn2an.Cn2an(target, "smart")
			if err != nil {
				return inputs
			}
			return strconv.FormatFloat(val, 'f', -1, 64) + "℃"

		case "number":
			val, err := t.cn2an.Cn2an(inputs, "smart")
			if err != nil {
				return inputs
			}
			return strconv.FormatFloat(val, 'f', -1, 64)
		}
	} else if method == "an2cn" {
		switch subMode {
		case "date":
			yearRe := regexp.MustCompile(`\d+年`)
			result := yearRe.ReplaceAllStringFunc(inputs, func(match string) string {
				digits := strings.TrimSuffix(match, "年")
				if digits == "" {
					return match
				}
				val, err := t.an2cn.An2cn(digits, "direct")
				if err != nil {
					return match
				}
				return val + "年"
			})
			// 月日用 low 模式
			numRe := regexp.MustCompile(`\d+`)
			return numRe.ReplaceAllStringFunc(result, func(match string) string {
				result, err := t.an2cn.An2cn(match, "low")
				if err != nil {
					return match
				}
				return result
			})

		case "fraction":
			numRe := regexp.MustCompile(`\d+`)
			result := numRe.ReplaceAllStringFunc(inputs, func(match string) string {
				cnNum, err := t.an2cn.An2cn(match, "low")
				if err != nil {
					return match
				}
				return cnNum
			})
			parts := strings.Split(result, "/")
			if len(parts) == 2 {
				return fmt.Sprintf("%s分之%s", parts[1], parts[0])
			}
			return inputs

		case "celsius":
			if strings.HasSuffix(inputs, "℃") {
				numPart := inputs[:len(inputs)-len("℃")]
				result, err := t.an2cn.An2cn(numPart, "low")
				if err != nil {
					return inputs
				}
				return result + "摄氏度"
			}
			return inputs

		case "percent":
			if strings.HasSuffix(inputs, "%") {
				numPart := inputs[:len(inputs)-1]
				result, err := t.an2cn.An2cn(numPart, "low")
				if err != nil {
					return inputs
				}
				return "百分之" + result
			}
			return inputs

		case "number":
			result, err := t.an2cn.An2cn(inputs, "low")
			if err != nil {
				return inputs
			}
			return result
		}
	}

	return inputs
}
