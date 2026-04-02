<div align="center">

# XTool

**一站式桌面工具集**

![Version](https://img.shields.io/badge/version-1.0.0-blue?style=flat-square)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-brightgreen?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square)
![Vue](https://img.shields.io/badge/Vue-3.4+-42b883?style=flat-square)

功能丰富 · 启动极快 · 界面现代 · 动效流畅 · 完全本地

</div>

---

## ✨ 特性

| 特性 | 说明 |
|------|------|
| 🚀 **极速启动** | 基于 Go + Wails，启动速度 < 800ms |
| 💎 **现代界面** | 云母/亚克力效果，深色/浅色主题，精致动效 |
| 🛠️ **功能丰富** | 40+ 实用工具，覆盖开发、系统、日常、网络 |
| 🔒 **完全本地** | 不联网、不收集数据、不留后门 |
| 📦 **轻量便携** | 打包体积 < 30MB，内存占用 < 120MB |
| 🌍 **跨平台** | 支持 Windows 10/11、Linux、macOS |

---

## 📦 下载

| 平台 | 下载 |
|------|------|
| Windows 10/11 | `XTool.exe` |
| Linux (amd64) | `xtool_linux_amd64.tar.gz` |
| macOS (Intel) | `xtool_macos_intel.dmg` |
| macOS (Apple Silicon) | `xtool_macos_arm64.dmg` |

> 前往 [Releases](https://github.com/MasterPick/universal-toolbox/releases) 下载最新版本

---

## 🛠️ 功能清单

### 📝 开发工具

| 工具 | 功能描述 |
|------|----------|
| **JSON 工具** | 格式化、压缩、校验、转义/反转义 |
| **JSON 对比** | 两个 JSON 差异对比，高亮显示 |
| **XML 工具** | XML 格式化美化 |
| **Base64** | 文本 Base64 编码/解码 |
| **URL 编解码** | URL 编码与解码 |
| **哈希计算** | MD5 / SHA1 / SHA256 / SHA512 |
| **加密解密** | AES / DES / RSA / Base64 编解码 |
| **文本工具** | 字符统计、大小写转换、查找替换、文本对比 |
| **二维码** | 生成二维码，自定义样式，导出 PNG/SVG |
| **UUID 生成** | 批量生成 UUID v1/v4 |
| **时间戳** | 时间戳 ↔ 日期时间互转，多时区支持 |
| **正则测试** | 正则表达式实时测试，匹配高亮，常用模板 |
| **代码片段** | 保存管理常用代码片段，支持语法高亮 |

### 💻 系统工具

| 工具 | 功能描述 |
|------|----------|
| **进程管理** | 查看/搜索/终止进程，资源占用排序 |
| **端口查看** | 查看端口占用，一键释放端口 |
| **系统信息** | CPU/内存/磁盘使用率，硬件详情 |
| **批量重命名** | 文件批量重命名，支持正则表达式 |
| **取色器** | 屏幕取色，颜色格式转换，历史记录 |
| **图片工具** | 格式转换(PNG/JPG/WebP)、压缩、尺寸调整 |
| **文件批量处理** | 文件搜索、批量复制/移动/删除 |

### 📅 日常工具

| 工具 | 功能描述 |
|------|----------|
| **计算器** | 标准 + 科学计算模式 |
| **单位换算** | 长度/重量/温度/速度/面积/体积 |
| **备忘录** | 便签式备忘管理，置顶、搜索、颜色分类 |

### 🌐 网络工具

| 工具 | 功能描述 |
|------|----------|
| **Ping 测试** | 主机连通性检测，延迟统计 |
| **内网扫描** | 局域网设备发现，端口扫描 |
| **HTTP 测试** | 接口请求测试，响应查看 |
| **DNS 查询** | 多种记录类型查询(A/AAAA/MX/TXT/NS) |
| **Hosts 编辑** | 系统 Hosts 文件编辑管理 |

### ⚙️ 设置

| 功能 | 描述 |
|------|------|
| **基础设置** | 主题切换、语言、窗口行为 |
| **主题编辑器** | 自定义颜色配置，实时预览，导入导出 |
| **快捷键管理** | 自定义快捷键，启用/禁用，搜索过滤 |

---

## 🏗️ 技术架构

```
XTool
├── 后端 (Go)
│   ├── Wails v2         # 跨平台桌面框架
│   ├── gopsutil         # 系统信息获取
│   ├── SQLite           # 本地数据存储
│   └── 模块化设计       # 各功能独立、可插拔
│
└── 前端 (Vue3)
    ├── Vue 3.4          # 响应式框架
    ├── TypeScript       # 类型安全
    ├── TailwindCSS      # 原子化 CSS
    ├── Pinia            # 状态管理
    ├── Vue Router       # 路由管理
    └── Lucide Icons     # 图标库
```

---

## 💻 本地开发

### 环境要求

- Go 1.22+
- Node.js 22+
- Wails CLI v2.12+

### 开发步骤

```bash
# 1. 克隆仓库
git clone https://github.com/MasterPick/universal-toolbox.git
cd universal-toolbox

# 2. 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 3. 安装前端依赖
cd frontend && npm install && cd ..

# 4. 开发模式运行（热重载）
wails dev

# 5. 构建生产版本
wails build
```

### 项目结构

```
universal-toolbox/
├── main.go                 # 应用入口
├── internal/               # Go 后端模块
│   ├── devtools/           # 开发工具模块
│   ├── sysinfo/            # 系统工具模块
│   ├── network/            # 网络工具模块
│   ├── daily/              # 日常工具模块
│   ├── config/             # 配置管理
│   └── db/                 # 数据库
├── frontend/               # Vue3 前端
│   ├── src/
│   │   ├── views/          # 页面组件
│   │   ├── components/     # 公共组件
│   │   ├── stores/         # 状态管理
│   │   └── router/         # 路由配置
│   └── wailsjs/            # Wails 绑定（自动生成）
└── build/                  # 构建输出
```

---

## 🎯 开发规范

- **Go 代码**：完整中文注释，错误全部捕获
- **Vue 组件**：Composition API，TypeScript 类型定义
- **样式**：TailwindCSS 原子类，CSS 变量主题系统
- **交互**：动效流畅，状态反馈清晰，空状态友好

---

## 📄 许可证

[MIT License](LICENSE) © 2026 MasterPick
