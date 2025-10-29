package gocn2an

import "strings"

var normalizeRuneMap = map[rune]rune{
	'兩': '两',
	'貳': '贰',
	'參': '叁',
	'陸': '陆',
	'億': '亿',
	'萬': '万',
	'點': '点',
	'圓': '圆',
	'負': '负',
	'廿': '廿', // handled separately where needed
}

// normalizeText performs common preprocessing:
//  1. full-width ASCII -> half-width
//  2. traditional numerals/symbols -> simplified equivalents
func normalizeText(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 0xFF01 && r <= 0xFF5E:
			// Full-width ASCII -> half-width
			r -= 0xFEE0
		case r == 0x3000:
			// Full-width space
			r = 0x20
		}

		if mapped, ok := normalizeRuneMap[r]; ok {
			r = mapped
		}

		result.WriteRune(r)
	}
	return result.String()
}
