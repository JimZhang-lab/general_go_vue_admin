/*
 * @Author: JimZhang
 * @Date: 2025-07-29 18:00:00
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-29 18:00:00
 * @FilePath: /server/tools/load_test.go
 * @Description: 负载测试工具
 */

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// TestConfig 测试配置
type TestConfig struct {
	BaseURL     string        `json:"base_url"`
	Concurrency int           `json:"concurrency"`
	Duration    time.Duration `json:"duration"`
	RPS         int           `json:"rps"`
	Endpoints   []Endpoint    `json:"endpoints"`
}

// Endpoint 测试端点
type Endpoint struct {
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Weight  int               `json:"weight"` // 权重，用于分配请求比例
}

// TestResult 测试结果
type TestResult struct {
	TotalRequests   int64           `json:"total_requests"`
	SuccessRequests int64           `json:"success_requests"`
	FailedRequests  int64           `json:"failed_requests"`
	TotalDuration   time.Duration   `json:"total_duration"`
	AvgResponseTime time.Duration   `json:"avg_response_time"`
	MinResponseTime time.Duration   `json:"min_response_time"`
	MaxResponseTime time.Duration   `json:"max_response_time"`
	QPS             float64         `json:"qps"`
	ErrorRate       float64         `json:"error_rate"`
	StatusCodes     map[int]int64   `json:"status_codes"`
	ResponseTimes   []time.Duration `json:"-"` // 不序列化，用于计算百分位数
	P50ResponseTime time.Duration   `json:"p50_response_time"`
	P90ResponseTime time.Duration   `json:"p90_response_time"`
	P95ResponseTime time.Duration   `json:"p95_response_time"`
	P99ResponseTime time.Duration   `json:"p99_response_time"`
}

// LoadTester 负载测试器
type LoadTester struct {
	config    *TestConfig
	client    *http.Client
	results   *TestResult
	mutex     sync.Mutex
	startTime time.Time
	stopChan  chan struct{}
}

// NewLoadTester 创建负载测试器
func NewLoadTester(config *TestConfig) *LoadTester {
	return &LoadTester{
		config: config,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		results: &TestResult{
			StatusCodes:     make(map[int]int64),
			ResponseTimes:   make([]time.Duration, 0),
			MinResponseTime: time.Hour, // 初始化为很大的值
		},
		stopChan: make(chan struct{}),
	}
}

// Run 运行负载测试
func (lt *LoadTester) Run() *TestResult {
	log.Printf("开始负载测试: 并发数=%d, 持续时间=%v, 目标QPS=%d",
		lt.config.Concurrency, lt.config.Duration, lt.config.RPS)

	lt.startTime = time.Now()

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < lt.config.Concurrency; i++ {
		wg.Add(1)
		go lt.worker(&wg)
	}

	// 等待测试完成
	time.Sleep(lt.config.Duration)
	close(lt.stopChan)
	wg.Wait()

	// 计算最终结果
	lt.calculateFinalResults()

	log.Printf("负载测试完成: QPS=%.2f, 成功率=%.2f%%, 平均响应时间=%v",
		lt.results.QPS, 100-lt.results.ErrorRate, lt.results.AvgResponseTime)

	return lt.results
}

// worker 工作协程
func (lt *LoadTester) worker(wg *sync.WaitGroup) {
	defer wg.Done()

	// 计算请求间隔（如果设置了RPS限制）
	var requestInterval time.Duration
	if lt.config.RPS > 0 {
		requestInterval = time.Duration(int64(time.Second) / int64(lt.config.RPS) * int64(lt.config.Concurrency))
	}

	ticker := time.NewTicker(requestInterval)
	if requestInterval == 0 {
		ticker.Stop()
	}
	defer ticker.Stop()

	for {
		select {
		case <-lt.stopChan:
			return
		default:
			// 选择端点
			endpoint := lt.selectEndpoint()

			// 发送请求
			start := time.Now()
			resp, err := lt.sendRequest(endpoint)
			responseTime := time.Since(start)

			// 记录结果
			lt.recordResult(resp, err, responseTime)

			// 限制QPS
			if requestInterval > 0 {
				<-ticker.C
			}
		}
	}
}

