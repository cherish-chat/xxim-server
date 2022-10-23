#!/bin/bash
SERVICE_NAME="im"
FILENAME="${SERVICE_NAME}.proto"

ln -s ../../../common/pb/${FILENAME} ./${FILENAME} || echo "link exists"
# shellcheck disable=SC2046
protoc --proto_path=$(dirname ../../../common/pb/${FILENAME}) ${FILENAME} --go_out=../../../common
# shellcheck disable=SC2013
# shellcheck disable=SC2006
for i in `grep package -rl ../../../common/pb/*.go` ; do
    sed -i "" "s#,omitempty##g" "${i}"
done
