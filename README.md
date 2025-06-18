# JSON2Image

一个将JSON数据转换为图片的Go库，支持自定义字体、颜色、样式和JSON裁剪功能。

## 功能特性

- 🎨 **多种字体支持**：内置Monaco、微软雅黑、苹方、王壬金石字体，也支持自定义字体
- 🌈 **可配置颜色方案**：支持自定义层级颜色和括号颜色
- ✂️ **JSON裁剪功能**：支持复杂的路径规则，可提取JSON的特定部分
- 🔧 **灵活配置**：通过选项模式轻松配置各种参数
- 📱 **多种输出格式**：支持保存为PNG文件或输出为Base64字符串
- 🔄 **向后兼容**：保持与旧版本API的兼容性

## 安装

```bash
go get github.com/BeCrafter/json2image
```

## 基本用法

### 使用默认配置

```go
import "github.com/BeCrafter/json2image"

jsonData := `{
    "name": "John Doe",
    "age": 30,
    "hobbies": ["reading", "coding"]
}`

// 保存为文件
_, err := json2image.Json2Image(jsonData, nil, "output.png")

// 或获取Base64字符串
base64Str, err := json2image.Json2Image(jsonData, nil)
```

## 高级配置

### 自定义字体和样式

```go
config := json2image.DefaultConfig().
    WithFont(json2image.FontTypePingFang).  // 使用苹方字体
    WithFontSize(16).                       // 字体大小
    WithLineHeight(24).                     // 行高
    WithPadding(30).                        // 内边距
    WithBackgroundColor(0.95, 0.95, 0.98)  // 背景色

_, err := json2image.Json2Image(jsonData, config, "custom.png")
```

### 自定义颜色方案

```go
// 定义层级颜色
levelColors := [][3]float64{
    {0.1, 0.2, 0.8}, // 深蓝色
    {0.8, 0.1, 0.2}, // 深红色
    {0.1, 0.8, 0.2}, // 深绿色
}

// 定义括号颜色
braceColors := [][3]float64{
    {0.4, 0.5, 0.9}, // 浅蓝色
    {0.9, 0.4, 0.5}, // 浅红色
    {0.4, 0.9, 0.5}, // 浅绿色
}

config := json2image.DefaultConfig().
    WithLevelColors(levelColors).
    WithBraceLevelColors(braceColors).
    WithDefaultTextColor(0.2, 0.2, 0.2)

_, err := json2image.Json2Image(jsonData, config, "colors.png")
```

### 使用自定义字体

```go
// 使用系统字体文件
config := json2image.DefaultConfig().
    WithCustomFont("/path/to/your/font.ttf")

_, err := json2image.Json2Image(jsonData, config, "custom_font.png")
```

## JSON裁剪功能

JSON裁剪允许你提取JSON中的特定部分，支持复杂的路径规则：

```go
jsonData := `{
    "users": [
        {
            "name": "Alice",
            "profile": {
                "email": "alice@example.com",
                "age": 25
            }
        },
        {
            "name": "Bob", 
            "profile": {
                "email": "bob@example.com",
                "age": 30
            }
        }
    ],
    "metadata": {
        "total": 2
    }
}`

// 定义裁剪规则
config := json2image.DefaultConfig().WithCropRules(
    "users[*].name",           // 所有用户的姓名
    "users[*].profile.email",  // 所有用户的邮箱
    "users[0].profile.age",    // 第一个用户的年龄
    "metadata.total",          // 总数
)

_, err := json2image.CropJson2Image(jsonData, config, "cropped.png")
```

### 裁剪规则语法

- `field.subfield` - 访问嵌套字段
- `array[*]` - 访问数组所有元素
- `array[0]` - 访问数组第0个元素
- `array[0,2]` - 访问数组第0和第2个元素
- `parent.*.field` - 使用通配符访问所有子对象的字段

## 配置选项

### 字体类型

```go
json2image.FontTypeMonaco     // Monaco字体
json2image.FontTypeMsyh       // 微软雅黑字体
json2image.FontTypePingFang   // 苹方字体
json2image.FontTypeWrjs       // 王壬金石字体
json2image.FontTypeCustom     // 自定义字体
```

### 链式配置方法

| 方法 | 说明 |
|------|------|
| `WithFont(fontType)` | 设置字体类型 |
| `WithCustomFont(path)` | 设置自定义字体路径 |
| `WithFontSize(size)` | 设置字体大小 |
| `WithLineHeight(height)` | 设置行高 |
| `WithPadding(padding)` | 设置内边距 |
| `WithBackgroundColor(r,g,b)` | 设置背景色 |
| `WithLevelColors(colors)` | 设置层级颜色 |
| `WithBraceLevelColors(colors)` | 设置括号颜色 |
| `WithDefaultTextColor(r,g,b)` | 设置默认文本颜色 |
| `WithCropRules(rules...)` | 设置裁剪规则 |

## 向后兼容

为了保持与旧版本的兼容性，我们提供了兼容的函数：

```go
// 旧版本兼容
_, err := json2image.Json2ImageDefault(jsonData, "output.png")

// 裁剪功能兼容
rules := []string{"field1", "field2"}
_, err := json2image.CropJson2ImageDefault(jsonData, rules, "output.png")
```

## API 参考

### 主要函数

#### Json2Image

```go
func Json2Image(jsonData string, config *Config, outputPath ...string) (string, error)
```

将JSON数据转换为图片。

**参数:**
- `jsonData`: JSON字符串
- `config`: 配置选项，传入nil使用默认配置
- `outputPath`: 可选的输出文件路径

**返回值:**
- 如果提供了`outputPath`，返回空字符串和可能的错误
- 如果未提供`outputPath`，返回Base64编码的图片数据

#### CropJson2Image

```go
func CropJson2Image(jsonData string, config *Config, outputPath ...string) (string, error)
```

对JSON进行裁剪后转换为图片。配置中必须包含裁剪规则。

#### DefaultConfig

```go
func DefaultConfig() *Config
```

返回默认配置实例。

## 示例项目

查看测试文件了解更多使用示例：

- `helper_test.go` - 配置和字体测试
- `json2image_test.go` - 图片生成测试  
- `jsoncrop_test.go` - JSON裁剪测试

## 注意事项

1. **字体兼容性**: 某些字体文件可能不受支持，建议使用标准的TTF或OTF格式
2. **内存使用**: 大型JSON数据可能消耗较多内存
3. **临时文件**: 内置字体会创建临时文件，函数执行完成后会自动清理
4. **自定义字体**: 使用自定义字体时，请确保字体文件存在且格式正确

## 许可证

本项目采用 [许可证名称] 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情
