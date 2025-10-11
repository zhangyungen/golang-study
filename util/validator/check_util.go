package validator

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// Validator 校验器结构体
type Validator struct {
	errors []string
}

// New 创建新的校验器实例
func New() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

// Errors 获取所有错误信息
func (v *Validator) Errors() []string {
	return v.errors
}

// HasErrors 检查是否有错误
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Error 返回错误信息字符串
func (v *Validator) Error() string {
	if !v.HasErrors() {
		return ""
	}
	return strings.Join(v.errors, "; ")
}

// Clear 清空错误信息
func (v *Validator) Clear() {
	v.errors = make([]string, 0)
}

// addError 添加错误信息
func (v *Validator) addError(field string, message string) {
	v.errors = append(v.errors, fmt.Sprintf("%s: %s", field, message))
}

// ==================== 基础类型校验 ====================

// Required 必填字段校验
func (v *Validator) Required(field string, value interface{}) *Validator {
	if isEmpty(value) {
		v.addError(field, "该字段为必填项")
	}
	return v
}

// NotNil 非空指针校验
func (v *Validator) NotNil(field string, value interface{}) *Validator {
	if value == nil {
		v.addError(field, "该字段不能为nil")
	}
	return v
}

// ==================== 字符串校验 ====================

// MinLength 最小长度校验
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if utf8.RuneCountInString(value) < min {
		v.addError(field, fmt.Sprintf("长度不能少于 %d 个字符", min))
	}
	return v
}

// MaxLength 最大长度校验
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if utf8.RuneCountInString(value) > max {
		v.addError(field, fmt.Sprintf("长度不能超过 %d 个字符", max))
	}
	return v
}

// LengthRange 长度范围校验
func (v *Validator) LengthRange(field, value string, min, max int) *Validator {
	length := utf8.RuneCountInString(value)
	if length < min || length > max {
		v.addError(field, fmt.Sprintf("长度必须在 %d 到 %d 个字符之间", min, max))
	}
	return v
}

// Matches 正则表达式匹配
func (v *Validator) Matches(field, value, pattern string) *Validator {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil || !matched {
		v.addError(field, "格式不符合要求")
	}
	return v
}

// IsEmail 邮箱格式校验
func (v *Validator) IsEmail(field, email string) *Validator {
	_, err := mail.ParseAddress(email)
	if err != nil {
		v.addError(field, "邮箱格式不正确")
	}
	return v
}

// IsURL URL格式校验
func (v *Validator) IsURL(field, urlStr string) *Validator {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		v.addError(field, "URL格式不正确")
		return v
	}

	u, err := url.Parse(urlStr)
	if err != nil || u.Scheme == "" || u.Host == "" {
		v.addError(field, "URL格式不正确")
	}
	return v
}

// IsPhoneCN 中国手机号校验
func (v *Validator) IsPhoneCN(field, phone string) *Validator {
	// 简单的中国手机号校验
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	if !matched {
		v.addError(field, "手机号格式不正确")
	}
	return v
}

// Contains 包含特定文本
func (v *Validator) Contains(field, value, substr string) *Validator {
	if !strings.Contains(value, substr) {
		v.addError(field, fmt.Sprintf("必须包含 '%s'", substr))
	}
	return v
}

// NotContains 不包含特定文本
func (v *Validator) NotContains(field, value, substr string) *Validator {
	if strings.Contains(value, substr) {
		v.addError(field, fmt.Sprintf("不能包含 '%s'", substr))
	}
	return v
}

// IsAlpha 只包含字母
func (v *Validator) IsAlpha(field, value string) *Validator {
	for _, r := range value {
		if !unicode.IsLetter(r) {
			v.addError(field, "只能包含字母")
			break
		}
	}
	return v
}

// IsAlphanumeric 只包含字母和数字
func (v *Validator) IsAlphanumeric(field, value string) *Validator {
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			v.addError(field, "只能包含字母和数字")
			break
		}
	}
	return v
}

// ==================== 数字校验 ====================

// MinInt 最小整数值
func (v *Validator) MinInt(field string, value, min int) *Validator {
	if value < min {
		v.addError(field, fmt.Sprintf("不能小于 %d", min))
	}
	return v
}

// MaxInt 最大整数值
func (v *Validator) MaxInt(field string, value, max int) *Validator {
	if value > max {
		v.addError(field, fmt.Sprintf("不能大于 %d", max))
	}
	return v
}

// RangeInt 整数范围
func (v *Validator) RangeInt(field string, value, min, max int) *Validator {
	if value < min || value > max {
		v.addError(field, fmt.Sprintf("必须在 %d 到 %d 之间", min, max))
	}
	return v
}

// MinFloat 最小浮点数值
func (v *Validator) MinFloat(field string, value, min float64) *Validator {
	if value < min {
		v.addError(field, fmt.Sprintf("不能小于 %.2f", min))
	}
	return v
}

