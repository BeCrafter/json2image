package json2image

import "testing"

func TestGetFontFile(t *testing.T) {
	fontFile, err := getFontFile()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	t.Logf("Font file: %s", fontFile)
}
