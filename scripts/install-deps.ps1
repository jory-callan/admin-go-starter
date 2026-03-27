# pkg 层依赖安装脚本 (PowerShell)
# 使用方法：.\scripts\install-deps.ps1

Write-Host "Installing pkg layer dependencies..." -ForegroundColor Green

# HTTP (Gin)
Write-Host "Installing github.com/gin-gonic/gin..." -ForegroundColor Cyan
go get github.com/gin-gonic/gin

# Rate Limiter
Write-Host "Installing golang.org/x/time/rate..." -ForegroundColor Cyan
go get golang.org/x/time/rate

# Service Discovery (Consul)
Write-Host "Installing github.com/hashicorp/consul/api..." -ForegroundColor Cyan
go get github.com/hashicorp/consul/api

# OpenTelemetry (Tracing)
Write-Host "Installing go.opentelemetry.io/otel..." -ForegroundColor Cyan
go get go.opentelemetry.io/otel
Write-Host "Installing go.opentelemetry.io/otel/exporters/jaeger..." -ForegroundColor Cyan
go get go.opentelemetry.io/otel/exporters/jaeger
Write-Host "Installing go.opentelemetry.io/otel/sdk..." -ForegroundColor Cyan
go get go.opentelemetry.io/otel/sdk
Write-Host "Installing go.opentelemetry.io/otel/semconv/v1.4.0..." -ForegroundColor Cyan
go get go.opentelemetry.io/otel/semconv/v1.4.0

# Tidy up
Write-Host "Running go mod tidy..." -ForegroundColor Yellow
go mod tidy

Write-Host "All dependencies installed successfully!" -ForegroundColor Green
