# ========================================================
# XTool Makefile
# 支持本地构建所有平台安装包
# ========================================================

APP_NAME := XTool
APP_NAME_EN := xtool
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0")
BUILD_DIR := build/bin

.PHONY: all dev build build-windows build-installer build-linux build-macos clean install-deps help

# 默认目标：显示帮助
help:
	@echo "XTool 构建系统 v$(VERSION)"
	@echo ""
	@echo "用法："
	@echo "  make dev              启动开发模式（热重载）"
	@echo "  make build            构建当前平台"
	@echo "  make build-windows    构建 Windows EXE"
	@echo "  make build-installer  构建 Windows 安装包（需要 NSIS）"
	@echo "  make build-linux      构建 Linux 二进制 + .deb + tar.gz"
	@echo "  make build-macos      构建 macOS .app + .dmg"
	@echo "  make all              构建当前平台版本"
	@echo "  make install-deps     安装构建依赖"
	@echo "  make clean            清理构建产物"

# 安装构建依赖
install-deps:
	@echo "📦 安装 Go 依赖..."
	go mod download
	@echo "📦 安装 Wails CLI..."
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	@echo "📦 安装前端依赖..."
	cd frontend && npm install
	@echo "✅ 依赖安装完成"

# 开发模式（热重载）
dev:
	@echo "🚀 启动开发模式..."
	wails dev

# 构建当前平台
build:
	@echo "🔨 构建当前平台..."
	wails build -ldflags "-s -w"
	@echo "✅ 构建完成：$(BUILD_DIR)/"

all: build

# Windows 构建
build-windows:
	@echo "🪟 构建 Windows amd64..."
	@echo "📦 清理旧文件..."
	rm -rf $(BUILD_DIR)
	wails build -platform windows/amd64 -ldflags "-s -w" -skipbindings
	@echo "✅ Windows 构建完成：$(BUILD_DIR)/$(APP_NAME).exe"
	ls -la $(BUILD_DIR)/

# Windows 安装包（需要安装 NSIS）
build-installer: build-windows
	@echo "📦 生成安装包..."
	@echo "⚠️  需要 NSIS 已安装并添加到 PATH"
	makensis /DVERSION=$(VERSION) installer.nsi
	@echo "✅ 安装包生成完成"
	ls -la XTool_Setup_*.exe 2>/dev/null || echo "❌ 安装包生成失败，请检查 NSIS 是否安装"

# Linux 构建 + .deb 打包
build-linux:
	@echo "🐧 构建 Linux amd64..."
	wails build -platform linux/amd64 -ldflags "-s -w" -o "$(APP_NAME_EN)"
	@echo "📦 打包 tar.gz..."
	cd $(BUILD_DIR) && tar -czf $(APP_NAME_EN)_linux_amd64.tar.gz $(APP_NAME_EN)
	@echo "📦 生成 .deb 包..."
	$(MAKE) _deb-package
	@echo "✅ Linux 构建完成"
	ls -la $(BUILD_DIR)/

# 内部：生成 .deb 包
_deb-package:
	@VERSION=$(VERSION) && \
	DEB_DIR=$(APP_NAME_EN)_$${VERSION#v}_amd64 && \
	rm -rf $$DEB_DIR && \
	mkdir -p $$DEB_DIR/DEBIAN && \
	mkdir -p $$DEB_DIR/usr/bin && \
	mkdir -p $$DEB_DIR/usr/share/applications && \
	mkdir -p $$DEB_DIR/usr/share/doc/$(APP_NAME_EN) && \
	cp $(BUILD_DIR)/$(APP_NAME_EN) $$DEB_DIR/usr/bin/ && \
	chmod +x $$DEB_DIR/usr/bin/$(APP_NAME_EN) && \
	cp README.md $$DEB_DIR/usr/share/doc/$(APP_NAME_EN)/ && \
	cp LICENSE $$DEB_DIR/usr/share/doc/$(APP_NAME_EN)/ && \
	printf '[Desktop Entry]\nVersion=%s\nName=XTool\nExec=/usr/bin/$(APP_NAME_EN)\nTerminal=false\nType=Application\nCategories=Utility;\n' "$${VERSION#v}" > $$DEB_DIR/usr/share/applications/$(APP_NAME_EN).desktop && \
	printf 'Package: $(APP_NAME_EN)\nVersion: %s\nSection: utils\nPriority: optional\nArchitecture: amd64\nMaintainer: MasterPick <masterpickself@outlook.com>\nDescription: XTool\n' "$${VERSION#v}" > $$DEB_DIR/DEBIAN/control && \
	fakeroot dpkg-deb --build $$DEB_DIR $(BUILD_DIR)/$(APP_NAME_EN)_$${VERSION#v}_amd64.deb && \
	rm -rf $$DEB_DIR

# macOS 构建 + .dmg 打包
build-macos:
	@echo "🍎 构建 macOS..."
	wails build -platform darwin/amd64 -ldflags "-s -w"
	@echo "📦 创建 DMG..."
	hdiutil create \
		-volname "$(APP_NAME)" \
		-srcfolder "$(BUILD_DIR)/$(APP_NAME).app" \
		-ov \
		-format UDZO \
		-imagekey zlib-level=9 \
		"$(BUILD_DIR)/$(APP_NAME_EN)_macos.dmg"
	@echo "✅ macOS 构建完成"
	ls -la $(BUILD_DIR)/

# 清理构建产物
clean:
	@echo "🧹 清理构建产物..."
	rm -rf $(BUILD_DIR)/*
	rm -f XTool_Setup_*.exe
	@echo "✅ 清理完成"
