package io

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"zyj.com/golang-study/util/timeutil"
)

// 文件复制示例
func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	bytesCopied, err := io.Copy(dst, src)
	if err != nil {
		return err
	}

	fmt.Printf("成功复制 %d 字节\n", bytesCopied)
	return nil
}

// 网络数据保存示例
func SaveHTTPResponse(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// 复制文件前N个字节（创建预览）
func CopyFirstNBytes(srcPath, dstPath string, n int64) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.CopyN(dst, src, n)
	return err
}

// 读取配置文件
func ReadToByte(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return data, err
	}
	// 解析配置数据...
	return data, nil
}

// 写入日志文件
func WriteLog(filepath, message string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	timestamp := timeutil.NowString(timeutil.FormatDateTime)
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	_, err = io.WriteString(file, logEntry)
	return err
}

// 向多个目标写入相同字符串
func BroadcastWrite(message string, writers ...io.Writer) error {
	multiWriter := io.MultiWriter(writers...)
	_, err := io.WriteString(multiWriter, message)
	return err
}

// 合并多个文件内容
func concatFiles(outputPath string, inputPaths ...string) error {
	readers := make([]io.Reader, len(inputPaths))

	for i, path := range inputPaths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		readers[i] = file
	}

	multiReader := io.MultiReader(readers...)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, multiReader)
	return err
}

// 分块处理大数据
func processInChunks(chunks ...[]byte) []byte {
	var readers []io.Reader
	for _, chunk := range chunks {
		readers = append(readers, bytes.NewReader(chunk))
	}

	multiReader := io.MultiReader(readers...)
	var result bytes.Buffer
	io.Copy(&result, multiReader)

	return result.Bytes()
}

// 日志多路输出：同时输出到文件和控制台
func setupLogging(filepath string) {
	logFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
}

// 数据备份：同时写入主存储和备份
func saveDataWithBackup(data []byte, primaryWriter, backupWriter io.Writer) error {
	multiWriter := io.MultiWriter(primaryWriter, backupWriter)
	_, err := multiWriter.Write(data)
	return err
}

// 实时数据监控
type Monitor struct {
	writers []io.Writer
}

func (m *Monitor) Write(p []byte) (n int, err error) {
	multiWriter := io.MultiWriter(m.writers...)
	return multiWriter.Write(p)
}

func (m *Monitor) AddWriter(w io.Writer) {
	m.writers = append(m.writers, w)
}

// 下载文件并计算哈希值
func downloadWithHashCheck(url, filepath string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	// 数据同时流向文件和哈希计算器
	teeReader := io.TeeReader(resp.Body, hasher)

	_, err = io.Copy(file, teeReader)
	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}

// 请求/响应日志记录
func logHTTPRequest(req *http.Request) {
	var reqBody bytes.Buffer
	tee := io.TeeReader(req.Body, &reqBody)

	// 读取请求体并记录
	body, _ := io.ReadAll(tee)
	req.Body = io.NopCloser(&reqBody) // 恢复请求体

	fmt.Printf("请求体: %s\n", string(body))
}

// 流式数据处理管道
func processDataStream(input io.Reader) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		// 数据处理逻辑
		buffer := make([]byte, 1024)
		for {
			n, err := input.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				pw.CloseWithError(err)
				return
			}

			// 处理数据（示例：转换为大写）
			processed := bytes.ToUpper(buffer[:n])
			if _, err := pw.Write(processed); err != nil {
				return
			}
		}
	}()

	return pr
}

// 实时数据转换
func jsonToCSVConverter(jsonData io.Reader) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		var data []map[string]interface{}
		decoder := json.NewDecoder(jsonData)
		if err := decoder.Decode(&data); err != nil {
			pw.CloseWithError(err)
			return
		}

		if len(data) > 0 {
			// 写入CSV头部
			var headers []string
			for key := range data[0] {
				headers = append(headers, key)
			}
			csvWriter := csv.NewWriter(pw)
			csvWriter.Write(headers)

			// 写入数据行
			for _, row := range data {
				var record []string
				for _, header := range headers {
					record = append(record, fmt.Sprintf("%v", row[header]))
				}
				csvWriter.Write(record)
			}
			csvWriter.Flush()
		}
	}()

	return pr
}

// 读取固定大小的数据块
func readFixedSizeBlock(r io.Reader, blockSize int) ([]byte, error) {
	buf := make([]byte, blockSize)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
