package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Go 语言时间操作完整示例 ===\n")

	// 1. 获取当前时间
	fmt.Println("1. 获取当前时间:")
	now := time.Now()
	fmt.Printf("当前时间: %v\n", now)
	fmt.Printf("时间戳(秒): %d\n", now.Unix())
	fmt.Printf("时间戳(毫秒): %d\n", now.UnixMilli())
	fmt.Printf("时间戳(微秒): %d\n", now.UnixMicro())
	fmt.Printf("时间戳(纳秒): %d\n", now.UnixNano())
	fmt.Println()

	// 2. 创建特定时间
	fmt.Println("2. 创建特定时间:")
	// 使用日期创建
	specificTime := time.Date(2023, time.December, 25, 15, 30, 45, 500000000, time.UTC)
	fmt.Printf("特定时间: %v\n", specificTime)
	fmt.Printf("本地时间: %v\n", specificTime.Local())

	// 从时间戳创建
	timestamp := time.Unix(1703044245, 0)
	fmt.Printf("从时间戳创建: %v\n", timestamp)

	// 从字符串创建（需要先解析）
	timeStr := "2023-12-25 15:30:45"
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	fmt.Printf("从字符串解析: %v\n", parsedTime)
	fmt.Println()

	// 3. 时间格式化 (Go 特有布局: 2006-01-02 15:04:05)
	fmt.Println("3. 时间格式化:")
	fmt.Printf("RFC3339格式: %v\n", now.Format(time.RFC3339))
	fmt.Printf("RFC1123格式: %v\n", now.Format(time.RFC1123))
	fmt.Printf("自定义格式: %v\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("日期格式: %v\n", now.Format("2006/01/02"))
	fmt.Printf("时间格式: %v\n", now.Format("15:04:05"))
	fmt.Printf("中文格式: %v\n", now.Format("2006年01月02日 15时04分05秒"))
	fmt.Printf("带毫秒: %v\n", now.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("带时区: %v\n", now.Format("2006-01-02 15:04:05 -0700 MST"))
	fmt.Println()

	// 4. 时间解析 (字符串转时间)
	fmt.Println("4. 时间解析:")

	// 解析标准格式
	timeStr1 := "2023-12-25 15:30:45"
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr1)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
	} else {
		fmt.Printf("解析成功: %v\n", t1)
	}

	// 解析RFC3339格式
	timeStr2 := "2023-12-25T15:30:45+08:00"
	t2, err := time.Parse(time.RFC3339, timeStr2)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
	} else {
		fmt.Printf("RFC3339解析: %v\n", t2)
	}

	// 解析自定义格式
	timeStr3 := "25/12/2023 15:30"
	t3, err := time.Parse("02/01/2006 15:04", timeStr3)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
	} else {
		fmt.Printf("自定义格式解析: %v\n", t3)
	}

	// 在特定时区解析
	timeStr4 := "2023-12-25 15:30:45"
	loc, _ := time.LoadLocation("America/New_York")
	t4, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr4, loc)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
	} else {
		fmt.Printf("纽约时区解析: %v\n", t4)
	}
	fmt.Println()

	// 5. 时间加减操作
	fmt.Println("5. 时间加减操作:")
	fmt.Printf("当前时间: %v\n", now)
	fmt.Printf("加1小时: %v\n", now.Add(time.Hour))
	fmt.Printf("加30分钟: %v\n", now.Add(30*time.Minute))
	fmt.Printf("加1天: %v\n", now.Add(24*time.Hour))
	fmt.Printf("加1周: %v\n", now.Add(7*24*time.Hour))
	fmt.Printf("减1小时: %v\n", now.Add(-time.Hour))
	fmt.Printf("减2天: %v\n", now.Add(-48*time.Hour))
	fmt.Println()

	// 6. 日期操作
	fmt.Println("6. 日期操作:")

	// 获取特定日期部分
	fmt.Printf("年: %d\n", now.Year())
	fmt.Printf("月: %d (%s)\n", now.Month(), now.Month())
	fmt.Printf("日: %d\n", now.Day())
	fmt.Printf("时: %d\n", now.Hour())
	fmt.Printf("分: %d\n", now.Minute())
	fmt.Printf("秒: %d\n", now.Second())
	fmt.Printf("纳秒: %d\n", now.Nanosecond())
	fmt.Printf("星期: %d (%s)\n", now.Weekday(), now.Weekday())
	fmt.Printf("一年中的第几天: %d\n", now.YearDay())
	fmt.Printf("一个月中的第几天: %d\n", now.Day())

	// 日期计算
	tomorrow := now.AddDate(0, 0, 1)
	nextMonth := now.AddDate(0, 1, 0)
	nextYear := now.AddDate(1, 0, 0)
	fmt.Printf("明天: %v\n", tomorrow.Format("2006-01-02"))
	fmt.Printf("下个月: %v\n", nextMonth.Format("2006-01-02"))
	fmt.Printf("明年: %v\n", nextYear.Format("2006-01-02"))

	// 获取月初和月末
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	fmt.Printf("本月第一天: %v\n", firstDayOfMonth.Format("2006-01-02"))
	fmt.Printf("本月最后一天: %v\n", lastDayOfMonth.Format("2006-01-02"))
	fmt.Println()

	// 7. 时间比较
	fmt.Println("7. 时间比较:")
	time1 := time.Now()
	time2 := time1.Add(5 * time.Minute)
	time3 := time1.Add(-5 * time.Minute)

	fmt.Printf("time1: %v\n", time1.Format("15:04:05"))
	fmt.Printf("time2: %v\n", time2.Format("15:04:05"))
	fmt.Printf("time3: %v\n", time3.Format("15:04:05"))
	fmt.Printf("time1 在 time2 之前: %t\n", time1.Before(time2))
	fmt.Printf("time2 在 time1 之后: %t\n", time2.After(time1))
	fmt.Printf("time1 等于 time2: %t\n", time1.Equal(time2))
	fmt.Printf("time1 在 time3 之后: %t\n", time1.After(time3))
	fmt.Println()

	// 8. 时间差计算
	fmt.Println("8. 时间差计算:")
	start := time.Now()
	time.Sleep(2 * time.Second) // 模拟耗时操作
	end := time.Now()
	duration := end.Sub(start)

	fmt.Printf("开始时间: %v\n", start.Format("15:04:05.000"))
	fmt.Printf("结束时间: %v\n", end.Format("15:04:05.000"))
	fmt.Printf("耗时: %v\n", duration)
	fmt.Printf("耗时(秒): %.3f\n", duration.Seconds())
	fmt.Printf("耗时(毫秒): %.0f\n", duration.Milliseconds())
	fmt.Printf("耗时(微秒): %.0f\n", duration.Microseconds())
	fmt.Printf("耗时(分钟): %.2f\n", duration.Minutes())
	fmt.Printf("耗时(小时): %.3f\n", duration.Hours())
	fmt.Println()

	// 9. 时区操作
	fmt.Println("9. 时区操作:")

	// 获取本地时区
	fmt.Printf("本地时间: %v\n", now.Local())
	fmt.Printf("本地时区: %v\n", now.Location())

	// 转换为UTC时间
	utcTime := now.UTC()
	fmt.Printf("UTC时间: %v\n", utcTime)

	// 指定时区
	nyLoc, _ := time.LoadLocation("America/New_York")
	nyTime := now.In(nyLoc)
	fmt.Printf("纽约时间: %v\n", nyTime)

	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
	shTime := now.In(shanghaiLoc)
	fmt.Printf("上海时间: %v\n", shTime)

	londonLoc, _ := time.LoadLocation("Europe/London")
	londonTime := now.In(londonLoc)
	fmt.Printf("伦敦时间: %v\n", londonTime)

	// 时区转换
	fmt.Printf("纽约 → 上海: %v\n", nyTime.In(shanghaiLoc))
	fmt.Println()

	// 10. 时间戳转换
	fmt.Println("10. 时间戳转换:")
	currentTimestamp := time.Now().Unix()

	fmt.Printf("当前时间戳(秒): %d\n", currentTimestamp)
	fmt.Printf("当前时间戳(毫秒): %d\n", time.Now().UnixMilli())

	// 时间戳转时间
	timeFromTimestamp := time.Unix(currentTimestamp, 0)
	fmt.Printf("时间戳转时间: %v\n", timeFromTimestamp.Format("2006-01-02 15:04:05"))

	// 带纳秒的时间戳
	timeFromNano := time.Unix(0, now.UnixNano())
	fmt.Printf("纳秒时间戳转时间: %v\n", timeFromNano.Format("2006-01-02 15:04:05.999999999"))
	fmt.Println()

	// 11. 定时器和休眠
	fmt.Println("11. 定时器操作:")
	fmt.Println("开始等待...")

	// 简单休眠
	time.Sleep(1 * time.Second)
	fmt.Println("1秒后")

	// 使用 After
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("2秒时间到!")
	}

	// 使用 Ticker
	fmt.Println("开始计时(3秒):")
	ticker := time.NewTicker(1 * time.Second)
	for i := 0; i < 3; i++ {
		<-ticker.C
		fmt.Printf("滴答 %d\n", i+1)
	}
	ticker.Stop()
	fmt.Println()

	// 12. 时间常量
	fmt.Println("12. 时间常量:")
	fmt.Printf("1纳秒: %v\n", time.Nanosecond)
	fmt.Printf("1微秒: %v\n", time.Microsecond)
	fmt.Printf("1毫秒: %v\n", time.Millisecond)
	fmt.Printf("1秒: %v\n", time.Second)
	fmt.Printf("1分钟: %v\n", time.Minute)
	fmt.Printf("1小时: %v\n", time.Hour)
	fmt.Printf("1天: %v\n", 24*time.Hour)
	fmt.Printf("1周: %v\n", 7*24*time.Hour)
	fmt.Println()

	// 13. 时间判断
	fmt.Println("13. 时间判断:")
	testTime := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	fmt.Printf("是否是周末: %t\n", isWeekend(testTime))
	fmt.Printf("是否是工作日: %t\n", isWeekday(testTime))
	fmt.Printf("是否是未来时间: %t\n", isFuture(testTime))
	fmt.Printf("是否是过去时间: %t\n", isPast(testTime))
}

// 判断是否是周末
func isWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

// 判断是否是工作日
func isWeekday(t time.Time) bool {
	return !isWeekend(t)
}

// 判断是否是未来时间
func isFuture(t time.Time) bool {
	return t.After(time.Now())
}

// 判断是否是过去时间
func isPast(t time.Time) bool {
	return t.Before(time.Now())
}
