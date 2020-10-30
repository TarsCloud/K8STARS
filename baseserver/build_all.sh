#!/bin/bash
set -e

repo="ccr.ccs.tencentyun.com/tarsbase"
version=$(git describe --tags)

make registry VERSION=$version IMG_REPO=$repo
docker tag $repo/tarsregistry:$version $repo/tarsregistry:latest
docker tag $repo/cppregistry:$version $repo/cppregistry:latest

for server in "tarsstat" "tarsconfig" "tarslog" "tarsnotify" "tarsproperty" "tarsquerystat" "tarsqueryproperty"; do
    make img SERVER=$server VERSION=$version IMG_REPO=$repo
    docker tag $repo/$server:$version $repo/$server:latest 
done 

# echo "push image:"
echo docker push $repo/tarsregistry:$version
echo docker push $repo/tarsregistry:latest
echo docker push $repo/cppregistry:$version
echo docker push $repo/cppregistry:latest
for server in "tarsstat" "tarsconfig" "tarslog" "tarsnotify" "tarsproperty" "tarsquerystat" "tarsqueryproperty"; do
    echo docker push $repo/$server:$version 
    echo docker push $repo/$server:latest 
done