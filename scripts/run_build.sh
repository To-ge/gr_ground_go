#!/bin/bash

# 引数がない場合のエラー処理
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <local|linux>"
  exit 1
fi

# ビルドモードの判定
MODE=$1
INPUT_FILE="./cmd/main.go"
OUTPUT_DIR="./bin"
OUTPUT_FILE="gr-ground-go"

# 出力ディレクトリの作成
mkdir -p "$OUTPUT_DIR"

# モードに応じたビルド
if [ "$MODE" == "local" ]; then
  echo "Building for local environment..."
  GOOS=$(go env GOOS)
  GOARCH=$(go env GOARCH)
elif [ "$MODE" == "linux" ]; then
  echo "Building for RaspberryPi4 ..."
  GOOS="linux"
  GOARCH="arm64"
else
  echo "Invalid mode: $MODE"
  echo "Usage: $0 <local|linux>"
  exit 1
fi

# ビルド実行
CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$OUTPUT_DIR/$OUTPUT_FILE" "$INPUT_FILE"

# ビルド結果の確認
if [ $? -eq 0 ]; then
  echo "Build successful! Output: $OUTPUT_DIR/$OUTPUT_FILE"
else
  echo "Build failed!"
  exit 1
fi
