package gocn2an

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Cn2An 中文数字转阿拉伯数字的转换器
type Cn2An struct {
	allNum          string
	allUnit         string
	strictCNNumber  map[string]string
	normalCNNumber  map[string]string
	checkKeyDict    map[string]string
	patternDict     map[string]map[string]*regexp.Regexp
	ac              *An2Cn
	modeList        []string
	yjfPattern      *regexp.Regexp
	pattern1        *regexp.Regexp
	ptnAllNum       *regexp.Regexp
	ptnSpeakingMode *regexp.Regexp
}

// NewCn2An 创建新的中文到阿拉伯数字转换器
func NewCn2An() *Cn2An {
	c := &Cn2An{
		strictCNNumber: StrictCNNumber,
		normalCNNumber: NormalCNNumber,
		modeList:       []string{"strict", "normal", "smart"},
	}

	// 构建 allNum 和 allUnit
	for r := range NumberCN2AN {
		c.allNum += string(r)
	}
	for r := range UnitCN2AN {
		c.allUnit += string(r)
	}

	// 构建 checkKeyDict
	c.checkKeyDict = map[string]string{
		"strict": c.buildCheckKey(c.strictCNNumber) + "点负",
		"normal": c.buildCheckKey(c.normalCNNumber) + "点负",
		"smart":  c.buildCheckKey(c.normalCNNumber) + "点负" + "01234567890.-",
	}

	// 构建模式字典
	c.patternDict = c.getPattern()

	// 创建 An2Cn 实例
	c.ac = NewAn2Cn()

	// 编译正则表达式
	c.yjfPattern = regexp.MustCompile(fmt.Sprintf(`^.*?[元圆][%s]角([%s]分)?$`, c.allNum, c.allNum))
	c.pattern1 = regexp.MustCompile(fmt.Sprintf(`^-?\d+(\.\d+)?[%s]?$`, c.allUnit))
	c.ptnAllNum = regexp.MustCompile(fmt.Sprintf(`^[%s]+$`, c.allNum))
	c.ptnSpeakingMode = regexp.MustCompile(fmt.Sprintf(`^([%s]{0,2}[%s])+[%s]$`, c.allNum, c.allUnit, c.allNum))

	return c
}

// buildCheckKey 构建检查键
func (c *Cn2An) buildCheckKey(cnNumber map[string]string) string {
	result := ""
	for _, v := range cnNumber {
		result += v
	}
	return result
}

// Cn2an 中文数字转阿拉伯数字的主函数
// inputs: 中文数字字符串
// mode: strict(严格), normal(正常), smart(智能)
func (c *Cn2An) Cn2an(inputs string, mode string) (float64, error) {
	if inputs == "" {
		return 0, errors.New("输入数据为空")
	}

	if !contains(c.modeList, mode) {
		return 0, fmt.Errorf("mode 仅支持 %v", c.modeList)
	}

	// 数据预处理
	inputs = c.preprocess(inputs)

	// 特殊转化 廿
	inputs = strings.ReplaceAll(inputs, "廿", "二十")

	// 检查输入数据是否有效
	sign, integerData, decimalData, isAllNum, specialValue, hasSpecialValue, err := c.checkInputDataIsValid(inputs, mode)
	if err != nil {
		return 0, err
	}

	if sign == 0 {
		return specialValue, nil
	}

	// smart 模式下的特殊情况（直接返回计算好的数值）
	if hasSpecialValue {
		return float64(sign) * specialValue, nil
	}

	var output float64
	if !isAllNum {
		if decimalData == "" {
			intVal, err := c.integerConvert(integerData)
			if err != nil {
				return 0, err
			}
			output = float64(intVal)
		} else {
			intVal, err := c.integerConvert(integerData)
			if err != nil {
				return 0, err
			}
			decVal, err := c.decimalConvert(decimalData)
			if err != nil {
				return 0, err
			}
			output = float64(intVal) + decVal
			// 修正精度问题
			output = roundToDecimal(output, len(decimalData))
		}
	} else {
		if decimalData == "" {
			intVal, err := c.directConvert(integerData)
			if err != nil {
				return 0, err
			}
			output = float64(intVal)
		} else {
			intVal, err := c.directConvert(integerData)
			if err != nil {
				return 0, err
			}
			decVal, err := c.decimalConvert(decimalData)
			if err != nil {
				return 0, err
			}
			output = float64(intVal) + decVal
			output = roundToDecimal(output, len(decimalData))
		}
	}

	return float64(sign) * output, nil
}

