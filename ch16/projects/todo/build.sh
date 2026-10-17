#!/bin/sh
# Build todo for the common platforms into dist/, with the version stamped in.
#   ./ch16/projects/todo/build.sh 1.0.0
set -eu
version=${1:-dev}
cd "$(dirname "$0")"
rm -rf dist && mkdir dist
for target in darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64; do
  os=${target%/*}; arch=${target#*/}
  out="dist/todo-$os-$arch"
  [ "$os" = windows ] && out="$out.exe"
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -trimpath -ldflags "-s -w -X main.version=$version" -o "$out" .
  echo "built $out"
done
(cd dist && shasum -a 256 * > SHA256SUMS)
