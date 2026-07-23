#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
linker_command="${LINKER_BIN:-linker}"
if [[ "$linker_command" == */* ]]; then
  linker_binary="$(cd "$(dirname "$linker_command")" && pwd)/$(basename "$linker_command")"
  if [[ ! -x "$linker_binary" ]]; then
    echo "无法验证 gRPC 脚手架：LINKER_BIN 不是可执行文件。" >&2
    exit 1
  fi
else
  linker_binary="$(command -v "$linker_command" || true)"
  if [[ -z "$linker_binary" ]]; then
    echo "无法验证 gRPC 脚手架：请先安装 linker，或通过 LINKER_BIN 指定可执行文件。" >&2
    exit 1
  fi
fi
if ! command -v qadc >/dev/null 2>&1; then
  echo "无法验证 gRPC 脚手架：当前环境缺少 qadc。" >&2
  exit 1
fi

stage="$(mktemp -d "${TMPDIR:-/tmp}/linker-v3-example-grpc.XXXXXX")"
finish() {
  status=$?
  trap - EXIT
  rm -rf "$stage"
  exit "$status"
}
trap finish EXIT

cp "$root/scaffold/grpc.yaml" "$stage/grpc.yaml"
(cd "$stage" && "$linker_binary" generate grpc grpc.yaml --dry-run --format json > plan.json)
(cd "$stage" && "$linker_binary" generate grpc grpc.yaml --format json > receipt.json)

project="$stage/example-grpc"
if [[ -n "$(cd "$project" && gofmt -l .)" ]]; then
  echo "gRPC 脚手架生成结果未通过 gofmt。" >&2
  exit 1
fi
(cd "$project" && GOWORK=off go mod tidy)
(cd "$project" && GOWORK=off go test -race ./...)
(cd "$project" && GOWORK=off go vet ./...)
(cd "$project" && GOWORK=off qadc gate --profile backend-service)

echo "gRPC 脚手架示例通过 dry-run、独立生成、race、vet 与 QADC 验证。"
