#!/bin/bash

# 根据api文件生成api代码
goctl api go --api ./api/pledge.api --dir ../api --home ../../../scripts/goctl/1.6.5

# 根据api文件生成swagger文档
goctl api plugin -p goctl-swagger="swagger -filename pledge-api.json" --api ./api/pledge.api --dir ./swagger

# 根据sql文件生成model文件
goctl model mysql ddl --src ./sql/pledge.sql --cache --dir ../model --home ../../../scripts/goctl/1.6.5
