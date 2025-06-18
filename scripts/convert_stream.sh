#!/bin/bash

work_path=$(dirname $(readlink -f $0))
root_path=$(dirname ${work_path})

font_name=$1

if [ -z "${font_name}" ]; then
    echo -e "\nUsage: $0 <font_name>\n"
    echo -e "Available fonts:"
    for item in $(ls -l ${work_path}/fonts | awk '{print $9}'); do
        echo "    ${item}"
    done
    echo ""
    exit 1
fi

font_name_lower=$(echo ${font_name} | tr '[:upper:]' '[:lower:]')
font_name_first=$(echo ${font_name} | sed 's/\b\(.\)/\u\1/g')

output_file=${root_path}/fonts/${font_name_lower}.go

# 单行字符串长度
line_len=150
font_name_len=${#font_name}
# 首行字符串长度
first_line_len=$(( ${line_len} - 13 - ${font_name_len} ))

# 直接生成Go文件头部
echo "package fonts" > ${output_file}
echo " const ${font_name_first}FontData = " >> ${output_file}

# 使用awk流式处理base64数据，避免加载整个字符串到内存
cat ${work_path}/fonts/${font_name}.ttf | base64 | tr -d '\n' | \
awk -v first_len=${first_line_len} -v line_len=${line_len} '
BEGIN {
    is_first = 1
}
{
    data = $0
    len = length(data)
    pos = 1
    
    # 处理第一行
    if (is_first) {
        first_line = substr(data, 1, first_len)
        print " `" first_line "` +"
        pos = first_len + 1
        is_first = 0
    }
    
    # 处理剩余部分
    while (pos <= len) {
        chunk = substr(data, pos, line_len)
        if (pos + line_len > len) {
            # 最后一行
            print " `" chunk "` ``"
        } else {
            print " `" chunk "` +"
        }
        pos += line_len
    }
}' >> ${output_file}

go fmt ${output_file}