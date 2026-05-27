#!/bin/bash

# 1. Định nghĩa các cấu hình cố định (ĐẢM BẢO KHÔNG CÓ DẤU "://" Ở ĐÂY)
REGISTRY="reg.attapps.com"
PROJECT="sinhnhatf168"
REPOSITORY="be_sinhnhatf168_new"

# 2. Kiểm tra thư mục Git
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo "❌ Lỗi: Thư mục này không phải là một Git repository!"
    exit 1
fi

# 3. Lấy tên nhánh và mã Commit hiện tại
GIT_BRANCH=$(git branch --show-current)
if [ -z "$GIT_BRANCH" ]; then
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
fi

# Lấy mã Commit ngắn (7 ký tự)
GIT_COMMIT=$(git rev-parse --short HEAD)

# Chuẩn hóa tên nhánh (thay thế dấu / bằng dấu - nếu có)
CLEAN_BRANCH=$(echo "$GIT_BRANCH" | sed 's/\//-/g')

# 4. Tạo Docker Tag kết hợp
DOCKER_TAG="${CLEAN_BRANCH}-${GIT_COMMIT}"

# SỬA TẠI ĐÂY: Đảm bảo nối chuỗi thuần túy dạng host/project/repo:tag
COMMIT_IMAGE_NAME="${REGISTRY}/${PROJECT}/${REPOSITORY}:${DOCKER_TAG}"
LATEST_IMAGE_NAME="${REGISTRY}/${PROJECT}/${REPOSITORY}:latest"

echo "=============================================="
echo "🌿 Nhánh Git: $GIT_BRANCH"
echo "📌 Mã Commit: $GIT_COMMIT"
echo "🏷️  Docker Tag: $DOCKER_TAG"
echo "🐳 Image mã định danh: $COMMIT_IMAGE_NAME"
echo "🐳 Image bản mới nhất: $LATEST_IMAGE_NAME"
echo "=============================================="

# 5. Tiến hành Build Docker Image
echo "🚀 Đang tiến hành build Docker image..."
docker build -t "$COMMIT_IMAGE_NAME" .

if [ $? -ne 0 ]; then
    echo "❌ Lỗi: Build Docker image thất bại!"
    exit 1
fi

# Gắn thêm tag latest cho bản build mới nhất này
docker tag "$COMMIT_IMAGE_NAME" "$LATEST_IMAGE_NAME"

# 6. Đẩy các Image lên Registry
echo "📤 Đang push các image lên registry..."
docker push "$COMMIT_IMAGE_NAME"
docker push "$LATEST_IMAGE_NAME"

if [ $? -eq 0 ]; then
    echo "=============================================="
    echo "✅ THÀNH CÔNG! Đã đẩy thành công 2 tags lên registry:"
    echo "   1. Bản cụ thể: $COMMIT_IMAGE_NAME"
    echo "   2. Bản mới nhất: $LATEST_IMAGE_NAME"
    echo "=============================================="
else
    echo "❌ Lỗi: Không thể push image. Vui lòng kiểm tra lại 'docker login $REGISTRY'."
    exit 1
fi
