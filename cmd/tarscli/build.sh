#/bin/bash
set -e

repo="ccr.ccs.tencentyun.com/tarsbase"
version=$(git describe --tags)

make img VERSION=$version
docker tag ccr.ccs.tencentyun.com/tarsbase/tarscli:$version ccr.ccs.tencentyun.com/tarsbase/tarscli:latest

echo "push image:"
echo docker push ccr.ccs.tencentyun.com/tarsbase/tarscli:$version
echo docker push ccr.ccs.tencentyun.com/tarsbase/tarscli:latest