// Package main XTool 主入口
// 使用 Wails v2 框架构建跨平台桌面应用
// 技术栈：Go + Vue3 + TypeScript + TailwindCSS + SQLite
package main

import (
	"embed"

	"xtool/internal/advanced"
	"xtool/internal/config"
	"xtool/internal/daily"
	"xtool/internal/db"
	"xtool/internal/devtools"
	applog "xtool/internal/log"
	"xtool/internal/network"
	"xtool/internal/sysinfo"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// 嵌入前端静态资源（编译时打包）
//
//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化日志系统
	logger := applog.NewLogger()
	logger.Info("XTool 启动中...")

	// 初始化数据库
	database, err := db.Init()
	if err != nil {
		logger.Error("数据库初始化失败: " + err.Error())
		panic(err)
	}
	defer database.Close()

	// 初始化配置中心
	cfg := config.NewConfig(database)

	// 初始化各功能模块
	devModule := devtools.NewDevTools(database)
	sysModule := sysinfo.NewSysInfo()
	dailyModule := daily.NewDailyTools(database)
	netModule := network.NewNetworkTools()
	advModule := advanced.NewAdvancedTools(cfg)

	// 创建 Wails 应用
	err = wails.Run(&options.App{
		Title:            "XTool",
		Width:            1280,
		Height:           800,
		MinWidth:         900,
		MinHeight:        600,
		DisableResize:    false,
		Fullscreen:       false,
		WindowStartState: options.Normal,
		Frameless:        true,  // 无边框窗口（使用自定义标题栏）

		// 嵌入前端资源
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		// 后台颜色（深色模式支持）
		BackgroundColour: &options.RGBA{R: 18, G: 18, B: 25, A: 255},

		// 绑定 Go 后端方法到前端
		Bind: []interface{}{
			devModule,
			sysModule,
			dailyModule,
			netModule,
			advModule,
			cfg,
		},

		// Windows 平台特定配置（亚克力/云母效果）
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			BackdropType:                      windows.Mica,     // Windows 11 云母效果
			DisablePinchZoom:                  false,
			DisableWindowIcon:                 false,
			IsZoomControlEnabled:              false,
			EnableSwipeGestures:               false,
			DisableFramelessWindowDecorations: true,   // 禁用系统标题栏装饰（避免双标题栏）
		},

		// 日志级别
		LogLevel: 2,
	})

	if err != nil {
		logger.Error("应用运行失败: " + err.Error())
		panic(err)
	}
}
