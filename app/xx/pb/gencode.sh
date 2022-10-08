#!/bin/bash
SERVICE_NAME="xx"
FILENAME="${SERVICE_NAME}.proto"
PBNAME="${SERVICE_NAME}.pb"

ln -s ../../../common/pb/${FILENAME} ./${FILENAME}
# shellcheck disable=SC2046
# shellcheck disable=SC2006
goctl rpc protoc -I=$(dirname ../../../common/pb/${FILENAME}) ${FILENAME} -v --go_out=../../../common --go-grpc_out=../../../common  --zrpc_out=.. --style=goZero
# shellcheck disable=SC2046
protoc --proto_path=$(dirname ../../../common/pb/${FILENAME}) --descriptor_set_out=../../../common/pb/${PBNAME} ${FILENAME}
# shellcheck disable=SC2013
# shellcheck disable=SC2006
for i in `grep package -rl ../../../common/pb/*.go` ; do
    sed -i "" "s#,omitempty##g" "${i}"
done
