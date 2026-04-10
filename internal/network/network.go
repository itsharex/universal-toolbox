// Package network 网络工具模块
// 提供 Ping 测试、IP 扫描、HTTP 接口测试、Hosts 编辑等网络工具
package network

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// NetworkTools 网络工具结构体（Wails 绑定到前端）
type NetworkTools struct {
	httpClient *http.Client // HTTP 客户端（可复用）
}

// PingResult Ping 测试结果
type PingResult struct {
	Host      string  `json:"host"`      // 目标主机
	Alive     bool    `json:"alive"`     // 是否在线
	LatencyMs float64 `json:"latencyMs"` // 延迟（毫秒）
	Method    string  `json:"method"`    // 检测方式（icmp/tcp/dns）
	Error     string  `json:"error"`     // 错误信息
}

// ScanResult 内网扫描结果
type ScanResult struct {
	IP      string `json:"ip"`      // IP 地址
	Online  bool   `json:"online"`  // 是否在线
	Latency int64  `json:"latency"` // 延迟（毫秒）
	Ports   string `json:"ports"`   // 开放的端口列表
}

// HTTPTestResult HTTP 接口测试结果
type HTTPTestResult struct {
	StatusCode  int               `json:"statusCode"`  // HTTP 状态码
	Status      string            `json:"status"`      // 状态描述
	LatencyMs   float64           `json:"latencyMs"`   // 请求耗时（毫秒）
	Body        string            `json:"body"`        // 响应体
	Headers     map[string]string `json:"headers"`     // 响应头
	ContentType string            `json:"contentType"` // 内容类型
	Error       string            `json:"error"`       // 错误信息
}

// DNSRecord DNS 记录
type DNSRecord struct {
	Name  string `json:"name"`  // 记录名称
	Type  string `json:"type"`  // 记录类型
	Value string `json:"value"` // 记录值
	TTL   uint32 `json:"ttl"`   // TTL（秒）
}

// DNSQueryResult DNS 查询结果
type DNSQueryResult struct {
	Domain  string      `json:"domain"`  // 查询域名
	Type    string      `json:"type"`    // 查询类型
	Records []DNSRecord `json:"records"` // 查询结果记录
	Error   string      `json:"error"`   // 错误信息
}

// NewNetworkTools 创建网络工具模块实例
func NewNetworkTools() *NetworkTools {
	return &NetworkTools{
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // 30 秒超时
		},
	}
}

// ============================================================
// Ping 工具
// ============================================================

// PingHost 对指定主机进行 Ping 测试
// 先尝试 ICMP（如果可用），再尝试 TCP 连接
func (n *NetworkTools) PingHost(host string) PingResult {
	result := PingResult{Host: host}

	// 方法 1：尝试 ICMP Ping（需要 root/管理员权限）
	icmpResult := n.tryICMPPing(host)
	if icmpResult.Alive {
		return icmpResult
	}

	// 方法 2：TCP 连接检测（不需要特殊权限）
	tcpResult := n.tryTCPPing(host)
	if tcpResult.Alive {
		return tcpResult
	}

	// 方法 3：DNS 解析检测（最后的手段）
	_, dnsErr := net.LookupHost(host)
	if dnsErr != nil {
		result.Alive = false
		result.Error = "主机不可达：" + icmpResult.Error
		return result
	}

	// DNS 可以解析但 TCP 不通
	result.Alive = true
	result.Method = "dns"
	result.Error = "主机存在但端口未开放（DNS 可解析）"

	return result
}

// tryICMPPing 尝试 ICMP Ping
func (n *NetworkTools) tryICMPPing(host string) PingResult {
	result := PingResult{Host: host, Method: "icmp"}

	// 解析域名
	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		result.Error = "DNS 解析失败"
		return result
	}

	// 构建 ICMP 连接
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		// ICMP 不可用（可能没有权限）
		result.Error = "ICMP 不可用（需要管理员权限）"
		return result
	}
	defer conn.Close()

	// 设置超时
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	// 构建 ICMP Echo Request
	var msg [32]byte
	msg[0] = 8  // Echo Request 类型
	msg[1] = 0  // Code
	msg[2] = 0  // Checksum（先置零）
	msg[3] = 0  // Checksum
	msg[4] = 0  // Identifier
	msg[5] = 1  // Sequence number
	// 计算校验和
	checksum := 0
	for i := 0; i < 32; i += 2 {
		checksum += int(msg[i])<<8 | int(msg[i+1])
	}
	checksum = (checksum >> 16) + (checksum & 0xffff)
	checksum = ^checksum
	msg[2] = byte(checksum >> 8)
	msg[3] = byte(checksum)

	start := time.Now()
	_, err = conn.Write(msg[:])
	if err != nil {
		result.Error = "ICMP 发送失败"
		return result
	}

	// 读取响应
	resp := make([]byte, 1024)
	_, err = conn.Read(resp)
	if err != nil {
		result.Error = "ICMP 响应超时"
		return result
	}

	result.Alive = true
	result.LatencyMs = float64(time.Since(start).Nanoseconds()) / 1e6
	return result
}

