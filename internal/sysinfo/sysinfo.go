// Package sysinfo 系统工具模块
// 提供进程管理、端口查看、系统信息、文件批量重命名等系统级工具
package sysinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// SysInfo 系统工具结构体（Wails 绑定到前端）
type SysInfo struct{}

// ProcessInfo 进程信息结构
type ProcessInfo struct {
	PID    int32   `json:"pid"`    // 进程 ID
	Name   string  `json:"name"`   // 进程名称
	CPU    float64 `json:"cpu"`    // CPU 使用率（%）
	Memory float32 `json:"memory"` // 内存使用率（%）
	Status string  `json:"status"` // 进程状态
}

// PortInfo 端口信息结构
type PortInfo struct {
	Port    uint32 `json:"port"`    // 端口号
	PID     int32  `json:"pid"`     // 占用端口的进程 ID
	Process string `json:"process"` // 进程名称
	Status  string `json:"status"`  // 连接状态
	Type    string `json:"type"`    // TCP/UDP
}

// SystemInfoResult 系统详细信息结构
type SystemInfoResult struct {
	OS          string  `json:"os"`          // 操作系统
	Arch        string  `json:"arch"`        // CPU 架构
	Hostname    string  `json:"hostname"`    // 主机名
	CPUModel    string  `json:"cpuModel"`    // CPU 型号
	CPUCores    int     `json:"cpuCores"`    // CPU 核心数
	CPUUsage    float64 `json:"cpuUsage"`    // CPU 使用率（%）
	MemTotal    uint64  `json:"memTotal"`    // 内存总量（字节）
	MemUsed     uint64  `json:"memUsed"`     // 已用内存（字节）
	MemPercent  float64 `json:"memPercent"`  // 内存使用率（%）
	DiskTotal   uint64  `json:"diskTotal"`   // 磁盘总量（字节）
	DiskUsed    uint64  `json:"diskUsed"`    // 已用磁盘（字节）
	DiskPercent float64 `json:"diskPercent"` // 磁盘使用率（%）
	Uptime      uint64  `json:"uptime"`      // 系统运行时间（秒）
}

// DiskInfo 磁盘信息结构
type DiskInfo struct {
	Device      string  `json:"device"`      // 设备路径
	Mountpoint  string  `json:"mountpoint"`  // 挂载点
	Fstype      string  `json:"fstype"`      // 文件系统类型
	Total       uint64  `json:"total"`       // 总容量（字节）
	Used        uint64  `json:"used"`        // 已用容量（字节）
	Free        uint64  `json:"free"`        // 可用容量（字节）
	UsedPercent float64 `json:"usedPercent"` // 使用率（%）
}

// BatchRenameResult 批量重命名结果
type BatchRenameResult struct {
	Success bool     `json:"success"` // 是否全部成功
	Results []RenameItem `json:"results"` // 每个文件的重命名结果
}

// RenameItem 单个重命名结果
type RenameItem struct {
	OldName string `json:"oldName"` // 原文件名
	NewName string `json:"newName"` // 新文件名
	Success bool   `json:"success"` // 是否成功
	Error   string `json:"error"`   // 错误信息
}

// NewSysInfo 创建系统工具模块实例
func NewSysInfo() *SysInfo {
	return &SysInfo{}
}

// GetSystemInfo 获取系统详细信息
func (s *SysInfo) GetSystemInfo() (*SystemInfoResult, error) {
	result := &SystemInfoResult{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// 获取主机信息
	if hostInfo, err := host.Info(); err == nil {
		result.Hostname = hostInfo.Hostname
		result.Uptime = hostInfo.Uptime
	}

	// 获取 CPU 信息
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		result.CPUModel = cpuInfo[0].ModelName
		result.CPUCores = int(cpuInfo[0].Cores)
	}

	// 获取 CPU 使用率（取 1 秒平均值）
	if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
		result.CPUUsage = cpuPercent[0]
	}

	// 获取内存信息
	if memInfo, err := mem.VirtualMemory(); err == nil {
		result.MemTotal = memInfo.Total
		result.MemUsed = memInfo.Used
		result.MemPercent = memInfo.UsedPercent
	}

	// 获取磁盘信息（根分区或 C 盘）
	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "C:\\"
	}
	if diskInfo, err := disk.Usage(diskPath); err == nil {
		result.DiskTotal = diskInfo.Total
		result.DiskUsed = diskInfo.Used
		result.DiskPercent = diskInfo.UsedPercent
	}

	return result, nil
}

