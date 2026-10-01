#!/usr/bin/env sh
# Cross-compile anon for every OS the team uses. Output goes to ../server/public/dl/,
# so `wrangler deploy` serves them to the installer scripts.
set -e
cd "$(dirname "$0")"
mkdir -p ../server/public/dl

for target in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64; do
  os=${target%/*}
  arch=${target#*/}
  ext=""
  [ "$os" = "windows" ] && ext=".exe"
  out="../server/public/dl/anon-$os-$arch$ext"
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -trimpath -ldflags="-s -w" -o "$out" .
  echo "built $out"
done
