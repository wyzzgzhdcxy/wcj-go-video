#requires -Version 5.1
<#
.SYNOPSIS
    本地打包本项目 Wails v2 GUI 应用，输出到 build/bin/，并部署到 D:\app_mgr\my

.DESCRIPTION
    用法:
        .\scripts\build-local.ps1                   # 标准打包 + 部署
        .\scripts\build-local.ps1 -Clean            # 清理 build/bin/ 和 frontend/dist
        .\scripts\build-local.ps1 -SkipFrontend     # 跳过前端构建（wails build -s）
        .\scripts\build-local.ps1 -NoDeploy         # 只打包不部署
        .\scripts\build-local.ps1 -DeployDir D:\xxx # 自定义部署目录（默认 $env:app_output_dir，兜底 D:\app_mgr\my）
        .\scripts\build-local.ps1 -OutputDir dist   # 自定义输出目录
        .\scripts\build-local.ps1 -DryRun           # 只预览不执行

    前置:
        go / wails CLI（v2.x）
        前端: node / pnpm（wails 自动调用 frontend:install + frontend:build；-SkipFrontend 时跳过）

    输出:
        build/bin/<wails.json#outputfilename>.exe                 <- Wails 默认产物
        D:\app_mgr\my\<wails.json#outputfilename>.exe                <- 部署副本（拷贝前自动停止同名进程）

    说明:
        - 产物名自动从本项目根目录的 wails.json#outputfilename 读取，无需修改脚本
        - 默认 DeployDir = D:\app_mgr\my；如需改为其他目录，调用时传 -DeployDir
        - 默认 OutputDir = build/bin（相对项目根）；如传绝对路径则直接使用
        - 部署前自动 Stop-Process 同名进程，避免文件占用
    部署目录（app_output_dir）:
        取值优先级：-DeployDir 参数 > 环境变量 app_output_dir > 内置兜底 D:\app_mgr\my
        查看当前值：$env:app_output_dir
        修改：[Environment]::SetEnvironmentVariable('app_output_dir', 'D:\app_mgr\my', 'User')
#>

[CmdletBinding()]
param(
    [switch]$Clean,
    [switch]$SkipFrontend,
    [switch]$NoDeploy,
    [string]$DeployDir,       # 部署目录；留空则取环境变量 app_output_dir，兜底 D:\app_mgr\my
    [string]$OutputDir = "build\bin",
    [switch]$DryRun
)

$ErrorActionPreference = 'Stop'

# --- resolve deploy dir ------------------------------------------------------
# 部署目录取值优先级：-DeployDir 参数 > 环境变量 app_output_dir > 内置兜底目录
$FallbackDeployDir = "D:\app_mgr\my"
$envDeployDir = [Environment]::GetEnvironmentVariable('app_output_dir')
$envDeployDir = if ($null -ne $envDeployDir) { $envDeployDir.Trim() } else { '' }

if ($DeployDir) {
    $DeployDirSource = '-DeployDir 参数'
} elseif ($envDeployDir) {
    $DeployDir = $envDeployDir
    $DeployDirSource = '环境变量 app_output_dir'
} else {
    $DeployDir = $FallbackDeployDir
    $DeployDirSource = '脚本内置默认值'
}

# 去掉结尾多余的分隔符，避免拼出 "D:\app_mgr\my\\app.exe"；盘符根目录保留
$DeployDir = $DeployDir.Trim().TrimEnd('\', '/')
if ($DeployDir -match '^[A-Za-z]:$') { $DeployDir += '\' }

function Write-Step($t) { Write-Host "`n==> $t" -ForegroundColor Cyan }
function Write-Ok($t)   { Write-Host $t -ForegroundColor Green }
function Write-Warn($t) { Write-Host $t -ForegroundColor Yellow }

function Require-Command($cmd) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        throw "未找到命令: $cmd，请先安装并加入 PATH"
    }
}

# 停止指定进程（不带 .exe 后缀），用于解锁目标 EXE
function Stop-ExeProcess([string]$name) {
    $procs = Get-Process -Name $name -ErrorAction SilentlyContinue
    if ($procs) {
        Write-Host "  停止进程: $name (PID: $(($procs.Id) -join ', '))" -ForegroundColor Yellow
        if ($DryRun) {
            Write-Host "  (DRYRUN) Stop-Process: $name"
        } else {
            $procs | Stop-Process -Force -ErrorAction SilentlyContinue
            Start-Sleep -Milliseconds 500
        }
    }
}

$root = Split-Path $PSScriptRoot -Parent

# --- 从 wails.json 读取 outputfilename ----------------------------------------
$wailsJsonPath = Join-Path $root "wails.json"
if (-not (Test-Path -LiteralPath $wailsJsonPath)) {
    throw "未找到 wails.json: $wailsJsonPath"
}
$wailsJson = Get-Content -LiteralPath $wailsJsonPath -Encoding UTF8 -Raw | ConvertFrom-Json
if (-not $wailsJson.outputfilename) {
    throw "wails.json 缺少 outputfilename 字段"
}
$exeName = "$($wailsJson.outputfilename).exe"
Write-Host "产物名（wails.json#outputfilename）: $exeName"

