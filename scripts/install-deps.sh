#!/bin/bash

# pkg 层依赖安装脚本
# 使用方法：./install-deps.sh

echo "Installing pkg layer dependencies..."

# HTTP (Gin)
echo "Installing github.com/gin-gonic/gin..."
go get github.com/gin-gonic/gin

# Rate Limiter
echo "Installing golang.org/x/time/rate..."
go get golang.org/x/time/rate

# Service Discovery (Consul)
echo "Installing github.com/hashicorp/consul/api..."
go get github.com/hashicorp/consul/api

# OpenTelemetry (Tracing)
echo "Installing go.opentelemetry.io/otel..."
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/jaeger
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/semconv/v1.4.0

# Tidy up
echo "Running go mod tidy..."
go mod tidy

echo "All dependencies installed successfully!"
