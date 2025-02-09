package json2image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fogleman/gg"
)

// 修改 measureText 和 jsonToImage 函数中的字体加载部分
func measureText(text string) (float64, float64) {
	lines := strings.Split(text, "\n")
	maxWidth := 0.0
	dc := gg.NewContext(1, 1)

	fontPath, err := getFontFile()
	if err != nil {
		log.Fatalf("警告: 加载字体失败: %v\n", err)
		return 0, 0
	}

	dc.LoadFontFace(fontPath, 14)

	for _, line := range lines {
		w, _ := dc.MeasureString(line)
		if w > maxWidth {
			maxWidth = w
		}
	}

	height := float64(len(lines)) * 20 // 每行高度20像素
	return maxWidth + 40, height + 40  // 添加边距
}

// 添加颜色配置
var levelColors = [][]float64{
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
}

// 添加括号颜色配置（比对应层级颜色更浅）
var braceLevelColors = [][]float64{
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
}

// 修改 parseJSONWithColor 函数中的括号处理逻辑
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

func parseJSONWithColor(text string) []ColoredLine {
	lines := strings.Split(text, "\n")
	coloredLines := make([]ColoredLine, len(lines))

	for i, line := range lines {
		level := (strings.Count(line, "    "))

		// 查找括号位置，同时记录括号类型
		bracePos := []int{}
		braceType := []rune{}
		// 直接在原始行中查找括号，不使用 trimmed
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

// 在 Json2Image 函数中修改绘制括号的代码
// 将数据按照Base64保存：JsonToImage(jsonData)
// 将数据按照图片格式输出：JsonToImage(jsonData, "output.png")
func Json2Image(jsonData string, outputPath ...string) (string, error) {
	// 格式化 JSON
	formattedJSON, err := formatJSON(jsonData)
	if err != nil {
		return "", fmt.Errorf("格式化 JSON 失败: %v", err)
	}

	// 解析带颜色信息的行
	coloredLines := parseJSONWithColor(formattedJSON)

	// 计算图片尺寸
	width, height := measureText(formattedJSON)

	// 创建画布
	dc := gg.NewContext(int(width), int(height))

	// 设置背景色
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// 绘制文本
	y := 20.0
	for _, line := range coloredLines {
		currentX := 20.0

		if line.isKey {
			colorIdx := line.level % len(levelColors)
			color := levelColors[colorIdx]
			braceColor := braceLevelColors[colorIdx]

			if line.hasBrace {
				// 绘制前导空格
				if line.startPos > 0 {
					text := line.text[:line.startPos]
					dc.SetRGB(0, 0, 0)
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
				dc.SetRGB(0, 0, 0)
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
					dc.SetRGB(0, 0, 0)
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
					dc.SetRGB(0, 0, 0)
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
						dc.SetRGB(0, 0, 0)
						text := line.text[lastPos:pos]
						dc.DrawString(text, currentX, y)
						width, _ := dc.MeasureString(text)
						currentX += width
					}

					// 使用浅色绘制括号
					colorIdx := line.level % len(braceLevelColors)
					braceColor := braceLevelColors[colorIdx]
					dc.SetRGB(braceColor[0], braceColor[1], braceColor[2])
					braceText := string(line.braceType[i])
					dc.DrawString(braceText, currentX, y)
					width, _ := dc.MeasureString(braceText)
					currentX += width
					lastPos = pos + 1
				}

				// 绘制最后剩余的文本
				if lastPos < len(line.text) {
					dc.SetRGB(0, 0, 0)
					dc.DrawString(line.text[lastPos:], currentX, y)
				}
			} else {
				dc.SetRGB(0, 0, 0)
				dc.DrawString(line.text, currentX, y)
			}
		}
		y += 20
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
func CropJson2Image(jsonData string, rules []string, outputPath ...string) (string, error) {
	var inputData map[string]interface{}
	if str, err := formatJSON(jsonData); err != nil {
		return "", fmt.Errorf("格式化JSON失败: %v", err)
	} else {
		if err := json.Unmarshal([]byte(str), &inputData); err != nil {
			return "", fmt.Errorf("解析输入JSON失败: %v", err)
		}
	}

	output, err := JsonCrop(inputData, rules)
	if err != nil {
		return "", err
	}

	return Json2Image(string(output), outputPath...)
}
