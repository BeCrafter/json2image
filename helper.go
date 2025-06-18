package json2image

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BeCrafter/json2image/fonts"
)

// FontType 表示字体类型
type FontType int

const (
	FontTypeMonaco   FontType = iota // FontTypeMonaco Monaco字体
	FontTypeMsyh                     // FontTypeMsyh 微软雅黑字体
	FontTypePingFang                 // FontTypePingFang 苹方字体
	FontTypeWrjs                     // FontTypeWrjs 王壬金石字体
	FontTypeCustom                   // FontTypeCustom 自定义字体
)

// Config 配置选项
type Config struct {
	Font      FontConfig  // Font 字体配置
	Image     ImageConfig // Image 图片配置
	Color     ColorConfig // Color 颜色配置
	CropRules []string    // CropRules 裁剪规则
}

// FontConfig 字体配置
type FontConfig struct {
	Type       FontType // Type 字体类型
	CustomPath string   // CustomPath 自定义字体文件路径（当Type为FontTypeCustom时使用）
	Size       float64  // Size 字体大小
	LineHeight float64  // LineHeight 行高
}

// ImageConfig 图片配置
type ImageConfig struct {
	Padding         float64    // Padding 内边距
	BackgroundColor [3]float64 // BackgroundColor 背景色
}

// ColorConfig 颜色配置
type ColorConfig struct {
	LevelColors      [][3]float64 // LevelColors 各层级的颜色
	BraceLevelColors [][3]float64 // BraceLevelColors 括号的颜色
	DefaultTextColor [3]float64   // DefaultTextColor 默认文本颜色
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Font: FontConfig{
			Type:       FontTypeMsyh,
			Size:       14,
			LineHeight: 20,
		},
		Image: ImageConfig{
			Padding:         20,
			BackgroundColor: [3]float64{1, 1, 1}, // 白色
		},
		Color: ColorConfig{
			LevelColors: [][3]float64{
				{0.2, 0.6, 0.9}, // 蓝色
				{0.8, 0.3, 0.3}, // 红色
				{0.3, 0.7, 0.3}, // 绿色
				{0.7, 0.3, 0.7}, // 紫色
				{0.9, 0.6, 0.2}, // 橙色
				{0.2, 0.7, 0.7}, // 青色
				{0.7, 0.7, 0.2}, // 黄色
				{0.5, 0.2, 0.8}, // 深紫色
				{0.8, 0.4, 0.6}, // 粉色
				{0.4, 0.5, 0.3}, // 橄榄绿
			},
			BraceLevelColors: [][3]float64{
				{0.5, 0.8, 1.0}, // 浅蓝色
				{1.0, 0.6, 0.6}, // 浅红色
				{0.6, 0.9, 0.6}, // 浅绿色
				{0.9, 0.6, 0.9}, // 浅紫色
				{1.0, 0.8, 0.5}, // 浅橙色
				{0.5, 0.9, 0.9}, // 浅青色
				{0.9, 0.9, 0.5}, // 浅黄色
				{0.7, 0.5, 0.9}, // 浅深紫色
				{0.9, 0.7, 0.8}, // 浅粉色
				{0.7, 0.8, 0.6}, // 浅橄榄绿
			},
			DefaultTextColor: [3]float64{0, 0, 0}, // 黑色
		},
	}
}

// WithFont 设置字体类型
func (c *Config) WithFont(fontType FontType) *Config {
	c.Font.Type = fontType
	return c
}

// WithCustomFont 设置自定义字体
func (c *Config) WithCustomFont(fontPath string) *Config {
	c.Font.Type = FontTypeCustom
	c.Font.CustomPath = fontPath
	return c
}

// WithFontSize 设置字体大小
func (c *Config) WithFontSize(size float64) *Config {
	c.Font.Size = size
	return c
}

// WithLineHeight 设置行高
func (c *Config) WithLineHeight(lineHeight float64) *Config {
	c.Font.LineHeight = lineHeight
	return c
}