// tryTCPPing 尝试 TCP 连接检测
func (n *NetworkTools) tryTCPPing(host string) PingResult {
	result := PingResult{Host: host, Method: "tcp"}

	// 依次尝试常见端口
	ports := []string{"80", "443", "22", "3389", "8080"}
	start := time.Now()

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 2*time.Second)
		if err == nil {
			conn.Close()
			result.Alive = true
			result.LatencyMs = float64(time.Since(start).Nanoseconds()) / 1e6
			return result
		}
	}

	result.Error = "TCP 连接超时（已尝试端口: 80, 443, 22, 3389, 8080）"
	return result
}

// PingMultiple 批量 Ping 多个主机
func (n *NetworkTools) PingMultiple(hosts []string) []PingResult {
	results := make([]PingResult, len(hosts))
	var wg sync.WaitGroup

	for i, host := range hosts {
		wg.Add(1)
		go func(idx int, h string) {
			defer wg.Done()
			results[idx] = n.PingHost(h)
		}(i, host)
	}

	wg.Wait()
	return results
}

// ============================================================
// 内网扫描工具
// ============================================================

// ScanLAN 扫描内网 IP 段，返回在线设备列表
// subnet: 如 "192.168.1"，将扫描 192.168.1.1-254
// 增加更多端口扫描：80, 443, 22, 3389, 8080
func (n *NetworkTools) ScanLAN(subnet string) []ScanResult {
	var results []ScanResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 并发扫描，限制并发数为 50
	semaphore := make(chan struct{}, 50)

	// 扫描端口列表
	scanPorts := []uint32{80, 443, 22, 3389, 8080}

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s.%d", subnet, i)
		wg.Add(1)
		semaphore <- struct{}{}

		go func(targetIP string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			start := time.Now()
			var openPorts []string

			// 尝试多个端口
			for _, port := range scanPorts {
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetIP, port), 500*time.Millisecond)
				if err == nil {
					conn.Close()
					openPorts = append(openPorts, fmt.Sprintf("%d", port))
				}
			}

			if len(openPorts) > 0 {
				mu.Lock()
				results = append(results, ScanResult{
					IP:      targetIP,
					Online:  true,
					Latency: time.Since(start).Milliseconds(),
					Ports:   strings.Join(openPorts, ","),
				})
				mu.Unlock()
			}
		}(ip)
	}

	wg.Wait()
	return results
}

// GetLocalSubnet 获取本机所在子网（用于扫描建议）
func (n *NetworkTools) GetLocalSubnet() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				// 返回前三段（如 192.168.1）
				parts := strings.Split(ip4.String(), ".")
				if len(parts) >= 3 {
					return strings.Join(parts[:3], "."), nil
				}
			}
		}
	}

	return "192.168.1", nil // 默认值
}

// ============================================================
// HTTP 接口测试工具
// ============================================================

// HTTPRequest 发送 HTTP 请求并返回结果
// method: GET/POST/PUT/DELETE/PATCH
// Content-Type 逻辑：先设置自定义请求头，再设置默认值（自定义优先）
func (n *NetworkTools) HTTPRequest(method, url, body string, headers map[string]string) HTTPTestResult {
	result := HTTPTestResult{}

	// 创建请求
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		result.Error = "请求创建失败：" + err.Error()
		return result
	}

	// 先设置自定义请求头（自定义优先）
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 再设置默认请求头（仅当自定义请求头未设置时才生效）
	req.Header.Set("User-Agent", "XTool/1.0")
	if body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求并计时
	start := time.Now()
	resp, err := n.httpClient.Do(req)
	if err != nil {
		result.Error = "请求失败：" + err.Error()
		return result
	}
	defer resp.Body.Close()

	result.LatencyMs = float64(time.Since(start).Nanoseconds()) / 1e6
	result.StatusCode = resp.StatusCode
	result.Status = resp.Status
	result.ContentType = resp.Header.Get("Content-Type")

	// 读取响应体（限制最多 1MB）
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		result.Error = "读取响应失败：" + err.Error()
		return result
	}
	result.Body = string(respBody)

	// 收集响应头
	result.Headers = make(map[string]string)
	for k, v := range resp.Header {
		result.Headers[k] = strings.Join(v, ", ")
	}

	return result
}