// preprocess 数据预处理（简化版，实际应该包括繁体转简体、全角转半角）
func (c *Cn2An) preprocess(s string) string {
	return normalizeText(s)
}

// getPattern 获取正则表达式模式
func (c *Cn2An) getPattern() map[string]map[string]*regexp.Regexp {
	// 整数严格检查
	_0 := "[零]"
	_1_9 := "[一二三四五六七八九]"
	_10_99 := fmt.Sprintf("%s?[十]%s?", _1_9, _1_9)
	_1_99 := fmt.Sprintf("(%s|%s)", _10_99, _1_9)
	_100_999 := fmt.Sprintf("(%s[百]([零]%s)?|%s[百]%s)", _1_9, _1_9, _1_9, _10_99)
	_1_999 := fmt.Sprintf("(%s|%s)", _100_999, _1_99)
	_1000_9999 := fmt.Sprintf("(%s[千]([零]%s)?|%s[千]%s)", _1_9, _1_99, _1_9, _100_999)
	_1_9999 := fmt.Sprintf("(%s|%s)", _1000_9999, _1_999)
	_10000_99999999 := fmt.Sprintf("(%s[万]([零]%s)?|%s[万]%s)", _1_9999, _1_999, _1_9999, _1000_9999)
	_1_99999999 := fmt.Sprintf("(%s|%s)", _10000_99999999, _1_9999)
	_100000000_9999999999999999 := fmt.Sprintf("(%s[亿]([零]%s)?|%s[亿]%s)", _1_99999999, _1_99999999, _1_99999999, _10000_99999999)
	_1_9999999999999999 := fmt.Sprintf("(%s|%s)", _100000000_9999999999999999, _1_99999999)
	strIntPattern := fmt.Sprintf("^(%s|%s)$", _0, _1_9999999999999999)
	norIntPattern := fmt.Sprintf("^(%s|%s)$", _0, _1_9999999999999999)

	strDecPattern := "^[零一二三四五六七八九]{0,15}[一二三四五六七八九]$"
	norDecPattern := "^[零一二三四五六七八九]{0,16}$"

	// 替换严格模式的字符
	for key, val := range c.strictCNNumber {
		strIntPattern = replacePattern(strIntPattern, key, val)
		strDecPattern = replacePattern(strDecPattern, key, val)
	}

	// 替换正常模式的字符
	for key, val := range c.normalCNNumber {
		norIntPattern = replacePattern(norIntPattern, key, val)
		norDecPattern = replacePattern(norDecPattern, key, val)
	}

	return map[string]map[string]*regexp.Regexp{
		"strict": {
			"int": regexp.MustCompile(strIntPattern),
			"dec": regexp.MustCompile(strDecPattern),
		},
		"normal": {
			"int": regexp.MustCompile(norIntPattern),
			"dec": regexp.MustCompile(norDecPattern),
		},
	}
}

// replacePattern 替换模式中的字符
func replacePattern(pattern, key, val string) string {
	return strings.ReplaceAll(pattern, key, val)
}

// copyNum 将数字字符串转换为中文
func (c *Cn2An) copyNum(num string) string {
	result := ""
	for _, ch := range num {
		n, _ := strconv.Atoi(string(ch))
		result += NumberLowAN2CN[n]
	}
	return result
}

