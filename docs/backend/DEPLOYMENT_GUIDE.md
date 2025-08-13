# 🚀 后端部署指南

Go后端生产环境部署完整指南，包含Docker部署、传统部署、性能调优等内容。

## 📋 部署概览

### 部署方式
- **Docker部署** (推荐): 容器化部署，环境一致性好
- **传统部署**: 直接在服务器上运行二进制文件
- **Kubernetes部署**: 大规模集群部署

### 环境要求
- **操作系统**: Linux (Ubuntu 20.04+ / CentOS 8+)
- **CPU**: 4核心以上
- **内存**: 8GB以上
- **磁盘**: 100GB以上 SSD
- **网络**: 带宽10Mbps以上

## 🐳 Docker部署 (推荐)

### 1. 创建Dockerfile
```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# 复制二进制文件和配置
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8080
CMD ["./main"]
```

### 2. 创建docker-compose.yml
```yaml
version: '3.8'

services:
  # Go后端服务
  go-backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
      - DB_HOST=mysql
      - REDIS_HOST=redis
      - RABBITMQ_HOST=rabbitmq
    depends_on:
      - mysql
      - redis
      - rabbitmq
    volumes:
      - ./logs:/root/logs
    restart: unless-stopped

  # MySQL数据库
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: admin1234
      MYSQL_DATABASE: go_vue_admin
      MYSQL_USER: admin
      MYSQL_PASSWORD: admin123
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./sql:/docker-entrypoint-initdb.d
    command: --default-authentication-plugin=mysql_native_password
    restart: unless-stopped

  # Redis缓存
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    restart: unless-stopped

  # RabbitMQ消息队列
  rabbitmq:
    image: rabbitmq:3-management
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: admin123
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    restart: unless-stopped

  # Nginx反向代理
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - go-backend
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
  rabbitmq_data:
```

### 3. 部署命令
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f go-backend

# 停止服务
docker-compose down

# 重启服务
docker-compose restart go-backend
```

## 🖥️ 传统部署

### 1. 环境准备
```bash
# 安装Go (如果需要编译)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装MySQL
sudo apt update
sudo apt install mysql-server
sudo mysql_secure_installation

# 安装Redis
sudo apt install redis-server
sudo systemctl enable redis-server

# 安装RabbitMQ
sudo apt install rabbitmq-server
sudo systemctl enable rabbitmq-server
sudo rabbitmq-plugins enable rabbitmq_management
```

### 2. 编译应用
```bash
# 克隆代码
git clone https://github.com/your-repo/go-vue-admin.git
cd go-vue-admin/server

# 安装依赖
go mod tidy

# 编译生产版本
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-vue-admin .

# 设置执行权限
chmod +x go-vue-admin
```

### 3. 配置文件
```yaml
# config.yaml
server:
  port: 8080
  mode: release

db:
  host: 127.0.0.1
  port: 3306
  username: admin
  password: admin123
  database: go_vue_admin
  maxIdleConns: 50
  maxOpenConns: 200
  setConnMaxLifetime: 3600
  connMaxIdleTime: 1800

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

rabbitmq:
  host: 127.0.0.1
  port: 5672
  username: admin
  password: admin123

jwt:
  secret: your-jwt-secret-key
  expire: 7200

log:
  level: info
  path: ./logs
```

### 4. 创建系统服务
```bash
# 创建服务文件
sudo tee /etc/systemd/system/go-vue-admin.service > /dev/null <<EOF
[Unit]
Description=Go Vue Admin Backend
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/go-vue-admin
ExecStart=/opt/go-vue-admin/go-vue-admin
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable go-vue-admin
sudo systemctl start go-vue-admin

# 查看服务状态
sudo systemctl status go-vue-admin
```

## 🌐 Nginx配置

### 1. 基础配置
```nginx
# /etc/nginx/sites-available/go-vue-admin
server {
    listen 80;
    server_name your-domain.com;

    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL配置
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;

    # 前端静态文件
    location / {
        root /var/www/go-vue-admin;
        try_files $uri $uri/ /index.html;
        
        # 缓存静态资源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # 后端API代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }

    # 健康检查
    location /health {
        proxy_pass http://127.0.0.1:8080;
        access_log off;
    }

    # 日志配置
    access_log /var/log/nginx/go-vue-admin.access.log;
    error_log /var/log/nginx/go-vue-admin.error.log;
}
```

### 2. 负载均衡配置
```nginx
upstream go-backend {
    server 127.0.0.1:8080 weight=1 max_fails=3 fail_timeout=30s;
    server 127.0.0.1:8081 weight=1 max_fails=3 fail_timeout=30s;
    server 127.0.0.1:8082 weight=1 max_fails=3 fail_timeout=30s;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    location /api/ {
        proxy_pass http://go-backend;
        # 其他配置...
    }
}
```

## 📊 监控配置

### 1. Prometheus配置
```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'go-vue-admin'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s
```

### 2. Grafana仪表板
- 导入Go应用监控模板
- 配置MySQL监控
- 配置Redis监控
- 配置系统资源监控

## 🔧 性能调优

### 1. 系统参数优化
```bash
# /etc/sysctl.conf
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_max_tw_buckets = 5000

# 应用配置
sudo sysctl -p
```

### 2. 数据库优化
```sql
-- MySQL配置优化
SET GLOBAL innodb_buffer_pool_size = 4294967296; -- 4GB
SET GLOBAL max_connections = 1000;
SET GLOBAL query_cache_size = 268435456; -- 256MB
```

### 3. Redis优化
```bash
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
```

## 🔒 安全配置

### 1. 防火墙设置
```bash
# UFW防火墙
sudo ufw enable
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw deny 8080/tcp  # 只允许内网访问
```

### 2. SSL证书
```bash
# 使用Let's Encrypt
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 3. 安全头配置
```nginx
# 安全头
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header X-Content-Type-Options "nosniff" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
add_header Content-Security-Policy "default-src 'self'" always;
```

## 📝 部署检查清单

### 部署前检查
- [ ] 代码已提交到生产分支
- [ ] 配置文件已更新
- [ ] 数据库迁移脚本已准备
- [ ] SSL证书已配置
- [ ] 监控系统已配置

### 部署后检查
- [ ] 服务正常启动
- [ ] 健康检查通过
- [ ] API接口正常响应
- [ ] 数据库连接正常
- [ ] 缓存服务正常
- [ ] 日志输出正常
- [ ] 监控指标正常

## 🆘 故障排查

### 常见问题
1. **服务启动失败**: 检查配置文件和依赖服务
2. **数据库连接失败**: 检查数据库服务和连接配置
3. **内存不足**: 调整JVM参数或增加服务器内存
4. **性能问题**: 查看监控指标，优化数据库查询

### 日志查看
```bash
# 应用日志
tail -f /opt/go-vue-admin/logs/app.log

# 系统日志
sudo journalctl -u go-vue-admin -f

# Nginx日志
tail -f /var/log/nginx/go-vue-admin.error.log
```

---

**最后更新**: 2025-07-29  
**维护者**: 运维团队  
**技术支持**: 查看故障排查文档或联系运维团队
