# json2image

json2image 是一个用于 JSON 数据可视化的 Go 语言工具库。它可以将 JSON 数据转换为带有语法高亮的图片，同时支持根据特定规则提取 JSON 数据并生成图片。

## 功能特点

- JSON 数据转图片，支持语法高亮
- 支持嵌套 JSON 结构的解析
- 支持按规则提取 JSON 数据
- 支持输出 PNG 图片或 Base64 编码
- 支持多层级的颜色区分
- 支持括号匹配的可视化

## 安装

```bash
go get github.com/BeCrafter/json2image@latest
```

## 主要功能

### 1. JSON 转图片 (Json2Image)

将 JSON 数据转换为图片，支持两种输出方式：

- 保存为 PNG 文件
- 返回 Base64 编码字符串

```go
// 方式一：保存为 PNG 文件
jsonData := `{
    "name": "John Doe",
    "age": 30,
    "address": {
        "street": "123 Main St",
        "city": "New York"
    }
}`

// 保存为文件
_, err := Json2Image(jsonData, "output.png")
if err != nil {
    log.Fatal(err)
}

// 方式二：获取 Base64 编码
base64Str, err := Json2Image(jsonData)
if err != nil {
    log.Fatal(err)
}
fmt.Println(base64Str)
```

### 2. JSON 数据提取并转图片 (CropJson2Image)

根据指定规则提取 JSON 数据并生成图片。支持复杂的数据提取规则：

```go
jsonData := `{
    "orders": {
        "order1": {
            "items": [
                {
                    "product": {
                        "name": "商品1",
                        "price": 100
                    }
                }
            ],
            "customer": {
                "name": "张三"
            }
        }
    }
}`

// 定义提取规则
rules := []string{
    "orders.order1.items[*].product.name",  // 提取所有商品名称
    "orders.*.customer.name",               // 提取所有客户名称
}

// 保存为文件
_, err := CropJson2Image(jsonData, rules, "output.png")
if err != nil {
    log.Fatal(err)
}

// 获取 Base64 编码
base64Str, err := CropJson2Image(jsonData, rules)
if err != nil {
    log.Fatal(err)
}
```

### 3. JSON 数据提取 (JsonCrop)

仅提取 JSON 数据而不生成图片：

```go
// 解析 JSON 数据
var inputData map[string]interface{}
if err := json.Unmarshal([]byte(jsonData), &inputData); err != nil {
    log.Fatal(err)
}

// 定义提取规则
rules := []string{
    "orders.order1.items[*].product.name",
    "orders.*.customer.name"
}

// 提取数据
output, err := JsonCrop(inputData, rules)
if err != nil {
    log.Fatal(err)
}
```

## 规则语法说明
提取规则支持以下语法：

1. 点号(`.`)：访问对象属性
   - 示例： `user.name`

2. 星号(`*`)：通配符，匹配所有键   
   - 示例： `orders.*.customer.name`

3. 数组索引：使用方括号   
   - 单个索引： `items[0]`
   - 多个索引： `items[0,1]`
   - 所有元素： `items[*]`

## 注意事项

1. 图片生成时会自动处理字体加载
2. 不同层级的 JSON 结构会使用不同的颜色标识
3. 支持处理嵌套的 JSON 字符串
4. Base64 输出适用于 Web 应用场景