// checkInputDataIsValid 检查输入数据是否有效
// 返回：sign(符号), integerData(整数部分), decimalData(小数部分), isAllNum(是否纯数字), specialValue(特殊值), hasSpecialValue(是否有特殊值), error
func (c *Cn2An) checkInputDataIsValid(checkData, mode string) (int, string, string, bool, float64, bool, error) {
	originalData := checkData
	hasDecimalPoint := false

	// 去除停用词
	stopWords := []string{"元整", "圆整", "元正", "圆正"}
	for _, word := range stopWords {
		if strings.HasSuffix(checkData, word) {
			checkData = checkData[:len(checkData)-len(word)]
		}
	}

	// 在非严格模式下去除元、圆
	if mode != "strict" {
		normalStopWords := []string{"圆", "元"}
		for _, word := range normalStopWords {
			checkData = strings.TrimSuffix(checkData, word)
		}
	}

	// 处理元角分
	if c.yjfPattern.MatchString(originalData) {
		checkData = strings.ReplaceAll(checkData, "元", "点")
		checkData = strings.ReplaceAll(checkData, "角", "")
		checkData = strings.ReplaceAll(checkData, "分", "")
	}

	// 处理特殊问法
	checkData = strings.ReplaceAll(checkData, "零十", "零一十")
	checkData = strings.ReplaceAll(checkData, "零百", "零一百")

	// 检查字符是否合法
	checkKeys := c.checkKeyDict[mode]
	for _, r := range checkData {
		if !strings.ContainsRune(checkKeys, r) {
			return 0, "", "", false, 0, false, fmt.Errorf("当前为%s模式，输入的数据不在转化范围内：%c", mode, r)
		}
	}

	// 确定正负号
	sign := 1
	if strings.HasPrefix(checkData, "负") {
		checkData = checkData[len("负"):]
		sign = -1
	}

	var integerData, decimalData string

	// 处理小数点
	if strings.Contains(checkData, "点") {
		hasDecimalPoint = true
		parts := strings.Split(checkData, "点")
		if len(parts) != 2 {
			return 0, "", "", false, 0, false, errors.New("数据中包含不止一个点")
		}
		integerData, decimalData = parts[0], parts[1]

		// smart 模式下转换阿拉伯数字
		if mode == "smart" {
			integerData = c.convertArabicInSmart(integerData)
			decimalData = c.convertArabicDecimalInSmart(decimalData)
			mode = "normal"
		}
	} else {
		integerData = checkData
		decimalData = ""

		// smart 模式处理
		if mode == "smart" {
			if c.pattern1.MatchString(integerData) {
				// 10.1万 或 10.1 这样的格式
				runes := []rune(integerData)
				if len(runes) > 0 {
					lastRune := runes[len(runes)-1]
					if val, ok := UnitCN2AN[lastRune]; ok {
						// 有单位
						numPart := string(runes[:len(runes)-1])
						if numPart == "" {
							return 0, "", "", false, 0, false, fmt.Errorf("不符合格式的数据：%s", checkData)
						}
						numVal, err := strconv.ParseFloat(numPart, 64)
						if err == nil {
							output := numVal * float64(val)
							return 0, "", "", false, output, true, nil
						}
					} else {
						// 没有单位，纯数字
						numVal, err := strconv.ParseFloat(integerData, 64)
						if err == nil {
							return 0, "", "", false, numVal, true, nil
						}
					}
				}
			}

			integerData = c.convertArabicInSmart(integerData)
			mode = "normal"
		}
	}

	// 验证整数部分
	if patterns, ok := c.patternDict[mode]; ok {
		if intPattern, ok := patterns["int"]; ok {
			if intPattern.MatchString(integerData) {
				if hasDecimalPoint {
					if decimalData == "" {
						return 0, "", "", false, 0, false, fmt.Errorf("不符合格式的数据：%s", checkData)
					}
					if decPattern, ok := patterns["dec"]; ok {
						if decPattern.MatchString(decimalData) {
							return sign, integerData, decimalData, false, 0, false, nil
						}
					}
				} else {
					return sign, integerData, decimalData, false, 0, false, nil
				}
			}
		}
	}

	// normal 模式的特殊处理
	if mode == "normal" {
		// 纯数模式：一二三
		if c.ptnAllNum.MatchString(integerData) {
			if hasDecimalPoint {
				if decimalData == "" {
					return 0, "", "", false, 0, false, fmt.Errorf("不符合格式的数据：%s", checkData)
				}
				if decPattern, ok := c.patternDict[mode]["dec"]; ok {
					if decPattern.MatchString(decimalData) {
						return sign, integerData, decimalData, true, 0, false, nil
					}
				}
			} else {
				return sign, integerData, decimalData, true, 0, false, nil
			}
		}

		// 口语模式：一万二
		if len(integerData) >= 3 && c.ptnSpeakingMode.MatchString(integerData) {
			// 找到最后一个单位字符
			runes := []rune(integerData)
			lastChar := runes[len(runes)-1]

			// 检查是否是数字字符（而非单位）
			if _, isNum := NumberCN2AN[lastChar]; isNum {
				// 倒数第二个字符应该是单位
				if len(runes) >= 2 {
					secondLastChar := runes[len(runes)-2]
					if val, ok := UnitCN2AN[secondLastChar]; ok {
						unit := UnitLowAN2CN[val/10]
						integerData = integerData + unit
						if hasDecimalPoint {
							if decimalData == "" {
								return 0, "", "", false, 0, false, fmt.Errorf("不符合格式的数据：%s", checkData)
							}
							if decPattern, ok := c.patternDict[mode]["dec"]; ok {
								if decPattern.MatchString(decimalData) {
									return sign, integerData, decimalData, false, 0, false, nil
								}
							}
						} else {
							return sign, integerData, decimalData, false, 0, false, nil
						}
					}
				}
			}
		}
	}

	return 0, "", "", false, 0, false, fmt.Errorf("不符合格式的数据：%s", checkData)
}