// GetProcessList 获取当前运行的进程列表（按 CPU 使用率降序排序）
func (s *SysInfo) GetProcessList() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败：%v", err)
	}

	var processList []ProcessInfo
	for _, p := range procs {
		info := ProcessInfo{PID: p.Pid}

		// 获取进程名（忽略错误）
		if name, err := p.Name(); err == nil {
			info.Name = name
		} else {
			continue // 跳过无法获取名称的进程
		}

		// 获取 CPU 使用率
		if cpu, err := p.CPUPercent(); err == nil {
			info.CPU = cpu
		}

		// 获取内存使用率
		if mem, err := p.MemoryPercent(); err == nil {
			info.Memory = mem
		}

		// 获取状态
		if status, err := p.Status(); err == nil {
			info.Status = strings.Join(status, ",")
		}

		processList = append(processList, info)
	}

	// 按进程名排序，再按 CPU 使用率降序排序
	sort.Slice(processList, func(i, j int) bool {
		if processList[i].Name != processList[j].Name {
			return processList[i].Name < processList[j].Name
		}
		return processList[i].CPU > processList[j].CPU
	})

	return processList, nil
}

// KillProcess 终止指定 PID 的进程
// 添加安全检查：防止杀 PID 1（init/systemd）和自身进程
func (s *SysInfo) KillProcess(pid int32) error {
	// 安全检查：禁止终止 PID 1（系统初始化进程）
	if pid == 1 {
		return fmt.Errorf("禁止终止 PID 1（系统初始化进程），此操作会导致系统不稳定")
	}

	// 安全检查：禁止终止自身进程
	selfPID := int32(os.Getpid())
	if pid == selfPID {
		return fmt.Errorf("禁止终止自身进程（PID: %d）", selfPID)
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("找不到进程 %d：%v", pid, err)
	}

	// 获取进程名用于日志提示
	procName := ""
	if name, err := p.Name(); err == nil {
		procName = name
	}

	if err := p.Kill(); err != nil {
		return fmt.Errorf("终止进程 %s (PID: %d) 失败：%v", procName, pid, err)
	}

	return nil
}

// ============================================================
// 磁盘清理工具
// ============================================================

// JunkItem 垃圾文件项
type JunkItem struct {
	Path  string `json:"path"`  // 文件路径
	Size  int64  `json:"size"`  // 文件大小（字节）
	Type  string `json:"type"`  // 类型分类（临时文件/缓存文件/日志文件等）
	MTime string `json:"mtime"` // 最后修改时间
}

// ScanDiskJunkResult 磁盘扫描结果
type ScanDiskJunkResult struct {
	Success   bool       `json:"success"`   // 是否成功
	TotalSize int64      `json:"totalSize"` // 总大小（字节）
	Count     int        `json:"count"`     // 文件数量
	Items     []JunkItem `json:"items"`     // 垃圾文件列表
	Error     string     `json:"error"`     // 错误信息
}

