package gocn2an

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// An2Cn 阿拉伯数字转中文数字的转换器
type An2Cn struct {
	allNum    string
	numberLow map[int]string
	numberUp  map[int]string
	modeList  []string
}

// NewAn2Cn 创建新的阿拉伯数字到中文转换器
func NewAn2Cn() *An2Cn {
	return &An2Cn{
		allNum:    "0123456789",
		numberLow: NumberLowAN2CN,
		numberUp:  NumberUpAN2CN,
		modeList:  []string{"low", "up", "rmb", "direct"},
	}
}

// An2cn 阿拉伯数字转中文数字的主函数
// inputs: 数字字符串或数字
// mode: low(小写), up(大写), rmb(人民币), direct(直接转换)
func (a *An2Cn) An2cn(inputs interface{}, mode string) (string, error) {
	if inputs == nil || inputs == "" {
		return "", errors.New("输入数据为空")
	}

	if !contains(a.modeList, mode) {
		return "", fmt.Errorf("mode 仅支持 %v", a.modeList)
	}

	// 转换为字符串
	var inputStr string
	switch v := inputs.(type) {
	case string:
		inputStr = v
	case int:
		inputStr = strconv.Itoa(v)
	case int64:
		inputStr = strconv.FormatInt(v, 10)
	case float64:
		inputStr = a.numberToString(v)
	case float32:
		inputStr = a.numberToString(float64(v))
	default:
		inputStr = fmt.Sprintf("%v", v)
	}

	// 数据预处理
	inputStr = a.preprocess(inputStr)

	// 检查输入是否有效
	if err := a.checkInputsIsValid(inputStr); err != nil {
		return "", err
	}

	// 判断正负
	sign := ""
	if strings.HasPrefix(inputStr, "-") {
		sign = "负"
		inputStr = inputStr[1:]
	}

	var output string

	if mode == "direct" {
		output = a.directConvert(inputStr)
	} else {
		// 切割整数和小数
		parts := strings.Split(inputStr, ".")
		if len(parts) == 1 {
			// 不包含小数
			integerData := parts[0]
			if mode == "rmb" {
				intOut, err := a.integerConvert(integerData, "up")
				if err != nil {
					return "", err
				}
				output = intOut + "元整"
			} else {
				intOut, err := a.integerConvert(integerData, mode)
				if err != nil {
					return "", err
				}
				output = intOut
			}
		} else if len(parts) == 2 {
			// 包含小数
			integerData, decimalData := parts[0], parts[1]
			if mode == "rmb" {
				intData, err := a.integerConvert(integerData, "up")
				if err != nil {
					return "", err
				}
				decData := a.decimalConvert(decimalData, "up")
				lenDecData := len([]rune(decData))

				if lenDecData == 0 {
					output = intData + "元整"
				} else if lenDecData == 1 {
					return "", fmt.Errorf("异常输出：%s", decData)
				} else if lenDecData == 2 {
					decRunes := []rune(decData)
					if decRunes[1] != '零' {
						if intData == "零" {
							output = string(decRunes[1]) + "角"
						} else {
							output = intData + "元" + string(decRunes[1]) + "角"
						}
					} else {
						output = intData + "元整"
					}
				} else {
					decRunes := []rune(decData)
					if decRunes[1] != '零' {
						if decRunes[2] != '零' {
							if intData == "零" {
								output = string(decRunes[1]) + "角" + string(decRunes[2]) + "分"
							} else {
								output = intData + "元" + string(decRunes[1]) + "角" + string(decRunes[2]) + "分"
							}
						} else {
							if intData == "零" {
								output = string(decRunes[1]) + "角"
							} else {
								output = intData + "元" + string(decRunes[1]) + "角"
							}
						}
					} else {
						if decRunes[2] != '零' {
							if intData == "零" {
								output = string(decRunes[2]) + "分"
							} else {
								output = intData + "元" + "零" + string(decRunes[2]) + "分"
							}
						} else {
							output = intData + "元整"
						}
					}
				}
			} else {
				intOut, err := a.integerConvert(integerData, mode)
				if err != nil {
					return "", err
				}
				decOut := a.decimalConvert(decimalData, mode)
				output = intOut + decOut
			}
		} else {
			return "", fmt.Errorf("输入格式错误：%s", inputStr)
		}
	}

	return sign + output, nil
}

