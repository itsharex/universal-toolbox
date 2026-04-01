<div align="center">

<img src="docs/logo.png" alt="XTool Logo" width="120" height="120">

# XTool - 万能工具箱

![Version](https://img.shields.io/badge/version-1.0.0-blue?style=flat-square)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-brightgreen?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)
![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square)
![Vue](https://img.shields.io/badge/Vue-3.4-42b883?style=flat-square)

**一站式桌面工具集 · 功能丰富 · 启动极快 · 界面现代**

[立即下载](#-下载安装) · [功能一览](#-功能一览) · [开发指南](#-开发指南)

</div>

---

## 📖 简介

XTool 是一款现代化的桌面工具集，专为开发者和高级用户设计。集成了 **40+ 实用工具**，覆盖开发调试、系统管理、日常办公、网络诊断等场景。

### ✨ 核心特性

| 特性 | 说明 |
|------|------|
| 🚀 **极速启动** | Go 原生编译，启动时间 < 500ms |
| 💎 **现代界面** | Windows 11 云母效果，深色/浅色主题 |
| 🛠️ **功能丰富** | 40+ 工具，持续更新 |
| 🔒 **隐私安全** | 纯本地运行，无网络请求，零数据上传 |
| 📦 **轻量便携** | 单文件 < 25MB，无需安装 |
| 🌍 **跨平台** | Windows / Linux / macOS |

---

## 📦 下载安装

### Windows

| 版本 | 下载 | 说明 |
|------|------|------|
| Windows 10/11 | [XTool.exe](https://github.com/MasterPick/universal-toolbox/releases/latest) | 推荐，支持云母效果 |
| Windows 7/8 | [XTool-win7.exe](https://github.com/MasterPick/universal-toolbox/releases/latest) | 兼容版本 |

### 其他平台

| 平台 | 文件 |
|------|------|
| Linux (amd64) | `xtool_linux_amd64.tar.gz` |
| Linux (arm64) | `xtool_linux_arm64.tar.gz` |
| macOS (Intel) | `xtool_macos_amd64.dmg` |
| macOS (Apple Silicon) | `xtool_macos_arm64.dmg` |

---

## 🛠️ 功能一览

### 📝 开发工具

| 工具 | 功能描述 |
|------|----------|
| **JSON 工具** | 格式化、压缩、校验、转义/反转义 |
| **JSON 对比** | 两个 JSON 差异对比，高亮显示 |
| **XML 工具** | XML 格式化美化 |
| **Base64** | 文本/文件 Base64 编码解码 |
| **URL 编解码** | URL 编码与解码 |
| **哈希计算** | MD5 / SHA1 / SHA256 / SHA512 |
| **加密解密** | AES / DES 加密解密，MD5/SHA 哈希 |
| **文本工具** | 字符统计、大小写转换、去重 |
| **二维码** | 生成/解析二维码，支持 Logo |
| **UUID 生成** | 批量生成 UUID v1/v4 |
| **时间戳** | 时间戳 ↔ 日期时间互转 |
| **正则测试** | 正则表达式实时匹配测试 |
| **代码片段** | 保存管理常用代码片段 |

### 💻 系统工具

| 工具 | 功能描述 |
|------|----------|
| **进程管理** | 查看进程列表，强制终止进程 |
| **端口查看** | 查看端口占用，一键释放端口 |
| **系统信息** | CPU/内存/磁盘实时监控 |
| **批量重命名** | 文件批量重命名，支持正则表达式 |
| **取色器** | 屏幕取色，8 种颜色格式转换 |
| **图片工具** | 格式转换(PNG/JPG/WebP)、压缩、调整尺寸 |
| **文件批量处理** | 批量复制/移动/删除文件 |

### 🔧 日常工具

| 工具 | 功能描述 |
|------|----------|
| **计算器** | 标准 + 科学计算器 |
| **单位换算** | 长度/重量/温度/速度/时间换算 |
| **备忘录** | 轻量级便签管理 |

### 🌐 网络工具

| 工具 | 功能描述 |
|------|----------|
| **Ping 测试** | 主机连通性检测，延迟统计 |
| **内网扫描** | 局域网设备发现 |
| **HTTP 测试** | HTTP 接口请求测试 |
| **DNS 查询** | DNS 记录查询(A/AAAA/MX/TXT/NS/CNAME) |
| **Hosts 编辑** | 系统 Hosts 文件可视化管理 |

### ⚙️ 设置

| 功能 | 描述 |
|------|------|
| **主题编辑器** | 自定义主题配色，实时预览 |
| **快捷键管理** | 自定义快捷键绑定 |

---

## 🎨 界面预览

<div align="center">

| 深色主题 | 浅色主题 |
|:--------:|:--------:|
| ![Dark Theme](docs/screenshot-dark.png) | ![Light Theme](docs/screenshot-light.png) |

</div>

---

## 🏗️ 技术架构

```
XTool
│
├── 后端 (Go)
│   ├── Wails v2          # 跨平台桌面应用框架
│   ├── gopsutil          # 系统信息采集
│   ├── SQLite            # 本地数据持久化
│   └── 模块化设计        # 各功能独立、可插拔
│
└── 前端 (Vue3)
    ├── Vue 3.4           # 响应式框架
    ├── TypeScript        # 类型安全
    ├── TailwindCSS       # 原子化 CSS
    ├── Pinia             # 状态管理
    ├── Vue Router        # 路由管理
    └── Lucide Icons      # 图标库
```

---

## 💻 开发指南

### 环境要求

| 工具 | 版本 |
|------|------|
| Go | 1.22+ |
| Node.js | 22+ |
| Wails CLI | v2.12+ |

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/MasterPick/universal-toolbox.git
cd universal-toolbox

# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式（热重载）
wails dev

# 构建生产版本
wails build

# 构建并打包
wails build -clean -upx
```

### 项目结构

```
universal-toolbox/
├── main.go                 # 应用入口
├── internal/               # Go 后端模块
│   ├── config/            # 配置管理
│   ├── db/                # 数据库
│   ├── devtools/          # 开发工具
│   ├── sysinfo/           # 系统工具
│   ├── network/           # 网络工具
│   ├── daily/             # 日常工具
│   └── advanced/          # 高级功能
├── frontend/              # Vue3 前端
│   ├── src/
│   │   ├── views/        # 页面组件
│   │   ├── components/   # 公共组件
│   │   ├── stores/       # 状态管理
│   │   ├── router/       # 路由配置
│   │   └── assets/       # 静态资源
│   └── package.json
├── build/                 # 构建输出
└── Makefile              # 构建脚本
```

---

## 🤝 参与贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

---

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给一个 Star ⭐**

Made with ❤️ by [MasterPick](https://github.com/MasterPick)

</div>
