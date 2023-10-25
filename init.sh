#!/bin/bash

# 定义要创建的目录
DIRECTORIES=(
  api
  build
  cmd
  configs
  deployments
  docs
  githooks
  init
  internal/pkg/code
  internal/pkg/logger
  internal/pkg/middleware
  internal/pkg/options
  internal/pkg/server
  internal/pkg/util
  internal/pkg/util/gormutil
  internal/pkg/validation
  pkg/app
  pkg/cli/genericclioptions
  pkg/db
  pkg/log/
  pkg/shutdown/
  pkg/storage
  pkg/util/flag
  pkg/util/genutil
  pkg/validator
  test
  third_party
  tools
)

# 创建目录
echo "Creating directories..."
for DIRECTORY in "${DIRECTORIES[@]}"
do
   echo "Creating ${DIRECTORY}/"
   mkdir -p "${DIRECTORY}"
done

# 创建基础文件
echo "Creating base files..."
touch CHANGELOG
echo "Please fill in here..." >> CHANGELOG

touch CONTRIBUTING.md
echo "Please fill in here..." >> CONTRIBUTING.md

touch LICENSE
echo "Please fill in here..." >> LICENSE

touch Makefile
echo "Please fill in here..." >> Makefile

touch README.md
echo "# Project Title" >> README.md
echo "Project description here..." >> README.md

# 创建internal/pkg下的Go文件
echo "Creating Go files in internal/pkg..."
# [在此处添加所有 internal/pkg 下的 touch 命令]

# 创建pkg下的Go文件
echo "Creating Go files in pkg..."
# [在此处添加所有 pkg 下的 touch 命令]

echo "Structure created successfully."