// MaxFloat 最大浮点数值
func (v *Validator) MaxFloat(field string, value, max float64) *Validator {
	if value > max {
		v.addError(field, fmt.Sprintf("不能大于 %.2f", max))
	}
	return v
}

// RangeFloat 浮点数范围
func (v *Validator) RangeFloat(field string, value, min, max float64) *Validator {
	if value < min || value > max {
		v.addError(field, fmt.Sprintf("必须在 %.2f 到 %.2f 之间", min, max))
	}
	return v
}

// IsPositive 正数校验
func (v *Validator) IsPositive(field string, value float64) *Validator {
	if value <= 0 {
		v.addError(field, "必须为正数")
	}
	return v
}

// ==================== 切片和数组校验 ====================

// MinSliceLength 切片最小长度
func (v *Validator) MinSliceLength(field string, slice interface{}, min int) *Validator {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		v.addError(field, "不是切片或数组类型")
		return v
	}

	if val.Len() < min {
		v.addError(field, fmt.Sprintf("元素数量不能少于 %d 个", min))
	}
	return v
}

// MaxSliceLength 切片最大长度
func (v *Validator) MaxSliceLength(field string, slice interface{}, max int) *Validator {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		v.addError(field, "不是切片或数组类型")
		return v
	}

	if val.Len() > max {
		v.addError(field, fmt.Sprintf("元素数量不能超过 %d 个", max))
	}
	return v
}

// SliceLengthRange 切片长度范围
func (v *Validator) SliceLengthRange(field string, slice interface{}, min, max int) *Validator {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		v.addError(field, "不是切片或数组类型")
		return v
	}

	length := val.Len()
	if length < min || length > max {
		v.addError(field, fmt.Sprintf("元素数量必须在 %d 到 %d 个之间", min, max))
	}
	return v
}

// ==================== 时间校验 ====================

// IsAfter 时间在指定时间之后
func (v *Validator) IsAfter(field string, t, after time.Time) *Validator {
	if !t.After(after) {
		v.addError(field, fmt.Sprintf("必须在 %s 之后", after.Format("2006-01-02 15:04:05")))
	}
	return v
}

// IsBefore 时间在指定时间之前
func (v *Validator) IsBefore(field string, t, before time.Time) *Validator {
	if !t.Before(before) {
		v.addError(field, fmt.Sprintf("必须在 %s 之前", before.Format("2006-01-02 15:04:05")))
	}
	return v
}

// IsBetweenTimes 时间在范围内
func (v *Validator) IsBetweenTimes(field string, t, start, end time.Time) *Validator {
	if t.Before(start) || t.After(end) {
		v.addError(field, fmt.Sprintf("必须在 %s 到 %s 之间",
			start.Format("2006-01-02 15:04:05"),
			end.Format("2006-01-02 15:04:05")))
	}
	return v
}

// ==================== 自定义校验 ====================

// Custom 自定义校验函数
func (v *Validator) Custom(field string, value interface{},
	validateFunc func(interface{}) (bool, string)) *Validator {

	if valid, message := validateFunc(value); !valid {
		v.addError(field, message)
	}
	return v
}

// ==================== 条件校验 ====================

// When 条件校验
func (v *Validator) When(condition bool, validateFunc func()) *Validator {
	if condition {
		validateFunc()
	}
	return v
}

// ==================== 辅助函数 ====================

// isEmpty 检查值是否为空
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Struct:
		// 对于time.Time特殊处理
		if t, ok := value.(time.Time); ok {
			return t.IsZero()
		}
		return false
	default:
		return false
	}
}

func IsEmpty(value interface{}, msg string) error {
	if isEmpty(value) {
		return errors.New(msg)
	}
	return nil
}

// ==================== 快捷函数 ====================

// ValidateStruct 校验结构体（简单版本）
func ValidateStruct(data interface{}) error {
	validator := New()

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("只能校验结构体")
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		// 获取字段标签
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		fieldName := field.Name

		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			switch {
			case rule == "required":
				validator.Required(fieldName, fieldValue.Interface())
			case strings.HasPrefix(rule, "min="):
				if min, err := strconv.Atoi(strings.TrimPrefix(rule, "min=")); err == nil {
					if fieldValue.Kind() == reflect.String {
						validator.MinLength(fieldName, fieldValue.String(), min)
					}
				}
			case strings.HasPrefix(rule, "max="):
				if max, err := strconv.Atoi(strings.TrimPrefix(rule, "max=")); err == nil {
					if fieldValue.Kind() == reflect.String {
						validator.MaxLength(fieldName, fieldValue.String(), max)
					}
				}
			case rule == "email":
				if fieldValue.Kind() == reflect.String {
					validator.IsEmail(fieldName, fieldValue.String())
				}
			}
		}
	}

	if validator.HasErrors() {
		return errors.New(validator.Error())
	}

	return nil
}
