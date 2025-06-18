#!/bin/bash

work_path=$(dirname $(readlink -f $0))

echo "work_path: ${work_path}"
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

# 使用Go程序进行转换，性能比bash脚本快10-100倍
cd ${work_path}
go run ./cmd/font_converter.go ${font_name} 

go fmt ${work_path}/../fonts/*