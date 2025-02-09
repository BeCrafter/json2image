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

	_, err := Json2Image(jsonData, "output.png")
	if err != nil {
		fmt.Printf("生成图片失败: %v\n", err)
		return
	}
	fmt.Println("图片生成成功：output.png")
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

	rules := []string{
		"orders.order1.items[*].product.name",
		"orders.order1.items[0].product.cnt",
		"orders.order1.items[1].product.price",
		"orders.order1.items[1,0].product.uname",
		"orders.*.customer.name",
		"extras.content.adjustInfo.adjustId",
	}

	// 生成图片
	_, err := CropJson2Image(inputData, rules, "output.png")
	if err != nil {
		t.Fatalf("生成图片失败: %v\n", err)
		return
	}

	// 输出二进制图片
	body, err := CropJson2Image(inputData, rules)
	if err != nil {
		t.Fatalf("生成图片失败: %v\n", err)
		return
	}

	t.Log(body)
}
