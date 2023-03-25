#!/bin/bash

# 获取最新tag版本号
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)
if [ -z "$latest_tag" ]; then
  latest_tag="v0.0.0"
fi

# 将版本号转化为数字并自增1
version_number=$(echo $latest_tag | tr -d 'v')
version_array=(${version_number//./ })
major=${version_array[0]}
minor=${version_array[1]}
patch=${version_array[2]}
patch=$((patch+1))
new_version="v$major.$minor.$patch"

# 输出新的版本号
echo "New version: $new_version"

# 提交代码并打tag
git add .
git commit -m "$1"
git tag $new_version
git push origin "$2" $new_version
