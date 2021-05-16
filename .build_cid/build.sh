#!/bin/bash
# The build script
# 2021-05-09 20:02:15 CST
build_cid_path=${BUILD_DIR_PREFIX}/.build_cid
registry=`grep "^registry" ${build_cid_path}/cid.conf |tr -d ' \t'|cut -d '=' -f 2-`
reg_project=`grep "^reg_project" ${build_cid_path}/cid.conf |tr -d ' \t'|cut -d '=' -f 2-`
images_name=`grep "^name" ${build_cid_path}/cid.conf |tr -d ' \t'|cut -d '=' -f 2-`
images_version=`grep "^version" ${build_cid_path}/cid.conf |tr -d ' \t'|cut -d '=' -f 2-`

images_tag=${registry}/${reg_project}/${images_name}:${images_version}

docker build --build-arg version="$images_version" -t ${images_tag} . -f ${BUILD_DIR_PREFIX}/Dockerfile