// ScanDiskJunk 扫描系统临时文件和缓存文件
// 纯扫描，不自动删除
func (s *SysInfo) ScanDiskJunk() ScanDiskJunkResult {
	result := ScanDiskJunkResult{Success: true}

	// 根据操作系统确定扫描路径
	var scanPaths []struct {
		path    string
		junkType string
		pattern string // glob 匹配模式
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}

	switch runtime.GOOS {
	case "windows":
		// Windows 系统临时目录
		windowsTemp := os.Getenv("TEMP")
		if windowsTemp == "" {
			windowsTemp = os.Getenv("TMP")
		}
		if windowsTemp != "" {
			scanPaths = append(scanPaths,
				struct {
					path     string
					junkType string
					pattern string
				}{windowsTemp, "系统临时文件", "*"},
			)
		}
		// Windows 预读取缓存
		scanPaths = append(scanPaths,
			struct {
				path     string
				junkType string
				pattern string
			}{"C:\\Windows\\Prefetch", "预读取缓存", "*.pf"},
		)
		// 用户缓存目录
		if homeDir != "" {
			scanPaths = append(scanPaths,
				struct {
					path     string
					junkType string
					pattern string
				}{filepath.Join(homeDir, "AppData", "Local", "Temp"), "用户临时文件", "*"},
			)
		}

	case "darwin":
		// macOS 系统临时目录
		scanPaths = append(scanPaths,
			struct {
				path     string
				junkType string
				pattern string
			}{"/tmp", "系统临时文件", "*"},
			struct {
				path     string
				junkType string
				pattern string
			}{"/private/tmp", "系统临时文件", "*"},
		)
		// macOS 用户缓存
		if homeDir != "" {
			scanPaths = append(scanPaths,
				struct {
					path     string
					junkType string
					pattern string
				}{filepath.Join(homeDir, "Library", "Caches"), "用户缓存文件", "*"},
			)
		}

	default:
		// Linux / 其他 Unix 系统
		scanPaths = append(scanPaths,
			struct {
				path     string
				junkType string
				pattern string
			}{"/tmp", "系统临时文件", "*"},
			struct {
				path     string
				junkType string
				pattern string
			}{"/var/tmp", "系统临时文件", "*"},
		)
		// 用户缓存
		if homeDir != "" {
			scanPaths = append(scanPaths,
				struct {
					path     string
					junkType string
					pattern string
				}{filepath.Join(homeDir, ".cache"), "用户缓存文件", "*"},
			)
		}
		// 包管理器缓存
		scanPaths = append(scanPaths,
			struct {
				path     string
				junkType string
				pattern string
			}{"/var/cache/apt/archives", "APT 缓存", "*.deb"},
		)
	}

	// 遍历扫描路径
	seen := make(map[string]bool) // 去重
	for _, sp := range scanPaths {
		matches, err := filepath.Glob(filepath.Join(sp.path, sp.pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			if seen[match] {
				continue
			}
			seen[match] = true

			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			// 跳过目录本身
			if info.IsDir() {
				continue
			}

			// 跳过最近 1 小时内修改的文件（避免误删正在使用的临时文件）
			if time.Since(info.ModTime()) < time.Hour {
				continue
			}

			result.Items = append(result.Items, JunkItem{
				Path:  match,
				Size:  info.Size(),
				Type:  sp.junkType,
				MTime: info.ModTime().Format("2006-01-02 15:04:05"),
			})
			result.TotalSize += info.Size()
		}
	}

	result.Count = len(result.Items)

	// 按大小降序排序
	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].Size > result.Items[j].Size
	})

	// 限制返回数量（最多 500 条）
	if len(result.Items) > 500 {
		result.Items = result.Items[:500]
	}

	return result
}

// CleanDiskItemsResult 清理结果
type CleanDiskItemsResult struct {
	Success   bool   `json:"success"`   // 是否成功
	TotalSize int64  `json:"totalSize"` // 清理的总大小（字节）
	Count     int    `json:"count"`     // 清理的文件数量
	Message   string `json:"message"`   // 结果描述
	Error     string `json:"error"`     // 错误信息
}