# --- preflight ---------------------------------------------------------------
Write-Step "检查构建环境"
$knownToolDirs = @(
    "E:\application\golang\go\bin",      # go 备选安装位
    "E:\application\nodejs",             # node / pnpm 备选安装位
    "$env:USERPROFILE\go\bin",           # wails（go install 默认 GOBIN）
    "$env:ProgramFiles\nodejs",          # node 标准安装位
    "$env:LOCALAPPDATA\Programs\Go\bin"  # go 标准安装位
) | Where-Object { $_ -and (Test-Path -LiteralPath $_) }

$pathEntries = @($env:Path -split ';' | Where-Object { $_ } | ForEach-Object { $_.TrimEnd('\') })
$added = @()
foreach ($d in $knownToolDirs) {
    $norm = $d.TrimEnd('\')
    if ($pathEntries -notcontains $norm) {
        $env:Path = "$d;$env:Path"
        $added += $d
    }
}
if ($added.Count -gt 0) {
    Write-Warn "已自动追加 PATH: $($added -join '; ')"
}

Require-Command go
Require-Command wails
Write-Host "  $(& go version)"
Write-Host "  wails $(& wails version)"

# --- output dir --------------------------------------------------------------
$binDir = if ([System.IO.Path]::IsPathRooted($OutputDir)) { $OutputDir } else { Join-Path $root $OutputDir }
if ($Clean) {
    foreach ($p in @($binDir, (Join-Path $root "frontend\dist"))) {
        if (Test-Path -LiteralPath $p) {
            Write-Step "清理: $p"
            if (-not $DryRun) { Remove-Item -LiteralPath $p -Recurse -Force }
        }
    }
}
if (-not (Test-Path -LiteralPath $binDir)) {
    if ($DryRun) {
        Write-Host "  (DRYRUN) 创建目录: $binDir"
    } else {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    }
}
Write-Host "输出目录: $binDir"

# --- build -------------------------------------------------------------------
$wailsArgs = @('build', '-platform', 'windows/amd64', '-trimpath')
if ($SkipFrontend) {
    $wailsArgs += '-s'
    if (-not (Test-Path -LiteralPath (Join-Path $root "frontend\dist\index.html"))) {
        throw "frontend/dist 不存在或为空，请去掉 -SkipFrontend 重新构建"
    }
} else {
    $wailsArgs += '-clean'   # wails 会自动执行 frontend:install + frontend:build
}

Write-Step "wails $($wailsArgs -join ' ')"
Push-Location $root
try {
    if ($DryRun) {
        Write-Host "  (DRYRUN) 在 $root 执行: wails $($wailsArgs -join ' ')"
    } else {
        & wails @wailsArgs
        if ($LASTEXITCODE -ne 0) { throw "wails build 失败 (exit=$LASTEXITCODE)" }
    }
} finally {
    Pop-Location
}

$exe = Join-Path $binDir $exeName
if (-not $DryRun -and -not (Test-Path -LiteralPath $exe)) {
    throw "未找到预期产物: $exe"
}
Write-Ok "  ✓ 已生成: $exe"

# --- deploy ------------------------------------------------------------------
if ($NoDeploy) {
    Write-Warn "已跳过部署 (-NoDeploy)"
} else {
    Write-Step "部署: $DeployDir  (来源: $DeployDirSource)"
    if (-not (Test-Path -LiteralPath $DeployDir)) {
        if ($DryRun) {
            Write-Host "  (DRYRUN) 创建目录: $DeployDir"
        } else {
            New-Item -ItemType Directory -Path $DeployDir -Force | Out-Null
        }
    }

    # 部署前停止同名进程（按 wails.json#outputfilename）
    $baseName = [System.IO.Path]::GetFileNameWithoutExtension($exeName)
    Stop-ExeProcess $baseName

    if ($DryRun) {
        Write-Host "  (DRYRUN) Copy-Item $exe -> $(Join-Path $DeployDir $exeName)"
    } else {
        Copy-Item -LiteralPath $exe -Destination (Join-Path $DeployDir $exeName) -Force
        Write-Ok "  ✓ 已部署: $(Join-Path $DeployDir $exeName)"
    }
}

# --- summary -----------------------------------------------------------------
Write-Step "打包完成"
if (Test-Path -LiteralPath $binDir) {
    Get-ChildItem -LiteralPath $binDir -File | Sort-Object Name |
        Select-Object Name, @{n='Size(MB)';e={[math]::Round($_.Length/1MB,2)}} |
        Format-Table -AutoSize | Out-String | Write-Host
}
Write-Host "本地产物: $exe" -ForegroundColor Green
if (-not $NoDeploy -and $DeployDir) {
    Write-Host "部署位置: $(Join-Path $DeployDir $exeName)" -ForegroundColor Green
}