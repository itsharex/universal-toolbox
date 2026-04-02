// Package sysinfo 系统工具模块
// 提供进程管理、端口查看、系统信息、文件批量重命名等系统级工具
package sysinfo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
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

// GetProcessList 获取当前运行的进程列表
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

	return processList, nil
}

// KillProcess 终止指定 PID 的进程
func (s *SysInfo) KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("找不到进程 %d：%v", pid, err)
	}

	if err := p.Kill(); err != nil {
		return fmt.Errorf("终止进程失败：%v", err)
	}

	return nil
}

// GetPortList 获取端口占用信息
func (s *SysInfo) GetPortList() ([]PortInfo, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("获取端口信息失败：%v", err)
	}

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

	return portList, nil
}

// ReleasePort 释放指定端口（通过终止占用进程）
func (s *SysInfo) ReleasePort(port uint32) error {
	connections, err := net.Connections("all")
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
	interfaces, err := net.Interfaces()
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
func (s *SysInfo) OpenFileManager(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// FileItem 文件信息结构（用于文件列表展示）
type FileItem struct {
	Name    string `json:"name"`    // 文件名（含扩展名）
	Path    string `json:"path"`    // 完整路径
	Size    int64  `json:"size"`    // 文件大小（字节）
	IsDir   bool   `json:"isDir"`   // 是否为目录
	ModTime int64  `json:"modTime"` // 修改时间戳
	Ext     string `json:"ext"`     // 扩展名
}

// RenameResult 重命名结果结构
type RenameResult struct {
	OldPath string `json:"oldPath"` // 原路径
	NewPath string `json:"newPath"` // 新路径
	Success bool   `json:"success"` // 是否成功
	Error   string `json:"error"`   // 错误信息
}

// ListDirectory 列出指定目录下的所有文件（不含子目录）
func (s *SysInfo) ListDirectory(dirPath string) ([]FileItem, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败：%v", err)
	}

	var files []FileItem
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // 跳过无法获取信息的文件
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		files = append(files, FileItem{
			Name:    entry.Name(),
			Path:    fullPath,
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime().Unix(),
			Ext:     filepath.Ext(entry.Name()),
		})
	}
	return files, nil
}

// BatchRenameFiles 批量重命名文件
// oldNewPairs: 形如 "原路径|新路径;原路径|新路径" 的字符串数组
func (s *SysInfo) BatchRenameFiles(pairs []string) ([]RenameResult, error) {
	var results []RenameResult

	for _, pair := range pairs {
		// 分割原路径和新路径
		parts := strings.Split(pair, "|")
		if len(parts) != 2 {
			results = append(results, RenameResult{
				OldPath: pair,
				NewPath: "",
				Success: false,
				Error:   "路径格式错误",
			})
			continue
		}

		oldPath := strings.TrimSpace(parts[0])
		newPath := strings.TrimSpace(parts[1])

		// 检查原文件是否存在
		if _, err := os.Stat(oldPath); err != nil {
			results = append(results, RenameResult{
				OldPath: oldPath,
				NewPath: newPath,
				Success: false,
				Error:   fmt.Sprintf("文件不存在或无权限：%v", err),
			})
			continue
		}

		// 检查目标文件是否已存在
		if _, err := os.Stat(newPath); err == nil {
			results = append(results, RenameResult{
				OldPath: oldPath,
				NewPath: newPath,
				Success: false,
				Error:   "目标文件已存在",
			})
			continue
		}

		// 执行重命名
		if err := os.Rename(oldPath, newPath); err != nil {
			results = append(results, RenameResult{
				OldPath: oldPath,
				NewPath: newPath,
				Success: false,
				Error:   fmt.Sprintf("重命名失败：%v", err),
			})
			continue
		}

		results = append(results, RenameResult{
			OldPath: oldPath,
			NewPath: newPath,
			Success: true,
			Error:   "",
		})
	}

	return results, nil
}

// BatchCopyFiles 批量复制文件（格式：源路径|目标路径）
func (s *SysInfo) BatchCopyFiles(pairs []string) error {
	for _, pair := range pairs {
		parts := strings.Split(pair, "|")
		if len(parts) != 2 {
			return fmt.Errorf("路径格式错误：%s", pair)
		}
		src, dst := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("读取文件失败 %s：%v", src, err)
		}

		if err := os.WriteFile(dst, data, 0644); err != nil {
			return fmt.Errorf("写入文件失败 %s：%v", dst, err)
		}
	}
	return nil
}

// BatchMoveFiles 批量移动文件（本质是重命名到新目录）
func (s *SysInfo) BatchMoveFiles(pairs []string) ([]RenameResult, error) {
	return s.BatchRenameFiles(pairs)
}

// BatchDeleteFiles 批量删除文件
func (s *SysInfo) BatchDeleteFiles(paths []string) error {
	for _, path := range paths {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("删除文件失败 %s：%v", path, err)
		}
	}
	return nil
}

// CopyFile 复制单个文件
func (s *SysInfo) CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取文件失败：%v", err)
	}
	return os.WriteFile(dst, data, 0644)
}

// MoveFile 移动单个文件（重命名）
func (s *SysInfo) MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// DeleteFile 删除单个文件
func (s *SysInfo) DeleteFile(path string) error {
	return os.Remove(path)
}