// CleanDiskItems 根据传入的文件路径列表执行清理
// itemsJSON: JSON 格式的文件路径数组，如 ["/tmp/file1", "/tmp/file2"]
func (s *SysInfo) CleanDiskItems(itemsJSON string) CleanDiskItemsResult {
	result := CleanDiskItemsResult{Success: true}

	if itemsJSON == "" {
		result.Error = "文件路径列表不能为空"
		result.Success = false
		return result
	}

	// 解析 JSON 路径列表
	var paths []string
	if err := json.Unmarshal([]byte(itemsJSON), &paths); err != nil {
		result.Error = "JSON 格式错误：" + err.Error()
		result.Success = false
		return result
	}

	if len(paths) == 0 {
		result.Error = "文件路径列表为空"
		result.Success = false
		return result
	}

	// 逐个删除文件
	var deletedCount int
	var totalSize int64
	var errors []string

	for _, filePath := range paths {
		// 安全检查：防止路径遍历攻击
		cleanPath := filepath.Clean(filePath)

		// 获取文件信息
		info, err := os.Stat(cleanPath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("无法访问 %s: %v", cleanPath, err))
			continue
		}

		// 记录大小
		fileSize := info.Size()

		// 执行删除
		if err := os.Remove(cleanPath); err != nil {
			errors = append(errors, fmt.Sprintf("删除失败 %s: %v", cleanPath, err))
			continue
		}

		deletedCount++
		totalSize += fileSize
	}

	result.Count = deletedCount
	result.TotalSize = totalSize

	// 构建结果消息
	if deletedCount > 0 {
		result.Message = fmt.Sprintf("成功清理 %d 个文件，释放 %.2f MB 空间",
			deletedCount, float64(totalSize)/(1024*1024))
	}

	if len(errors) > 0 {
		if result.Message != "" {
			result.Message += "\n"
		}
		result.Message += fmt.Sprintf("其中 %d 个文件清理失败", len(errors))
	}

	return result
}

// GetPortList 获取端口占用信息（添加去重逻辑）
func (s *SysInfo) GetPortList() ([]PortInfo, error) {
	connections, err := psnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取端口信息失败：%v", err)
	}

	// 使用 map 进行去重（相同端口+PID+类型只保留一条）
	seen := make(map[string]bool)
	var portList []PortInfo

	for _, conn := range connections {
		// 只显示监听状态和已建立连接的端口
		if conn.Laddr.Port == 0 {
			continue
		}

		// conn.Type 在 gopsutil v3.24+ 中为 uint32，需要手动转换为可读字符串
		typeStr := map[uint32]string{
			1: "TCP", 2: "UDP", 3: "TCPv6", 4: "UDPv6",
		}
		connTypeName, ok := typeStr[conn.Type]
		if !ok {
			connTypeName = fmt.Sprintf("UNKNOWN(%d)", conn.Type)
		}

		// 去重键：端口 + 类型
		dedupeKey := fmt.Sprintf("%d-%s", conn.Laddr.Port, connTypeName)
		if seen[dedupeKey] {
			continue
		}
		seen[dedupeKey] = true

		info := PortInfo{
			Port:   conn.Laddr.Port,
			PID:    conn.Pid,
			Status: conn.Status,
			Type:   connTypeName,
		}

		// 获取进程名
		if conn.Pid > 0 {
			if p, err := process.NewProcess(conn.Pid); err == nil {
				if name, err := p.Name(); err == nil {
					info.Process = name
				}
			}
		}

		portList = append(portList, info)
	}

	// 按端口号排序
	sort.Slice(portList, func(i, j int) bool {
		return portList[i].Port < portList[j].Port
	})

	return portList, nil
}

// ReleasePort 释放指定端口（通过终止占用进程）
func (s *SysInfo) ReleasePort(port uint32) error {
	connections, err := psnet.Connections("all")
	if err != nil {
		return fmt.Errorf("获取端口信息失败：%v", err)
	}

	for _, conn := range connections {
		if conn.Laddr.Port == port && conn.Pid > 0 {
			p, err := process.NewProcess(conn.Pid)
			if err != nil {
				continue
			}
			if err := p.Kill(); err != nil {
				return fmt.Errorf("释放端口 %d 失败：%v", port, err)
			}
			return nil
		}
	}

	return fmt.Errorf("端口 %d 未被占用", port)
}

