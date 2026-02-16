# build.ps1 - 使用 goreleaser 构建项目
# 用法:
#   .\build.ps1              # snapshot 构建（不需要 git tag）
#   .\build.ps1 -Release     # 正式发布构建（需要 git tag）
#   .\build.ps1 -Clean       # 清理 dist 目录

param(
    [switch]$Release,
    [switch]$Clean
)

$ErrorActionPreference = "Stop"

# 清理
if ($Clean) {
    Write-Host "[*] Cleaning dist/ ..." -ForegroundColor Yellow
    if (Test-Path "dist") {
        Remove-Item -Recurse -Force "dist"
    }
    Write-Host "[+] Done." -ForegroundColor Green
    exit 0
}

# 检查 goreleaser 是否安装
if (-not (Get-Command "goreleaser" -ErrorAction SilentlyContinue)) {
    Write-Host "[-] goreleaser not found. Install: go install github.com/goreleaser/goreleaser/v2@latest" -ForegroundColor Red
    exit 1
}

# 构建
if ($Release) {
    Write-Host "[*] Running goreleaser release ..." -ForegroundColor Cyan
    goreleaser release --clean
} else {
    Write-Host "[*] Running goreleaser snapshot ..." -ForegroundColor Cyan
    goreleaser release --snapshot --clean
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "[-] Build failed." -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "[+] Build succeeded. Output in dist/" -ForegroundColor Green
Get-ChildItem dist/*.tar.gz, dist/*.zip 2>$null | ForEach-Object {
    Write-Host "    $_" -ForegroundColor Gray
}