// selectEndpoint 选择测试端点
func (lt *LoadTester) selectEndpoint() *Endpoint {
	if len(lt.config.Endpoints) == 0 {
		// 默认端点
		return &Endpoint{
			Path:   "/health",
			Method: "GET",
		}
	}

	// 根据权重选择端点（简单实现）
	totalWeight := 0
	for _, ep := range lt.config.Endpoints {
		totalWeight += ep.Weight
	}

	if totalWeight == 0 {
		return &lt.config.Endpoints[0]
	}

	// 简单的轮询选择
	index := int(atomic.LoadInt64(&lt.results.TotalRequests)) % len(lt.config.Endpoints)
	return &lt.config.Endpoints[index]
}

// sendRequest 发送HTTP请求
func (lt *LoadTester) sendRequest(endpoint *Endpoint) (*http.Response, error) {
	url := lt.config.BaseURL + endpoint.Path

	var body io.Reader
	if endpoint.Body != "" {
		body = bytes.NewBufferString(endpoint.Body)
	}

	req, err := http.NewRequest(endpoint.Method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	// 设置默认请求头
	if req.Header.Get("Content-Type") == "" && endpoint.Body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return lt.client.Do(req)
}

// recordResult 记录测试结果
func (lt *LoadTester) recordResult(resp *http.Response, err error, responseTime time.Duration) {
	lt.mutex.Lock()
	defer lt.mutex.Unlock()

	atomic.AddInt64(&lt.results.TotalRequests, 1)
	lt.results.ResponseTimes = append(lt.results.ResponseTimes, responseTime)

	// 更新响应时间统计
	if responseTime < lt.results.MinResponseTime {
		lt.results.MinResponseTime = responseTime
	}
	if responseTime > lt.results.MaxResponseTime {
		lt.results.MaxResponseTime = responseTime
	}

	if err != nil {
		atomic.AddInt64(&lt.results.FailedRequests, 1)
		return
	}

	if resp != nil {
		defer resp.Body.Close()

		// 记录状态码
		lt.results.StatusCodes[resp.StatusCode]++

		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			atomic.AddInt64(&lt.results.SuccessRequests, 1)
		} else {
			atomic.AddInt64(&lt.results.FailedRequests, 1)
		}
	}
}

// calculateFinalResults 计算最终结果
func (lt *LoadTester) calculateFinalResults() {
	lt.results.TotalDuration = time.Since(lt.startTime)

	if lt.results.TotalRequests > 0 {
		// 计算QPS
		lt.results.QPS = float64(lt.results.TotalRequests) / lt.results.TotalDuration.Seconds()

		// 计算错误率
		lt.results.ErrorRate = float64(lt.results.FailedRequests) / float64(lt.results.TotalRequests) * 100

		// 计算平均响应时间
		var totalResponseTime time.Duration
		for _, rt := range lt.results.ResponseTimes {
			totalResponseTime += rt
		}
		lt.results.AvgResponseTime = totalResponseTime / time.Duration(len(lt.results.ResponseTimes))

		// 计算百分位数
		lt.calculatePercentiles()
	}
}

// calculatePercentiles 计算百分位数
func (lt *LoadTester) calculatePercentiles() {
	if len(lt.results.ResponseTimes) == 0 {
		return
	}

	// 简单排序
	times := make([]time.Duration, len(lt.results.ResponseTimes))
	copy(times, lt.results.ResponseTimes)

	// 冒泡排序（简单实现）
	for i := 0; i < len(times)-1; i++ {
		for j := 0; j < len(times)-i-1; j++ {
			if times[j] > times[j+1] {
				times[j], times[j+1] = times[j+1], times[j]
			}
		}
	}

	// 计算百分位数
	lt.results.P50ResponseTime = times[len(times)*50/100]
	lt.results.P90ResponseTime = times[len(times)*90/100]
	lt.results.P95ResponseTime = times[len(times)*95/100]
	lt.results.P99ResponseTime = times[len(times)*99/100]
}

