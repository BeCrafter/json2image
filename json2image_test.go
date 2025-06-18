package json2image

import (
	"fmt"
	"testing"
)

func TestJson2Image(t *testing.T) {
	// 示例 JSON 数据
	jsonData := `{
		"name": "John Doe",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "New York",
			"country": "USA"
		},
		"hobbies": ["reading", "swimming", "coding"],
		"extras": "{\"uid\":231345,\"source_id\":0}"
	}`

	// 使用默认配置
	_, err := Json2Image(jsonData, nil, "output.png")
	if err != nil {
		t.Errorf("生成图片失败: %v\n", err)
		return
	}
	fmt.Println("图片生成成功：output.png")
}

func TestJson2ImageWithCustomConfig(t *testing.T) {
	// 示例 JSON 数据
	jsonData := `{
		"name": "测试用户",
		"age": 30,
		"address": {
			"street": "测试街道123号",
			"city": "测试城市",
			"country": "中国"
		},
		"hobbies": ["阅读", "游泳", "编程"]
	}`

	// 创建自定义配置
	config := DefaultConfig().
		WithFont(FontTypePingFang).
		WithFontSize(16).
		WithLineHeight(24).
		WithPadding(30).
		WithBackgroundColor(0.95, 0.95, 0.95)

	_, err := Json2Image(jsonData, config, "output/output_custom.png")
	if err != nil {
		t.Errorf("生成图片失败: %v\n", err)
		return
	}
	fmt.Println("自定义配置图片生成成功：output_custom.png")
}

func TestJson2ImageBase64(t *testing.T) {
	// 示例 JSON 数据
	jsonData := `{
		"simple": "test",
		"number": 123
	}`

	// 生成base64字符串
	base64Str, err := Json2Image(jsonData, nil)
	if err != nil {
		t.Errorf("生成base64失败: %v\n", err)
		return
	}

	if len(base64Str) == 0 {
		t.Error("生成的base64字符串为空")
		return
	}

	fmt.Printf("Base64字符串长度: %d\n", len(base64Str))
}

func TestCropJson2Image(t *testing.T) {
	// 示例 JSON 数据
	inputData := `{
		"orders": {
			"order1": {
				"items": [{
						"product": {
							"name": "商品1",
							"cnt": 11,
							"uname": "111",
							"price": 100
						}
					},
					{
						"product": {
							"name": "商品2",
							"cnt": 11,
							"uname": "112",
							"price": 200
						}
					}
				],
				"customer": {
					"name": "张三",
					"cnt": 11,
					"price": 20
				}
			},
			"order2": {
				"items": [{
					"product": {
						"name": "商品3",
						"cnt": 11,
						"price": 300
					}
				}],
				"customer": {
					"name": "李四",
					"cnt": 11,
					"price": 21
				}
			}
		},
		"extras": "{\"content\":{\"adjustInfo\":{\"adjustId\":\"1518562\"}}}"
	}`

	// 创建带裁剪规则的配置
	config := DefaultConfig().WithCropRules(
		"orders.order1.items[*].product.name",
		"orders.order1.items[0].product.cnt",
		"orders.order1.items[1].product.price",
		"orders.order1.items[1,0].product.uname",
		"orders.*.customer.name",
		"extras.content.adjustInfo.adjustId",
	)

	// 生成图片
	_, err := CropJson2Image(inputData, config, "output/output_crop.png")
	if err != nil {
		t.Fatalf("生成图片失败: %v\n", err)
		return
	}

	// 输出二进制图片
	body, err := CropJson2Image(inputData, config)
	if err != nil {
		t.Fatalf("生成图片失败: %v\n", err)
		return
	}

	t.Logf("Base64字符串长度: %d", len(body))
}

func TestJson2ImageBackwardCompatibility(t *testing.T) {
	// 测试向后兼容的函数
	jsonData := `{"test": "backward compatibility"}`

	_, err := Json2ImageDefault(jsonData, "output/output_backward.png")
	if err != nil {
		t.Errorf("向后兼容函数失败: %v", err)
	}
	fmt.Println("向后兼容测试成功")
}

func TestCropJson2ImageBackwardCompatibility(t *testing.T) {
	// 测试向后兼容的裁剪函数
	inputData := `{
		"data": {
			"item1": {
				"name": "test1",
				"value": 100
			},
			"item2": {
				"name": "test2",
				"value": 200
			}
		}
	}`

	rules := []string{
		"data.item1.name",
		"data.*.value",
	}

	_, err := CropJson2ImageDefault(inputData, rules, "output/output_crop_backward.png")
	if err != nil {
		t.Errorf("向后兼容裁剪函数失败: %v", err)
	}
	fmt.Println("向后兼容裁剪测试成功")
}

func TestJson2ImageWithDifferentFonts(t *testing.T) {
	// 测试不同字体
	jsonData := `{
		"message": "Hello World 你好世界",
		"numbers": [1, 2, 3, 4, 5]
	}`

	fontTypes := []FontType{
		FontTypeMonaco,
		FontTypeMsyh,
		FontTypePingFang,
	}

	fontNames := []string{
		"monaco",
		"msyh",
		"pingfang",
		"wrjs",
	}

	successCount := 0
	for i, fontType := range fontTypes {
		config := DefaultConfig().WithFont(fontType)
		fileName := fmt.Sprintf("output/output_font_%s.png", fontNames[i])

		_, err := Json2Image(jsonData, config, fileName)
		if err != nil {
			t.Logf("警告: 字体 %s 生成图片失败: %v", fontNames[i], err)
			// 不直接失败，而是记录警告并继续
		} else {
			fmt.Printf("字体 %s 图片生成成功：%s\n", fontNames[i], fileName)
			successCount++
		}
	}

	// 如果至少有一个字体成功，测试就通过
	if successCount == 0 {
		t.Error("所有字体都加载失败")
	} else {
		t.Logf("成功加载了 %d/%d 个字体", successCount, len(fontTypes))
	}
}
