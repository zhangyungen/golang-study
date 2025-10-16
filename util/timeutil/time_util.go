package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//// timeUtil 时间日期工具类
//type TimeUtil struct{}
//
//// NewtimeUtil 创建时间工具实例
//func GetTimeUtil() *TimeUtil {
//	once.Do(func() {
//		timeUtil = &TimeUtil{}
//	})
//	return timeUtil
//}

// 预定义时间格式
const (
	FormatDateTime   = "2006-01-02 15:04:05"
	FormatDate       = "2006-01-02"
	FormatTime       = "15:04:05"
	FormatDateTimeMs = "2006-01-02 15:04:05.000"
	FormatDateTimeCN = "2006年01月02日 15时04分05秒"
	FormatDateCN     = "2006年01月02日"
	FormatTimeCN     = "15时04分05秒"
	FormatISO8601    = "2006-01-02T15:04:05Z07:00"
	FormatRFC3339    = time.RFC3339
	FormatRFC1123    = time.RFC1123
)

// ==================== 时间获取相关方法 ====================

// Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

// NowString 获取当前时间字符串
func NowString(format string) string {
	return time.Now().Format(format)
}

// Today 获取今天日期
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// TodayString 获取今天日期字符串
func TodayString() string {
	return time.Now().Format(FormatDate)
}

// Timestamp 获取当前时间戳（秒）
func TimestampUnix() int64 {
	return time.Now().Unix()
}

// TimestampMilli 获取当前时间戳（毫秒）
func TimestampMilli() int64 {
	return time.Now().UnixMilli()
}

// ==================== 时间格式化方法 ====================

// Format 时间格式化
func Format(t time.Time, format string) string {
	return t.Format(format)
}

// FormatDefault 默认格式化 (2006-01-02 15:04:05)
func FormatDefault(t time.Time) string {
	return t.Format(FormatDateTime)
}

// FormatDateOnly 只格式化日期
func FormatDateOnly(t time.Time) string {
	return t.Format(FormatDate)
}

// FormatTimeOnly 只格式化时间
func FormatTimeOnly(t time.Time) string {
	return t.Format(FormatTime)
}

// FormatChinese 中文格式格式化
func FormatChinese(t time.Time) string {
	return t.Format(FormatDateTimeCN)
}

// ==================== 时间解析方法 ====================

// Parse 解析时间字符串
func Parse(timeStr, format string) (time.Time, error) {
	return time.Parse(format, timeStr)
}

// ParseInLocation 在指定时区解析时间字符串
func ParseInLocation(timeStr, format string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(format, timeStr, loc)
}

// ParseDefault 解析默认格式时间字符串
func ParseDefault(timeStr string) (time.Time, error) {
	return time.Parse(FormatDateTime, timeStr)
}

// ParseDate 解析日期字符串
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(FormatDate, dateStr)
}

