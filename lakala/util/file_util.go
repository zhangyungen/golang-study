package util

import (
	"os"
	"path/filepath"
	"zyj.com/golang-study/lakala"
)

// FileUtils 文件工具类
type FileUtils struct{}

// GetProperties 读取Properties文件
func (f *FileUtils) GetProperties(filePath string) (map[string]string, error) {
	content, err := f.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return f.parseProperties(content)
}

// ReadFile 读取文件内容
func (f *FileUtils) ReadFile(filePath string) (string, error) {
	// 检查文件是否存在
	if !f.FileExists(filePath) {
		return "", lakala.NewSDKExceptionFromEnumsWithInfo(
			lakala.INITIALIZE_KEYSTORE_ERROR,
			"找不到指定的文件["+filePath+"]",
		)
	}

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", lakala.NewSDKExceptionFromEnumsWithInfo(
			lakala.FILE_READ_FAIL_EXCEPTION,
			"读取文件内容出错: "+err.Error(),
		)
	}

	return string(data), nil
}

// ReadFileBytes 读取文件内容为字节数组
func (f *FileUtils) ReadFileBytes(filePath string) ([]byte, error) {
	if !f.FileExists(filePath) {
		return nil, lakala.NewSDKExceptionFromEnumsWithInfo(
			lakala.INITIALIZE_KEYSTORE_ERROR,
			"找不到指定的文件["+filePath+"]",
		)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, lakala.NewSDKExceptionFromEnumsWithInfo(
			lakala.FILE_READ_FAIL_EXCEPTION,
			"读取文件内容出错: "+err.Error(),
		)
	}

	return data, nil
}

// WriteFile 写入文件内容
func (f *FileUtils) WriteFile(filePath, content string) error {
	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return lakala.NewSDKException("创建目录失败", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return lakala.NewSDKException("写入文件失败", err)
	}

	return nil
}

// FileExists 检查文件是否存在
func (f *FileUtils) FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirectoryExists 检查目录是否存在
func (f *FileUtils) DirectoryExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// parseProperties 解析Properties格式内容
func (f *FileUtils) parseProperties(content string) (map[string]string, error) {
	result := make(map[string]string)

	// 简单的properties解析实现
	// 实际应用中可能需要更复杂的解析逻辑
	lines := splitLines(content)
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			continue // 跳过空行和注释
		}

		if idx := indexOf(line, '='); idx != -1 {
			key := trimSpace(line[:idx])
			value := trimSpace(line[idx+1:])
			if key != "" {
				result[key] = value
			}
		}
	}

	return result, nil
}

// 辅助函数：按行分割
func splitLines(content string) []string {
	var lines []string
	start := 0
	for i, ch := range content {
		if ch == '\n' {
			lines = append(lines, content[start:i])
			start = i + 1
		}
	}
	if start < len(content) {
		lines = append(lines, content[start:])
	}
	return lines
}

// 辅助函数：去除空格
func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// 辅助函数：查找字符位置
func indexOf(s string, ch byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			return i
		}
	}
	return -1
}

// 全局工具实例
var fileUtils = &FileUtils{}

// GetProperties 全局函数：读取Properties文件
func GetProperties(filePath string) (map[string]string, error) {
	return fileUtils.GetProperties(filePath)
}

// ReadFile 全局函数：读取文件内容
func ReadFile(filePath string) (string, error) {
	return fileUtils.ReadFile(filePath)
}

// ReadFileBytes 全局函数：读取文件内容为字节数组
func ReadFileBytes(filePath string) ([]byte, error) {
	return fileUtils.ReadFileBytes(filePath)
}

// WriteFile 全局函数：写入文件内容
func WriteFile(filePath, content string) error {
	return fileUtils.WriteFile(filePath, content)
}

// FileExists 全局函数：检查文件是否存在
func FileExists(filePath string) bool {
	return fileUtils.FileExists(filePath)
}

// DirectoryExists 全局函数：检查目录是否存在
func DirectoryExists(dirPath string) bool {
	return fileUtils.DirectoryExists(dirPath)
}