// GetIPAddresses 获取本机 IP 地址列表
func (s *SysInfo) GetIPAddresses() ([]string, error) {
	interfaces, err := psnet.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, iface := range interfaces {
		for _, addr := range iface.Addrs {
			if addr.Addr != "" {
				ips = append(ips, fmt.Sprintf("%s (%s)", addr.Addr, iface.Name))
			}
		}
	}

	return ips, nil
}

// OpenFileManager 在系统文件管理器中打开指定路径
// 添加路径校验防止命令注入
func (s *SysInfo) OpenFileManager(path string) error {
	// 路径安全校验：防止命令注入
	if path == "" {
		return fmt.Errorf("路径不能为空")
	}

	// 检查路径中是否包含危险字符（防止命令注入）
	dangerousChars := []string{"|", "&", ";", "$", "`", "(", ")", "<", ">", "!", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("路径包含非法字符: %q，可能存在命令注入风险", char)
		}
	}

	// 清理路径，防止路径遍历
	cleanPath := filepath.Clean(path)

	// 验证路径是否存在（仅验证目录）
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在: %s", cleanPath)
		}
		return fmt.Errorf("访问路径失败: %v", err)
	}
	if !info.IsDir() {
		// 如果是文件，打开其所在目录
		cleanPath = filepath.Dir(cleanPath)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", cleanPath)
	case "darwin":
		cmd = exec.Command("open", cleanPath)
	default:
		cmd = exec.Command("xdg-open", cleanPath)
	}
	return cmd.Start()
}

// BatchRename 批量文件重命名
// 支持模式：prefix（添加前缀）、suffix（添加后缀）、sequence（序号重命名）、regex（正则替换）
// directory: 目标目录
// pattern: 文件匹配模式（如 "*.txt"）
// mode: 重命名模式
// value: 模式参数（前缀/后缀文本、正则表达式、序号起始值等）
// replace: 替换文本（仅 regex 模式使用）
func (s *SysInfo) BatchRename(directory, filePattern, mode, value, replace string) BatchRenameResult {
	result := BatchRenameResult{
		Success: true,
	}

	// 校验目录是否存在
	info, err := os.Stat(directory)
	if err != nil || !info.IsDir() {
		result.Success = false
		result.Results = []RenameItem{{
			OldName: directory,
			Error:   "目录不存在或无法访问",
		}}
		return result
	}

	// 匹配文件
	matches, err := filepath.Glob(filepath.Join(directory, filePattern))
	if err != nil {
		result.Success = false
		result.Results = []RenameItem{{
			OldName: filePattern,
			Error:   "文件匹配模式错误：" + err.Error(),
		}}
		return result
	}

	if len(matches) == 0 {
		result.Success = false
		result.Results = []RenameItem{{
			OldName: filePattern,
			Error:   "未找到匹配的文件",
		}}
		return result
	}

	result.Results = make([]RenameItem, 0, len(matches))

	switch mode {
	case "prefix":
		// 添加前缀模式
		for _, match := range matches {
			baseName := filepath.Base(match)
			newName := value + baseName
			newPath := filepath.Join(directory, newName)
			item := RenameItem{OldName: baseName, NewName: newName}

			if err := os.Rename(match, newPath); err != nil {
				item.Error = err.Error()
				result.Success = false
			} else {
				item.Success = true
			}
			result.Results = append(result.Results, item)
		}

	case "suffix":
		// 添加后缀模式（在扩展名之前插入）
		for _, match := range matches {
			baseName := filepath.Base(match)
			ext := filepath.Ext(baseName)
			nameWithoutExt := strings.TrimSuffix(baseName, ext)
			newName := nameWithoutExt + value + ext
			newPath := filepath.Join(directory, newName)
			item := RenameItem{OldName: baseName, NewName: newName}

			if err := os.Rename(match, newPath); err != nil {
				item.Error = err.Error()
				result.Success = false
			} else {
				item.Success = true
			}
			result.Results = append(result.Results, item)
		}

	case "sequence":
		// 序号重命名模式
		startNum := 1
		if value != "" {
			if n, err := strconv.Atoi(value); err == nil {
				startNum = n
			}
		}
		// 获取扩展名
		ext := filepath.Ext(matches[0])
		for i, match := range matches {
			baseName := filepath.Base(match)
			fileExt := filepath.Ext(baseName)
			if fileExt != "" {
				ext = fileExt
			}
			newName := fmt.Sprintf("%s%04d%s", replace, startNum+i, ext)
			newPath := filepath.Join(directory, newName)
			item := RenameItem{OldName: baseName, NewName: newName}

			if err := os.Rename(match, newPath); err != nil {
				item.Error = err.Error()
				result.Success = false
			} else {
				item.Success = true
			}
			result.Results = append(result.Results, item)
		}

	case "regex":
		// 正则替换模式
		re, err := regexp.Compile(value)
		if err != nil {
			result.Success = false
			result.Results = []RenameItem{{
				Error: "正则表达式错误：" + err.Error(),
			}}
			return result
		}
		for _, match := range matches {
			baseName := filepath.Base(match)
			newName := re.ReplaceAllString(baseName, replace)
			if newName == baseName {
				// 没有变化，跳过
				result.Results = append(result.Results, RenameItem{
					OldName: baseName, NewName: newName, Success: true,
				})
				continue
			}
			newPath := filepath.Join(directory, newName)
			item := RenameItem{OldName: baseName, NewName: newName}

			if err := os.Rename(match, newPath); err != nil {
				item.Error = err.Error()
				result.Success = false
			} else {
				item.Success = true
			}
			result.Results = append(result.Results, item)
		}

	default:
		result.Success = false
		result.Results = []RenameItem{{
			Error: fmt.Sprintf("不支持的重命名模式: %s（支持: prefix/suffix/sequence/regex）", mode),
		}}
	}

	return result
}