// HttpTestWithHeaders 支持完整请求头编辑的 HTTP 测试
// headersJSON: JSON 格式的请求头，如 {"Content-Type":"application/json","Authorization":"Bearer xxx"}
// 返回完整的请求和响应信息
func (n *NetworkTools) HttpTestWithHeaders(method, url, body, headersJSON string) HTTPTestResult {
	result := HTTPTestResult{}

	// 解析请求头 JSON
	var headers map[string]string
	if headersJSON != "" {
		if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
			result.Error = "请求头 JSON 解析失败：" + err.Error()
			return result
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		result.Error = "请求创建失败：" + err.Error()
		return result
	}

	// 设置自定义请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 设置默认 User-Agent（仅当未自定义时）
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "XTool/1.0")
	}

	// 发送请求并计时
	start := time.Now()
	resp, err := n.httpClient.Do(req)
	if err != nil {
		result.Error = "请求失败：" + err.Error()
		return result
	}
	defer resp.Body.Close()

	result.LatencyMs = float64(time.Since(start).Nanoseconds()) / 1e6
	result.StatusCode = resp.StatusCode
	result.Status = resp.Status
	result.ContentType = resp.Header.Get("Content-Type")

	// 读取响应体（限制最多 2MB）
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		result.Error = "读取响应失败：" + err.Error()
		return result
	}
	result.Body = string(respBody)

	// 收集响应头
	result.Headers = make(map[string]string)
	for k, v := range resp.Header {
		result.Headers[k] = strings.Join(v, ", ")
	}

	return result
}

// ============================================================
// DNS 查询工具
// ============================================================

// DNSQuery 查询 DNS 记录
// domain: 要查询的域名
// recordType: 记录类型（A/AAAA/MX/CNAME/NS/TXT）
func (n *NetworkTools) DNSQuery(domain, recordType string) DNSQueryResult {
	result := DNSQueryResult{
		Domain: domain,
		Type:   strings.ToUpper(recordType),
	}

	// 将记录类型转换为 DNS 类型
	var qtype uint16
	switch strings.ToUpper(recordType) {
	case "A":
		qtype = dns.TypeA
	case "AAAA":
		qtype = dns.TypeAAAA
	case "MX":
		qtype = dns.TypeMX
	case "CNAME":
		qtype = dns.TypeCNAME
	case "NS":
		qtype = dns.TypeNS
	case "TXT":
		qtype = dns.TypeTXT
	default:
		result.Error = fmt.Sprintf("不支持的记录类型: %s（支持: A/AAAA/MX/CNAME/NS/TXT）", recordType)
		return result
	}

	// 创建 DNS 客户端
	client := &dns.Client{
		Timeout: 5 * time.Second,
	}

	// 构建 DNS 查询消息
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), qtype)
	msg.RecursionDesired = true

	// 使用公共 DNS 服务器查询
	dnsServers := []string{"8.8.8.8:53", "1.1.1.1:53"}
	var r *dns.Msg
	var err error

	for _, server := range dnsServers {
		r, _, err = client.Exchange(msg, server)
		if err == nil {
			break
		}
	}

	if err != nil {
		result.Error = fmt.Sprintf("DNS 查询失败: %v", err)
		return result
	}

	// 解析响应
	result.Records = make([]DNSRecord, 0)
	for _, ans := range r.Answer {
		record := DNSRecord{
			Name: ans.Header().Name,
			Type: dns.TypeToString[ans.Header().Rrtype],
			TTL:  ans.Header().Ttl,
		}

		switch v := ans.(type) {
		case *dns.A:
			record.Value = v.A.String()
		case *dns.AAAA:
			record.Value = v.AAAA.String()
		case *dns.MX:
			record.Value = fmt.Sprintf("%s (优先级: %d)", v.Mx, v.Preference)
		case *dns.CNAME:
			record.Value = v.Target
		case *dns.NS:
			record.Value = v.Ns
		case *dns.TXT:
			record.Value = strings.Join(v.Txt, " ")
		default:
			record.Value = ans.String()
		}

		result.Records = append(result.Records, record)
	}

	if len(result.Records) == 0 {
		result.Error = "未找到匹配的 DNS 记录"
	}

	return result
}

