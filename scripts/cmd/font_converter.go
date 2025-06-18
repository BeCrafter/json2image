package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("\nUsage: go run font_converter.go <font_name>")
		fmt.Println("\nAvailable fonts:")
		showAvailableFonts()
		os.Exit(1)
	}

	fontName := os.Args[1]

	// 获取当前工作目录
	workPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting work path: %v\n", err)
		os.Exit(1)
	}
	rootPath := filepath.Dir(workPath)

	// 生成文件名
	fontNameLower := strings.ToLower(fontName)
	fontNameFirst := titleCase(fontName)

	inputFile := filepath.Join(workPath, "fonts", fontName+".ttf")
	outputFile := filepath.Join(rootPath, "fonts", fontNameLower+".go")

	// 转换字体文件
	err = convertFont(inputFile, outputFile, fontNameFirst)
	if err != nil {
		fmt.Printf("Error converting font: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func showAvailableFonts() {
	workPath, err := os.Getwd()
	if err != nil {
		return
	}
	fontsDir := filepath.Join(workPath, "fonts")
	files, err := os.ReadDir(fontsDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".ttf" {
			fmt.Printf("    %s\n", strings.TrimSuffix(file.Name(), ".ttf"))
		}
	}
}

func convertFont(inputFile, outputFile, fontName string) error {
	// 读取字体文件
	fontData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read font file: %v", err)
	}

	// 创建输出文件
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// 编码为base64
	encoded := base64.StdEncoding.EncodeToString(fontData)

	// 对于大文件，拆分成多个常量避免编译器嵌套深度限制
	const maxConstSize = 50000 // 每个常量最大50KB

	// 写入包头
	writer.WriteString("package fonts\n\n")

	if len(encoded) <= maxConstSize {
		// 小文件，使用单个常量
		writer.WriteString(fmt.Sprintf("const %sFontData = ", fontName))
		processEncodedData(writer, encoded, 150, 150-13-len(fontName))
		writer.WriteString("``\n")
	} else {
		// 大文件，拆分成多个常量
		parts := (len(encoded) + maxConstSize - 1) / maxConstSize

		// 生成各个部分的常量
		for i := 0; i < parts; i++ {
			start := i * maxConstSize
			end := start + maxConstSize
			if end > len(encoded) {
				end = len(encoded)
			}

			fontnameVar := fmt.Sprintf("%sFontData%d", fontName, i)
			writer.WriteString(fmt.Sprintf("const %s = `", fontnameVar))
			// 直接写入数据，不使用processEncodedData
			chunk := encoded[start:end]
			chunkLen := 150
			pos := 0
			for pos < len(chunk) {
				lineEnd := pos + chunkLen - len(fontnameVar) - 5
				if pos > 0 {
					lineEnd = pos + chunkLen
				}
				if lineEnd > len(chunk) {
					lineEnd = len(chunk)
				}
				if pos > 0 {
					writer.WriteString(" +\n `")
				}
				writer.WriteString(chunk[pos:lineEnd])
				writer.WriteString("`")
				pos = lineEnd
			}
			writer.WriteString("\n\n")
		}

		// 生成主常量，连接所有部分
		writer.WriteString(fmt.Sprintf("const %sFontData = ", fontName))
		for i := 0; i < parts; i++ {
			if i > 0 {
				writer.WriteString(" + ")
			}
			writer.WriteString(fmt.Sprintf("%sFontData%d", fontName, i))
		}
		writer.WriteString("\n")
	}

	return nil
}

func processEncodedData(writer *bufio.Writer, encoded string, lineLen, firstLineLen int) {
	pos := 0
	dataLen := len(encoded)
	isFirst := true

	for pos < dataLen {
		var chunkLen int
		if isFirst {
			chunkLen = firstLineLen
			isFirst = false
		} else {
			chunkLen = lineLen
		}

		end := pos + chunkLen
		if end > dataLen {
			end = dataLen
		}

		chunk := encoded[pos:end]
		if end >= dataLen {
			// 最后一块，不加 +
			writer.WriteString(fmt.Sprintf(" `%s` ", chunk))
		} else {
			writer.WriteString(fmt.Sprintf(" `%s` +\n", chunk))
		}
		pos = end
	}
}
