package gocn2an

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Transform 句子转换器
type Transform struct {
	allNum                 string
	allUnit                string
	cn2an                  *Cn2An
	an2cn                  *An2Cn
	cnPattern              string
	smartCnPattern         string
	cnPatternRe            *regexp.Regexp
	smartCnPatternRe       *regexp.Regexp
	mathSymbolReplacer     *strings.Replacer
	binaryMinusPlaceholder string
}

var exponentPattern = regexp.MustCompile(`([^\s\^]+)\s*\^\s*([^\s\^]+)`)

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

	mathPairs := []string{
		"<=", "小于等于",
		">=", "大于等于",
		"≤", "小于等于",
		"≦", "小于等于",
		"≥", "大于等于",
		"≧", "大于等于",
		"!=", "不等于",
		"≠", "不等于",
		"=", "等于",
		"＝", "等于",
		"<", "小于",
		"＜", "小于",
		">", "大于",
		"＞", "大于",
		"+", "加",
		"＋", "加",
		"×", "乘",
		"✕", "乘",
		"✖", "乘",
		"⋅", "乘",
		"·", "乘",
		"÷", "除以",
		"±", "正负",
		"∓", "负正",
		"∴", "所以",
		"∵", "因为",
		"∪", "并集",
		"∩", "交集",
		"∑", "求和",
		"Σ", "求和",
		"Sigma", "求和",
		"sigma", "求和",
		"∫", "积分",
		"∞", "无穷大",
		"π", "派",
		"√", "根号",
		"∂", "偏导",
	}
	t.mathSymbolReplacer = strings.NewReplacer(mathPairs...)
	t.binaryMinusPlaceholder = "@@__CNAN_MINUS__@@"

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
		inputs = t.preprocessAn2cnMathSymbols(inputs)

		// 日期
		dateRe := regexp.MustCompile(`(?:\d{2,4}\s*年\s*(?:\d{1,2}\s*月\s*)?(?:\d{1,2}\s*日)?)|(?:\d{1,2}\s*月\s*(?:\d{1,2}\s*日)?)|(?:\d{1,2}\s*日)`)
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

		output = t.postprocessAn2cnMathSymbols(output)

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
			normalized := strings.Join(strings.Fields(inputs), "")
			if normalized == "" {
				return inputs
			}
			yearRe := regexp.MustCompile(`\d+年`)
			result := yearRe.ReplaceAllStringFunc(normalized, func(match string) string {
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
				prefix := ""
				trimmed := result
				if strings.HasPrefix(result, "负") {
					prefix = "-"
					trimmed = strings.TrimPrefix(result, "负")
				}
				return prefix + trimmed + "摄氏度"
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

func (t *Transform) preprocessAn2cnMathSymbols(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	for i, r := range runes {
		runes[i] = normalizeMinusRune(r)
	}

	return t.markBinaryMinus(string(runes))
}

func normalizeMinusRune(r rune) rune {
	switch r {
	case '−', '﹣', '－':
		return '-'
	default:
		return r
	}
}

func (t *Transform) markBinaryMinus(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	var builder strings.Builder
	builder.Grow(len(runes) * 2)
	for i, r := range runes {
		if r == '-' {
			if t.isBinaryMinus(runes, i) {
				builder.WriteString(t.binaryMinusPlaceholder)
			} else {
				builder.WriteRune(r)
			}
			continue
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

func (t *Transform) isBinaryMinus(runes []rune, idx int) bool {
	prevIdx, prev := previousNonSpaceRune(runes, idx)
	nextIdx, next := nextNonSpaceRune(runes, idx)
	if prevIdx == -1 || nextIdx == -1 {
		return false
	}

	if !isBinaryMinusPrevOperand(prev) || !isBinaryMinusNextOperand(next) {
		return false
	}

	if isHyphenInsideWord(runes, idx, prevIdx, nextIdx) {
		return false
	}

	return true
}

func (t *Transform) postprocessAn2cnMathSymbols(s string) string {
	if s == "" {
		return s
	}

	s = strings.ReplaceAll(s, t.binaryMinusPlaceholder, "减")
	s = t.replaceEmbeddedNegativeBetweenOperands(s)
	s = t.replaceExponentNotation(s)
	s = t.replaceAbsoluteValue(s)
	if t.mathSymbolReplacer != nil {
		s = t.mathSymbolReplacer.Replace(s)
	}
	s = t.replaceSlashSymbols(s)
	s = t.replaceAsteriskSymbols(s)
	s = t.convertRemainingMinus(s)

	return s
}

func (t *Transform) replaceEmbeddedNegativeBetweenOperands(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(runes) * 2)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '负' {
			prevIdx, prev := previousNonSpaceRune(runes, i)
			nextIdx, next := nextNonSpaceRune(runes, i)
			if prevIdx != -1 && nextIdx != -1 && isMinusContextRune(prev) && isMinusContextRune(next) {
				builder.WriteString("减")
				continue
			}
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

func (t *Transform) replaceSlashSymbols(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(runes) * 2)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '/' {
			if t.isDivisionSlash(runes, i) {
				builder.WriteString("除以")
			} else {
				builder.WriteRune(r)
			}
			continue
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

func (t *Transform) replaceAsteriskSymbols(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(runes) * 2)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '*' {
			if t.isMultiplicationAsterisk(runes, i) {
				builder.WriteString("乘")
			} else {
				builder.WriteRune(r)
			}
			continue
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

func (t *Transform) replaceExponentNotation(s string) string {
	if !strings.ContainsRune(s, '^') {
		return s
	}

	for strings.ContainsRune(s, '^') {
		replaced := exponentPattern.ReplaceAllStringFunc(s, func(match string) string {
			subs := exponentPattern.FindStringSubmatch(match)
			if len(subs) != 3 {
				return match
			}
			base := strings.TrimSpace(subs[1])
			exponent := strings.TrimSpace(subs[2])
			if base == "" || exponent == "" {
				return match
			}
			return base + "的" + exponent + "次方"
		})
		if replaced == s {
			break
		}
		s = replaced
	}

	return s
}

func (t *Transform) replaceAbsoluteValue(s string) string {
	if !strings.ContainsRune(s, '|') {
		return s
	}

	runes := []rune(s)
	var builder strings.Builder
	builder.Grow(len(runes) * 2)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r != '|' {
			builder.WriteRune(r)
			continue
		}
		if (i+1 < len(runes) && runes[i+1] == '|') || (i-1 >= 0 && runes[i-1] == '|') {
			builder.WriteRune(r)
			continue
		}
		j := i + 1
		for j < len(runes) && runes[j] != '|' {
			j++
		}
		if j >= len(runes) {
			builder.WriteRune(r)
			continue
		}
		inner := strings.TrimSpace(string(runes[i+1 : j]))
		if inner == "" {
			builder.WriteRune(r)
			continue
		}
		builder.WriteString(inner)
		builder.WriteString("的绝对值")
		i = j
	}

	return builder.String()
}

func (t *Transform) convertRemainingMinus(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(runes) * 2)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '-' {
			prevIdx, prev := previousNonSpaceRune(runes, i)
			nextIdx, next := nextNonSpaceRune(runes, i)
			if prevIdx != -1 && nextIdx != -1 {
				if isHyphenInsideWord(runes, i, prevIdx, nextIdx) {
					builder.WriteRune(r)
					continue
				}
				if isBinaryMinusPrevOperand(prev) && isBinaryMinusNextOperand(next) {
					builder.WriteString("减")
					continue
				}
			}
			if nextIdx != -1 && isUnaryMinusOperand(next) {
				if prevIdx == -1 || isUnaryMinusPrefixRune(prev) {
					builder.WriteString("负")
					continue
				}
			}
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

func (t *Transform) isDivisionSlash(runes []rune, idx int) bool {
	prevIdx, prev := previousNonSpaceRune(runes, idx)
	nextIdx, next := nextNonSpaceRune(runes, idx)
	if prevIdx == -1 || nextIdx == -1 {
		return false
	}
	if runes[prevIdx] == '/' || runes[nextIdx] == '/' {
		return false
	}
	if prev == ':' {
		return false
	}
	if isURLSlash(runes, idx) {
		return false
	}
	if !isSlashOperand(prev) || !isSlashOperand(next) {
		return false
	}
	return true
}

func (t *Transform) isMultiplicationAsterisk(runes []rune, idx int) bool {
	prevIdx, prev := previousNonSpaceRune(runes, idx)
	nextIdx, next := nextNonSpaceRune(runes, idx)
	if prevIdx == -1 || nextIdx == -1 {
		return false
	}
	if runes[prevIdx] == '*' || runes[nextIdx] == '*' {
		return false
	}
	if !isSlashOperand(prev) || !isSlashOperand(next) {
		return false
	}
	return true
}

func previousNonSpaceRune(runes []rune, idx int) (int, rune) {
	for i := idx - 1; i >= 0; i-- {
		if !unicode.IsSpace(runes[i]) {
			return i, runes[i]
		}
	}
	return -1, 0
}

func nextNonSpaceRune(runes []rune, idx int) (int, rune) {
	for i := idx + 1; i < len(runes); i++ {
		if !unicode.IsSpace(runes[i]) {
			return i, runes[i]
		}
	}
	return -1, 0
}

func isAsciiLetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func isChineseNumeralRune(r rune) bool {
	if _, ok := NumberCN2AN[r]; ok {
		return true
	}
	switch r {
	case '点':
		return true
	}
	return false
}

func isGreekRune(r rune) bool {
	return unicode.In(r, unicode.Greek)
}

func isBinaryMinusPrevOperand(r rune) bool {
	if unicode.IsDigit(r) {
		return true
	}
	if isAsciiLetter(r) {
		return true
	}
	if isChineseNumeralRune(r) {
		return true
	}
	if isGreekRune(r) {
		return true
	}
	switch r {
	case ')', ']', '}', '℃', '％', '%', '°', '∞', 'π', '∑', 'Σ', '∫', '∂':
		return true
	}
	return false
}

func isBinaryMinusNextOperand(r rune) bool {
	if unicode.IsDigit(r) {
		return true
	}
	if isAsciiLetter(r) {
		return true
	}
	if isChineseNumeralRune(r) {
		return true
	}
	if isGreekRune(r) {
		return true
	}
	switch r {
	case '(', '[', '{':
		return true
	}
	return false
}

func isUnaryMinusOperand(r rune) bool {
	return isBinaryMinusNextOperand(r)
}

func isMinusContextRune(r rune) bool {
	if unicode.IsDigit(r) {
		return true
	}
	if isAsciiLetter(r) {
		return true
	}
	if isChineseNumeralRune(r) {
		return true
	}
	if isGreekRune(r) {
		return true
	}
	return false
}

func isSlashOperand(r rune) bool {
	if unicode.IsDigit(r) {
		return true
	}
	if isAsciiLetter(r) {
		return true
	}
	if isChineseNumeralRune(r) {
		return true
	}
	if isGreekRune(r) {
		return true
	}
	switch r {
	case ')', ']', '}', '(', '[', '{':
		return true
	}
	return false
}

func isURLSlash(runes []rune, idx int) bool {
	if idx > 0 && runes[idx-1] == '/' {
		return true
	}
	if idx+1 < len(runes) && runes[idx+1] == '/' {
		return true
	}
	return false
}

func isUnaryMinusPrefixRune(r rune) bool {
	switch r {
	case '+', '-', '−', '*', '×', '✕', '✖', '⋅', '·', '/', '÷', '=', '(', '[', '{', ',', '，', '。', '、', '；', '：':
		return true
	}
	return false
}

func isHyphenInsideWord(runes []rune, minusIdx, prevIdx, nextIdx int) bool {
	if prevIdx != minusIdx-1 || nextIdx != minusIdx+1 {
		return false
	}
	prev := runes[prevIdx]
	next := runes[nextIdx]
	if !isAsciiLetter(prev) || !isAsciiLetter(next) {
		return false
	}
	if prevIdx-1 >= 0 && isAsciiLetter(runes[prevIdx-1]) {
		return true
	}
	if nextIdx+1 < len(runes) && isAsciiLetter(runes[nextIdx+1]) {
		return true
	}
	return false
}
