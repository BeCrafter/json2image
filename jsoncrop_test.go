package json2image

import (
	"encoding/json"
	"testing"
)

func TestJsonCrop(t *testing.T) {
	inputJSON := `{
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

	var inputData map[string]interface{}
	if str, err := formatJSON(inputJSON); err != nil {
		t.Fatalf("解析输入JSON失败: %v", err)
	} else {
		if err := json.Unmarshal([]byte(str), &inputData); err != nil {
			t.Fatalf("解析输入JSON失败: %v", err)
		}
	}

	rules := []string{
		"orders.order1.items[*].product.name",
		"orders.order1.items[0].product.cnt",
		"orders.order1.items[1].product.price",
		"orders.order1.items[1,0].product.uname",
		"orders.*.customer.name",
		"extras.content.adjustInfo.adjustId",
	}

	output, _ := JsonCrop(inputData, rules)

	var res map[string]interface{}
	json.Unmarshal(output, &res)

	outputJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		t.Fatalf("生成输出JSON失败: %v", err)
	}
	t.Log(string(outputJSON))
}

func TestJsonCropWithComplexRules(t *testing.T) {
	// 测试更复杂的裁剪规则
	inputJSON := `{
		"data": {
			"users": [
				{
					"id": 1,
					"name": "Alice",
					"profile": {
						"email": "alice@example.com",
						"age": 25,
						"preferences": {
							"theme": "dark",
							"language": "en"
						}
					}
				},
				{
					"id": 2,
					"name": "Bob",
					"profile": {
						"email": "bob@example.com",
						"age": 30,
						"preferences": {
							"theme": "light",
							"language": "zh"
						}
					}
				}
			],
			"metadata": {
				"total": 2,
				"page": 1
			}
		}
	}`

	var inputData map[string]interface{}
	if str, err := formatJSON(inputJSON); err != nil {
		t.Fatalf("解析输入JSON失败: %v", err)
	} else {
		if err := json.Unmarshal([]byte(str), &inputData); err != nil {
			t.Fatalf("解析输入JSON失败: %v", err)
		}
	}

	rules := []string{
		"data.users[*].name",
		"data.users[*].profile.email",
		"data.users[0].profile.preferences.theme",
		"data.metadata.total",
	}

	output, err := JsonCrop(inputData, rules)
	if err != nil {
		t.Fatalf("JsonCrop失败: %v", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(output, &res); err != nil {
		t.Fatalf("解析输出JSON失败: %v", err)
	}

	outputJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		t.Fatalf("生成输出JSON失败: %v", err)
	}
	t.Log("Complex crop result:", string(outputJSON))
}

func TestJsonCropWithWildcards(t *testing.T) {
	// 测试通配符规则
	inputJSON := `{
		"services": {
			"web": {
				"port": 80,
				"status": "running"
			},
			"db": {
				"port": 3306,
				"status": "stopped"
			},
			"cache": {
				"port": 6379,
				"status": "running"
			}
		}
	}`

	var inputData map[string]interface{}
	if str, err := formatJSON(inputJSON); err != nil {
		t.Fatalf("解析输入JSON失败: %v", err)
	} else {
		if err := json.Unmarshal([]byte(str), &inputData); err != nil {
			t.Fatalf("解析输入JSON失败: %v", err)
		}
	}

	rules := []string{
		"services.*.port",
		"services.*.status",
	}

	output, err := JsonCrop(inputData, rules)
	if err != nil {
		t.Fatalf("JsonCrop失败: %v", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(output, &res); err != nil {
		t.Fatalf("解析输出JSON失败: %v", err)
	}

	outputJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		t.Fatalf("生成输出JSON失败: %v", err)
	}
	t.Log("Wildcard crop result:", string(outputJSON))
}
