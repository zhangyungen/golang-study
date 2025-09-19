package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ContentType 常量
const (
	ContentTypeJSON      = "application/json"
	ContentTypeXML       = "application/xml"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeMultipart = "multipart/form-data"
	ContentTypeText      = "text/plain"
	ContentTypeHTML      = "text/html"
)

// HTTPMethod HTTP 方法
type HTTPMethod string

const (
	MethodGet     HTTPMethod = "GET"
	MethodPost    HTTPMethod = "POST"
	MethodPut     HTTPMethod = "PUT"
	MethodDelete  HTTPMethod = "DELETE"
	MethodPatch   HTTPMethod = "PATCH"
	MethodHead    HTTPMethod = "HEAD"
	MethodOptions HTTPMethod = "OPTIONS"
)

// RequestOptions 请求选项
type RequestOptions struct {
	Headers    map[string]string
	Query      map[string]string
	Timeout    time.Duration
	Context    context.Context
	RetryTimes int
	RetryDelay time.Duration
}

// Response 响应结构
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
}

// FilePart 文件上传部分
type FilePart struct {
	FieldName string
	FileName  string
	FilePath  string
	Reader    io.Reader
}

// HTTPClient HTTP 客户端
type HTTPClient struct {
	baseURL     string
	client      *http.Client
	defaultOpts RequestOptions
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(baseURL string, opts ...func(*HTTPClient)) *HTTPClient {
	client := &HTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		defaultOpts: RequestOptions{
			Headers: map[string]string{
				"User-Agent": "Go-HTTP-Client/1.0",
			},
			Timeout:    30 * time.Second,
			RetryTimes: 0,
			RetryDelay: 1 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithTimeout 设置默认超时时间
func WithTimeout(timeout time.Duration) func(*HTTPClient) {
	return func(hc *HTTPClient) {
		hc.defaultOpts.Timeout = timeout
		hc.client.Timeout = timeout
	}
}

// WithDefaultHeaders 设置默认请求头
func WithDefaultHeaders(headers map[string]string) func(*HTTPClient) {
	return func(hc *HTTPClient) {
		for k, v := range headers {
			hc.defaultOpts.Headers[k] = v
		}
	}
}

// WithRetry 设置重试策略
func WithRetry(times int, delay time.Duration) func(*HTTPClient) {
	return func(hc *HTTPClient) {
		hc.defaultOpts.RetryTimes = times
		hc.defaultOpts.RetryDelay = delay
	}
}

// request 基础请求方法
func (hc *HTTPClient) request(method HTTPMethod, path string, body io.Reader, opts *RequestOptions) (*Response, error) {
	start := time.Now()

	// 合并选项
	options := hc.mergeOptions(opts)

	// 构建完整URL
	fullURL, err := hc.buildURL(path, options.Query)
	if err != nil {
		return nil, err
	}

	// 创建请求
	req, err := http.NewRequestWithContext(
		hc.getContext(&options),
		string(method),
		fullURL,
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	hc.setHeaders(req, options.Headers)

	var resp *http.Response
	var lastErr error

	// 重试机制
	for attempt := 0; attempt <= options.RetryTimes; attempt++ {
		if attempt > 0 {
			time.Sleep(options.RetryDelay)
		}

		resp, lastErr = hc.client.Do(req)
		if lastErr == nil && resp.StatusCode < 500 {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("请求失败: %w", lastErr)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       bodyBytes,
		Duration:   time.Since(start),
	}, nil
}

// buildURL 构建完整URL
func (hc *HTTPClient) buildURL(path string, query map[string]string) (string, error) {
	u, err := url.Parse(hc.baseURL)
	if err != nil {
		return "", err
	}

	u.Path = strings.TrimSuffix(u.Path, "/") + "/" + strings.TrimPrefix(path, "/")

	// 添加查询参数
	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

// setHeaders 设置请求头
func (hc *HTTPClient) setHeaders(req *http.Request, headers map[string]string) {
	// 设置默认头
	for k, v := range hc.defaultOpts.Headers {
		req.Header.Set(k, v)
	}

	// 设置自定义头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

// getContext 获取上下文
func (hc *HTTPClient) getContext(opts *RequestOptions) context.Context {
	if opts.Context != nil {
		return opts.Context
	}
	return context.Background()
}

// mergeOptions 合并选项
func (hc *HTTPClient) mergeOptions(opts *RequestOptions) RequestOptions {
	merged := hc.defaultOpts

	if opts != nil {
		if opts.Timeout > 0 {
			merged.Timeout = opts.Timeout
		}
		if opts.Context != nil {
			merged.Context = opts.Context
		}
		if opts.RetryTimes > 0 {
			merged.RetryTimes = opts.RetryTimes
		}
		if opts.RetryDelay > 0 {
			merged.RetryDelay = opts.RetryDelay
		}

		// 合并头
		for k, v := range opts.Headers {
			merged.Headers[k] = v
		}

		// 合并查询参数
		if merged.Query == nil {
			merged.Query = make(map[string]string)
		}
		for k, v := range opts.Query {
			merged.Query[k] = v
		}
	}

	return merged
}

// Get GET 请求
func (hc *HTTPClient) Get(path string, opts *RequestOptions) (*Response, error) {
	return hc.request(MethodGet, path, nil, opts)
}

// Post POST 请求
func (hc *HTTPClient) Post(path string, body interface{}, opts *RequestOptions) (*Response, error) {
	return hc.withBody(MethodPost, path, body, opts)
}

// Put PUT 请求
func (hc *HTTPClient) Put(path string, body interface{}, opts *RequestOptions) (*Response, error) {
	return hc.withBody(MethodPut, path, body, opts)
}

// Delete DELETE 请求
func (hc *HTTPClient) Delete(path string, opts *RequestOptions) (*Response, error) {
	return hc.request(MethodDelete, path, nil, opts)
}

// Patch PATCH 请求
func (hc *HTTPClient) Patch(path string, body interface{}, opts *RequestOptions) (*Response, error) {
	return hc.withBody(MethodPatch, path, body, opts)
}

// Head HEAD 请求
func (hc *HTTPClient) Head(path string, opts *RequestOptions) (*Response, error) {
	return hc.request(MethodHead, path, nil, opts)
}

// Options OPTIONS 请求
func (hc *HTTPClient) Options(path string, opts *RequestOptions) (*Response, error) {
	return hc.request(MethodOptions, path, nil, opts)
}

// withBody 处理带请求体的方法
func (hc *HTTPClient) withBody(method HTTPMethod, path string, body interface{}, opts *RequestOptions) (*Response, error) {
	var reader io.Reader
	var contentType string

	if body != nil {
		switch v := body.(type) {
		case []byte:
			reader = bytes.NewReader(v)
			contentType = ContentTypeText
		case string:
			reader = strings.NewReader(v)
			contentType = ContentTypeText
		case url.Values:
			reader = strings.NewReader(v.Encode())
			contentType = ContentTypeForm
		case io.Reader:
			reader = v
		default:
			// 默认为 JSON
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("JSON序列化失败: %w", err)
			}
			reader = bytes.NewReader(jsonData)
			contentType = ContentTypeJSON
		}
	}

	// 设置内容类型
	if opts == nil {
		opts = &RequestOptions{}
	}
	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}
	if _, exists := opts.Headers["Content-Type"]; !exists && contentType != "" {
		opts.Headers["Content-Type"] = contentType
	}

	return hc.request(method, path, reader, opts)
}

// PostJSON POST JSON 数据
func (hc *HTTPClient) PostJSON(path string, data interface{}, opts *RequestOptions) (*Response, error) {
	if opts == nil {
		opts = &RequestOptions{}
	}
	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}
	opts.Headers["Content-Type"] = ContentTypeJSON

	return hc.Post(path, data, opts)
}

// PostForm POST 表单数据
func (hc *HTTPClient) PostForm(path string, formData map[string]string, opts *RequestOptions) (*Response, error) {
	values := url.Values{}
	for k, v := range formData {
		values.Add(k, v)
	}

	if opts == nil {
		opts = &RequestOptions{}
	}
	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}
	opts.Headers["Content-Type"] = ContentTypeForm

	return hc.Post(path, values, opts)
}

// PostMultipart POST 多部分表单数据（文件上传）
func (hc *HTTPClient) PostMultipart(path string, formData map[string]string, files []FilePart, opts *RequestOptions) (*Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加表单字段
	for key, value := range formData {
		writer.WriteField(key, value)
	}

	// 添加文件
	for _, file := range files {
		var part io.Writer
		var err error

		if file.Reader != nil {
			part, err = writer.CreateFormFile(file.FieldName, file.FileName)
		} else if file.FilePath != "" {
			fileReader, err := os.Open(file.FilePath)
			if err != nil {
				return nil, fmt.Errorf("打开文件失败: %w", err)
			}
			defer fileReader.Close()

			part, err = writer.CreateFormFile(file.FieldName, filepath.Base(file.FilePath))
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, fileReader)
		}

		if err != nil {
			return nil, fmt.Errorf("创建表单文件失败: %w", err)
		}
	}

	writer.Close()

	if opts == nil {
		opts = &RequestOptions{}
	}
	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}
	opts.Headers["Content-Type"] = writer.FormDataContentType()

	return hc.Post(path, body, opts)
}

// DownloadFile 下载文件
func (hc *HTTPClient) DownloadFile(url, filepath string, opts *RequestOptions) error {
	resp, err := hc.Get(url, opts)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入文件
	_, err = file.Write(resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// JSON 解析响应为 JSON
func (r *Response) JSON(target interface{}) error {
	return json.Unmarshal(r.Body, target)
}

// XML 解析响应为 XML
func (r *Response) XML(target interface{}) error {
	return xml.Unmarshal(r.Body, target)
}

// Text 获取响应文本
func (r *Response) Text() string {
	return string(r.Body)
}

// Bytes 获取响应字节
func (r *Response) Bytes() []byte {
	return r.Body
}

// IsSuccess 检查是否成功
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// GetHeader 获取响应头
func (r *Response) GetHeader(key string) string {
	return r.Headers.Get(key)
}
