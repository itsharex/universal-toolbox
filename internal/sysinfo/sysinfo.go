// Package sysinfo 系统工具模块
// 提供进程管理、端口查看、系统信息、文件批量重命名等系统级工具
package sysinfo

import (
	"fmt"
	"os/exec"
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
