package json2image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

// ColoredLine 带颜色信息的行
type ColoredLine struct {
	text      string
	level     int
	isKey     bool
	startPos  int
	keyLength int
	hasBrace  bool   // 是否包含括号
	bracePos  []int  // 括号的位置
	braceType []rune // 括号的类型
}

// measureText 测量文本尺寸
func measureText(text string, config *Config) (float64, float64) {
	lines := strings.Split(text, "\n")
	maxWidth := 0.0
	dc := gg.NewContext(1, 1)

	fontPath, err := getFontFile(config)
	if err != nil {
		log.Printf("警告: 加载字体失败: %v\n", err)
		return 0, 0
	}

	// 确保在函数结束时清理临时字体文件
	defer func() {
		// 只有当字体不是自定义字体时才删除临时文件
		if config.Font.Type != FontTypeCustom {
			if err := os.Remove(fontPath); err != nil {
				log.Printf("警告: 清理临时字体文件失败: %v", err)
			}
		}
	}()

	if err := dc.LoadFontFace(fontPath, config.Font.Size); err != nil {
		log.Printf("警告: 设置字体失败: %v", err)
		return 0, 0
	}

	for _, line := range lines {
		w, _ := dc.MeasureString(line)
		if w > maxWidth {
			maxWidth = w
		}
	}

	height := float64(len(lines)) * config.Font.LineHeight
	return maxWidth + config.Image.Padding*2, height + config.Image.Padding*2
}

// parseJSONWithColor 解析JSON并添加颜色信息
func parseJSONWithColor(text string) []ColoredLine {
	lines := strings.Split(text, "\n")
	coloredLines := make([]ColoredLine, len(lines))

	for i, line := range lines {
		level := (strings.Count(line, "    "))

		// 查找括号位置，同时记录括号类型
		bracePos := []int{}
		braceType := []rune{}
		for pos, char := range line {
			if char == '{' || char == '}' || char == '[' || char == ']' {
				bracePos = append(bracePos, pos)
				braceType = append(braceType, char)
			}
		}

		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			keyLength := len(strings.Trim(parts[0], `" `))
			startPos := strings.Index(line, `"`)

			coloredLines[i] = ColoredLine{
				text:      line,
				level:     level,
				isKey:     true,
				startPos:  startPos,
				keyLength: keyLength,
				hasBrace:  len(bracePos) > 0,
				bracePos:  bracePos,
				braceType: braceType,
			}
		} else {
			coloredLines[i] = ColoredLine{
				text:      line,
				level:     level,
				isKey:     false,
				hasBrace:  len(bracePos) > 0,
				bracePos:  bracePos,
				braceType: braceType,
			}
		}
	}
	return coloredLines
}