// ============================================================
// Hosts 文件编辑工具
// ============================================================

// GetHostsContent 读取 hosts 文件内容
func (n *NetworkTools) GetHostsContent() (string, error) {
	hostsPath := n.getHostsPath()
	content, err := os.ReadFile(hostsPath)
	if err != nil {
		return "", fmt.Errorf("读取 hosts 文件失败：%v", err)
	}
	return string(content), nil
}

// SaveHostsContent 保存 hosts 文件内容（需要管理员权限）
// 修复备份逻辑：检查备份操作中的错误
func (n *NetworkTools) SaveHostsContent(content string) error {
	hostsPath := n.getHostsPath()

	// 创建备份（检查错误）
	backupPath := hostsPath + ".bak"
	original, err := os.ReadFile(hostsPath)
	if err != nil {
		return fmt.Errorf("读取原始 hosts 文件失败（无法创建备份）：%v", err)
	}
	if err := os.WriteFile(backupPath, original, 0644); err != nil {
		return fmt.Errorf("创建 hosts 备份文件失败：%v", err)
	}

	// 写入新内容
	if err := os.WriteFile(hostsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("保存 hosts 文件失败（可能需要管理员权限）：%v", err)
	}
	return nil
}

// getHostsPath 获取当前系统的 hosts 文件路径
func (n *NetworkTools) getHostsPath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

// ============================================================
// WHOIS 查询工具
// ============================================================

// WhoisResult WHOIS 查询结果
type WhoisResult struct {
	Success    bool   `json:"success"`    // 是否成功
	Domain     string `json:"domain"`     // 查询域名
	Registrar  string `json:"registrar"`  // 注册商
	CreateTime string `json:"createTime"` // 注册时间
	ExpireTime string `json:"expireTime"` // 到期时间
	UpdateTime string `json:"updateTime"` // 更新时间
	NameServer string `json:"nameServer"` // DNS 服务器
	DomainStatus string `json:"domainStatus"` // 域名状态
	Registrant string `json:"registrant"` // 注册人
	RawData    string `json:"rawData"`    // 原始 WHOIS 数据
	Error      string `json:"error"`      // 错误信息
}

// WhoisQuery 查询域名的 WHOIS 注册信息
// 使用 net.Dial 连接 WHOIS 服务器
func (n *NetworkTools) WhoisQuery(domain string) WhoisResult {
	result := WhoisResult{Domain: domain}

	// 去掉可能的协议前缀和路径
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.SplitN(domain, "/", 2)[0]
	domain = strings.TrimSpace(domain)

	if domain == "" {
		result.Error = "域名不能为空"
		return result
	}

	// 根据顶级域名选择 WHOIS 服务器
	whoisServer := getWhoisServer(domain)

	// 连接 WHOIS 服务器
	conn, err := net.DialTimeout("tcp", whoisServer, 10*time.Second)
	if err != nil {
		result.Error = fmt.Sprintf("连接 WHOIS 服务器失败: %v", err)
		return result
	}
	defer conn.Close()

	// 设置读写超时
	conn.SetDeadline(time.Now().Add(15 * time.Second))

	// 发送查询请求
	_, err = fmt.Fprintf(conn, "%s\r\n", domain)
	if err != nil {
		result.Error = fmt.Sprintf("发送 WHOIS 查询失败: %v", err)
		return result
	}

	// 读取响应
	var rawData strings.Builder
	scanner := bufio.NewScanner(conn)
	// 增大缓冲区以处理长行
	scanner.Buffer(make([]byte, 1024*64), 1024*64)
	for scanner.Scan() {
		rawData.WriteString(scanner.Text())
		rawData.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		result.Error = fmt.Sprintf("读取 WHOIS 响应失败: %v", err)
		return result
	}

	raw := rawData.String()
	result.RawData = raw
	result.Success = true

	// 解析关键字段
	result.Registrar = extractWhoisField(raw, []string{"Registrar:", "registrar:"})
	result.CreateTime = extractWhoisField(raw, []string{"Creation Date:", "Created On:", "Registry Creation Date:", "created:", "Registration Time:"})
	result.ExpireTime = extractWhoisField(raw, []string{"Registry Expiry Date:", "Expiry Date:", "Expiration Time:", "expires:", "paid-till:"})
	result.UpdateTime = extractWhoisField(raw, []string{"Updated Date:", "Last Updated:", "Modified:"})
	result.Registrant = extractWhoisField(raw, []string{"Registrant:", "Registrant Name:"})

	// 提取 Name Server（可能有多行）
	result.NameServer = extractWhoisFieldMulti(raw, []string{"Name Server:", "nserver:"})

	// 提取 Domain Status（可能有多行）
	result.DomainStatus = extractWhoisFieldMulti(raw, []string{"Domain Status:", "status:"})

	return result
}

