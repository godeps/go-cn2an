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
