#!/bin/bash
cd $(dirname $0)
basedir=`pwd`
find k8s | grep Dockerfile | while read filename
do
  path=`echo $filename | sed "s#/Dockerfile##g"`
  cd $path
  pwd
  imageName=`cat Dockerfile | head -n 1 | sed "s/# //g"`
  echo $imageName
  docker build . -t $imageName $1
  cd $basedir
done
