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

ret="package fonts\n"
ret="${ret} const ${font_name_first}FontData = "


body_str=$(cat ${work_path}/fonts/${font_name}.ttf | base64)

# 第一行使用189字符
first_line=${body_str:0:${first_line_len}}
ret="${ret} \`${first_line}\` +\n"

# 剩余部分每200字符一行
remaining=${body_str:${first_line_len}}
len=${#remaining}
pos=0

while [ $pos -lt $len ]; do
    ret="${ret} \`${remaining:$pos:${line_len}}\` +\n"
    pos=$(( pos + ${line_len} ))
done

echo -e "${ret} \`\`" > ${output_file}

go fmt ${output_file}