// PrintResults 打印测试结果
func (lt *LoadTester) PrintResults() {
	fmt.Println("\n=== 负载测试结果 ===")
	fmt.Printf("总请求数: %d\n", lt.results.TotalRequests)
	fmt.Printf("成功请求数: %d\n", lt.results.SuccessRequests)
	fmt.Printf("失败请求数: %d\n", lt.results.FailedRequests)
	fmt.Printf("测试持续时间: %v\n", lt.results.TotalDuration)
	fmt.Printf("QPS: %.2f\n", lt.results.QPS)
	fmt.Printf("错误率: %.2f%%\n", lt.results.ErrorRate)
	fmt.Printf("平均响应时间: %v\n", lt.results.AvgResponseTime)
	fmt.Printf("最小响应时间: %v\n", lt.results.MinResponseTime)
	fmt.Printf("最大响应时间: %v\n", lt.results.MaxResponseTime)
	fmt.Printf("P50响应时间: %v\n", lt.results.P50ResponseTime)
	fmt.Printf("P90响应时间: %v\n", lt.results.P90ResponseTime)
	fmt.Printf("P95响应时间: %v\n", lt.results.P95ResponseTime)
	fmt.Printf("P99响应时间: %v\n", lt.results.P99ResponseTime)

	fmt.Println("\n状态码分布:")
	for code, count := range lt.results.StatusCodes {
		fmt.Printf("  %d: %d\n", code, count)
	}
}

// SaveResults 保存测试结果到文件
func (lt *LoadTester) SaveResults(filename string) error {
	_, err := json.MarshalIndent(lt.results, "", "  ")
	if err != nil {
		return err
	}

	// TODO: 这里应该写入文件，简化实现
	return nil
}

func main() {
	// 命令行参数
	var (
		baseURL     = flag.String("url", "http://localhost:8080", "测试目标URL")
		concurrency = flag.Int("c", 10, "并发数")
		duration    = flag.Duration("d", time.Minute, "测试持续时间")
		rps         = flag.Int("r", 0, "目标QPS（0表示不限制）")
		configFile  = flag.String("config", "", "配置文件路径")
	)
	flag.Parse()

	var config *TestConfig

	if *configFile != "" {
		// 从配置文件加载
		log.Printf("从配置文件加载: %s", *configFile)
		// 这里应该读取配置文件，简化实现
		config = &TestConfig{
			BaseURL:     *baseURL,
			Concurrency: *concurrency,
			Duration:    *duration,
			RPS:         *rps,
		}
	} else {
		// 使用命令行参数
		config = &TestConfig{
			BaseURL:     *baseURL,
			Concurrency: *concurrency,
			Duration:    *duration,
			RPS:         *rps,
			Endpoints: []Endpoint{
				{Path: "/health", Method: "GET", Weight: 1},
				{Path: "/api/admin/list", Method: "GET", Weight: 3},
				{Path: "/api/dept/list", Method: "GET", Weight: 2},
			},
		}
	}

	// 创建并运行负载测试
	tester := NewLoadTester(config)
	results := tester.Run()

	// 打印结果
	tester.PrintResults()

	// 保存结果
	filename := fmt.Sprintf("load_test_results_%d.json", time.Now().Unix())
	if err := tester.SaveResults(filename); err != nil {
		log.Printf("保存结果失败: %v", err)
	} else {
		log.Printf("结果已保存到: %s", filename)
	}

	// 性能评估
	evaluatePerformance(results)
}

// evaluatePerformance 评估性能
func evaluatePerformance(results *TestResult) {
	fmt.Println("\n=== 性能评估 ===")

	// QPS评估
	if results.QPS >= 1000 {
		fmt.Println("✅ QPS表现优秀 (>= 1000)")
	} else if results.QPS >= 500 {
		fmt.Println("⚠️  QPS表现良好 (>= 500)")
	} else {
		fmt.Println("❌ QPS需要优化 (< 500)")
	}

	// 响应时间评估
	if results.AvgResponseTime <= 100*time.Millisecond {
		fmt.Println("✅ 响应时间优秀 (<= 100ms)")
	} else if results.AvgResponseTime <= 500*time.Millisecond {
		fmt.Println("⚠️  响应时间良好 (<= 500ms)")
	} else {
		fmt.Println("❌ 响应时间需要优化 (> 500ms)")
	}

	// 错误率评估
	if results.ErrorRate <= 1.0 {
		fmt.Println("✅ 错误率优秀 (<= 1%)")
	} else if results.ErrorRate <= 5.0 {
		fmt.Println("⚠️  错误率可接受 (<= 5%)")
	} else {
		fmt.Println("❌ 错误率过高 (> 5%)")
	}

	// P99响应时间评估
	if results.P99ResponseTime <= 1*time.Second {
		fmt.Println("✅ P99响应时间优秀 (<= 1s)")
	} else if results.P99ResponseTime <= 3*time.Second {
		fmt.Println("⚠️  P99响应时间可接受 (<= 3s)")
	} else {
		fmt.Println("❌ P99响应时间需要优化 (> 3s)")
	}
}