// Json2Image 将JSON数据转换为图片
// 参数：
// - jsonData: JSON字符串
// - config: 配置选项，如果为nil则使用默认配置
// - outputPath: 输出路径（可选），如果不提供则返回base64字符串
func Json2Image(jsonData string, config *Config, outputPath ...string) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 格式化 JSON
	formattedJSON, err := formatJSON(jsonData)
	if err != nil {
		return "", fmt.Errorf("格式化 JSON 失败: %v", err)
	}

	// 解析带颜色信息的行
	coloredLines := parseJSONWithColor(formattedJSON)

	// 计算图片尺寸
	width, height := measureText(formattedJSON, config)

	// 创建画布
	dc := gg.NewContext(int(width), int(height))

	// 加载字体
	fontPath, err := getFontFile(config)
	if err != nil {
		return "", fmt.Errorf("加载字体失败: %v", err)
	}

	// 确保在函数结束时清理临时字体文件
	defer func() {
		// 只有当字体不是自定义字体时才删除临时文件
		if config.Font.Type != FontTypeCustom {
			if err := os.Remove(fontPath); err != nil {
				log.Printf("警告: 清理临时字体文件失败: %v", err)
			}
		}
	}()

	if err := dc.LoadFontFace(fontPath, config.Font.Size); err != nil {
		return "", fmt.Errorf("设置字体失败: %v", err)
	}

	// 设置背景色
	dc.SetRGB(config.Image.BackgroundColor[0], config.Image.BackgroundColor[1], config.Image.BackgroundColor[2])
	dc.Clear()

	// 绘制文本
	y := config.Image.Padding
	for _, line := range coloredLines {
		currentX := config.Image.Padding

		if line.isKey {
			colorIdx := line.level % len(config.Color.LevelColors)
			color := config.Color.LevelColors[colorIdx]
			braceColor := config.Color.BraceLevelColors[colorIdx]

			if line.hasBrace {
				// 绘制前导空格
				if line.startPos > 0 {
					text := line.text[:line.startPos]
					dc.SetRGB(config.Color.DefaultTextColor[0], config.Color.DefaultTextColor[1], config.Color.DefaultTextColor[2])
					dc.DrawString(text, currentX, y)
					width, _ := dc.MeasureString(text)
					currentX += width
				}

				// 绘制字段名（使用正常颜色）
				dc.SetRGB(color[0], color[1], color[2])
				keyText := line.text[line.startPos : line.startPos+line.keyLength+2]
				dc.DrawString(keyText, currentX, y)
				width, _ := dc.MeasureString(keyText)
				currentX += width

				// 绘制冒号和空格
				colonPos := line.startPos + line.keyLength + 2
				text := line.text[colonPos:line.bracePos[0]]
				dc.SetRGB(config.Color.DefaultTextColor[0], config.Color.DefaultTextColor[1], config.Color.DefaultTextColor[2])
				dc.DrawString(text, currentX, y)
				width, _ = dc.MeasureString(text)
				currentX += width

				// 绘制括号（使用浅色）
				dc.SetRGB(braceColor[0], braceColor[1], braceColor[2])
				braceText := string(line.braceType[0])
				dc.DrawString(braceText, currentX, y)
			} else {
				// 原有的键值绘制逻辑保持不变
				if line.startPos > 0 {
					text := line.text[:line.startPos]
					dc.SetRGB(config.Color.DefaultTextColor[0], config.Color.DefaultTextColor[1], config.Color.DefaultTextColor[2])
					dc.DrawString(text, currentX, y)
					width, _ := dc.MeasureString(text)
					currentX += width
				}

				dc.SetRGB(color[0], color[1], color[2])
				keyText := line.text[line.startPos : line.startPos+line.keyLength+2]
				dc.DrawString(keyText, currentX, y)
				width, _ := dc.MeasureString(keyText)
				currentX += width

				if line.startPos+line.keyLength+2 < len(line.text) {
					remainingText := line.text[line.startPos+line.keyLength+2:]
					dc.SetRGB(config.Color.DefaultTextColor[0], config.Color.DefaultTextColor[1], config.Color.DefaultTextColor[2])
					dc.DrawString(remainingText, currentX, y)
				}
			}
		} else {
			// 处理非键值行（可能包含括号）
			if line.hasBrace {
				lastPos := 0
				for i, pos := range line.bracePos {
					// 绘制括号前的文本
					if pos > lastPos {
						dc.SetRGB(config.Color.DefaultTextColor[0], config.Color.DefaultTextColor[1], config.Color.DefaultTextColor[2])
						text := line.text[lastPos:pos]
						dc.DrawString(text, currentX, y)
						width, _ := dc.MeasureString(text)
						currentX += width
					}

					// 使用浅色绘制括号
					colorIdx := line.level % len(config.Color.BraceLevelColors)
					braceColor := config.Color.BraceLevelColors[colorIdx]
					dc.SetRGB(braceColor[0], braceColor[1], braceColor[2])
					braceText := string(line.braceType[i])
					dc.DrawString(braceText, currentX, y)
					width, _ := dc.MeasureString(braceText)
					currentX += width
					lastPos = pos + 1
				}

				// 绘制最后剩余的文本
				if lastPos < len(line.text) {
					dc.SetRGB(config.Color.DefaultTextColor[0], config.Color.DefaultTextColor[1], config.Color.DefaultTextColor[2])
					dc.DrawString(line.text[lastPos:], currentX, y)
				}
			} else {
				dc.SetRGB(config.Color.DefaultTextColor[0], config.Color.DefaultTextColor[1], config.Color.DefaultTextColor[2])
				dc.DrawString(line.text, currentX, y)
			}
		}
		y += config.Font.LineHeight
	}

	if len(outputPath) > 0 {
		// 使用提供的路径保存图片
		return "", dc.SavePNG(outputPath[0])
	}

	// 保存为 base64
	var buf bytes.Buffer
	if err := dc.EncodePNG(&buf); err != nil {
		return "", fmt.Errorf("保存图片为 base64 失败: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// CropJson2Image 将裁剪后的数据转换为图片
func CropJson2Image(jsonData string, config *Config, outputPath ...string) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	var inputData map[string]interface{}
	if str, err := formatJSON(jsonData); err != nil {
		return "", fmt.Errorf("格式化JSON失败: %v", err)
	} else {
		if err := json.Unmarshal([]byte(str), &inputData); err != nil {
			return "", fmt.Errorf("解析输入JSON失败: %v", err)
		}
	}

	if len(config.CropRules) == 0 {
		return "", fmt.Errorf("裁剪规则不能为空")
	}

	output, err := JsonCrop(inputData, config.CropRules)
	if err != nil {
		return "", err
	}

	return Json2Image(string(output), config, outputPath...)
}

// 以下是为了向后兼容而保留的函数，它们使用默认配置

// Json2ImageDefault 使用默认配置将JSON转换为图片（向后兼容）
func Json2ImageDefault(jsonData string, outputPath ...string) (string, error) {
	return Json2Image(jsonData, DefaultConfig(), outputPath...)
}

// CropJson2ImageDefault 使用默认配置将裁剪后的JSON转换为图片（向后兼容）
func CropJson2ImageDefault(jsonData string, rules []string, outputPath ...string) (string, error) {
	config := DefaultConfig().WithCropRules(rules...)
	return CropJson2Image(jsonData, config, outputPath...)
}