// GetDiskList 获取所有磁盘分区信息
func (s *SysInfo) GetDiskList() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(true) // true 表示只显示物理设备
	if err != nil {
		return nil, fmt.Errorf("获取磁盘分区列表失败：%v", err)
	}

	var diskList []DiskInfo
	for _, part := range partitions {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			// 某些分区可能无法访问（如权限不足），跳过
			continue
		}

		diskList = append(diskList, DiskInfo{
			Device:      part.Device,
			Mountpoint:  part.Mountpoint,
			Fstype:      part.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	return diskList, nil
}

// FileOperation 执行文件操作（复制/移动/删除）
// operation: copy/move/delete
// src: 源文件或目录路径
// dst: 目标路径（复制和移动时使用，删除时忽略）
func (s *SysInfo) FileOperation(operation, src, dst string) error {
	// 路径安全校验
	if src == "" {
		return fmt.Errorf("源路径不能为空")
	}

	// 清理路径
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	switch operation {
	case "copy":
		if dst == "" {
			return fmt.Errorf("复制操作需要指定目标路径")
		}
		// 检查源是否存在
		srcInfo, err := os.Stat(src)
		if err != nil {
			return fmt.Errorf("源路径不存在: %v", err)
		}
		if srcInfo.IsDir() {
			// 目录复制
			return copyDir(src, dst)
		}
		// 文件复制
		return copyFile(src, dst)

	case "move":
		if dst == "" {
			return fmt.Errorf("移动操作需要指定目标路径")
		}
		return os.Rename(src, dst)

	case "delete":
		// 检查路径是否存在
		info, err := os.Stat(src)
		if err != nil {
			return fmt.Errorf("路径不存在: %v", err)
		}
		if info.IsDir() {
			return os.RemoveAll(src)
		}
		return os.Remove(src)

	default:
		return fmt.Errorf("不支持的操作类型: %s（支持: copy/move/delete）", operation)
	}
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	// 读取源文件
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取源文件失败: %v", err)
	}

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 获取源文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %v", err)
	}

	// 写入目标文件
	if err := os.WriteFile(dst, data, srcInfo.Mode()); err != nil {
		return fmt.Errorf("写入目标文件失败: %v", err)
	}

	return nil
}

// copyDir 递归复制目录
func copyDir(src, dst string) error {
	// 创建目标目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
