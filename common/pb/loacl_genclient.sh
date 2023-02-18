# brew install swift-protobuf
rm -rf local_swift
mkdir -p local_swift/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
protoc -I . sdk.proto --swift_out=./local_swift/pb
cp sdk.proto local_swift/pb
rm -rf ~/Desktop/xxim.pb.core.swift.zip || true
zip -r ~/Desktop/xxim.pb.core.swift.zip local_swift
rm -rf local_swift

mkdir -p local_swift/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
# shellcheck disable=SC2006
protoFiles=(
    "common.proto"
    "conn.proto"
    "group.proto"
    "im.proto"
    "msg.proto"
    "notice.proto"
    "validate.proto"
    "relation.proto"
    "user.proto"
)
# shellcheck disable=SC2128
for file in $protoFiles
do
    protoc -I . "$file" --swift_out=./local_swift/pb
    cp "$file" local_swift/pb
done
rm -rf ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/swift
cp -r local_swift/pb ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/swift
rm -rf ~/Desktop/xxim.api.pb.swift.zip || true
zip -r ~/Desktop/xxim.api.pb.swift.zip local_swift

#################################################################################
rm -rf local_objc
mkdir -p local_objc/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
protoc -I . sdk.proto --objc_out=./local_objc/pb
cp sdk.proto local_objc/pb
rm -rf ~/Desktop/xxim.pb.core.objc.zip || true
zip -r ~/Desktop/xxim.pb.core.objc.zip local_objc
rm -rf local_objc

mkdir -p local_objc/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
# shellcheck disable=SC2006
protoFiles=(
    "common.proto"
    "conn.proto"
    "group.proto"
    "im.proto"
    "msg.proto"
    "notice.proto"
    "validate.proto"
    "relation.proto"
    "user.proto"
)
# shellcheck disable=SC2128
for file in $protoFiles
do
    protoc -I . "$file" --objc_out=./local_objc/pb
    cp "$file" local_objc/pb
done
rm -rf ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/objc
cp -r local_objc/pb ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/objc
rm -rf ~/Desktop/xxim.api.pb.objc.zip || true
zip -r ~/Desktop/xxim.api.pb.objc.zip local_objc
rm -rf local_objc

#################################################################################
rm -rf local_java
mkdir -p local_java/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
protoc -I . sdk.proto --java_out=./local_java/pb
cp sdk.proto local_java/pb
rm -rf ~/Desktop/xxim.pb.core.java.zip || true
zip -r ~/Desktop/xxim.pb.core.java.zip local_java
rm -rf local_java

mkdir -p local_java/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
# shellcheck disable=SC2006
protoFiles=(
    "common.proto"
    "conn.proto"
    "group.proto"
    "im.proto"
    "msg.proto"
    "notice.proto"
    "validate.proto"
    "relation.proto"
    "user.proto"
)
# shellcheck disable=SC2128
for file in $protoFiles
do
    protoc -I . "$file" --java_out=./local_java/pb
    cp "$file" local_java/pb
done
rm -rf ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/java
cp -r local_java/pb ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/java
rm -rf ~/Desktop/xxim.api.pb.java.zip || true
zip -r ~/Desktop/xxim.api.pb.java.zip local_java
rm -rf local_java

#################################################################################
rm -rf local_dart
mkdir -p local_dart/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
protoc -I . sdk.proto --dart_out=./local_dart/pb
cp sdk.proto local_dart/pb
rm -rf ~/Desktop/xxim.pb.dart.zip || true
zip -r ~/Desktop/xxim.pb.dart.zip local_dart
rm -rf local_dart

mkdir -p local_dart/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
# shellcheck disable=SC2006
protoFiles=(
    "common.proto"
    "conn.proto"
    "group.proto"
    "im.proto"
    "msg.proto"
    "notice.proto"
    "validate.proto"
    "relation.proto"
    "user.proto"
)
# shellcheck disable=SC2128
for file in $protoFiles
do
    protoc -I . "$file" --dart_out=./local_dart/pb
    cp "$file" local_dart/pb
done
rm -rf ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/dart
cp -r local_dart/pb ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/dart
rm -rf ~/Desktop/xxim.api.pb.dart.zip || true
zip -r ~/Desktop/xxim.api.pb.dart.zip local_dart
#rm -rf local_dart

#################################################################################
rm -rf local_cpp
mkdir -p local_cpp/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
protoc -I . sdk.proto --cpp_out=./local_cpp/pb
cp sdk.proto local_cpp/pb
rm -rf ~/Desktop/xxim.pb.core.cpp.zip || true
zip -r ~/Desktop/xxim.pb.core.cpp.zip local_cpp
rm -rf local_cpp

mkdir -p local_cpp/pb || true
export PATH=$PATH:$HOME/.pub-cache/bin
# shellcheck disable=SC2006
protoFiles=(
    "common.proto"
    "conn.proto"
    "group.proto"
    "im.proto"
    "msg.proto"
    "notice.proto"
    "validate.proto"
    "relation.proto"
    "user.proto"
)
# shellcheck disable=SC2128
for file in $protoFiles
do
    protoc -I . "$file" --cpp_out=./local_cpp/pb
    cp "$file" local_cpp/pb
done
rm -rf ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/cpp
cp -r local_cpp/pb ~/go/src/eyq7369.xyz/showurl/xxim-protobuf/cpp
rm -rf ~/Desktop/xxim.api.pb.cpp.zip || true
zip -r ~/Desktop/xxim.api.pb.cpp.zip local_cpp
