package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("=== Go 语言数据类型转换示例 ===\n")

	// 1. 整数与浮点数之间的转换
	fmt.Println("1. 整数与浮点数转换:")
	var i int = 42
	var f float64 = float64(i)
	fmt.Printf("int(%d) → float64(%f)\n", i, f)

	var f2 float64 = 3.14159
	var i2 int = int(f2)
	fmt.Printf("float64(%f) → int(%d)\n", f2, i2)
	fmt.Println()

	// 2. 不同整数类型之间的转换
	fmt.Println("2. 不同整数类型转换:")
	var i32 int32 = 12345
	var i64 int64 = int64(i32)
	fmt.Printf("int32(%d) → int64(%d)\n", i32, i64)

	var u8 uint8 = 255
	var i16 int16 = int16(u8)
	fmt.Printf("uint8(%d) → int16(%d)\n", u8, i16)
	fmt.Println()

	// 3. 各种类型转换为字符串
	fmt.Println("3. 各种类型转换为字符串:")

	// 整数转字符串
	num := 42
	str1 := strconv.Itoa(num)
	fmt.Printf("int(%d) → string(%q)\n", num, str1)

	// 使用 fmt.Sprintf 转换
	str2 := fmt.Sprintf("%d", num)
	fmt.Printf("fmt.Sprintf: int(%d) → string(%q)\n", num, str2)

	// 浮点数转字符串
	pi := 3.14159
	str3 := strconv.FormatFloat(pi, 'f', 2, 64)
	fmt.Printf("float64(%f) → string(%q)\n", pi, str3)

	// 布尔值转字符串
	b := true
	str4 := strconv.FormatBool(b)
	fmt.Printf("bool(%t) → string(%q)\n", b, str4)

	// 字节切片转字符串
	bytes := []byte{72, 101, 108, 108, 111}
	str5 := string(bytes)
	fmt.Printf("[]byte(%v) → string(%q)\n", bytes, str5)

	// rune 转字符串
	r := '世'
	str6 := string(r)
	fmt.Printf("rune(%U) → string(%q)\n", r, str6)
	fmt.Println()

	// 4. 字符串转换为其他类型
	fmt.Println("4. 字符串转换为其他类型:")

	// 字符串转整数
	s := "123"
	num2, _ := strconv.Atoi(s)
	fmt.Printf("string(%q) → int(%d)\n", s, num2)

	// 字符串转浮点数
	s2 := "3.14"
	f3, _ := strconv.ParseFloat(s2, 64)
	fmt.Printf("string(%q) → float64(%f)\n", s2, f3)

	// 字符串转布尔值
	s3 := "true"
	b2, _ := strconv.ParseBool(s3)
	fmt.Printf("string(%q) → bool(%t)\n", s3, b2)

	// 字符串转字节切片
	s4 := "Hello"
	bytes2 := []byte(s4)
	fmt.Printf("string(%q) → []byte(%v)\n", s4, bytes2)
	fmt.Println()

	// 5. 接口类型转换（类型断言）
	fmt.Println("5. 接口类型转换:")
	var val interface{} = "Hello, Go!"
	if s, ok := val.(string); ok {
		fmt.Printf("interface{}(%v) → string(%q)\n", val, s)
	}

	var val2 interface{} = 42
	if num, ok := val2.(int); ok {
		fmt.Printf("interface{}(%v) → int(%d)\n", val2, num)
	}
	fmt.Println()

	// 6. 自定义类型转换
	fmt.Println("6. 自定义类型转换:")
	type Celsius float64
	type Fahrenheit float64

	c := Celsius(100)
	fahrenheit := Fahrenheit(c*9/5 + 32)
	fmt.Printf("Celsius(%g) → Fahrenheit(%g)\n", c, fahrenheit)

	// 7. 使用 strings 包进行转换和操作
	fmt.Println("\n7. 字符串大小写转换:")
	original := "Hello, World!"
	upper := strings.ToUpper(original)
	lower := strings.ToLower(original)
	fmt.Printf("原始: %q\n大写: %q\n小写: %q\n", original, upper, lower)
}