// getWhoisServer 根据域名顶级域选择 WHOIS 服务器
func getWhoisServer(domain string) string {
	parts := strings.Split(domain, ".")
	tld := ""
	if len(parts) >= 2 {
		tld = strings.ToLower(parts[len(parts)-1])
	}

	// 常见顶级域的 WHOIS 服务器
	servers := map[string]string{
		"com":  "whois.verisign-grs.com:43",
		"net":  "whois.verisign-grs.com:43",
		"org":  "whois.pir.org:43",
		"info": "whois.afilias.net:43",
		"cn":   "whois.cnnic.cn:43",
		"uk":   "whois.nic.uk:43",
		"de":   "whois.denic.de:43",
		"jp":   "whois.jprs.jp:43",
		"kr":   "whois.kr:43",
		"ru":   "whois.tcinet.ru:43",
		"fr":   "whois.nic.fr:43",
		"au":   "whois.auda.org.au:43",
		"ca":   "whois.cira.ca:43",
		"br":   "whois.registro.br:43",
		"in":   "whois.registry.in:43",
		"io":   "whois.nic.io:43",
		"dev":  "whois.nic.google:43",
		"app":  "whois.nic.google:43",
	}

	if server, ok := servers[tld]; ok {
		return server
	}

	// 默认使用 IANA WHOIS 服务器
	return "whois.iana.org:43"
}

// extractWhoisField 从 WHOIS 数据中提取指定字段值（取第一个匹配）
func extractWhoisField(raw string, fieldNames []string) string {
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		for _, fieldName := range fieldNames {
			if idx := strings.Index(strings.ToLower(line), strings.ToLower(fieldName)); idx >= 0 {
				value := strings.TrimSpace(line[idx+len(fieldName):])
				// 去掉 URL 括号
				value = strings.TrimPrefix(value, "https://")
				value = strings.TrimPrefix(value, "http://")
				// 去掉常见的前缀标签
				value = strings.TrimSpace(value)
				if value != "" {
					return value
				}
			}
		}
	}
	return "-"
}

// extractWhoisFieldMulti 从 WHOIS 数据中提取指定字段值（合并多行）
func extractWhoisFieldMulti(raw string, fieldNames []string) string {
	lines := strings.Split(raw, "\n")
	var values []string
	for _, line := range lines {
		for _, fieldName := range fieldNames {
			if idx := strings.Index(strings.ToLower(line), strings.ToLower(fieldName)); idx >= 0 {
				value := strings.TrimSpace(line[idx+len(fieldName):])
				value = strings.TrimSpace(value)
				if value != "" {
					values = append(values, value)
				}
			}
		}
	}
	if len(values) == 0 {
		return "-"
	}
	// 最多返回前 5 个
	if len(values) > 5 {
		values = values[:5]
	}
	return strings.Join(values, ", ")
}

// ============================================================
// 网络测速工具
// ============================================================

// SpeedTestResult 测速结果
type SpeedTestResult struct {
	Success      bool    `json:"success"`      // 是否成功
	DownloadMbps float64 `json:"downloadMbps"` // 下载速度（Mbps）
	UploadMbps   float64 `json:"uploadMbps"`   // 上传速度（Mbps）
	LatencyMs    float64 `json:"latencyMs"`    // 延迟（毫秒）
	JitterMs     float64 `json:"jitterMs"`     // 抖动（毫秒）
	PacketLoss   float64 `json:"packetLoss"`   // 丢包率（%）
	TestServer   string  `json:"testServer"`   // 测试服务器
	TestDuration string  `json:"testDuration"` // 测试耗时
	Error        string  `json:"error"`        // 错误信息
}

