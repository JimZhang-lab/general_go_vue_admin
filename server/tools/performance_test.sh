#!/bin/bash

# 性能验证脚本
# Author: JimZhang
# Date: 2025-07-29 18:10:00
# Description: 自动化性能测试和验证脚本

set -e

# 配置变量
SERVER_URL="http://localhost:8080"
TEST_DURATION="2m"
CONCURRENCY_LEVELS=(10 50 100 200 500)
RESULTS_DIR="./performance_results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    # 检查Go是否安装
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go"
        exit 1
    fi
    
    # 检查curl是否安装
    if ! command -v curl &> /dev/null; then
        log_error "curl未安装，请先安装curl"
        exit 1
    fi
    
    # 检查jq是否安装（用于JSON处理）
    if ! command -v jq &> /dev/null; then
        log_warning "jq未安装，建议安装jq以获得更好的JSON处理体验"
    fi
    
    log_success "依赖检查完成"
}

# 检查服务器状态
check_server() {
    log_info "检查服务器状态..."
    
    if curl -s -f "${SERVER_URL}/health" > /dev/null; then
        log_success "服务器运行正常"
    else
        log_error "服务器未运行或无法访问: ${SERVER_URL}"
        exit 1
    fi
}

# 预热服务器
warmup_server() {
    log_info "预热服务器..."
    
    for i in {1..10}; do
        curl -s "${SERVER_URL}/health" > /dev/null
        curl -s "${SERVER_URL}/api/admin/list" > /dev/null || true
        sleep 0.1
    done
    
    log_success "服务器预热完成"
}

# 运行单个负载测试
run_load_test() {
    local concurrency=$1
    local output_file="${RESULTS_DIR}/load_test_c${concurrency}_${TIMESTAMP}.json"
    
    log_info "运行负载测试: 并发数=${concurrency}, 持续时间=${TEST_DURATION}"
    
    # 编译并运行负载测试工具
    go run tools/load_test.go \
        -url="${SERVER_URL}" \
        -c="${concurrency}" \
        -d="${TEST_DURATION}" \
        -r=0 > "${output_file}.log" 2>&1
    
    if [ $? -eq 0 ]; then
        log_success "负载测试完成: 并发数=${concurrency}"
    else
        log_error "负载测试失败: 并发数=${concurrency}"
        return 1
    fi
}

# 运行所有负载测试
run_all_tests() {
    log_info "开始运行所有负载测试..."
    
    mkdir -p "${RESULTS_DIR}"
    
    for concurrency in "${CONCURRENCY_LEVELS[@]}"; do
        run_load_test "${concurrency}"
        sleep 10  # 测试间隔，让服务器恢复
    done
    
    log_success "所有负载测试完成"
}

# 收集系统指标
collect_system_metrics() {
    log_info "收集系统指标..."
    
    local metrics_file="${RESULTS_DIR}/system_metrics_${TIMESTAMP}.json"
    
    # 获取系统指标
    curl -s "${SERVER_URL}/metrics?type=system" > "${metrics_file}" || true
    
    # 获取HTTP指标
    curl -s "${SERVER_URL}/metrics?type=http" > "${RESULTS_DIR}/http_metrics_${TIMESTAMP}.json" || true
    
    # 获取QPS信息
    curl -s "${SERVER_URL}/qps" > "${RESULTS_DIR}/qps_metrics_${TIMESTAMP}.json" || true
    
    log_success "系统指标收集完成"
}

