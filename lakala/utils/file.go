package utils

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// ReadFileString 读取文本文件全部内容
func ReadFileString(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	builder := &strings.Builder{}
	reader := bufio.NewReader(file)
	for {
		chunk, err := reader.ReadString('\n')
		builder.WriteString(chunk)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return builder.String(), nil
}

// ReadProperties 解析 properties 文件到 map
func ReadProperties(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	props := make(map[string]string)
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		if idx := strings.IndexAny(line, "=:"); idx >= 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			value = decodeUnicodeEscapes(value)
			props[key] = value
		}
	}
	if err := reader.Err(); err != nil {
		return nil, err
	}
	return props, nil
}

func decodeUnicodeEscapes(input string) string {
	if !strings.Contains(input, `\u`) {
		return input
	}
	builder := strings.Builder{}
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' && i+1 < len(input) && input[i+1] == 'u' && i+5 < len(input) {
			hex := input[i+2 : i+6]
			if val, err := strconv.ParseInt(hex, 16, 32); err == nil {
				builder.WriteRune(rune(val))
				i += 5
				continue
			}
		}
		builder.WriteByte(input[i])
	}
	return builder.String()
}