// SpeedTest 执行网络测速
// 包含下载测速、上传测速、延迟测试
func (n *NetworkTools) SpeedTest() SpeedTestResult {
	result := SpeedTestResult{Success: true}
	startTime := time.Now()

	// 测试服务器列表（公共 CDN 节点）
	testServers := []struct {
		name string
		url  string
	}{
		{"Cloudflare", "https://speed.cloudflare.com/__down?bytes=5000000"},
		{"Cloudflare (alt)", "https://speed.cloudflare.com/__down?bytes=10000000"},
	}

	// 延迟测试
	latency, jitter, packetLoss := n.testLatency()
	result.LatencyMs = latency
	result.JitterMs = jitter
	result.PacketLoss = packetLoss

	// 下载测速
	var downloadSpeeds []float64
	for _, server := range testServers {
		speed := n.testDownloadSpeed(server.name, server.url)
		if speed > 0 {
			downloadSpeeds = append(downloadSpeeds, speed)
			result.TestServer = server.name
		}
	}
	if len(downloadSpeeds) > 0 {
		result.DownloadMbps = avgFloat64(downloadSpeeds)
	}

	// 上传测速
	uploadSpeed := n.testUploadSpeed()
	result.UploadMbps = uploadSpeed

	result.TestDuration = time.Since(startTime).String()

	return result
}

// testLatency 延迟测试：多次 TCP Ping 取平均值
// 返回：平均延迟(ms)、抖动(ms)、丢包率(%)
func (n *NetworkTools) testLatency() (avgLatency, jitter, packetLoss float64) {
	// 测试目标：公共 DNS 服务器
	targets := []string{"8.8.8.8:53", "1.1.1.1:53", "223.5.5.5:53"}
	pingCount := 10

	var latencies []float64
	var lostPackets int

	for _, target := range targets {
		for i := 0; i < pingCount; i++ {
			start := time.Now()
			conn, err := net.DialTimeout("tcp", target, 3*time.Second)
			if err != nil {
				lostPackets++
				continue
			}
			latency := float64(time.Since(start).Nanoseconds()) / 1e6
			conn.Close()
			latencies = append(latencies, latency)
		}
		// 只要有一个目标成功就不再尝试其他
		if len(latencies) >= pingCount {
			break
		}
	}

	if len(latencies) == 0 {
		return 0, 0, 100
	}

	totalPings := len(latencies) + lostPackets

	// 计算平均延迟
	sum := 0.0
	for _, l := range latencies {
		sum += l
	}
	avgLatency = sum / float64(len(latencies))

	// 计算抖动（相邻延迟差的平均值）
	if len(latencies) >= 2 {
		jitterSum := 0.0
		for i := 1; i < len(latencies); i++ {
			diff := latencies[i] - latencies[i-1]
			if diff < 0 {
				diff = -diff
			}
			jitterSum += diff
		}
		jitter = jitterSum / float64(len(latencies)-1)
	}

	// 计算丢包率
	packetLoss = float64(lostPackets) / float64(totalPings) * 100

	return avgLatency, jitter, packetLoss
}

// testDownloadSpeed 下载测速
// 从指定 URL 下载测试文件，计算下载速度
func (n *NetworkTools) testDownloadSpeed(serverName, url string) float64 {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("User-Agent", "XTool/1.0 SpeedTest")

	start := time.Now()
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	// 读取全部数据
	written, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return 0
	}

	elapsed := time.Since(start).Seconds()
	if elapsed <= 0 {
		return 0
	}

	// 计算速度（Mbps = bytes * 8 / 1024 / 1024 / seconds）
	speedMbps := float64(written) * 8.0 / (1024.0 * 1024.0) / elapsed

	return speedMbps
}

// testUploadSpeed 上传测速
// 向测试端点 POST 随机数据，计算上传速度
func (n *NetworkTools) testUploadSpeed() float64 {
	// 生成 2MB 随机测试数据
	uploadData := make([]byte, 2*1024*1024)
	rand.Read(uploadData)

	// 使用 Cloudflare 的测速端点
	testURL := "https://speed.cloudflare.com/__up"

	req, err := http.NewRequest("POST", testURL, bytes.NewReader(uploadData))
	if err != nil {
		return 0
	}
	req.Header.Set("User-Agent", "XTool/1.0 SpeedTest")
	req.Header.Set("Content-Type", "application/octet-stream")

	start := time.Now()
	resp, err := n.httpClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	// 读取响应（确保请求完成）
	io.Copy(io.Discard, resp.Body)

	elapsed := time.Since(start).Seconds()
	if elapsed <= 0 {
		return 0
	}

	// 计算速度（Mbps）
	speedMbps := float64(len(uploadData)) * 8.0 / (1024.0 * 1024.0) / elapsed

	return speedMbps
}

// avgFloat64 计算浮点数切片的平均值
func avgFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
