#!/bin/bash
HOST_NAME=$1
OPTS=$2
HADOLINT_OPTS=$3
cd $(dirname $0)
cd ..
basedir=`pwd`
find k8s | grep Dockerfile | while read filename
do
  path=`echo $filename | sed "s#/Dockerfile##g"`
  cd $path
  pwd
  hadolint Dockerfile $HADOLINT_OPTS
  imageName=`cat Dockerfile | head -n 1 | sed "s/# //g"`
  echo $HOST_NAME/$imageName
  docker build . -t $HOST_NAME/$imageName $OPTS
  cd $basedir
done
