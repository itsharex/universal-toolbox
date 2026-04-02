// Package advanced 高级扩展功能模块
// 提供主题切换、窗口控制、快捷键管理等高级功能
package advanced

import (
	"xtool/internal/config"
)

// AdvancedTools 高级工具结构体（Wails 绑定到前端）
type AdvancedTools struct {
	config *config.Config // 配置中心引用
}

// ShortcutConfig 快捷键配置结构
type ShortcutConfig struct {
	Action  string `json:"action"`  // 功能名称
	Hotkey  string `json:"hotkey"`  // 快捷键组合（如 Ctrl+Shift+T）
	Enabled bool   `json:"enabled"` // 是否启用
}

// ThemeOption 主题选项结构
type ThemeOption struct {
	ID          string `json:"id"`          // 主题 ID
	Name        string `json:"name"`        // 主题名称
	Description string `json:"description"` // 主题描述
	Preview     string `json:"preview"`     // 预览颜色（CSS 渐变）
}

// NewAdvancedTools 创建高级功能模块实例
func NewAdvancedTools(cfg *config.Config) *AdvancedTools {
	return &AdvancedTools{config: cfg}
}

// GetAvailableThemes 获取所有可用主题列表
func (a *AdvancedTools) GetAvailableThemes() []ThemeOption {
	return []ThemeOption{
		{
			ID:          "dark",
			Name:        "深色主题",
			Description: "护眼深色配色，适合夜间使用",
			Preview:     "linear-gradient(135deg, #1a1a2e, #16213e)",
		},
		{
			ID:          "light",
			Name:        "浅色主题",
			Description: "清爽浅色配色，适合日间使用",
			Preview:     "linear-gradient(135deg, #f0f2f5, #ffffff)",
		},
		{
			ID:          "auto",
			Name:        "跟随系统",
			Description: "自动跟随系统深色/浅色模式",
			Preview:     "linear-gradient(135deg, #667eea, #764ba2)",
		},
		{
			ID:          "blue",
			Name:        "深海蓝",
			Description: "以蓝色为主调的专业风格",
			Preview:     "linear-gradient(135deg, #0f2027, #203a43, #2c5364)",
		},
		{
			ID:          "green",
			Name:        "森林绿",
			Description: "清新自然的绿色主题",
			Preview:     "linear-gradient(135deg, #134e5e, #71b280)",
		},
		{
			ID:          "purple",
			Name:        "紫罗兰",
			Description: "优雅紫色风格主题",
			Preview:     "linear-gradient(135deg, #4776e6, #8e54e9)",
		},
		{
			ID:          "orange",
			Name:        "活力橙",
			Description: "活力充沛的暖色主题",
			Preview:     "linear-gradient(135deg, #f7971e, #ffd200)",
		},
		{
			ID:          "mica",
			Name:        "云母效果",
			Description: "Windows 11 原生云母半透明效果",
			Preview:     "linear-gradient(135deg, rgba(255,255,255,0.1), rgba(255,255,255,0.3))",
		},
	}
}

// SwitchTheme 切换应用主题
func (a *AdvancedTools) SwitchTheme(themeID string) error {
	return a.config.SetTheme(themeID)
}

// GetDefaultShortcuts 获取默认快捷键配置
func (a *AdvancedTools) GetDefaultShortcuts() []ShortcutConfig {
	return []ShortcutConfig{
		{Action: "搜索工具", Hotkey: "Ctrl+K", Enabled: true},
		{Action: "切换主题", Hotkey: "Ctrl+Shift+T", Enabled: true},
		{Action: "窗口置顶", Hotkey: "Ctrl+Shift+P", Enabled: true},
		{Action: "调小字体", Hotkey: "Ctrl+-", Enabled: true},
		{Action: "调大字体", Hotkey: "Ctrl+=", Enabled: true},
		{Action: "JSON 格式化", Hotkey: "Ctrl+Shift+F", Enabled: true},
		{Action: "快速备忘", Hotkey: "Ctrl+Shift+N", Enabled: true},
		{Action: "最小化到托盘", Hotkey: "Ctrl+M", Enabled: true},
	}
}

// GetAppVersion 获取应用版本信息
func (a *AdvancedTools) GetAppVersion() map[string]string {
	return map[string]string{
		"version":   "1.1.0",
		"buildDate": "2026-04-02",
		"author":    "MasterPick",
		"repo":      "https://github.com/MasterPick/xtool",
	}
}
