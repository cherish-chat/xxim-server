#!/bin/zsh

# 此shell脚本用于生成pb文件

## 1.生成java文件
rm -rf sdk/java
mkdir -p sdk/java/pb
protoc --java_out=./sdk/java/pb/ *.proto
cp ./*.proto sdk/java/pb/

## 2.生成go文件
rm -rf sdk/go
mkdir -p sdk/go/pb
protoc --go_out=./sdk/go/pb/ *.proto
cp ./*.proto sdk/go/pb/

## 3.生成python文件
rm -rf sdk/python
mkdir -p sdk/python/pb
protoc --python_out=./sdk/python/pb/ *.proto
cp ./*.proto sdk/python/pb/

## 4.生成c++文件
rm -rf sdk/cpp
mkdir -p sdk/cpp/pb
protoc --cpp_out=./sdk/cpp/pb/ *.proto
cp ./*.proto sdk/cpp/pb/

## 5.生成c#文件
rm -rf sdk/csharp
mkdir -p sdk/csharp/pb
protoc --csharp_out=./sdk/csharp/pb/ *.proto
cp ./*.proto sdk/csharp/pb/

## 6.生成objc文件
rm -rf sdk/objc
mkdir -p sdk/objc/pb
protoc --objc_out=./sdk/objc/pb/ *.proto
cp ./*.proto sdk/objc/pb/

## 7.生成swift文件
rm -rf sdk/swift
mkdir -p sdk/swift/pb
protoc --swift_out=./sdk/swift/pb/ *.proto
cp ./*.proto sdk/swift/pb/

## 8.生成ts文件
rm -rf sdk/ts
mkdir -p sdk/ts/pb
npm install -g ts-protoc-gen
protoc --ts_out=./sdk/ts/pb/ *.proto
cp ./*.proto sdk/ts/pb/

## 9.生成rust文件
rm -rf sdk/rust
mkdir -p sdk/rust/pb
#cargo install protobuf-codegen
cp ./*.proto sdk/rust/pb/
protoc --rust_out=./sdk/rust/pb/ *.proto

## 10.生成dart文件
rm -rf sdk/dart
mkdir -p sdk/dart/pb
#dart pub global activate protoc_plugin
protoc --dart_out=./sdk/dart/pb/ *.proto
cp ./*.proto sdk/dart/pb/


### 每一个struct都加上#[derive(serde::Serialize,serde::Deserialize)]，否则无法序列化和反序列化
### 在#[derive(PartialEq,Clone,Default,Debug)]的上面加
# shellcheck disable=SC2164
#cd sdk/rust/pb
#for file in ./*.rs
#do
#  #判断系统 如果是mac则不需要加""，如果是linux则需要加""
#      if [[ `uname` == "Darwin" ]]; then
#        sed -i '' 's#fn from_i32(value: i32)#pub fn from_i32(value: i32)#g' "$file"
#        #sed -i '' 's#fn value(&self) -> i32#pub fn value(&self) -> i32#g' "$file"
#      elif [[ `uname` == "Linux" ]]; then
#        sed -i 's#fn from_i32(value: i32)#pub fn from_i32(value: i32)#g' "$file"
#        #sed -i 's#fn value(&self) -> i32#pub fn value(&self) -> i32#g' "$file"
#      else
#        echo "未知系统，请手动对pb.rs文件增加json序列化支持"
#      fi
#
#done
## shellcheck disable=SC2164
#cd -
