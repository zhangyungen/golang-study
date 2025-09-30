package str

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"log"
	"math/rand"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// ==================== 基础判断工具 ====================

// IsEmpty 判断字符串是否为空（包括空格）
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty 判断字符串是否非空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// IsBlank 判断字符串是否全是空白字符
func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// IsNotBlank 判断字符串是否包含非空白字符
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// ==================== 转换工具 ====================

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ToCamelCase 下划线转驼峰
func ToCamelCase(s string) string {
	if IsBlank(s) {
		return s
	}

	words := strings.Split(s, "_")
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	return strings.Join(words, "")
}

// ToSnakeCase 驼峰转下划线
func ToSnakeCase(s string) string {
	if IsBlank(s) {
		return s
	}

	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// Truncate 截断字符串
func Truncate(s string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}

	return string(runes[:maxLength]) + "..."
}

// Abbreviate 缩写字符串
func Abbreviate(s string, maxWidth int) string {
	if maxWidth < 4 {
		return s // 宽度太小无法缩写
	}

	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}

	return string(runes[:maxWidth-3]) + "..."
}

// ==================== 提取工具 ====================

// SubstringBetween 提取两个标记之间的字符串
func SubstringBetween(str, start, end string) string {
	startIndex := strings.Index(str, start)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(start)
	endIndex := strings.Index(str[startIndex:], end)
	if endIndex == -1 {
		return ""
	}

	return str[startIndex : startIndex+endIndex]
}

// ExtractNumbers 提取字符串中的所有数字
func ExtractNumbers(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ExtractLetters 提取字符串中的所有字母
func ExtractLetters(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ==================== 格式化工具 ====================

// Center 居中显示字符串
func Center(s string, width int, padChar rune) string {
	if width <= len(s) {
		return s
	}

	padding := width - len(s)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return strings.Repeat(string(padChar), leftPadding) + s +
		strings.Repeat(string(padChar), rightPadding)
}

// PadLeft 左填充
func PadLeft(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(string(padChar), length-len(s)) + s
}

// PadRight 右填充
func PadRight(s string, length int, padChar rune) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(string(padChar), length-len(s))
}

// FormatPhone 格式化手机号
func FormatPhone(phone string) string {
	// 只保留数字
	numbers := ExtractNumbers(phone)

	if len(numbers) == 11 {
		return numbers[:3] + "-" + numbers[3:7] + "-" + numbers[7:]
	}
	return numbers
}

// ==================== 编码/解码工具 ====================

// SafeHTML 转义HTML特殊字符
func SafeHTML(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		switch r {
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '&':
			buf.WriteString("&amp;")
		case '"':
			buf.WriteString("&quot;")
		case '\'':
			buf.WriteString("&#39;")
		default:
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// UnescapeHTML 反转义HTML
func UnescapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	return s
}

// ==================== 验证工具 ====================

// IsEmail 验证邮箱格式
func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// IsPhone 验证手机号格式（中国）
func IsPhone(s string) bool {
	// 只保留数字验证
	numbers := ExtractNumbers(s)
	if len(numbers) != 11 {
		return false
	}

	// 简单的手机号格式验证
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, numbers)
	return matched
}

// IsNumeric 判断是否全是数字
func IsNumeric(s string) bool {
	if IsBlank(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlpha 判断是否全是字母
func IsAlpha(s string) bool {
	if IsBlank(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphaNumeric 判断是否只包含字母和数字
func IsAlphaNumeric(s string) bool {
	if IsBlank(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// ==================== 生成工具 ====================

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// RandomNumber 生成随机数字字符串
func RandomNumber(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// UUID 生成简化的UUID
func UUID() string {
	timestamp := time.Now().UnixNano()
	randomPart := RandomString(8)
	return fmt.Sprintf("%x-%s", timestamp, randomPart)
}

// ==================== 高级功能 ====================

// CountMatches 统计子串出现次数
func CountMatches(str, sub string) int {
	if IsEmpty(str) || IsEmpty(sub) {
		return 0
	}
	return (len(str) - len(strings.ReplaceAll(str, sub, ""))) / len(sub)
}

// 字符串切片中是否startWith某字符串
func StringsStartWith(slice []string, target string) bool {
	targetLen := len(target)
	if targetLen == 0 {
		return false
	}
	for i := 0; i < len(slice); i++ {
		s := slice[i]
		if len(s) >= targetLen && s[:targetLen] == target {
			return true
		}
	}
	return false
}

func StartWith(str string, target string) bool {
	targetLen := len(target)
	if targetLen == 0 {
		return false
	}
	if len(str) >= targetLen && str[:targetLen] == target {
		return true
	}
	return false
}

// Remove 移除所有指定字符
func Remove(str string, remove string) string {
	if IsEmpty(str) || IsEmpty(remove) {
		return str
	}
	return strings.ReplaceAll(str, remove, "")
}

// RemoveWhitespace 移除所有空白字符
func RemoveWhitespace(str string) string {
	var result strings.Builder
	for _, r := range str {
		if !unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ContainsAny 判断是否包含任意指定字符串
func ContainsAny(str string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}

// ContainsAll 判断是否包含所有指定字符串
func ContainsAll(str string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(str, sub) {
			return false
		}
	}
	return true
}

// ==================== 性能优化工具 ====================

// StringToBytes 高性能字符串转字节切片（无内存分配）
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString 高性能字节切片转字符串（无内存分配）
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Concat 高性能字符串拼接
func Concat(strs ...string) string {
	var totalLen int
	for _, s := range strs {
		totalLen += len(s)
	}

	var builder strings.Builder
	builder.Grow(totalLen)
	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()
}

// ==================== 编码相关 ====================

// IsUTF8 判断是否为有效的UTF-8编码
func IsUTF8(s string) bool {
	return utf8.ValidString(s)
}

// RuneLength 获取字符数量（非字节长度）
func RuneLength(s string) int {
	return utf8.RuneCountInString(s)
}

// SubstringByRune 按字符数截取子串
func SubstringByRune(s string, start, length int) string {
	runes := []rune(s)
	if start < 0 || start >= len(runes) || length <= 0 {
		return ""
	}

	end := start + length
	if end > len(runes) {
		end = len(runes)
	}

	return string(runes[start:end])
}

// ObjToJson 对象转JSON字符串
func ObjToJson[T any](obj T) string {
	meta, err := json.Marshal(obj)
	if err != nil {
		log.Printf("ObjToJson", err)
		return ""
	}
	return string(meta)
}

// JsonToObj 字符串转对象
func JsonToObj[T any](str string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(str), &result)
	return result, err
}
