#!/bin/bash

# 根据api文件生成api代码
goctl api go --api ./api/reward.api --dir ../api --home ../../../scripts/goctl/1.6.5

# 根据api文件生成swagger文档
goctl api plugin -p goctl-swagger="swagger -filename reward-api.json" --api ./api/reward.api --dir ./swagger

# 根据sql文件生成model文件
goctl model mysql ddl --src ./sql/reward.sql --cache --dir ../model --home ../../../scripts/goctl/1.6.5
goctl model mysql ddl --src ./sql/other.sql --dir ../model