# 生成性能报告
generate_report() {
    log_info "生成性能报告..."
    
    local report_file="${RESULTS_DIR}/performance_report_${TIMESTAMP}.md"
    
    cat > "${report_file}" << EOF
# 性能测试报告

**测试时间**: $(date)
**服务器地址**: ${SERVER_URL}
**测试持续时间**: ${TEST_DURATION}

## 测试概览

| 并发数 | QPS | 平均响应时间 | P99响应时间 | 错误率 | 状态 |
|--------|-----|-------------|-------------|--------|------|
EOF

    # 分析每个测试结果
    for concurrency in "${CONCURRENCY_LEVELS[@]}"; do
        local log_file="${RESULTS_DIR}/load_test_c${concurrency}_${TIMESTAMP}.json.log"
        
        if [ -f "${log_file}" ]; then
            # 从日志中提取关键指标（简化实现）
            local qps=$(grep "QPS=" "${log_file}" | tail -1 | sed 's/.*QPS=\([0-9.]*\).*/\1/' || echo "N/A")
            local success_rate=$(grep "成功率=" "${log_file}" | tail -1 | sed 's/.*成功率=\([0-9.]*\)%.*/\1/' || echo "N/A")
            local avg_time=$(grep "平均响应时间=" "${log_file}" | tail -1 | sed 's/.*平均响应时间=\([^,]*\).*/\1/' || echo "N/A")
            
            # 判断状态
            local status="❌"
            if (( $(echo "${qps} >= 500" | bc -l 2>/dev/null || echo 0) )); then
                status="✅"
            elif (( $(echo "${qps} >= 200" | bc -l 2>/dev/null || echo 0) )); then
                status="⚠️"
            fi
            
            echo "| ${concurrency} | ${qps} | ${avg_time} | N/A | $((100-${success_rate:-0}))% | ${status} |" >> "${report_file}"
        else
            echo "| ${concurrency} | N/A | N/A | N/A | N/A | ❌ |" >> "${report_file}"
        fi
    done

    cat >> "${report_file}" << EOF

## 性能分析

### QPS性能
- **目标**: >= 1000 QPS (优秀), >= 500 QPS (良好)
- **实际**: 请查看上表

### 响应时间
- **目标**: <= 100ms (优秀), <= 500ms (良好)
- **实际**: 请查看上表

### 错误率
- **目标**: <= 1% (优秀), <= 5% (可接受)
- **实际**: 请查看上表

## 优化建议

1. **如果QPS不达标**:
   - 检查数据库连接池配置
   - 启用Redis缓存
   - 优化数据库查询
   - 考虑使用RabbitMQ异步处理

2. **如果响应时间过长**:
   - 添加数据库索引
   - 优化SQL查询
   - 启用缓存机制
   - 减少不必要的数据传输

3. **如果错误率过高**:
   - 检查错误日志
   - 优化错误处理逻辑
   - 增加系统资源
   - 调整限流和熔断参数

## 系统资源使用

$(curl -s "${SERVER_URL}/metrics?type=system" 2>/dev/null | head -20 || echo "无法获取系统指标")

## 测试文件

- 详细日志: ${RESULTS_DIR}/
- 系统指标: ${RESULTS_DIR}/system_metrics_${TIMESTAMP}.json
- HTTP指标: ${RESULTS_DIR}/http_metrics_${TIMESTAMP}.json
- QPS指标: ${RESULTS_DIR}/qps_metrics_${TIMESTAMP}.json

EOF

    log_success "性能报告生成完成: ${report_file}"
}

# 验证性能目标
verify_performance_goals() {
    log_info "验证性能目标..."
    
    local passed=0
    local total=0
    
    # 检查QPS目标
    log_info "检查QPS目标 (>= 500)..."
    for concurrency in "${CONCURRENCY_LEVELS[@]}"; do
        local log_file="${RESULTS_DIR}/load_test_c${concurrency}_${TIMESTAMP}.json.log"
        if [ -f "${log_file}" ]; then
            local qps=$(grep "QPS=" "${log_file}" | tail -1 | sed 's/.*QPS=\([0-9.]*\).*/\1/' || echo "0")
            total=$((total + 1))
            
            if (( $(echo "${qps} >= 500" | bc -l 2>/dev/null || echo 0) )); then
                log_success "并发数${concurrency}: QPS=${qps} ✅"
                passed=$((passed + 1))
            else
                log_warning "并发数${concurrency}: QPS=${qps} ❌"
            fi
        fi
    done
    
    # 输出验证结果
    echo ""
    log_info "性能目标验证结果:"
    log_info "通过测试: ${passed}/${total}"
    
    if [ "${passed}" -eq "${total}" ] && [ "${total}" -gt 0 ]; then
        log_success "🎉 所有性能目标均已达成！"
        return 0
    else
        log_warning "⚠️  部分性能目标未达成，请查看详细报告"
        return 1
    fi
}

# 清理函数
cleanup() {
    log_info "清理临时文件..."
    # 这里可以添加清理逻辑
}

# 主函数
main() {
    echo "========================================="
    echo "         性能测试与验证工具"
    echo "========================================="
    echo ""
    
    # 设置清理陷阱
    trap cleanup EXIT
    
    # 执行测试流程
    check_dependencies
    check_server
    warmup_server
    run_all_tests
    collect_system_metrics
    generate_report
    
    # 验证性能目标
    if verify_performance_goals; then
        exit 0
    else
        exit 1
    fi
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--url)
            SERVER_URL="$2"
            shift 2
            ;;
        -d|--duration)
            TEST_DURATION="$2"
            shift 2
            ;;
        -o|--output)
            RESULTS_DIR="$2"
            shift 2
            ;;
        -h|--help)
            echo "用法: $0 [选项]"
            echo ""
            echo "选项:"
            echo "  -u, --url URL        服务器地址 (默认: http://localhost:8080)"
            echo "  -d, --duration TIME  测试持续时间 (默认: 2m)"
            echo "  -o, --output DIR     结果输出目录 (默认: ./performance_results)"
            echo "  -h, --help           显示帮助信息"
            echo ""
            echo "示例:"
            echo "  $0 -u http://localhost:8080 -d 5m"
            exit 0
            ;;
        *)
            log_error "未知参数: $1"
            exit 1
            ;;
    esac
done

# 运行主函数
main
