# go-cn2an

**go-cn2an** is a Go library that mirrors the behaviour of the Python project [cn2an](https://github.com/Ailln/cn2an). It converts numbers between Chinese numerals and Arabic numerals and supports sentence level transformations.

> 这是 go-cn2an 的英文文档。阅读中文版请参见 `README.md`。

## Features

- **Chinese → Arabic**: strict, normal, and smart modes, including mixed Chinese/Arabic inputs.
- **Arabic → Chinese**: lowercase, uppercase, RMB uppercase, and direct digit mapping modes.
- **Sentence transformation**:
  - Dates
  - Fractions
  - Percentages
  - Temperature (℃ / 摄氏度)
- Handles negatives and decimals (up to 16 fractional digits).
- Supports values in the range `[10^-16, 10^16]`.

### Supported Capabilities

| Category | Direction | Details / Examples |
| --- | --- | --- |
| Core numerals | Chinese → Arabic | Strict / normal / smart modes; uppercase numerals; colloquial forms; negatives and decimals; range `10^-16` ~ `10^16`. |
| Core numerals | Arabic → Chinese | low / up / rmb / direct modes; negatives; decimals; RMB wording. |
| Sentence transform | Chinese → Arabic | Automatically recognises dates, fractions, percentages, Celsius expressions, and colloquial numbers. |
| Sentence transform | Arabic → Chinese | Handles dates, fractions, percentages, Celsius. |
| Math symbol reading | Symbols → Chinese wording | `+`→`加` (plus), `-`→`减` (minus), `*`/`×`→`乘` (times), `÷`/`/`→`除以` (divide), `=`→`等于` (equals), `<`/`≤`/`>`/`≥` comparators, `!=`/`≠`→`不等于` (not equal), `±`/`∓`→`正负`/`负正`, `^`→`…次方`, `√`→`根号` (square root), `|x|`→`x的绝对值`, `∑`/`Sigma`→`求和`, `∫`→`积分`, `∞`→`无穷大`, `π`→`派`, `∂`→`偏导`, `∪`→`并集`, `∩`→`交集`, `∵`→`因为`, `∴`→`所以`. |

## Installation

```bash
go get github.com/godeps/go-cn2an
```

## Quick Start

```go
package main

import (
	"fmt"
	gocn2an "github.com/godeps/go-cn2an"
)

func main() {
	c := gocn2an.NewCn2An()
	res, _ := c.Cn2an("一百二十三点四五", "strict")
	fmt.Println(res) // 123.45

	a := gocn2an.NewAn2Cn()
	out, _ := a.An2cn(123.45, "rmb")
	fmt.Println(out) // 壹佰贰拾叁元肆角伍分
}
```

### Sentence Transform

```go
t := gocn2an.NewTransform()

cn, _ := t.Transform("小王的生日是2001年3月4日", "an2cn")
fmt.Println(cn) // 小王的生日是二零零一年三月四日

an, _ := t.Transform("创业板指九月九日早盘低开百分之一点五七", "cn2an")
fmt.Println(an) // 创业板指9月9日早盘低开1.57%
```

## Testing

```bash
go test ./...
```

## Example Program

```bash
go run example/main.go
```

## License

MIT

## Credits

- Original Python project: [Ailln/cn2an](https://github.com/Ailln/cn2an)
