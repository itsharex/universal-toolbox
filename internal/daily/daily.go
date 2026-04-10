// Package daily 日常工具模块
// 提供计算器、单位换算、汇率换算、备忘录等日常实用工具
package daily

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unicode"
	"xtool/internal/db"
)

// DailyTools 日常工具结构体（Wails 绑定到前端）
type DailyTools struct {
	db *db.Database // 数据库连接
}

// Note 备忘录结构
type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Color     string `json:"color"`
	Pinned    bool   `json:"pinned"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ConversionResult 单位转换结果
type ConversionResult struct {
	Value   float64 `json:"value"`   // 转换后的值
	Unit    string  `json:"unit"`    // 目标单位
	Formula string  `json:"formula"` // 转换公式说明
}

// NewDailyTools 创建日常工具模块实例
func NewDailyTools(database *db.Database) *DailyTools {
	return &DailyTools{db: database}
}

// ============================================================
// 计算器（标准模式 + 科学模式）
// ============================================================

// CalcBasic 基础四则运算
// op: +, -, *, /
func (d *DailyTools) CalcBasic(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("不支持的运算符: %s", op)
	}
}

// CalcScientific 科学计算函数
// func: sqrt/pow/log/ln/sin/cos/tan/abs/ceil/floor/round
func (d *DailyTools) CalcScientific(value, param float64, fn string) (float64, error) {
	switch fn {
	case "sqrt":
		if value < 0 {
			return 0, fmt.Errorf("负数不能开平方根")
		}
		return math.Sqrt(value), nil
	case "pow":
		return math.Pow(value, param), nil
	case "log":
		if value <= 0 {
			return 0, fmt.Errorf("对数的真数必须大于 0")
		}
		return math.Log10(value), nil
	case "ln":
		if value <= 0 {
			return 0, fmt.Errorf("自然对数的真数必须大于 0")
		}
		return math.Log(value), nil
	case "sin":
		return math.Sin(value * math.Pi / 180), nil // 角度转弧度
	case "cos":
		return math.Cos(value * math.Pi / 180), nil
	case "tan":
		// tan 函数在 90 度和 270 度附近趋向无穷大，添加边界处理
		// 将角度归一化到 0-360 范围
		normalized := math.Mod(value, 360)
		if normalized < 0 {
			normalized += 360
		}
		// 检查是否接近 90 度或 270 度（奇数倍的 90 度）
		if math.Mod(normalized, 180) == 90 {
			return 0, fmt.Errorf("tan(%.1f 度) 无定义（趋向无穷大）", value)
		}
		// 接近 90 度或 270 度时，检查是否会导致极大值
		radians := normalized * math.Pi / 180
		result := math.Tan(radians)
		if math.IsInf(result, 0) || math.IsNaN(result) {
			return 0, fmt.Errorf("tan(%.1f 度) 无定义（趋向无穷大）", value)
		}
		return result, nil
	case "abs":
		return math.Abs(value), nil
	case "ceil":
		return math.Ceil(value), nil
	case "floor":
		return math.Floor(value), nil
	case "round":
		return math.Round(value), nil
	case "pi":
		return math.Pi * value, nil
	default:
		return 0, fmt.Errorf("不支持的函数: %s", fn)
	}
}

// CalcExpression 安全表达式求值（替代前端 Function()）
// 支持基本四则运算、括号、数学常量（PI、E）
// 不使用 eval/exec，通过安全的词法分析实现
func (d *DailyTools) CalcExpression(expr string) (float64, error) {
	// 预处理：去除空格
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, fmt.Errorf("表达式不能为空")
	}

	// 替换数学常量
	expr = strings.ReplaceAll(expr, "PI", fmt.Sprintf("%v", math.Pi))
	expr = strings.ReplaceAll(expr, "E", fmt.Sprintf("%v", math.E))

	// 安全检查：只允许数字、运算符、括号和小数点
	for _, ch := range expr {
		if !unicode.IsDigit(ch) && ch != '+' && ch != '-' && ch != '*' && ch != '/' &&
			ch != '(' && ch != ')' && ch != '.' && ch != ' ' {
			return 0, fmt.Errorf("表达式包含非法字符: %c（仅支持数字和 +-*/()）", ch)
		}
	}

	// 使用安全的方式解析表达式
	result, err := safeEval(expr)
	if err != nil {
		return 0, fmt.Errorf("表达式计算失败: %v", err)
	}
	return result, nil
}

// safeEval 安全的表达式求值器
// 支持加减乘除和括号，使用递归下降解析
func safeEval(expr string) (float64, error) {
	// 词法分析：将表达式拆分为 token
	tokens := tokenize(expr)
	pos := 0

	// 解析加法/减法（最低优先级）
	parseAddSub := func() (float64, error) {
		left, err := parseMulDiv(tokens, &pos)
		if err != nil {
			return 0, err
		}
		for pos < len(tokens) && (tokens[pos] == "+" || tokens[pos] == "-") {
			op := tokens[pos]
			pos++
			right, err := parseMulDiv(tokens, &pos)
			if err != nil {
				return 0, err
			}
			if op == "+" {
				left += right
			} else {
				left -= right
			}
		}
		return left, nil
	}

	result, err := parseAddSub()
	if err != nil {
		return 0, err
	}
	if pos != len(tokens) {
		return 0, fmt.Errorf("表达式解析不完整，剩余部分: %s", strings.Join(tokens[pos:], " "))
	}
	return result, nil
}

// parseMulDiv 解析乘除法
func parseMulDiv(tokens []string, pos *int) (float64, error) {
	left, err := parsePrimary(tokens, pos)
	if err != nil {
		return 0, err
	}
	for *pos < len(tokens) && (tokens[*pos] == "*" || tokens[*pos] == "/") {
		op := tokens[*pos]
		*pos++
		right, err := parsePrimary(tokens, pos)
		if err != nil {
			return 0, err
		}
		if op == "*" {
			left *= right
		} else {
			if right == 0 {
				return 0, fmt.Errorf("除数不能为零")
			}
			left /= right
		}
	}
	return left, nil
}

// parsePrimary 解析基本元素（数字或括号表达式）
func parsePrimary(tokens []string, pos *int) (float64, error) {
	if *pos >= len(tokens) {
		return 0, fmt.Errorf("表达式不完整")
	}

	token := tokens[*pos]

	// 处理负号（一元运算符）
	if token == "-" {
		*pos++
		val, err := parsePrimary(tokens, pos)
		if err != nil {
			return 0, err
		}
		return -val, nil
	}

	// 处理正号（一元运算符）
	if token == "+" {
		*pos++
		return parsePrimary(tokens, pos)
	}

	// 处理括号
	if token == "(" {
		*pos++ // 跳过 '('
		val, err := parseAddSub(tokens, pos)
		if err != nil {
			return 0, err
		}
		if *pos >= len(tokens) || tokens[*pos] != ")" {
			return 0, fmt.Errorf("缺少右括号")
		}
		*pos++ // 跳过 ')'
		return val, nil
	}

	// 处理数字
	val, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, fmt.Errorf("无法解析数字: %s", token)
	}
	*pos++
	return val, nil
}

// tokenize 将表达式字符串拆分为 token 列表
func tokenize(expr string) []string {
	var tokens []string
	var current strings.Builder

	for _, ch := range expr {
		if ch == ' ' {
			continue
		}
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '(' || ch == ')' {
			// 如果之前有累积的数字，先添加
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		} else {
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// ============================================================
// 单位换算工具
// ============================================================

// ConvertLength 长度单位换算
// fromUnit/toUnit: mm, cm, m, km, inch, foot, yard, mile
func (d *DailyTools) ConvertLength(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为米（基准单位）
	toMeter := map[string]float64{
		"mm":   0.001,
		"cm":   0.01,
		"m":    1.0,
		"km":   1000.0,
		"inch": 0.0254,
		"foot": 0.3048,
		"yard": 0.9144,
		"mile": 1609.344,
	}

	fromFactor, ok1 := toMeter[fromUnit]
	toFactor, ok2 := toMeter[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的长度单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertWeight 重量单位换算
// fromUnit/toUnit: mg, g, kg, t, oz, lb
func (d *DailyTools) ConvertWeight(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为克（基准单位）
	toGram := map[string]float64{
		"mg": 0.001,
		"g":  1.0,
		"kg": 1000.0,
		"t":  1000000.0,
		"oz": 28.3495,
		"lb": 453.592,
	}

	fromFactor, ok1 := toGram[fromUnit]
	toFactor, ok2 := toGram[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的重量单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertTemperature 温度换算
// fromUnit/toUnit: C（摄氏度）、F（华氏度）、K（开尔文）
func (d *DailyTools) ConvertTemperature(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 先转换为摄氏度
	var celsius float64
	switch strings.ToUpper(fromUnit) {
	case "C":
		celsius = value
	case "F":
		celsius = (value - 32) * 5 / 9
	case "K":
		celsius = value - 273.15
	default:
		return ConversionResult{}, fmt.Errorf("不支持的温度单位（支持: C/F/K）")
	}

	// 再从摄氏度转为目标单位
	var result float64
	switch strings.ToUpper(toUnit) {
	case "C":
		result = celsius
	case "F":
		result = celsius*9/5 + 32
	case "K":
		result = celsius + 273.15
	default:
		return ConversionResult{}, fmt.Errorf("不支持的温度单位（支持: C/F/K）")
	}

	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g%s = %g%s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertSpeed 速度单位换算
// fromUnit/toUnit: ms（米/秒）、kmh（千米/时）、mph（英里/时）、knot（节）
func (d *DailyTools) ConvertSpeed(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为米/秒
	toMPS := map[string]float64{
		"ms":   1.0,
		"kmh":  1.0 / 3.6,
		"mph":  0.44704,
		"knot": 0.514444,
	}

	fromFactor, ok1 := toMPS[fromUnit]
	toFactor, ok2 := toMPS[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的速度单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertArea 面积单位换算
// fromUnit/toUnit: mm2, cm2, m2, km2, ha（公顷）, acre（英亩）, mu（亩）
func (d *DailyTools) ConvertArea(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为平方米（基准单位）
	toSqm := map[string]float64{
		"mm2":  1e-6,
		"cm2":  1e-4,
		"m2":   1.0,
		"km2":  1e6,
		"ha":   1e4,      // 公顷
		"acre": 4046.8564, // 英亩
		"mu":   666.6667,  // 亩
	}

	fromFactor, ok1 := toSqm[fromUnit]
	toFactor, ok2 := toSqm[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的面 积单位（支持: mm2/cm2/m2/km2/ha/acre/mu）")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertVolume 体积单位换算
// fromUnit/toUnit: ml, l（升）, m3, cm3, gallon（加仑）, oz_fl（液量盎司）
func (d *DailyTools) ConvertVolume(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为升（基准单位）
	toLiter := map[string]float64{
		"ml":    0.001,
		"l":     1.0,
		"m3":    1000.0,
		"cm3":   0.001,
		"gallon": 3.78541, // 美制加仑
		"oz_fl": 0.0295735, // 液量盎司
	}

	fromFactor, ok1 := toLiter[fromUnit]
	toFactor, ok2 := toLiter[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的体积单位（支持: ml/l/m3/cm3/gallon/oz_fl）")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertDataSize 数据存储单位换算
// fromUnit/toUnit: b（字节）, kb, mb, gb, tb, pb, bit（位）
func (d *DailyTools) ConvertDataSize(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为字节（基准单位）
	toByte := map[string]float64{
		"bit": 0.125,  // 1 字节 = 8 位
		"b":   1.0,
		"kb":  1024.0,
		"mb":  1024 * 1024,
		"gb":  1024 * 1024 * 1024,
		"tb":  1024.0 * 1024 * 1024 * 1024,
		"pb":  1024.0 * 1024 * 1024 * 1024 * 1024,
	}

	fromFactor, ok1 := toByte[fromUnit]
	toFactor, ok2 := toByte[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的数据存储单位（支持: bit/b/kb/mb/gb/tb/pb）")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:   result,
		Unit:    strings.ToUpper(toUnit),
		Formula: fmt.Sprintf("%g %s = %g %s", value, strings.ToUpper(fromUnit), result, strings.ToUpper(toUnit)),
	}, nil
}

// ============================================================
// 时间日期工具
// ============================================================

// GetCurrentTime 获取当前时间的多种格式表示
func (d *DailyTools) GetCurrentTime() map[string]string {
	now := time.Now()
	return map[string]string{
		"datetime":  now.Format("2006-01-02 15:04:05"),
		"date":      now.Format("2006-01-02"),
		"time":      now.Format("15:04:05"),
		"weekday":   now.Weekday().String(),
		"timestamp": fmt.Sprintf("%d", now.Unix()),
		"utc":       now.UTC().Format("2006-01-02 15:04:05 UTC"),
		"timezone":  now.Format("MST"),
	}
}

// ============================================================
// 备忘录工具
// ============================================================

// GetNotes 获取所有备忘录
func (d *DailyTools) GetNotes() ([]Note, error) {
	rows, err := d.db.DB.Query(
		"SELECT id, title, content, color, pinned, created_at, updated_at FROM notes ORDER BY pinned DESC, updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		var pinnedInt int
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Color, &pinnedInt, &n.CreatedAt, &n.UpdatedAt); err != nil {
			continue
		}
		n.Pinned = pinnedInt == 1
		notes = append(notes, n)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历备忘录结果集失败: %w", err)
	}

	return notes, nil
}

// SaveNote 保存备忘录（新增，支持同时设置 pinned 状态）
func (d *DailyTools) SaveNote(title, content, color string, pinned bool) (int64, error) {
	pinnedInt := 0
	if pinned {
		pinnedInt = 1
	}
	result, err := d.db.DB.Exec(
		"INSERT INTO notes (title, content, color, pinned) VALUES (?, ?, ?, ?)",
		title, content, color, pinnedInt,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateNote 更新备忘录内容
func (d *DailyTools) UpdateNote(id int64, title, content, color string) error {
	_, err := d.db.DB.Exec(
		"UPDATE notes SET title=?, content=?, color=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		title, content, color, id,
	)
	return err
}

// PinNote 切换备忘录置顶状态
func (d *DailyTools) PinNote(id int64) error {
	_, err := d.db.DB.Exec(
		"UPDATE notes SET pinned = CASE WHEN pinned = 1 THEN 0 ELSE 1 END WHERE id = ?",
		id,
	)
	return err
}

// DeleteNote 删除备忘录
func (d *DailyTools) DeleteNote(id int64) error {
	_, err := d.db.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}

// ============================================================
// 密码生成工具
// ============================================================

// PasswordResult 密码生成结果
type PasswordResult struct {
	Success  bool   `json:"success"`  // 是否成功
	Password string `json:"password"` // 生成的密码
	Strength string `json:"strength"` // 密码强度（弱/中/强/极强）
	Error    string `json:"error"`    // 错误信息
}

// GeneratePassword 生成随机密码
// length: 密码长度 4-128
// useUpper: 是否包含大写字母
// useLower: 是否包含小写字母
// useNumbers: 是否包含数字
// useSymbols: 是否包含特殊符号
// excludeChars: 需要排除的字符
func (d *DailyTools) GeneratePassword(length int, useUpper, useLower, useNumbers, useSymbols bool, excludeChars string) PasswordResult {
	result := PasswordResult{Success: true}

	// 校验密码长度
	if length < 4 {
		length = 4
	}
	if length > 128 {
		length = 128
	}

	// 至少选择一种字符类型
	if !useUpper && !useLower && !useNumbers && !useSymbols {
		return PasswordResult{Success: false, Error: "至少需要选择一种字符类型"}
	}

	// 构建字符池
	var pool string
	var requiredChars []string // 每种选中类型至少一个字符

	if useUpper {
		upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		pool += upperChars
		// 随机选一个大写字母作为必需字符
		idx := cryptoRandInt(len(upperChars))
		requiredChars = append(requiredChars, string(upperChars[idx]))
	}
	if useLower {
		lowerChars := "abcdefghijklmnopqrstuvwxyz"
		pool += lowerChars
		idx := cryptoRandInt(len(lowerChars))
		requiredChars = append(requiredChars, string(lowerChars[idx]))
	}
	if useNumbers {
		numberChars := "0123456789"
		pool += numberChars
		idx := cryptoRandInt(len(numberChars))
		requiredChars = append(requiredChars, string(numberChars[idx]))
	}
	if useSymbols {
		symbolChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
		pool += symbolChars
		idx := cryptoRandInt(len(symbolChars))
		requiredChars = append(requiredChars, string(symbolChars[idx]))
	}

	// 从字符池中移除排除字符
	if excludeChars != "" {
		filtered := strings.Builder{}
		for _, ch := range pool {
			if !strings.ContainsRune(excludeChars, ch) {
				filtered.WriteRune(ch)
			}
		}
		pool = filtered.String()
		if pool == "" {
			return PasswordResult{Success: false, Error: "排除字符后可用字符池为空"}
		}
		// 重新生成必需字符（确保不在排除列表中）
		var newRequired []string
		for _, rc := range requiredChars {
			if !strings.ContainsAny(rc, excludeChars) {
				newRequired = append(newRequired, rc)
			} else {
				// 从剩余字符池中重新选取
				if len(pool) > 0 {
					idx := cryptoRandInt(len(pool))
					newRequired = append(newRequired, string(pool[idx]))
				}
			}
		}
		requiredChars = newRequired
	}

	// 生成密码
	password := make([]byte, 0, length)

	// 先添加必需字符（确保每种类型至少一个）
	password = append(password, []byte(strings.Join(requiredChars, ""))...)

	// 填充剩余长度
	remaining := length - len(password)
	for i := 0; i < remaining; i++ {
		idx := cryptoRandInt(len(pool))
		password = append(password, pool[idx])
	}

	// 打乱密码字符顺序（Fisher-Yates 洗牌算法）
	for i := len(password) - 1; i > 0; i-- {
		j := cryptoRandInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	result.Password = string(password)

	// 评估密码强度
	result.Strength = evaluatePasswordStrength(result.Password)

	return result
}

// evaluatePasswordStrength 评估密码强度
func evaluatePasswordStrength(password string) string {
	length := len(password)
	score := 0

	// 长度评分
	if length >= 8 {
		score++
	}
	if length >= 12 {
		score++
	}
	if length >= 16 {
		score++
	}

	// 字符类型评分
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSymbol := false
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		default:
			hasSymbol = true
		}
	}

	types := 0
	if hasUpper {
		types++
	}
	if hasLower {
		types++
	}
	if hasNumber {
		types++
	}
	if hasSymbol {
		types++
	}
	score += types

	// 评定强度
	switch {
	case score >= 6:
		return "极强"
	case score >= 4:
		return "强"
	case score >= 3:
		return "中"
	default:
		return "弱"
	}
}

// cryptoRandInt 生成安全的随机整数 [0, max)
func cryptoRandInt(max int) int {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// 降级到非安全随机（不应该发生）
		return 0
	}
	return int(n.Int64())
}
