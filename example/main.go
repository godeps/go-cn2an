package main

import (
	"fmt"
	gocn2an "github.com/godeps/go-cn2an"
)

func main() {
	fmt.Println("=== go-cn2an 示例 ===")
	fmt.Println()

	// 1. 中文数字转阿拉伯数字
	fmt.Println("1. 中文数字 => 阿拉伯数字")
	c := gocn2an.NewCn2An()

	// strict 模式
	result1, _ := c.Cn2an("一百二十三", "strict")
	fmt.Printf("  一百二十三 (strict) => %.0f\n", result1)

	result2, _ := c.Cn2an("一千零一", "strict")
	fmt.Printf("  一千零一 (strict) => %.0f\n", result2)

	result3, _ := c.Cn2an("负一百二十三点四五", "strict")
	fmt.Printf("  负一百二十三点四五 (strict) => %f\n", result3)

	// normal 模式
	result4, _ := c.Cn2an("一二三", "normal")
	fmt.Printf("  一二三 (normal) => %.0f\n", result4)

	result5, _ := c.Cn2an("一万二", "normal")
	fmt.Printf("  一万二 (normal) => %.0f\n", result5)

	// smart 模式
	result6, _ := c.Cn2an("1百23", "smart")
	fmt.Printf("  1百23 (smart) => %.0f\n", result6)

	result7, _ := c.Cn2an("10.1万", "smart")
	fmt.Printf("  10.1万 (smart) => %.0f\n", result7)

	fmt.Println()

	// 2. 阿拉伯数字转中文数字
	fmt.Println("2. 阿拉伯数字 => 中文数字")
	a := gocn2an.NewAn2Cn()

	// low 模式
	result8, _ := a.An2cn(123, "low")
	fmt.Printf("  123 (low) => %s\n", result8)

	result9, _ := a.An2cn(1001, "low")
	fmt.Printf("  1001 (low) => %s\n", result9)

	result10, _ := a.An2cn(-123.45, "low")
	fmt.Printf("  -123.45 (low) => %s\n", result10)

	// up 模式
	result11, _ := a.An2cn(123, "up")
	fmt.Printf("  123 (up) => %s\n", result11)

	// rmb 模式
	result12, _ := a.An2cn(123, "rmb")
	fmt.Printf("  123 (rmb) => %s\n", result12)

	result13, _ := a.An2cn(123.45, "rmb")
	fmt.Printf("  123.45 (rmb) => %s\n", result13)

	result14, _ := a.An2cn(0.5, "rmb")
	fmt.Printf("  0.5 (rmb) => %s\n", result14)

	// direct 模式
	result15, _ := a.An2cn(123, "direct")
	fmt.Printf("  123 (direct) => %s\n", result15)

	fmt.Println()

	// 3. 句子转换
	fmt.Println("3. 句子转换")
	t := gocn2an.NewTransform()

	// cn2an
	result16, _ := t.Transform("小王捡了一百块钱", "cn2an")
	fmt.Printf("  小王捡了一百块钱 (cn2an) => %s\n", result16)

	result17, _ := t.Transform("小王的生日是二零零一年三月四日", "cn2an")
	fmt.Printf("  小王的生日是二零零一年三月四日 (cn2an) => %s\n", result17)

	result18, _ := t.Transform("抛出去的硬币为正面的概率是二分之一", "cn2an")
	fmt.Printf("  抛出去的硬币为正面的概率是二分之一 (cn2an) => %s\n", result18)

	// an2cn
	result19, _ := t.Transform("小王捡了100块钱", "an2cn")
	fmt.Printf("  小王捡了100块钱 (an2cn) => %s\n", result19)

	result20, _ := t.Transform("小王的生日是2001年3月4日", "an2cn")
	fmt.Printf("  小王的生日是2001年3月4日 (an2cn) => %s\n", result20)

	result21, _ := t.Transform("抛出去的硬币为正面的概率是1/2", "an2cn")
	fmt.Printf("  抛出去的硬币为正面的概率是1/2 (an2cn) => %s\n", result21)

	fmt.Println()
	fmt.Println("=== 示例完成 ===")
}
