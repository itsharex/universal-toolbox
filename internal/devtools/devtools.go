// Package devtools 开发工具模块
// 提供 JSON/XML/YAML 格式化、Base64、URL 编解码、哈希计算、UUID 生成等开发常用工具
package devtools

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
	"xtool/internal/db"

	"github.com/google/uuid"
)

// DevTools 开发工具结构体（Wails 绑定到前端）
type DevTools struct {
	db *db.Database // 数据库连接，用于存储代码片段和历史记录
}

// Snippet 代码片段结构
type Snippet struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Language  string `json:"language"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ToolResult 工具执行结果通用结构
type ToolResult struct {
	Success bool   `json:"success"` // 是否成功
	Data    string `json:"data"`    // 结果数据
	Error   string `json:"error"`   // 错误信息（失败时填充）
}

// NewDevTools 创建开发工具模块实例
func NewDevTools(database *db.Database) *DevTools {
	return &DevTools{db: database}
}

// ============================================================
// JSON 工具
// ============================================================

// FormatJSON 格式化 JSON 字符串（美化输出，带缩进）
func (d *DevTools) FormatJSON(input string) ToolResult {
	// 解析 JSON 验证合法性
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ToolResult{Success: false, Error: "JSON 格式错误：" + err.Error()}
	}

	// 重新序列化（带 2 空格缩进）
	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ToolResult{Success: false, Error: "格式化失败：" + err.Error()}
	}

	return ToolResult{Success: true, Data: string(formatted)}
}

// CompressJSON 压缩 JSON 字符串（去除多余空格和换行）
func (d *DevTools) CompressJSON(input string) ToolResult {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ToolResult{Success: false, Error: "JSON 格式错误：" + err.Error()}
	}

	compressed, err := json.Marshal(obj)
	if err != nil {
		return ToolResult{Success: false, Error: "压缩失败：" + err.Error()}
	}

	return ToolResult{Success: true, Data: string(compressed)}
}

// ValidateJSON 校验 JSON 是否合法
func (d *DevTools) ValidateJSON(input string) ToolResult {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ToolResult{Success: false, Error: "不合法的 JSON：" + err.Error()}
	}
	return ToolResult{Success: true, Data: "JSON 格式正确"}
}

// EscapeJSON 转义 JSON 字符串（用于嵌入到字符串字面量中）
// 添加长度检查防止空字符串时越界
func (d *DevTools) EscapeJSON(input string) ToolResult {
	escaped, err := json.Marshal(input)
	if err != nil {
		return ToolResult{Success: false, Error: "转义失败：" + err.Error()}
	}
	// 去掉首尾引号（json.Marshal 会为字符串加上引号）
	result := string(escaped)
	if len(result) >= 2 {
		result = result[1 : len(result)-1]
	} else {
		// 空字符串情况：json.Marshal("") 返回 `""`，长度为 2
		// 如果长度不足 2，返回空字符串
		result = ""
	}
	return ToolResult{Success: true, Data: result}
}

// UnescapeJSON 反转义 JSON 字符串
func (d *DevTools) UnescapeJSON(input string) ToolResult {
	// 加上引号以构成合法的 JSON 字符串值
	quoted := `"` + input + `"`
	var result string
	if err := json.Unmarshal([]byte(quoted), &result); err != nil {
		return ToolResult{Success: false, Error: "反转义失败：" + err.Error()}
	}
	return ToolResult{Success: true, Data: result}
}

// JsonCompare 结构化 JSON 对比
// 解析两个 JSON 并递归比较差异，返回结构化的差异报告
func (d *DevTools) JsonCompare(json1, json2 string) ToolResult {
	var obj1, obj2 interface{}

	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return ToolResult{Success: false, Error: "JSON 1 格式错误：" + err.Error()}
	}
	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return ToolResult{Success: false, Error: "JSON 2 格式错误：" + err.Error()}
	}

	var sb strings.Builder
	diffCount := 0
	compareValues("", obj1, obj2, &sb, &diffCount)

	if diffCount == 0 {
		return ToolResult{Success: true, Data: "两个 JSON 完全相同"}
	}

	return ToolResult{
		Success: true,
		Data:    fmt.Sprintf("共发现 %d 处差异：\n\n%s", diffCount, sb.String()),
	}
}

// compareValues 递归比较两个 JSON 值并记录差异
func compareValues(path string, v1, v2 interface{}, sb *strings.Builder, diffCount *int) {
	switch val1 := v1.(type) {
	case map[string]interface{}:
		val2, ok := v2.(map[string]interface{})
		if !ok {
			*diffCount++
			sb.WriteString(fmt.Sprintf("[路径: %s] 类型不同: object -> %T\n", path, v2))
			return
		}
		// 收集所有键
		allKeys := make(map[string]bool)
		for k := range val1 {
			allKeys[k] = true
		}
		for k := range val2 {
			allKeys[k] = true
		}
		for k := range allKeys {
			childPath := k
			if path != "" {
				childPath = path + "." + k
			}
			if _, ok1 := val1[k]; !ok1 {
				*diffCount++
				sb.WriteString(fmt.Sprintf("[路径: %s] 仅在 JSON2 中存在\n", childPath))
				continue
			}
			if _, ok2 := val2[k]; !ok2 {
				*diffCount++
				sb.WriteString(fmt.Sprintf("[路径: %s] 仅在 JSON1 中存在\n", childPath))
				continue
			}
			compareValues(childPath, val1[k], val2[k], sb, diffCount)
		}
	case []interface{}:
		val2, ok := v2.([]interface{})
		if !ok {
			*diffCount++
			sb.WriteString(fmt.Sprintf("[路径: %s] 类型不同: array -> %T\n", path, v2))
			return
		}
		if len(val1) != len(val2) {
			*diffCount++
			sb.WriteString(fmt.Sprintf("[路径: %s] 数组长度不同: %d vs %d\n", path, len(val1), len(val2)))
		}
		minLen := len(val1)
		if len(val2) < minLen {
			minLen = len(val2)
		}
		for i := 0; i < minLen; i++ {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			compareValues(childPath, val1[i], val2[i], sb, diffCount)
		}
	default:
		// 基本类型比较
		if fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2) {
			*diffCount++
			sb.WriteString(fmt.Sprintf("[路径: %s] 值不同: %v -> %v\n", path, v1, v2))
		}
	}
}

// ============================================================
// XML 工具
// ============================================================

// FormatXML 格式化 XML 字符串
func (d *DevTools) FormatXML(input string) ToolResult {
	// 解码后重新编码（带缩进）
	decoder := xml.NewDecoder(strings.NewReader(input))
	decoder.Strict = false

	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		if err := encoder.EncodeToken(token); err != nil {
			return ToolResult{Success: false, Error: "XML 编码失败：" + err.Error()}
		}
	}

	if err := encoder.Flush(); err != nil {
		return ToolResult{Success: false, Error: "格式化失败：" + err.Error()}
	}

	result := buf.String()
	if result == "" {
		return ToolResult{Success: false, Error: "XML 格式错误，无法解析"}
	}

	return ToolResult{Success: true, Data: result}
}

// ============================================================
// Base64 工具
// ============================================================

// Base64Encode 将字符串编码为 Base64
func (d *DevTools) Base64Encode(input string) ToolResult {
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return ToolResult{Success: true, Data: encoded}
}

// Base64Decode 将 Base64 字符串解码
func (d *DevTools) Base64Decode(input string) ToolResult {
	// 去除可能的空格和换行
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")

	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		// 尝试 URL 安全的 Base64
		decoded, err = base64.URLEncoding.DecodeString(input)
		if err != nil {
			return ToolResult{Success: false, Error: "Base64 解码失败：" + err.Error()}
		}
	}

	// 检查是否为合法 UTF-8 文本
	if !utf8.Valid(decoded) {
		// 以十六进制显示
		return ToolResult{Success: true, Data: fmt.Sprintf("(二进制数据) 十六进制: %x", decoded)}
	}

	return ToolResult{Success: true, Data: string(decoded)}
}

// ============================================================
// URL 编解码工具
// ============================================================

// URLEncode URL 编码字符串
func (d *DevTools) URLEncode(input string) ToolResult {
	encoded := url.QueryEscape(input)
	return ToolResult{Success: true, Data: encoded}
}

// URLDecode URL 解码字符串
func (d *DevTools) URLDecode(input string) ToolResult {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return ToolResult{Success: false, Error: "URL 解码失败：" + err.Error()}
	}
	return ToolResult{Success: true, Data: decoded}
}

// ============================================================
// 哈希计算工具
// ============================================================

// CalcMD5 计算字符串的 MD5 哈希值
func (d *DevTools) CalcMD5(input string) ToolResult {
	hash := md5.Sum([]byte(input))
	return ToolResult{Success: true, Data: fmt.Sprintf("%x", hash)}
}

// CalcSHA1 计算字符串的 SHA1 哈希值
func (d *DevTools) CalcSHA1(input string) ToolResult {
	hash := sha1.Sum([]byte(input))
	return ToolResult{Success: true, Data: fmt.Sprintf("%x", hash)}
}

// CalcSHA256 计算字符串的 SHA256 哈希值
func (d *DevTools) CalcSHA256(input string) ToolResult {
	hash := sha256.Sum256([]byte(input))
	return ToolResult{Success: true, Data: fmt.Sprintf("%x", hash)}
}

// CalcSHA512 计算字符串的 SHA512 哈希值
func (d *DevTools) CalcSHA512(input string) ToolResult {
	hash := sha512.Sum512([]byte(input))
	return ToolResult{Success: true, Data: fmt.Sprintf("%x", hash)}
}

// CalcFileHash 计算文件的哈希值（支持 MD5/SHA1/SHA256/SHA512）
// algorithm: md5, sha1, sha256, sha512
func (d *DevTools) CalcFileHash(filePath, algorithm string) ToolResult {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return ToolResult{Success: false, Error: "无法打开文件：" + err.Error()}
	}
	defer file.Close()

	// 根据算法选择哈希函数
	var hashWriter interface {
		io.Writer
		Sum([]byte) []byte
	}

	switch strings.ToLower(algorithm) {
	case "md5":
		hashWriter = md5.New()
	case "sha1":
		hashWriter = sha1.New()
	case "sha256":
		hashWriter = sha256.New()
	case "sha512":
		hashWriter = sha512.New()
	default:
		return ToolResult{Success: false, Error: "不支持的哈希算法（支持: md5/sha1/sha256/sha512）"}
	}

	// 读取文件内容并计算哈希
	if _, err := io.Copy(hashWriter, file); err != nil {
		return ToolResult{Success: false, Error: "读取文件失败：" + err.Error()}
	}

	hashBytes := hashWriter.Sum(nil)
	return ToolResult{
		Success: true,
		Data:    fmt.Sprintf("%s: %s\n文件: %s", strings.ToUpper(algorithm), hex.EncodeToString(hashBytes), filePath),
	}
}

// ============================================================
// 文本工具
// ============================================================

// TextCompare 对比两个文本，返回差异行
func (d *DevTools) TextCompare(text1, text2 string) ToolResult {
	lines1 := strings.Split(text1, "\n")
	lines2 := strings.Split(text2, "\n")

	var sb strings.Builder
	maxLen := len(lines1)
	if len(lines2) > maxLen {
		maxLen = len(lines2)
	}

	diffCount := 0
	for i := 0; i < maxLen; i++ {
		var l1, l2 string
		if i < len(lines1) {
			l1 = lines1[i]
		}
		if i < len(lines2) {
			l2 = lines2[i]
		}
		if l1 != l2 {
			diffCount++
			sb.WriteString(fmt.Sprintf("第 %d 行:\n  原文: %s\n  新文: %s\n\n", i+1, l1, l2))
		}
	}

	if diffCount == 0 {
		return ToolResult{Success: true, Data: "两段文本完全相同"}
	}

	return ToolResult{
		Success: true,
		Data:    fmt.Sprintf("共发现 %d 处差异：\n\n%s", diffCount, sb.String()),
	}
}

// TextReplace 在文本中执行查找替换
func (d *DevTools) TextReplace(text, search, replace string, useRegex bool) ToolResult {
	if text == "" || search == "" {
		return ToolResult{Success: false, Error: "文本和搜索内容不能为空"}
	}

	var result string
	if useRegex {
		// 使用正则表达式替换
		re, err := regexp.Compile(search)
		if err != nil {
			return ToolResult{Success: false, Error: "正则表达式错误：" + err.Error()}
		}
		result = re.ReplaceAllString(text, replace)
	} else {
		// 普通字符串替换
		result = strings.ReplaceAll(text, search, replace)
	}

	return ToolResult{Success: true, Data: result}
}

// TextStats 统计文本的字符数、单词数、行数等信息
func (d *DevTools) TextStats(text string) ToolResult {
	charCount := len([]rune(text))           // 字符数（支持中文）
	charNoSpace := len([]rune(strings.ReplaceAll(text, " ", ""))) // 字符数（不含空格）
	byteCount := len(text)                   // 字节数
	lineCount := strings.Count(text, "\n") + 1 // 行数
	wordCount := len(strings.Fields(text))    // 英文单词数
	paragraphCount := 0                       // 段落数（非空行）

	// 统计中文字符数
	chineseCount := 0
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			chineseCount++
		}
	}

	// 统计段落数（以连续空行分隔）
	paragraphs := strings.Split(text, "\n\n")
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			paragraphCount++
		}
	}

	result := fmt.Sprintf(
		"字符数（含空格）: %d\n字符数（不含空格）: %d\n字节数: %d\n行数: %d\n英文单词数: %d\n中文字符数: %d\n段落数: %d",
		charCount, charNoSpace, byteCount, lineCount, wordCount, chineseCount, paragraphCount,
	)

	return ToolResult{Success: true, Data: result}
}

// TextStatsJSON 统计文本信息并返回 JSON 结构化数据
func (d *DevTools) TextStatsJSON(text string) ToolResult {
	charCount := len([]rune(text))
	charNoSpace := len([]rune(strings.ReplaceAll(text, " ", "")))
	byteCount := len(text)
	lineCount := strings.Count(text, "\n") + 1
	wordCount := len(strings.Fields(text))
	paragraphCount := 0

	// 统计中文字符数
	chineseCount := 0
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			chineseCount++
		}
	}

	// 统计段落数
	paragraphs := strings.Split(text, "\n\n")
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			paragraphCount++
		}
	}

	// 构建结构化结果
	stats := map[string]interface{}{
		"charCount":       charCount,
		"charNoSpace":     charNoSpace,
		"byteCount":       byteCount,
		"lineCount":       lineCount,
		"wordCount":       wordCount,
		"chineseCount":    chineseCount,
		"paragraphCount":  paragraphCount,
	}

	jsonBytes, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return ToolResult{Success: false, Error: "序列化统计结果失败：" + err.Error()}
	}

	return ToolResult{Success: true, Data: string(jsonBytes)}
}

// ============================================================
// UUID 生成工具
// ============================================================

// GenerateUUID 生成一个 UUID v4
func (d *DevTools) GenerateUUID() ToolResult {
	return ToolResult{Success: true, Data: uuid.New().String()}
}

// GenerateUUIDs 批量生成多个 UUID
// 返回提示信息，告知生成了多少个 UUID
func (d *DevTools) GenerateUUIDs(count int) ToolResult {
	if count <= 0 || count > 100 {
		count = 10 // 默认 10 个，最多 100 个
	}

	var ids []string
	for i := 0; i < count; i++ {
		ids = append(ids, uuid.New().String())
	}

	result := strings.Join(ids, "\n")
	// 添加返回提示信息
	result = fmt.Sprintf("已生成 %d 个 UUID：\n%s", count, result)

	return ToolResult{Success: true, Data: result}
}

// ============================================================
// 时间戳工具
// ============================================================

// TimestampToDatetime 时间戳（秒/毫秒）转日期时间字符串
func (d *DevTools) TimestampToDatetime(timestamp int64) ToolResult {
	var t time.Time

	// 自动判断是秒还是毫秒级时间戳
	if timestamp > 1e12 {
		// 毫秒级时间戳
		t = time.Unix(timestamp/1000, (timestamp%1000)*int64(time.Millisecond)).Local()
	} else {
		// 秒级时间戳
		t = time.Unix(timestamp, 0).Local()
	}

	result := fmt.Sprintf(
		"本地时间: %s\nUTC 时间: %s\n时区: %s",
		t.Format("2006-01-02 15:04:05"),
		t.UTC().Format("2006-01-02 15:04:05 UTC"),
		t.Format("MST"),
	)

	return ToolResult{Success: true, Data: result}
}

// DatetimeToTimestamp 日期时间字符串转时间戳
func (d *DevTools) DatetimeToTimestamp(datetime string) ToolResult {
	// 支持多种格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.ParseInLocation(format, datetime, time.Local)
		if err == nil {
			break
		}
	}

	if err != nil {
		return ToolResult{Success: false, Error: "日期格式不正确，请使用 YYYY-MM-DD HH:MM:SS 格式"}
	}

	result := fmt.Sprintf(
		"秒级时间戳: %d\n毫秒级时间戳: %d",
		t.Unix(),
		t.UnixMilli(),
	)

	return ToolResult{Success: true, Data: result}
}

// GetCurrentTimestamp 获取当前时间戳
func (d *DevTools) GetCurrentTimestamp() ToolResult {
	now := time.Now()
	result := fmt.Sprintf(
		"当前时间: %s\n秒级时间戳: %d\n毫秒级时间戳: %d",
		now.Format("2006-01-02 15:04:05"),
		now.Unix(),
		now.UnixMilli(),
	)
	return ToolResult{Success: true, Data: result}
}

// ============================================================
// 正则表达式工具
// ============================================================

// RegexTest 测试正则表达式匹配结果
func (d *DevTools) RegexTest(pattern, text string) ToolResult {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ToolResult{Success: false, Error: "正则表达式语法错误：" + err.Error()}
	}

	matches := re.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return ToolResult{Success: true, Data: "未找到匹配项"}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("找到 %d 个匹配项：\n\n", len(matches)))

	for i, match := range matches {
		sb.WriteString(fmt.Sprintf("匹配 #%d: %s\n", i+1, match[0]))
		for j, group := range match[1:] {
			sb.WriteString(fmt.Sprintf("  捕获组 %d: %s\n", j+1, group))
		}
	}

	return ToolResult{Success: true, Data: sb.String()}
}

// ============================================================
// 代码片段管理
// ============================================================

// GetSnippets 获取所有代码片段
func (d *DevTools) GetSnippets() ([]Snippet, error) {
	rows, err := d.db.DB.Query(
		"SELECT id, title, content, language, tags, created_at, updated_at FROM snippets ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Language, &s.Tags, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		snippets = append(snippets, s)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历代码片段结果集失败: %w", err)
	}

	return snippets, nil
}

// ============================================================
// 进制转换工具
// ============================================================

// ConvertBase 进制转换，支持 2/8/10/16 进制互转，支持小数
// input: 输入字符串
// fromBase: 源进制（2/8/10/16）
// toBase: 目标进制（2/8/10/16）
func (d *DevTools) ConvertBase(input string, fromBase, toBase int) ToolResult {
	// 校验进制参数
	validBases := map[int]bool{2: true, 8: true, 10: true, 16: true}
	if !validBases[fromBase] || !validBases[toBase] {
		return ToolResult{Success: false, Error: "不支持的进制（支持: 2/8/10/16）"}
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return ToolResult{Success: false, Error: "输入不能为空"}
	}

	// 判断是否包含小数点
	if strings.Contains(input, ".") {
		return convertBaseFloat(input, fromBase, toBase)
	}
	return convertBaseInt(input, fromBase, toBase)
}

// convertBaseInt 整数部分进制转换
func convertBaseInt(input string, fromBase, toBase int) ToolResult {
	// 将源进制字符串转为十进制整数值
	var value *big.Int
	switch fromBase {
	case 2:
		v := new(big.Int)
		_, ok := v.SetString(input, 2)
		if !ok {
			return ToolResult{Success: false, Error: fmt.Sprintf("无效的%d进制数: %s", fromBase, input)}
		}
		value = v
	case 8:
		v := new(big.Int)
		_, ok := v.SetString(input, 8)
		if !ok {
			return ToolResult{Success: false, Error: fmt.Sprintf("无效的%d进制数: %s", fromBase, input)}
		}
		value = v
	case 10:
		v := new(big.Int)
		_, ok := v.SetString(input, 10)
		if !ok {
			return ToolResult{Success: false, Error: fmt.Sprintf("无效的%d进制数: %s", fromBase, input)}
		}
		value = v
	case 16:
		v := new(big.Int)
		_, ok := v.SetString(input, 16)
		if !ok {
			return ToolResult{Success: false, Error: fmt.Sprintf("无效的%d进制数: %s", fromBase, input)}
		}
		value = v
	}

	// 转换为目标进制
	var result string
	switch toBase {
	case 2:
		result = value.Text(2)
	case 8:
		result = value.Text(8)
	case 10:
		result = value.Text(10)
	case 16:
		result = strings.ToUpper(value.Text(16))
	}

	// 生成转换过程说明
	process := fmt.Sprintf("转换过程：\n  输入: %s（%d进制）\n  十进制值: %s\n  输出: %s（%d进制）",
		input, fromBase, value.Text(10), result, toBase)

	return ToolResult{Success: true, Data: fmt.Sprintf("%s\n\n%s", result, process)}
}

// convertBaseFloat 小数部分进制转换
func convertBaseFloat(input string, fromBase, toBase int) ToolResult {
	parts := strings.SplitN(input, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = parts[1]
	}

	// 整数部分转换
	intResult := convertBaseInt(intPart, fromBase, toBase)
	if !intResult.Success {
		return intResult
	}

	// 小数部分转换：先转为十进制小数
	fracDecimal := 0.0
	for i, ch := range fracPart {
		digitVal := digitToValue(ch, fromBase)
		if digitVal < 0 {
			return ToolResult{Success: false, Error: fmt.Sprintf("无效字符 '%c' 在%d进制中", ch, fromBase)}
		}
		fracDecimal += float64(digitVal) * powFloat(1.0/float64(fromBase), float64(i+1))
	}

	// 十进制小数转目标进制（最多 20 位精度）
	if fracDecimal == 0 {
		// 没有小数部分
		return ToolResult{Success: true, Data: intResult.Data}
	}

	fracResult := ""
	seen := make(map[float64]bool)
	for i := 0; i < 20; i++ {
		fracDecimal *= float64(toBase)
		intVal := int(fracDecimal)
		fracDecimal -= float64(intVal)

		if intVal >= toBase {
			intVal = toBase - 1
		}

		fracResult += valueToDigit(intVal)

		if fracDecimal == 0 {
			break
		}
		// 检测循环小数
		if seen[fracDecimal] {
			fracResult += "..."
			break
		}
		seen[fracDecimal] = true
	}

	// 提取整数部分的纯结果（去掉转换过程说明）
	intPure := strings.SplitN(intResult.Data, "\n", 2)[0]

	finalResult := intPure + "." + fracResult
	process := fmt.Sprintf("转换过程：\n  输入: %s（%d进制）\n  整数部分: %s\n  小数部分: %s\n  输出: %s（%d进制）",
		input, fromBase, intPure, fracResult, finalResult, toBase)

	return ToolResult{Success: true, Data: fmt.Sprintf("%s\n\n%s", finalResult, process)}
}

// digitToValue 将字符转换为对应进制的数值
func digitToValue(ch rune, base int) int {
	if ch >= '0' && ch <= '9' {
		v := int(ch - '0')
		if v < base {
			return v
		}
	}
	if ch >= 'a' && ch <= 'f' {
		v := int(ch-'a') + 10
		if v < base {
			return v
		}
	}
	if ch >= 'A' && ch <= 'F' {
		v := int(ch-'A') + 10
		if v < base {
			return v
		}
	}
	return -1
}

// valueToDigit 将数值转换为对应字符
func valueToDigit(v int) string {
	if v < 10 {
		return fmt.Sprintf("%d", v)
	}
	return string(rune('A' + v - 10))
}

// powFloat 浮点数幂运算
func powFloat(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

// ============================================================
// 占位文本生成工具
// ============================================================

// GenerateDummyData 生成占位/测试数据
// dataType: name/phone/email/address/company/ipv4/ipv6/mac/date/uuid/sentence/paragraph
// locale: zh/en
// count: 生成数量 1-100
func (d *DevTools) GenerateDummyData(dataType string, count int, locale string) ToolResult {
	dataType = strings.ToLower(dataType)
	locale = strings.ToLower(locale)

	if count < 1 {
		count = 1
	}
	if count > 100 {
		count = 100
	}

	var results []string
	for i := 0; i < count; i++ {
		result := generateOneData(dataType, locale)
		results = append(results, result)
	}

	data := strings.Join(results, "\n")
	return ToolResult{Success: true, Data: data}
}

// generateOneData 生成一条指定类型的数据
func generateOneData(dataType, locale string) string {
	switch dataType {
	case "name":
		if locale == "en" {
			return generateEnglishName()
		}
		return generateChineseName()
	case "phone":
		if locale == "en" {
			return generateEnglishPhone()
		}
		return generateChinesePhone()
	case "email":
		return generateEmail(locale)
	case "address":
		if locale == "en" {
			return generateEnglishAddress()
		}
		return generateChineseAddress()
	case "company":
		if locale == "en" {
			return generateEnglishCompany()
		}
		return generateChineseCompany()
	case "ipv4":
		return generateIPv4()
	case "ipv6":
		return generateIPv6()
	case "mac":
		return generateMAC()
	case "date":
		return generateRandomDate()
	case "uuid":
		return uuid.New().String()
	case "sentence":
		if locale == "en" {
			return generateEnglishSentence()
		}
		return generateChineseSentence()
	case "paragraph":
		if locale == "en" {
			return generateEnglishParagraph()
		}
		return generateChineseParagraph()
	default:
		return fmt.Sprintf("不支持的数据类型: %s", dataType)
	}
}

// --- 中文姓名库 ---
var chineseSurnames = []string{
	"王", "李", "张", "刘", "陈", "杨", "赵", "黄", "周", "吴",
	"徐", "孙", "胡", "朱", "高", "林", "何", "郭", "马", "罗",
	"梁", "宋", "郑", "谢", "韩", "唐", "冯", "于", "董", "萧",
	"程", "曹", "袁", "邓", "许", "傅", "沈", "曾", "彭", "吕",
	"苏", "卢", "蒋", "蔡", "贾", "丁", "魏", "薛", "叶", "阎",
	"余", "潘", "杜", "戴", "夏", "钟", "汪", "田", "任", "姜",
	"范", "方", "石", "姚", "谭", "廖", "邹", "熊", "金", "陆",
	"郝", "孔", "白", "崔", "康", "毛", "邱", "秦", "江", "史",
}

var chineseMaleNames = []string{
	"伟", "强", "磊", "洋", "勇", "军", "杰", "涛", "明", "辉",
	"鹏", "华", "飞", "刚", "平", "建", "超", "志", "国", "海",
	"波", "亮", "斌", "龙", "宏", "峰", "毅", "浩", "宇", "轩",
	"泽", "昊", "天", "翔", "晨", "睿", "博", "文", "俊", "达",
}

var chineseFemaleNames = []string{
	"芳", "娜", "秀英", "敏", "静", "丽", "强", "洁", "娟", "艳",
	"秀兰", "霞", "燕", "玲", "桂英", "红", "梅", "莉", "婷", "慧",
	"雪", "琳", "倩", "颖", "露", "瑶", "欣", "蕾", "薇", "洁",
	"梦", "怡", "佳", "馨", "月", "云", "蓉", "汐", "妍", "彤",
}

// generateChineseName 生成中文姓名
func generateChineseName() string {
	surname := chineseSurnames[randInt(len(chineseSurnames))]
	if randInt(2) == 0 {
		// 两字名
		given := chineseMaleNames[randInt(len(chineseMaleNames))]
		return surname + given
	}
	// 三字名（姓氏 + 单字）
	allNames := append(chineseMaleNames, chineseFemaleNames...)
	given := allNames[randInt(len(allNames))]
	return surname + given
}

// --- 英文姓名库 ---
var englishFirstNames = []string{
	"James", "John", "Robert", "Michael", "William", "David", "Richard", "Joseph", "Thomas", "Charles",
	"Mary", "Patricia", "Jennifer", "Linda", "Barbara", "Elizabeth", "Susan", "Jessica", "Sarah", "Karen",
	"Daniel", "Matthew", "Anthony", "Mark", "Steven", "Andrew", "Joshua", "Kevin", "Brian", "Edward",
	"Emma", "Olivia", "Ava", "Sophia", "Isabella", "Mia", "Charlotte", "Amelia", "Harper", "Evelyn",
}

var englishLastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
	"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
}

func generateEnglishName() string {
	first := englishFirstNames[randInt(len(englishFirstNames))]
	last := englishLastNames[randInt(len(englishLastNames))]
	return first + " " + last
}

// --- 手机号生成 ---
func generateChinesePhone() string {
	// 中国手机号前缀
	prefixes := []string{"130", "131", "132", "133", "134", "135", "136", "137", "138", "139",
		"150", "151", "152", "153", "155", "156", "157", "158", "159",
		"170", "171", "172", "173", "175", "176", "177", "178",
		"180", "181", "182", "183", "184", "185", "186", "187", "188", "189",
		"191", "198", "199"}
	prefix := prefixes[randInt(len(prefixes))]
	suffix := randIntN(100000000)
	return fmt.Sprintf("%s%08d", prefix, suffix)
}

func generateEnglishPhone() string {
	// 美国手机号格式: (XXX) XXX-XXXX
	area := randIntN(900) + 100
	mid := randIntN(900) + 100
	last := randIntN(9000) + 1000
	return fmt.Sprintf("(%d) %d-%d", area, mid, last)
}

// --- 邮箱生成 ---
var emailDomains = []string{"gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "qq.com", "163.com", "126.com", "foxmail.com"}

func generateEmail(locale string) string {
	var name string
	if locale == "en" {
		first := strings.ToLower(englishFirstNames[randInt(len(englishFirstNames))])
		last := strings.ToLower(englishLastNames[randInt(len(englishLastNames))])
		// 随机选择格式
		switch randInt(3) {
		case 0:
			name = first + "." + last
		case 1:
			name = first + last
		default:
			name = first + randString(3)
		}
	} else {
		pinyin := []string{"zhangsan", "lisi", "wangwu", "zhaoliu", "qianqi", "sunba", "zhoujiu", "wushi",
			"chenxiaoming", "lihua", "wangfang", "zhangwei", "liuyang", "chenjie", "huanglei"}
		name = pinyin[randInt(len(pinyin))]
		if randInt(2) == 0 {
			name += fmt.Sprintf("%d", randIntN(999))
		}
	}
	domain := emailDomains[randInt(len(emailDomains))]
	return name + "@" + domain
}

// --- 地址生成 ---
var chineseCities = []string{"北京市", "上海市", "广州市", "深圳市", "杭州市", "成都市", "武汉市", "南京市", "重庆市", "西安市",
	"苏州市", "天津市", "长沙市", "郑州市", "东莞市", "青岛市", "沈阳市", "宁波市", "昆明市", "大连市"}

var chineseDistricts = []string{"朝阳区", "海淀区", "浦东新区", "南山区", "天河区", "武侯区", "武昌区", "鼓楼区", "渝中区", "雁塔区",
	"西湖区", "滨江区", "雨花台区", "江干区", "江汉区", "洪山区"}

var chineseStreets = []string{"中山路", "人民路", "建设路", "和平路", "解放路", "长安路", "文化路", "科技路", "创新路", "幸福路",
	"光明路", "学院路", "花园路", "金融街", "商业街", "工业路", "滨河路", "湖滨路"}

func generateChineseAddress() string {
	city := chineseCities[randInt(len(chineseCities))]
	district := chineseDistricts[randInt(len(chineseDistricts))]
	street := chineseStreets[randInt(len(chineseStreets))]
	number := randIntN(200) + 1
	building := randIntN(30) + 1
	unit := randIntN(5) + 1
	room := randIntN(2000) + 101
	return fmt.Sprintf("%s%s%s%d号%d栋%d单元%d室", city, district, street, number, building, unit, room)
}

func generateEnglishAddress() string {
	streets := []string{"Main St", "Oak Ave", "Maple Dr", "Cedar Ln", "Pine Rd", "Elm St", "Washington Blvd", "Park Ave", "Lake Dr", "Hill Rd"}
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	states := []string{"NY", "CA", "IL", "TX", "AZ", "PA", "FL", "OH", "GA", "NC"}
	number := randIntN(9999) + 1
	zip := randIntN(90000) + 10000
	return fmt.Sprintf("%d %s, %s, %s %d", number, streets[randInt(len(streets))], cities[randInt(len(cities))], states[randInt(len(states))], zip)
}

// --- 公司生成 ---
var chineseCompanyPrefixes = []string{"华", "中", "国", "天", "金", "鑫", "盛", "宏", "伟", "新", "创", "智", "汇", "达", "恒", "信", "安", "博", "瑞", "嘉"}
var chineseCompanyMiddles = []string{"联", "通", "泰", "科", "达", "创", "信", "和", "润", "源", "正", "方", "同", "合", "利", "升", "启", "明", "辰", "宇"}
var chineseCompanySuffixes = []string{"科技有限公司", "信息技术有限公司", "网络科技有限公司", "数据服务有限公司", "软件开发有限公司", "咨询有限公司", "电子商务有限公司", "文化传媒有限公司"}

func generateChineseCompany() string {
	prefix := chineseCompanyPrefixes[randInt(len(chineseCompanyPrefixes))]
	middle := chineseCompanyMiddles[randInt(len(chineseCompanyMiddles))]
	suffix := chineseCompanySuffixes[randInt(len(chineseCompanySuffixes))]
	return prefix + middle + suffix
}

func generateEnglishCompany() string {
	prefixes := []string{"Global", "Pacific", "Atlantic", "National", "United", "Advanced", "Digital", "Smart", "Innovative", "Prime"}
	middles := []string{"Tech", "Solutions", "Systems", "Dynamics", "Networks", "Services", "Labs", "Group", "Industries", "Ventures"}
	suffixes := []string{"Inc.", "LLC", "Corp.", "Ltd.", "Co.", "Group"}
	return fmt.Sprintf("%s %s %s", prefixes[randInt(len(prefixes))], middles[randInt(len(middles))], suffixes[randInt(len(suffixes))])
}

// --- 网络地址生成 ---
func generateIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", randIntN(223)+1, randIntN(256), randIntN(256), randIntN(254)+1)
}

func generateIPv6() string {
	groups := make([]string, 8)
	for i := range groups {
		groups[i] = fmt.Sprintf("%04x", randIntN(0x10000))
	}
	return strings.Join(groups, ":")
}

func generateMAC() string {
	groups := make([]string, 6)
	for i := range groups {
		groups[i] = fmt.Sprintf("%02x", randIntN(256))
	}
	return strings.Join(groups, ":")
}

// --- 日期生成 ---
func generateRandomDate() string {
	year := randIntN(10) + 2015
	month := randIntN(12) + 1
	day := randIntN(28) + 1
	hour := randIntN(24)
	minute := randIntN(60)
	second := randIntN(60)
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}

// --- 句子/段落生成 ---
var chineseSentences = []string{
	"今天天气真不错，适合出去散步。",
	"学习使人进步，知识改变命运。",
	"科技的发展改变了人们的生活方式。",
	"坚持就是胜利，不要轻易放弃。",
	"时间如流水，一去不复返。",
	"读书破万卷，下笔如有神。",
	"千里之行，始于足下。",
	"学而不思则罔，思而不学则殆。",
	"生活不止眼前的苟且，还有诗和远方。",
	"实践是检验真理的唯一标准。",
	"创新是一个民族进步的灵魂。",
	"数据是新时代的石油。",
	"人工智能正在改变各行各业。",
	"云计算让计算资源触手可及。",
	"网络安全是数字时代的重要基石。",
}

var englishSentences = []string{
	"The quick brown fox jumps over the lazy dog.",
	"Technology is best when it brings people together.",
	"Innovation distinguishes between a leader and a follower.",
	"The only way to do great work is to love what you do.",
	"Data is the new oil of the digital economy.",
	"Artificial intelligence is the future of computing.",
	"Cloud computing has revolutionized the IT industry.",
	"Security is not a product, but a process.",
	"The best way to predict the future is to invent it.",
	"Code is like humor. When you have to explain it, it is bad.",
}

func generateChineseSentence() string {
	return chineseSentences[randInt(len(chineseSentences))]
}

func generateEnglishSentence() string {
	return englishSentences[randInt(len(englishSentences))]
}

func generateChineseParagraph() string {
	count := randIntN(3) + 3 // 3-5 句
	var sentences []string
	for i := 0; i < count; i++ {
		sentences = append(sentences, chineseSentences[randInt(len(chineseSentences))])
	}
	return strings.Join(sentences, "")
}

func generateEnglishParagraph() string {
	count := randIntN(3) + 3
	var sentences []string
	for i := 0; i < count; i++ {
		sentences = append(sentences, englishSentences[randInt(len(englishSentences))])
	}
	return strings.Join(sentences, " ")
}

// --- 随机工具函数 ---
func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func randIntN(max int) int {
	return randInt(max)
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[randInt(len(charset))]
	}
	return string(result)
}

// ============================================================
// 接口文档生成工具
// ============================================================

// GenerateAPIDoc 根据输入信息生成 Markdown 格式接口文档
// apiName: 接口名称
// apiDesc: 接口描述
// method: 请求方法（GET/POST/PUT/DELETE/PATCH）
// url: 请求地址
// headersJSON: 请求头 JSON
// bodyJSON: 请求体 JSON
// responseJSON: 响应体 JSON
func (d *DevTools) GenerateAPIDoc(apiName, apiDesc, method, url, headersJSON, bodyJSON, responseJSON string) ToolResult {
	method = strings.ToUpper(method)
	if method == "" {
		method = "GET"
	}

	var sb strings.Builder

	// 文档标题
	sb.WriteString(fmt.Sprintf("# %s\n\n", apiName))
	if apiDesc != "" {
		sb.WriteString(fmt.Sprintf("> %s\n\n", apiDesc))
	}

	// 基本信息
	sb.WriteString("## 基本信息\n\n")
	sb.WriteString(fmt.Sprintf("| 项目 | 值 |\n| --- | --- |\n"))
	sb.WriteString(fmt.Sprintf("| 请求方法 | `%s` |\n", method))
	sb.WriteString(fmt.Sprintf("| 请求地址 | `%s` |\n", url))
	sb.WriteString(fmt.Sprintf("| Content-Type | `application/json` |\n\n"))

	// 请求头
	if headersJSON != "" {
		sb.WriteString("## 请求头\n\n")
		sb.WriteString("| 参数名 | 值 | 说明 |\n| --- | --- | --- |\n")
		var headers map[string]string
		if err := json.Unmarshal([]byte(headersJSON), &headers); err == nil {
			for k, v := range headers {
				sb.WriteString(fmt.Sprintf("| `%s` | `%s` | - |\n", k, v))
			}
		}
		sb.WriteString("\n")
	}

	// 请求参数
	if bodyJSON != "" {
		sb.WriteString("## 请求参数\n\n")
		// 格式化 JSON
		var obj interface{}
		if err := json.Unmarshal([]byte(bodyJSON), &obj); err == nil {
			formatted, _ := json.MarshalIndent(obj, "", "  ")
			sb.WriteString(fmt.Sprintf("```json\n%s\n```\n\n", string(formatted)))

			// 生成参数表
			if m, ok := obj.(map[string]interface{}); ok {
				sb.WriteString("### 参数说明\n\n")
				sb.WriteString("| 参数名 | 类型 | 必填 | 说明 |\n| --- | --- | --- | --- |\n")
				for k, v := range m {
					typeName := "string"
					switch v.(type) {
					case float64:
						typeName = "number"
					case bool:
						typeName = "boolean"
					case []interface{}:
						typeName = "array"
					case map[string]interface{}:
						typeName = "object"
					case nil:
						typeName = "null"
					}
					sb.WriteString(fmt.Sprintf("| `%s` | %s | 是 | - |\n", k, typeName))
				}
				sb.WriteString("\n")
			}
		}
	}

	// 响应示例
	sb.WriteString("## 响应示例\n\n")
	if responseJSON != "" {
		var obj interface{}
		if err := json.Unmarshal([]byte(responseJSON), &obj); err == nil {
			formatted, _ := json.MarshalIndent(obj, "", "  ")
			sb.WriteString(fmt.Sprintf("```json\n%s\n```\n\n", string(formatted)))
		}
	} else {
		sb.WriteString("```json\n{\n  \"code\": 200,\n  \"message\": \"success\",\n  \"data\": null\n}\n```\n\n")
	}

	// 错误码
	sb.WriteString("## 错误码\n\n")
	sb.WriteString("| 错误码 | 说明 |\n| --- | --- |\n")
	sb.WriteString("| 200 | 请求成功 |\n")
	sb.WriteString("| 400 | 请求参数错误 |\n")
	sb.WriteString("| 401 | 未授权，需要登录 |\n")
	sb.WriteString("| 403 | 禁止访问 |\n")
	sb.WriteString("| 404 | 资源不存在 |\n")
	sb.WriteString("| 500 | 服务器内部错误 |\n\n")

	return ToolResult{Success: true, Data: sb.String()}
}