// WithPadding 设置内边距
func (c *Config) WithPadding(padding float64) *Config {
	c.Image.Padding = padding
	return c
}

// WithBackgroundColor 设置背景色
func (c *Config) WithBackgroundColor(r, g, b float64) *Config {
	c.Image.BackgroundColor = [3]float64{r, g, b}
	return c
}

// WithLevelColors 设置层级颜色
func (c *Config) WithLevelColors(colors [][3]float64) *Config {
	c.Color.LevelColors = colors
	return c
}

// WithBraceLevelColors 设置括号颜色
func (c *Config) WithBraceLevelColors(colors [][3]float64) *Config {
	c.Color.BraceLevelColors = colors
	return c
}

// WithDefaultTextColor 设置默认文本颜色
func (c *Config) WithDefaultTextColor(r, g, b float64) *Config {
	c.Color.DefaultTextColor = [3]float64{r, g, b}
	return c
}

// WithCropRules 设置裁剪规则
func (c *Config) WithCropRules(rules ...string) *Config {
	c.CropRules = rules
	return c
}

// formatJSON 格式化JSON字符串
func formatJSON(data string) (string, error) {
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(data), &jsonObj); err != nil {
		return "", err
	}

	// 递归处理 JSON 对象
	processedObj := processNestedJSON(jsonObj)

	// 重新格式化整个 JSON
	prettyJSON, err := json.MarshalIndent(processedObj, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

// processNestedJSON 递归处理嵌套的 JSON 结构
func processNestedJSON(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		// 处理对象
		m := make(map[string]interface{})
		for key, value := range v {
			m[key] = processNestedJSON(value)
		}
		return m
	case []interface{}:
		// 处理数组
		a := make([]interface{}, len(v))
		for i, value := range v {
			a[i] = processNestedJSON(value)
		}
		return a
	case string:
		// 尝试解析字符串值是否为 JSON
		var nestedJSON interface{}
		if err := json.Unmarshal([]byte(v), &nestedJSON); err == nil {
			// 如果是有效的 JSON，则递归处理
			return processNestedJSON(nestedJSON)
		}
		return v
	default:
		return v
	}
}

// getFontFile 根据配置获取字体文件路径
func getFontFile(config *Config) (string, error) {
	var fontData string
	var err error

	switch config.Font.Type {
	case FontTypeMonaco:
		fontData = fonts.MonacoFontData
	case FontTypeMsyh:
		fontData = fonts.MsyhFontData
	case FontTypePingFang:
		fontData = fonts.PingfangscFontData
	case FontTypeWrjs:
		fontData = fonts.WrjsFontData
	case FontTypeCustom:
		if config.Font.CustomPath == "" {
			return "", fmt.Errorf("自定义字体路径不能为空")
		}

		// 检查文件是否存在
		if _, err := os.Stat(config.Font.CustomPath); os.IsNotExist(err) {
			return "", fmt.Errorf("字体文件不存在: %s", config.Font.CustomPath)
		}

		// 检查文件扩展名
		ext := filepath.Ext(config.Font.CustomPath)
		if ext != ".ttf" && ext != ".otf" {
			return "", fmt.Errorf("不支持的字体文件格式: %s", ext)
		}

		return config.Font.CustomPath, nil
	default:
		return "", fmt.Errorf("未知的字体类型: %d", config.Font.Type)
	}

	if fontData == "" {
		return "", fmt.Errorf("字体数据为空")
	}

	// 将 base64 字体数据解码为字节
	decodedData, err := base64.StdEncoding.DecodeString(fontData)
	if err != nil {
		return "", fmt.Errorf("解码字体数据失败: %v", err)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "font-*.ttf")
	if err != nil {
		return "", fmt.Errorf("创建临时字体文件失败: %v", err)
	}

	// 写入字体数据
	if _, err := tmpFile.Write(decodedData); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("写入字体数据失败: %v", err)
	}

	// 关闭文件
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("关闭临时文件失败: %v", err)
	}

	return tmpFile.Name(), nil
}
