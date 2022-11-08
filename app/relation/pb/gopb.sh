#!/bin/bash
SERVICE_NAME="relation"
FILENAME="${SERVICE_NAME}.proto"

ln -s ../../../common/pb/${FILENAME} ${FILENAME} || echo "link exists"
ln -s ../../../common/pb/common.proto ./common.proto || echo "link exists"
# shellcheck disable=SC2046
protoc --proto_path=$(dirname ../../../common/pb/common.proto) common.proto --go_out=../../../common
# shellcheck disable=SC2046
goctl rpc protoc -I=$(dirname ../../../common/pb/${FILENAME}) ${FILENAME} -v --go_out=../../../common --go-grpc_out=../../../common  --zrpc_out=.. --style=goZero
# shellcheck disable=SC2013
# shellcheck disable=SC2006
for i in `grep package -rl ../../../common/pb/*.go` ; do
    sed -i "" "s#,omitempty##g" "${i}"
done
