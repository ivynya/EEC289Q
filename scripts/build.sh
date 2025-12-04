OS_ARCHS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
)
for os_arch in "${OS_ARCHS[@]}"; do
  IFS=/ read -r GOOS GOARCH <<< "$os_arch"
  BIN_NAME="tsp_solver_${GOOS}_${GOARCH}"
  if [ "$GOOS" = "windows" ]; then
    BIN_NAME+=".exe"
  fi
  env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "bin/$BIN_NAME" ./cmd/solver
done
