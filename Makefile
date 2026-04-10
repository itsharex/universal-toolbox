# ========================================================
# XTool Makefile
# 支持本地构建所有平台安装包
# ========================================================

APP_NAME := XTool
APP_NAME_EN := xtool
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := build/bin
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 版本号注入 ldflags
LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.buildDate=$(BUILD_DATE) \
	-X main.gitCommit=$(GIT_COMMIT)

.PHONY: all dev build-windows build-linux build-macos build-macos-arm64 \
       package lint clean install-deps help _deb-package

# 默认目标：显示帮助
help:
	@echo "XTool 构建系统 v$(VERSION)"
	@echo ""
	@echo "用法："
	@echo "  make dev              启动开发模式（热重载）"
	@echo "  make build-windows    构建 Windows .exe（需要 Windows 环境）"
	@echo "  make build-linux      构建 Linux 二进制 + .deb + tar.gz"
	@echo "  make build-macos      构建 macOS .app + .dmg (Intel amd64)"
	@echo "  make build-macos-arm64 构建 macOS .app + .dmg (Apple Silicon)"
	@echo "  make package          构建所有平台安装包（当前平台）"
	@echo "  make lint             运行代码检查（go vet + 前端类型检查）"
	@echo "  make all              构建当前平台版本"
	@echo "  make install-deps     安装构建依赖"
	@echo "  make clean            清理构建产物"

# 安装构建依赖
install-deps:
	@echo "安装 Go 依赖..."
	go mod download
	@echo "安装 Wails CLI..."
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	@echo "安装前端依赖..."
	cd frontend && npm install
	@echo "依赖安装完成"

# 代码检查：Go vet + 前端类型检查
lint:
	@echo "运行 Go vet..."
	go vet ./...
	@echo "Go vet 通过"
	@echo "运行前端类型检查..."
	cd frontend && npx vue-tsc --noEmit
	@echo "前端类型检查通过"

# 开发模式（热重载）
dev:
	@echo "启动开发模式..."
	wails dev

# 构建当前平台
all:
	@echo "构建当前平台..."
	CGO_ENABLED=1 wails build -ldflags "$(LDFLAGS)"
	@echo "构建完成：$(BUILD_DIR)/"

# Windows 构建（需要 Windows 环境或交叉编译工具链）
build-windows:
	@echo "构建 Windows amd64..."
	CGO_ENABLED=1 wails build -platform windows/amd64 -ldflags "$(LDFLAGS)" -o "$(APP_NAME).exe"
	@echo "Windows 构建完成：$(BUILD_DIR)/$(APP_NAME).exe"
	@echo ""
	@echo "如需生成 NSIS 安装包，请在 Windows 环境下运行："
	@echo "   makensis installer.nsi"

# Linux 构建 + .deb 打包
build-linux:
	@echo "构建 Linux amd64..."
	CGO_ENABLED=1 wails build -platform linux/amd64 -ldflags "$(LDFLAGS)" -o "$(APP_NAME_EN)"
	@echo "打包 tar.gz..."
	cd $(BUILD_DIR) && tar -czf $(APP_NAME_EN)_linux_amd64.tar.gz $(APP_NAME_EN)
	@echo "生成 .deb 包..."
	$(MAKE) _deb-package
	@echo "Linux 构建完成"
	ls -la $(BUILD_DIR)/

# macOS 构建 + .dmg 打包 (Intel amd64)
build-macos:
	@echo "构建 macOS amd64 (Intel)..."
	CGO_ENABLED=1 wails build -platform darwin/amd64 -ldflags "$(LDFLAGS)"
	@echo "创建 DMG (Intel)..."
	hdiutil create \
		-volname "$(APP_NAME)" \
		-srcfolder "$(BUILD_DIR)/$(APP_NAME).app" \
		-ov \
		-format UDZO \
		-imagekey zlib-level=9 \
		"$(BUILD_DIR)/$(APP_NAME_EN)_macos_amd64.dmg"
	@echo "macOS Intel 构建完成"
	ls -la $(BUILD_DIR)/

# macOS 构建 + .dmg 打包 (Apple Silicon arm64)
build-macos-arm64:
	@echo "构建 macOS arm64 (Apple Silicon)..."
	CGO_ENABLED=1 wails build -platform darwin/arm64 -ldflags "$(LDFLAGS)"
	@echo "创建 DMG (Apple Silicon)..."
	hdiutil create \
		-volname "$(APP_NAME)" \
		-srcfolder "$(BUILD_DIR)/$(APP_NAME).app" \
		-ov \
		-format UDZO \
		-imagekey zlib-level=9 \
		"$(BUILD_DIR)/$(APP_NAME_EN)_macos_arm64.dmg"
	@echo "macOS Apple Silicon 构建完成"
	ls -la $(BUILD_DIR)/

# 统一打包目标：构建当前平台所有产物
package:
	@echo "========================================="
	@echo "  XTool 统一打包 v$(VERSION)"
	@echo "========================================="
	@echo ""
	@$(MAKE) all
	@echo ""
	@echo "打包完成，产物目录：$(BUILD_DIR)/"
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
	printf 'Package: $(APP_NAME_EN)\nVersion: %s\nSection: utils\nPriority: optional\nArchitecture: amd64\nMaintainer: MasterPick <masterpickself@outlook.com>\nDepends: libgtk-3-0, libwebkit2gtk-4.0-37\nDescription: XTool - One-stop Desktop Toolbox\n' "$${VERSION#v}" > $$DEB_DIR/DEBIAN/control && \
	fakeroot dpkg-deb --build $$DEB_DIR $(BUILD_DIR)/$(APP_NAME_EN)_$${VERSION#v}_amd64.deb && \
	rm -rf $$DEB_DIR

# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -rf $(BUILD_DIR)/*
	@echo "清理完成"