// convertArabicInSmart 在 smart 模式下转换阿拉伯数字为中文
func (c *Cn2An) convertArabicInSmart(s string) string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		result, err := c.ac.An2cn(match, "low")
		if err != nil {
			return match
		}
		return result
	})
}

// convertArabicDecimalInSmart 在 smart 模式下转换小数部分的阿拉伯数字
func (c *Cn2An) convertArabicDecimalInSmart(s string) string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		return c.copyNum(match)
	})
}

// integerConvert 转换整数部分
func (c *Cn2An) integerConvert(integerData string) (int64, error) {
	var output int64 = 0
	var unit int64 = 1
	var tenThousandUnit int64 = 1

	runes := []rune(integerData)
	for i := len(runes) - 1; i >= 0; i-- {
		cnNum := runes[i]

		// 数值
		if num, ok := NumberCN2AN[cnNum]; ok {
			output += int64(num) * unit
		} else if unitVal, ok := UnitCN2AN[cnNum]; ok {
			// 单位
			unit = unitVal

			// 判断万、亿
			if unit%10000 == 0 {
				if unit > tenThousandUnit {
					tenThousandUnit = unit
				} else {
					tenThousandUnit = unit * tenThousandUnit
					unit = tenThousandUnit
				}
			}

			if unit < tenThousandUnit {
				unit = unit * tenThousandUnit
			}

			// 如果是最后一个字符且是单位，需要加上单位值
			if i == 0 {
				output += unit
			}
		} else {
			return 0, fmt.Errorf("%c 不在转化范围内", cnNum)
		}
	}

	return output, nil
}

// decimalConvert 转换小数部分
func (c *Cn2An) decimalConvert(decimalData string) (float64, error) {
	lenDecimal := len([]rune(decimalData))
	if lenDecimal > 16 {
		decimalData = string([]rune(decimalData)[:16])
		lenDecimal = 16
	}

	if lenDecimal == 0 {
		return 0, nil
	}

	var builder strings.Builder
	builder.WriteString("0.")

	for _, r := range decimalData {
		num, ok := NumberCN2AN[r]
		if !ok {
			return 0, fmt.Errorf("%c 不在转化范围内", r)
		}
		builder.WriteByte(byte('0' + num))
	}

	val, err := strconv.ParseFloat(builder.String(), 64)
	if err != nil {
		return 0, err
	}

	return roundToDecimal(val, lenDecimal), nil
}

// directConvert 直接转换（纯数字模式：一二三 => 123）
func (c *Cn2An) directConvert(data string) (int64, error) {
	var output int64 = 0
	runes := []rune(data)
	lenData := len(runes)

	for i := lenData - 1; i >= 0; i-- {
		if num, ok := NumberCN2AN[runes[i]]; ok {
			output += int64(num) * int64(math.Pow(10, float64(lenData-i-1)))
		} else {
			return 0, fmt.Errorf("%c 不在转化范围内", runes[i])
		}
	}

	return output, nil
}

// 辅助函数

// contains 检查字符串切片是否包含某个字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// roundToDecimal 四舍五入到指定小数位数
func roundToDecimal(val float64, precision int) float64 {
	if precision < 0 {
		return val
	}
	format := fmt.Sprintf("%%.%df", precision)
	str := fmt.Sprintf(format, val)
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return val
	}
	return result
}
