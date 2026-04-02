# XTool 构建指南

## 快速开始

### 开发模式

```bash
make dev
```

### 构建 Windows 版本

```bash
# 构建 EXE
make build-windows

# 构建安装包（需要 NSIS）
make build-installer
```

## 安装包问题排查

### 问题 1: 桌面出现 `.link` 文件夹

**原因**: 安装脚本错误地创建了文件夹而非快捷方式

**解决方案**: 使用本项目提供的 `installer.nsi` 脚本：
- 使用 `CreateShortCut` 创建 `.lnk` 快捷方式
- 不是 `.link` 文件夹

### 问题 2: 二次安装新功能未生效

**原因**:
1. 浏览器缓存了旧的安装包
2. 前端代码未重新编译
3. 数据库缓存

**解决方案**:

```bash
# 1. 完全清理
make clean
rm -rf frontend/node_modules/.vite
rm -rf frontend/dist

# 2. 重新构建
make build-windows

# 3. 生成新安装包
make build-installer

# 4. 安装前卸载旧版本
# 控制面板 → 程序 → 卸载 XTool
```

## NSIS 安装说明

1. 下载 NSIS: https://nsis.sourceforge.io/Download
2. 安装后添加到 PATH: `C:\Program Files (x86)\NSIS`
3. 运行 `make build-installer`

## 版本更新流程

```bash
# 1. 更新版本号
# 编辑 wails.json 中的 productVersion
# 编辑 installer.nsi 中的 VERSION

# 2. 清理并构建
make clean
make build-installer

# 3. 测试安装包
# 运行 XTool_Setup_*.exe

# 4. 发布
# 上传到 GitHub Releases
```

## 常见问题

### Q: 启动后白屏

A: 检查前端是否正确构建：
```bash
cd frontend
npm run build
```

### Q: 功能不显示

A: 检查路由配置 `frontend/src/router/index.ts`

### Q: 数据丢失

A: 数据存储在 `%APPDATA%/xtool/` 目录