// ParseMultiple 尝试多种格式解析时间字符串
func ParseMultiple(timeStr string) (time.Time, error) {
	formats := []string{
		FormatDateTime,
		FormatDate,
		FormatISO8601,
		FormatRFC3339,
		FormatRFC1123,
		"2006/01/02 15:04:05",
		"2006/01/02",
		"2006-01-02",
		"20060102150405",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间字符串: %s", timeStr)
}

// ==================== 时间计算相关方法 ====================

// AddDays 添加天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours 添加小时
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 添加分钟
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 添加秒数
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// AddMonths 添加月数
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears 添加年数
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// ==================== 日期操作相关方法 ====================

// BeginningOfDay 获取某天的开始时间 (00:00:00)
func BeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取某天的结束时间 (23:59:59)
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// BeginningOfMonth 获取月初时间
func BeginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 获取月末时间
func EndOfMonth(t time.Time) time.Time {
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return firstDay.AddDate(0, 1, -1)
}

// BeginningOfYear 获取年初时间
func BeginningOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 获取年末时间
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

// ==================== 时间比较相关方法 ====================

// IsBefore 判断时间是否在另一个时间之前
func IsBefore(t1, t2 time.Time) bool {
	return t1.Before(t2)
}

// IsAfter 判断时间是否在另一个时间之后
func IsAfter(t1, t2 time.Time) bool {
	return t1.After(t2)
}

// IsEqual 判断两个时间是否相等
func IsEqual(t1, t2 time.Time) bool {
	return t1.Equal(t2)
}

// IsBetween 判断时间是否在两个时间之间
func IsBetween(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}

// DiffDays 计算两个时间相差的天数
func DiffDays(t1, t2 time.Time) int {
	t1 = BeginningOfDay(t1)
	t2 = BeginningOfDay(t2)
	diff := t2.Sub(t1)
	return int(diff.Hours() / 24)
}

// DiffHours 计算两个时间相差的小时数
func DiffHours(t1, t2 time.Time) int {
	diff := t2.Sub(t1)
	return int(diff.Hours())
}

// DiffMinutes 计算两个时间相差的分钟数
func DiffMinutes(t1, t2 time.Time) int {
	diff := t2.Sub(t1)
	return int(diff.Minutes())
}

// ==================== 时间判断相关方法 ====================

// IsToday 判断是否是今天
func IsToday(t time.Time) bool {
	today := Today()
	return t.Year() == today.Year() && t.Month() == today.Month() && t.Day() == today.Day()
}

// IsWeekend 判断是否是周末
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday 判断是否是工作日
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// IsLeapYear 判断是否是闰年
func IsLeapYear(t time.Time) bool {
	year := t.Year()
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// ==================== 时区相关方法 ====================

// ToUTC 转换为UTC时间
func ToUTC(t time.Time) time.Time {
	return t.UTC()
}

// ToLocal 转换为本地时间
func ToLocal(t time.Time) time.Time {
	return t.Local()
}

// ToTimezone 转换为指定时区时间
func ToTimezone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// GetTimezoneOffset 获取时区偏移量（小时）
func GetTimezoneOffset(t time.Time) int {
	_, offset := t.Zone()
	return offset / 3600
}

// ==================== 时间戳相关方法 ====================

// FromTimestamp 从时间戳创建时间
func FromTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// FromTimestampMilli 从毫秒时间戳创建时间
func FromTimestampMilli(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

// ToTimestamp 转换为时间戳（秒）
func ToTimestamp(t time.Time) int64 {
	return t.Unix()
}

// ToTimestampMilli 转换为时间戳（毫秒）
func ToTimestampMilli(t time.Time) int64 {
	return t.UnixMilli()
}

// ==================== 其他实用方法 ====================

// Age 计算年龄
func Age(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()

	// 如果生日还没过，年龄减1
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}

// WeekNumber 获取是一年中的第几周
func WeekNumber(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}

// DaysInMonth 获取月份的天数
func DaysInMonth(t time.Time) int {
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth.Day()
}

// FormatDuration 格式化时间间隔
func FormatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分%d秒", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分%d秒", minutes, seconds)
	} else {
		return fmt.Sprintf("%d秒", seconds)
	}
}

// Humanize 人性化时间显示
func Humanize(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "刚刚"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d分钟前", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d小时前", int(diff.Hours()))
	} else if diff < 30*24*time.Hour {
		return fmt.Sprintf("%d天前", int(diff.Hours()/24))
	} else if diff < 365*24*time.Hour {
		return fmt.Sprintf("%d个月前", int(diff.Hours()/(24*30)))
	} else {
		return fmt.Sprintf("%d年前", int(diff.Hours()/(24*365)))
	}
}

// ParseDuration 解析时间间隔字符串
func ParseDuration(durationStr string) (time.Duration, error) {
	// 支持格式: 1h30m, 2d5h, 1天2小时等
	durationStr = strings.ReplaceAll(durationStr, "天", "d")
	durationStr = strings.ReplaceAll(durationStr, "小时", "h")
	durationStr = strings.ReplaceAll(durationStr, "分钟", "m")
	durationStr = strings.ReplaceAll(durationStr, "秒", "s")

	var totalDuration time.Duration
	var currentNum string

	for _, char := range durationStr {
		if char >= '0' && char <= '9' {
			currentNum += string(char)
		} else {
			if currentNum == "" {
				continue
			}

			num, _ := strconv.Atoi(currentNum)
			var unit time.Duration

			switch char {
			case 'd':
				unit = 24 * time.Hour
			case 'h':
				unit = time.Hour
			case 'm':
				unit = time.Minute
			case 's':
				unit = time.Second
			default:
				continue
			}

			totalDuration += time.Duration(num) * unit
			currentNum = ""
		}
	}

	return totalDuration, nil
}

// generateTimestamp13 生成13位时间戳(毫秒级)
func generateTimestamp13() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
}

func timestamp13() int64 {
	return time.Now().UnixNano() / 1e6
}
