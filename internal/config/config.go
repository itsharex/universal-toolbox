// Package config 配置中心模块
// 统一管理应用主题、快捷键、用户偏好等配置项
// 所有配置持久化存储到 SQLite 数据库
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"xtool/internal/db"
)

// AppConfig 应用全局配置结构
type AppConfig struct {
	Theme        string  `json:"theme"`        // 主题：dark/light/auto/blue/green/purple
	Language     string  `json:"language"`     // 语言：zh-CN/en-US
	FontSize     int     `json:"fontSize"`     // 字体大小（px）
	Density      string  `json:"density"`      // 布局密度：compact/normal/spacious
	Transparency float64 `json:"transparency"` // 窗口透明度（0-1）
	AlwaysOnTop  bool    `json:"alwaysOnTop"`  // 窗口置顶
	MinToTray    bool    `json:"minToTray"`    // 最小化到托盘
	AutoBackup   bool    `json:"autoBackup"`   // 自动备份
}

// defaultConfig 返回默认配置
func defaultConfig() *AppConfig {
	return &AppConfig{
		Theme:        "dark",    // 默认深色主题
		Language:     "zh-CN",   // 默认简体中文
		FontSize:     14,        // 默认字体大小
		Density:      "normal",  // 默认正常密度
		Transparency: 0.95,      // 默认透明度
		AlwaysOnTop:  false,     // 默认不置顶
		MinToTray:    true,      // 默认最小化到托盘
		AutoBackup:   true,      // 默认自动备份
	}
}

// Config 配置管理器
type Config struct {
	mu     sync.RWMutex      // 读写锁，保护配置并发访问
	db     *db.Database      // 数据库实例
	config *AppConfig        // 当前配置
}

// NewConfig 创建配置管理器并加载已保存的配置
func NewConfig(database *db.Database) *Config {
	c := &Config{
		db:     database,
		config: defaultConfig(),
	}

	// 从数据库加载已保存的配置（错误不再静默忽略）
	if err := c.loadFromDB(); err != nil {
		log.Printf("[WARN] 从数据库加载配置失败，使用默认值: %v\n", err)
	}
	return c
}

// GetConfig 获取当前配置（供前端调用）
// 使用读锁保护并发访问
func (c *Config) GetConfig() *AppConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// SaveConfig 保存配置（供前端调用）
// 包含字段校验：FontSize 范围 10-30，Transparency 范围 0.1-1.0
func (c *Config) SaveConfig(configJSON string) error {
	// 解析前端传来的 JSON 配置
	var newConfig AppConfig
	if err := json.Unmarshal([]byte(configJSON), &newConfig); err != nil {
		return fmt.Errorf("配置 JSON 解析失败: %w", err)
	}

	// 字段校验
	if err := validateConfig(&newConfig); err != nil {
		return err
	}

	// 使用写锁保护配置更新
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新内存中的配置
	c.config = &newConfig

	// 持久化到数据库
	return c.saveToDB()
}

// GetTheme 获取当前主题（供前端调用）
func (c *Config) GetTheme() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.Theme
}

// SetTheme 设置主题（供前端调用）
func (c *Config) SetTheme(theme string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Theme = theme
	return c.saveToDB()
}

// ResetConfig 重置所有配置为默认值
func (c *Config) ResetConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 重置为默认配置
	c.config = defaultConfig()

	// 持久化到数据库
	return c.saveToDB()
}

// ExportConfig 导出配置为 JSON 字符串
func (c *Config) ExportConfig() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	jsonBytes, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("导出配置失败: %w", err)
	}
	return string(jsonBytes), nil
}

// ImportConfig 从 JSON 字符串导入配置
func (c *Config) ImportConfig(configJSON string) error {
	// 解析导入的 JSON 配置
	var newConfig AppConfig
	if err := json.Unmarshal([]byte(configJSON), &newConfig); err != nil {
		return fmt.Errorf("导入配置解析失败: %w", err)
	}

	// 字段校验
	if err := validateConfig(&newConfig); err != nil {
		return fmt.Errorf("导入配置校验失败: %w", err)
	}

	// 使用写锁保护配置更新
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新内存中的配置
	c.config = &newConfig

	// 持久化到数据库
	return c.saveToDB()
}

// validateConfig 校验配置字段的合法性
func validateConfig(cfg *AppConfig) error {
	// 字体大小校验：范围 10-30
	if cfg.FontSize < 10 || cfg.FontSize > 30 {
		return fmt.Errorf("字体大小必须在 10-30 之间，当前值: %d", cfg.FontSize)
	}

	// 透明度校验：范围 0.1-1.0
	if cfg.Transparency < 0.1 || cfg.Transparency > 1.0 {
		return fmt.Errorf("透明度必须在 0.1-1.0 之间，当前值: %.2f", cfg.Transparency)
	}

	return nil
}

// loadFromDB 从数据库加载配置
func (c *Config) loadFromDB() error {
	row := c.db.DB.QueryRow("SELECT value FROM settings WHERE key = 'app_config'")

	var jsonStr string
	if err := row.Scan(&jsonStr); err != nil {
		// 配置不存在，使用默认值（这不是错误）
		return fmt.Errorf("未找到已保存的配置: %w", err)
	}

	return json.Unmarshal([]byte(jsonStr), c.config)
}

// saveToDB 将配置序列化后存入数据库
func (c *Config) saveToDB() error {
	jsonBytes, err := json.Marshal(c.config)
	if err != nil {
		return err
	}

	// 使用 UPSERT 语法（INSERT OR REPLACE）
	_, err = c.db.DB.Exec(
		"INSERT OR REPLACE INTO settings (key, value, updated_at) VALUES ('app_config', ?, CURRENT_TIMESTAMP)",
		string(jsonBytes),
	)
	return err
}
