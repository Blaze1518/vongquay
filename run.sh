#!/bin/bash
set -e

echo "🔍 Kiểm tra công cụ swag..."
if ! command -v swag &> /dev/null; then
    echo "📥 Cài đặt swag..."
    go install ://github.com
fi

echo "📝 Đang tự động tạo tài liệu Swagger docs..."
swag init -g cmd/server/main.go --parseDependency --dir ./

echo "🚀 Khởi chạy server Golang..."
go run cmd/server/main.go
