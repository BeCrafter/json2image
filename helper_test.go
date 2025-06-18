package json2image

import (
	"os"
	"testing"
)

func TestGetFontFile(t *testing.T) {
	// 测试默认字体
	config := DefaultConfig()
	fontFile, err := getFontFile(config)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	t.Logf("Font file: %s", fontFile)

	// 清理临时文件
	os.Remove(fontFile)
}

func TestGetFontFileWithDifferentTypes(t *testing.T) {
	// 测试不同的字体类型
	fontTypes := []FontType{
		FontTypeMonaco,
		FontTypeMsyh,
		FontTypePingFang,
		FontTypeWrjs,
	}

	for _, fontType := range fontTypes {
		config := DefaultConfig().WithFont(fontType)
		fontFile, err := getFontFile(config)
		if err != nil {
			t.Errorf("Failed to load font type %v: %v", fontType, err)
			continue
		}
		t.Logf("Font type %v loaded: %s", fontType, fontFile)

		// 清理临时文件
		os.Remove(fontFile)
	}
}

func TestGetFontFileWithCustomFont(t *testing.T) {
	// 测试自定义字体（使用不存在的路径）
	config := DefaultConfig().WithCustomFont("/path/to/nonexistent/font.ttf")
	_, err := getFontFile(config)
	if err == nil {
		t.Error("Expected error for nonexistent font file, got nil")
	}
	t.Logf("Expected error for nonexistent font: %v", err)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// 验证默认配置
	if config.Font.Type != FontTypeMsyh {
		t.Errorf("Expected default font type to be FontTypeMsyh, got %v", config.Font.Type)
	}

	if config.Font.Size != 14 {
		t.Errorf("Expected default font size to be 14, got %v", config.Font.Size)
	}

	if config.Font.LineHeight != 20 {
		t.Errorf("Expected default line height to be 20, got %v", config.Font.LineHeight)
	}

	if config.Image.Padding != 20 {
		t.Errorf("Expected default padding to be 20, got %v", config.Image.Padding)
	}
}

func TestConfigChaining(t *testing.T) {
	// 测试配置链式调用
	config := DefaultConfig().
		WithFont(FontTypeMonaco).
		WithFontSize(16).
		WithLineHeight(24).
		WithPadding(30).
		WithBackgroundColor(0.9, 0.9, 0.9).
		WithCropRules("test.rule1", "test.rule2")

	if config.Font.Type != FontTypeMonaco {
		t.Errorf("Expected font type to be FontTypeMonaco, got %v", config.Font.Type)
	}

	if config.Font.Size != 16 {
		t.Errorf("Expected font size to be 16, got %v", config.Font.Size)
	}

	if config.Font.LineHeight != 24 {
		t.Errorf("Expected line height to be 24, got %v", config.Font.LineHeight)
	}

	if config.Image.Padding != 30 {
		t.Errorf("Expected padding to be 30, got %v", config.Image.Padding)
	}

	if len(config.CropRules) != 2 {
		t.Errorf("Expected 2 crop rules, got %v", len(config.CropRules))
	}
}