// ============================================================
// 代码混淆工具
// ============================================================

// ObfuscateCode 对代码进行简单混淆
// language: javascript
// 实现：变量重命名、字符串编码、控制流混淆
func (d *DevTools) ObfuscateCode(code string, language string) ToolResult {
	language = strings.ToLower(language)
	if language != "javascript" {
		return ToolResult{Success: false, Error: "目前仅支持 JavaScript 混淆"}
	}

	if strings.TrimSpace(code) == "" {
		return ToolResult{Success: false, Error: "代码不能为空"}
	}

	var result strings.Builder
	result.WriteString("// 混淆后的代码\n")

	// 步骤1：变量重命名
	obfuscated := obfuscateJSVariables(code)

	// 步骤2：字符串编码
	obfuscated = obfuscateJSStrings(obfuscated)

	// 步骤3：控制流混淆（包裹在自执行函数中）
	obfuscated = obfuscateJSControlFlow(obfuscated)

	result.WriteString(obfuscated)

	return ToolResult{Success: true, Data: result.String()}
}

// obfuscateJSVariables JavaScript 变量重命名
// 将常见的变量名替换为短变量名
func obfuscateJSVariables(code string) string {
	// 匹配 var/let/const 声明的变量名
	re := regexp.MustCompile(`\b(var|let|const)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)`)
	matches := re.FindAllStringSubmatchIndex(code, -1)

	if len(matches) == 0 {
		return code
	}

	// 收集变量名并建立映射
	nameMap := make(map[string]string)
	varCount := 0

	for _, match := range matches {
		nameStart := match[4]
		nameEnd := match[5]
		name := code[nameStart:nameEnd]

		if _, exists := nameMap[name]; !exists {
			// 生成混淆变量名 _0x + 十六进制
			obfName := fmt.Sprintf("_0x%x", varCount)
			nameMap[name] = obfName
			varCount++
		}
	}

	// 替换变量名（从长到短排序，避免部分替换）
	var names []string
	for name := range nameMap {
		names = append(names, name)
	}
	sortStringsByLength(names, true)

	for _, name := range names {
		obfName := nameMap[name]
		// 使用单词边界匹配替换
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`)
		code = re.ReplaceAllString(code, obfName)
	}

	return code
}

// obfuscateJSStrings JavaScript 字符串编码
// 将字符串字面量替换为十六进制编码的解码调用
func obfuscateJSStrings(code string) string {
	// 匹配双引号和单引号字符串（简单匹配，不处理转义）
	re := regexp.MustCompile(`"(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'`)

	return re.ReplaceAllStringFunc(code, func(s string) string {
		// 去掉引号
		inner := s[1 : len(s)-1]
		if len(inner) == 0 {
			return s
		}

		// 转为十六进制编码字符串
		var hexParts []string
		for _, ch := range inner {
			hexParts = append(hexParts, fmt.Sprintf("0x%x", ch))
		}

		// 生成解码表达式
		return fmt.Sprintf("(function(){var _s=[%s];var _r='';for(var _i=0;_i<_s.length;_i++){_r+=String.fromCharCode(_s[_i]);}return _r;})()",
			strings.Join(hexParts, ","))
	})
}

// obfuscateJSControlFlow JavaScript 控制流混淆
// 将代码包裹在自执行函数中，添加控制流平坦化
func obfuscateJSControlFlow(code string) string {
	// 添加控制流平坦化外壳
	wrapper := fmt.Sprintf(`(function(_0xentry){
var _0xstate=_0xentry===1?0:1;
switch(_0xstate){
case 0:
%s
break;
default:
%s
break;
}
})(1);`, code, code)

	return wrapper
}

// sortStringsByLength 按字符串长度排序
func sortStringsByLength(strs []string, desc bool) {
	sort.Slice(strs, func(i, j int) bool {
		if desc {
			return len(strs[i]) >= len(strs[j])
		}
		return len(strs[i]) < len(strs[j])
	})
}

// SaveSnippet 保存代码片段（支持新增和更新）
// 如果 id > 0 则更新已有记录，否则新增
func (d *DevTools) SaveSnippet(id int64, title, content, language, tags string) (int64, error) {
	if id > 0 {
		// 更新已有代码片段
		_, err := d.db.DB.Exec(
			"UPDATE snippets SET title=?, content=?, language=?, tags=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
			title, content, language, tags, id,
		)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	// 新增代码片段
	result, err := d.db.DB.Exec(
		"INSERT INTO snippets (title, content, language, tags) VALUES (?, ?, ?, ?)",
		title, content, language, tags,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// DeleteSnippet 删除指定代码片段
func (d *DevTools) DeleteSnippet(id int64) error {
	_, err := d.db.DB.Exec("DELETE FROM snippets WHERE id = ?", id)
	return err
}

// SearchSnippets 搜索代码片段（标题、内容、标签）
func (d *DevTools) SearchSnippets(keyword string) ([]Snippet, error) {
	query := "%" + keyword + "%"
	rows, err := d.db.DB.Query(
		"SELECT id, title, content, language, tags, created_at, updated_at FROM snippets WHERE title LIKE ? OR content LIKE ? OR tags LIKE ? ORDER BY updated_at DESC",
		query, query, query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Language, &s.Tags, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		snippets = append(snippets, s)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历搜索结果集失败: %w", err)
	}

	return snippets, nil
}