// numberToString 将数字转换为字符串（处理科学记数法）
func (a *An2Cn) numberToString(number float64) string {
	str := strconv.FormatFloat(number, 'f', -1, 64)

	// 在 Go 中，FormatFloat 对于 12.0 会返回 "12"，为了与 Python 保持一致，
	// 需要保留 .0 结尾的小数信息。
	if !strings.Contains(str, ".") && math.Trunc(number) == number {
		str += ".0"
	}

	return str
}

// preprocess 数据预处理
func (a *An2Cn) preprocess(s string) string {
	return normalizeText(s)
}

// checkInputsIsValid 检查输入数据是否有效
func (a *An2Cn) checkInputsIsValid(checkData string) error {
	allCheckKeys := a.allNum + ".-"
	for _, r := range checkData {
		if !strings.ContainsRune(allCheckKeys, r) {
			return fmt.Errorf("输入的数据不在转化范围内：%c", r)
		}
	}
	return nil
}

// integerConvert 转换整数部分
func (a *An2Cn) integerConvert(integerData, mode string) (string, error) {
	var numeralList map[int]string
	var unitList []string

	if mode == "low" {
		numeralList = NumberLowAN2CN
		unitList = UnitLowOrderAN2CN
	} else if mode == "up" {
		numeralList = NumberUpAN2CN
		unitList = UnitUpOrderAN2CN
	} else {
		return "", fmt.Errorf("error mode: %s", mode)
	}

	// 去除前导零
	intVal, err := strconv.ParseInt(integerData, 10, 64)
	if err != nil {
		return "", err
	}
	integerData = strconv.FormatInt(intVal, 10)

	lenInteger := len(integerData)
	if lenInteger > len(unitList) {
		return "", fmt.Errorf("超出数据范围，最长支持 %d 位", len(unitList))
	}

	outputAn := ""
	for i, ch := range integerData {
		d := int(ch - '0')
		if d != 0 {
			outputAn += numeralList[d] + unitList[lenInteger-i-1]
		} else {
			// 在万、亿位置，即使是0也要加单位
			if (lenInteger-i-1)%4 == 0 {
				outputAn += numeralList[d] + unitList[lenInteger-i-1]
			}
			// 如果前面不是零，加零
			if i > 0 && !strings.HasSuffix(outputAn, "零") {
				outputAn += numeralList[d]
			}
		}
	}

	// 清理多余的零和单位
	outputAn = strings.ReplaceAll(outputAn, "零零", "零")
	outputAn = strings.ReplaceAll(outputAn, "零万", "万")
	outputAn = strings.ReplaceAll(outputAn, "零亿", "亿")
	outputAn = strings.ReplaceAll(outputAn, "亿万", "亿")
	outputAn = strings.Trim(outputAn, "零")

	// 解决「一十几」问题
	if strings.HasPrefix(outputAn, "一十") {
		outputAn = outputAn[len("一"):]
	}

	// 0-1 之间的小数
	if outputAn == "" {
		outputAn = "零"
	}

	return outputAn, nil
}

// decimalConvert 转换小数部分
func (a *An2Cn) decimalConvert(decimalData, mode string) string {
	lenDecimal := len(decimalData)

	if lenDecimal > 16 {
		decimalData = decimalData[:16]
		lenDecimal = 16
	}

	outputAn := ""
	if lenDecimal > 0 {
		outputAn = "点"
	}

	var numeralList map[int]string
	if mode == "low" {
		numeralList = NumberLowAN2CN
	} else if mode == "up" {
		numeralList = NumberUpAN2CN
	} else {
		return ""
	}

	for _, ch := range decimalData {
		d := int(ch - '0')
		outputAn += numeralList[d]
	}

	return outputAn
}

// directConvert 直接转换（每位数字单独转换）
func (a *An2Cn) directConvert(inputs string) string {
	output := ""
	for _, ch := range inputs {
		if ch == '.' {
			output += "点"
		} else {
			d := int(ch - '0')
			output += a.numberLow[d]
		}
	}
	return output
}

// roundFloat 四舍五入浮点数
func roundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
