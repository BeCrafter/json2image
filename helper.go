package json2image

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/BeCrafter/json2image/fonts"
)

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

// 递归处理嵌套的 JSON 结构
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

// getFontFile 获取字体文件
func getFontFile() (string, error) {
	// 将 base64 字体数据解码为字节
	fontData, err := base64.StdEncoding.DecodeString(fonts.MonacoFontData)
	if err != nil {
		return "", fmt.Errorf("解码字体数据失败: %v", err)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "monaco-*.ttf")
	if err != nil {
		return "", fmt.Errorf("创建临时字体文件失败: %v", err)
	}

	// 写入字体数据
	if _, err := tmpFile.Write(fontData); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("写入字体数据失败: %v", err)
	}

	// 关闭文件
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("关闭临时文件失败: %v", err)
	}

	return tmpFile.Name(), nil
}
