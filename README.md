# go-cn2an

📦 **go-cn2an** 是一个快速转化 `中文数字` 和 `阿拉伯数字` 的 Golang 工具包！

这是 Python 版本 [cn2an](https://github.com/Ailln/cn2an) 的 Golang 完整复刻实现。

> For the English README please refer to `README_EN.md`.

## 功能特性

### 1. `中文数字` => `阿拉伯数字`

- 支持 `中文数字` => `阿拉伯数字`
- 支持 `大写中文数字` => `阿拉伯数字`
- 支持 `中文数字和阿拉伯数字混合` => `阿拉伯数字`

### 2. `阿拉伯数字` => `中文数字`

- 支持 `阿拉伯数字` => `中文数字`
- 支持 `阿拉伯数字` => `大写中文数字`
- 支持 `阿拉伯数字` => `大写人民币`

### 3. 句子转化

- 支持 `中文数字` => `阿拉伯数字`
  - 支持 `日期`
  - 支持 `分数`
  - 支持 `百分比`
  - 支持 `摄氏度`

- 支持 `阿拉伯数字` => `中文数字`
  - 支持 `日期`
  - 支持 `分数`
  - 支持 `百分比`
  - 支持 `摄氏度`

### 4. 其他

- 支持 `小数`
- 支持 `负数`
- 支持最大到 `10^16`（千万亿）
- 支持最小到 `10^-16`

### 支持能力一览

| 分类 | 转换方向 | 支持内容 / 示例 |
| --- | --- | --- |
| 基础数字 | 中文 → 阿拉伯 | 严格 / normal / smart 模式；支持大写、口语化、负数、小数，范围 `10^-16` ~ `10^16` |
| 基础数字 | 阿拉伯 → 中文 | low / up / rmb / direct 模式；支持负数、小数、人民币金额描述 |
| 句子转换 | 中文 → 阿拉伯 | 自动识别日期、分数、百分比、摄氏度、口语化表达 |
| 句子转换 | 阿拉伯 → 中文 | 自动识别日期、分数、百分比、摄氏度 |
| 数学符号读法 | 符号 → 中文描述 | `+`→`加`、`-`→`减`、`*`/`×`→`乘`、`/`/`÷`→`除以`、`=`→`等于`、`<`/`≤`/`>`/`≥`→比较、`!=`/`≠`→`不等于`、`±`/`∓`→`正负`/`负正`、`^`→`…次方`、`√`→`根号`、`|x|`→`x的绝对值`、`∑`/`Sigma`→`求和`、`∫`→`积分`、`∞`→`无穷大`、`π`→`派`、`∂`→`偏导`、`∪`→`并集`、`∩`→`交集`、`∵`→`因为`、`∴`→`所以` |

## 安装

```bash
go get github.com/godeps/go-cn2an
```

## 使用方法

### 1. 中文数字 => 阿拉伯数字

```go
package main

import (
    "fmt"
    gocn2an "github.com/godeps/go-cn2an"
)

func main() {
    c := gocn2an.NewCn2An()

    // strict 模式（默认）：只有严格符合数字拼写的才可以进行转化
    result, _ := c.Cn2an("一百二十三", "strict")
    fmt.Println(result) // 123

    // normal 模式：可以将 一二三 进行转化
    result, _ = c.Cn2an("一二三", "normal")
    fmt.Println(result) // 123

    // smart 模式：可以将混合拼写的 1百23 进行转化
    result, _ = c.Cn2an("1百23", "smart")
    fmt.Println(result) // 123

    // 支持负数
    result, _ = c.Cn2an("负一百二十三", "strict")
    fmt.Println(result) // -123

    // 支持小数
    result, _ = c.Cn2an("一点二三", "strict")
    fmt.Println(result) // 1.23

    // 支持口语化表达
    result, _ = c.Cn2an("一万二", "normal")
    fmt.Println(result) // 12000
}
```

### 2. 阿拉伯数字 => 中文数字

```go
package main

import (
    "fmt"
    gocn2an "github.com/godeps/go-cn2an"
)

func main() {
    a := gocn2an.NewAn2Cn()

    // low 模式（默认）：数字转化为小写的中文数字
    result, _ := a.An2cn(123, "low")
    fmt.Println(result) // 一百二十三

    // up 模式：数字转化为大写的中文数字
    result, _ = a.An2cn(123, "up")
    fmt.Println(result) // 壹佰贰拾叁

    // rmb 模式：数字转化为人民币专用的描述
    result, _ = a.An2cn(123, "rmb")
    fmt.Println(result) // 壹佰贰拾叁元整

    // direct 模式：直接转换每一位数字
    result, _ = a.An2cn(123, "direct")
    fmt.Println(result) // 一二三

    // 支持负数
    result, _ = a.An2cn(-123, "low")
    fmt.Println(result) // 负一百二十三

    // 支持小数
    result, _ = a.An2cn(1.23, "low")
    fmt.Println(result) // 一点二三

    // 支持人民币小数
    result, _ = a.An2cn(123.45, "rmb")
    fmt.Println(result) // 壹佰贰拾叁元肆角伍分
}
```

### 3. 句子转化

```go
package main

import (
    "fmt"
    gocn2an "github.com/godeps/go-cn2an"
)

func main() {
    t := gocn2an.NewTransform()

    // cn2an 方法：将句子中的中文数字转成阿拉伯数字
    result, _ := t.Transform("小王捡了一百块钱", "cn2an")
    fmt.Println(result) // 小王捡了100块钱

    // an2cn 方法：将句子中的阿拉伯数字转成中文数字
    result, _ = t.Transform("小王捡了100块钱", "an2cn")
    fmt.Println(result) // 小王捡了一百块钱

    // 支持日期
    result, _ = t.Transform("小王的生日是二零零一年三月四日", "cn2an")
    fmt.Println(result) // 小王的生日是2001年3月4日

    result, _ = t.Transform("小王的生日是2001年3月4日", "an2cn")
    fmt.Println(result) // 小王的生日是二零零一年三月四日

    // 支持分数
    result, _ = t.Transform("抛出去的硬币为正面的概率是二分之一", "cn2an")
    fmt.Println(result) // 抛出去的硬币为正面的概率是1/2

    result, _ = t.Transform("抛出去的硬币为正面的概率是1/2", "an2cn")
    fmt.Println(result) // 抛出去的硬币为正面的概率是二分之一
}
```

## API 说明

### Cn2an 中文转阿拉伯数字

```go
func (c *Cn2An) Cn2an(inputs string, mode string) (float64, error)
```

**参数：**
- `inputs`: 中文数字字符串
- `mode`: 转换模式
  - `"strict"`: 严格模式，只支持标准的中文数字表达
  - `"normal"`: 普通模式，支持口语化表达（如"一万二"）
  - `"smart"`: 智能模式，支持中文数字和阿拉伯数字混合

**返回：**
- `float64`: 转换后的阿拉伯数字
- `error`: 错误信息

### An2cn 阿拉伯数字转中文

```go
func (a *An2Cn) An2cn(inputs interface{}, mode string) (string, error)
```

**参数：**
- `inputs`: 阿拉伯数字（可以是 int, int64, float64, float32 或 string）
- `mode`: 转换模式
  - `"low"`: 小写中文数字
  - `"up"`: 大写中文数字
  - `"rmb"`: 人民币大写
  - `"direct"`: 直接转换（每位数字单独转换）

**返回：**
- `string`: 转换后的中文数字
- `error`: 错误信息

### Transform 句子转换

```go
func (t *Transform) Transform(inputs string, method string) (string, error)
```

**参数：**
- `inputs`: 输入句子
- `method`: 转换方法
  - `"cn2an"`: 中文数字转阿拉伯数字
  - `"an2cn"`: 阿拉伯数字转中文数字

**返回：**
- `string`: 转换后的句子
- `error`: 错误信息

## 运行测试

```bash
go test -v
```

## 运行示例

```bash
cd example
go run main.go
```

## 项目结构

```
go-cn2an/
├── config.go          # 配置文件，包含所有常量和映射
├── cn2an.go           # 中文数字转阿拉伯数字
├── an2cn.go           # 阿拉伯数字转中文数字
├── transform.go       # 句子转换
├── cn2an_test.go      # cn2an 测试
├── an2cn_test.go      # an2cn 测试
├── transform_test.go  # transform 测试
├── example/
│   └── main.go        # 使用示例
├── go.mod             # Go 模块文件
└── README.md          # 说明文档
```

## 特性支持

### 支持的中文数字

- 小写：零、一、二、三、四、五、六、七、八、九、十、百、千、万、亿
- 大写：零、壹、贰、叁、肆、伍、陆、柒、捌、玖、拾、佰、仟、万、亿
- 特殊：〇、幺、两、廿（二十）

### 支持的数字范围

- 整数：0 到 9999999999999999（千万亿）
- 小数：支持最多 16 位小数精度

### 特殊功能

1. **口语化支持**（normal 模式）
   - "一万二" => 12000
   - "三千五" => 3500
   - "两百三" => 230

2. **混合模式支持**（smart 模式）
   - "1百23" => 123
   - "10.1万" => 101000
   - "35.1亿" => 3510000000

3. **人民币格式**（rmb 模式）
   - 123.45 => "壹佰贰拾叁元肆角伍分"
   - 0.5 => "伍角"
   - 0.05 => "伍分"

## 许可证

MIT License

## 致谢

本项目是 Python 版本 [cn2an](https://github.com/Ailln/cn2an) 的 Golang 完整复刻。感谢原作者 [Ailln](https://github.com/Ailln) 的优秀设计和实现。

## 参考

- [cn2an - Python 版本](https://github.com/Ailln/cn2an)
- [中文数字 - Wikipedia](https://zh.wikipedia.org/zh-sg/%E4%B8%AD%E6%96%87%E6%95%B0%E5%AD%97)
