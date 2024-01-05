#!/bin/bash
cd $(dirname $0)
basedir=`pwd`
find k8s | grep Dockerfile | while read filename
do
  pwd
  path=`echo $filename | sed "s#/Dockerfile##g"`
  cd $path
  imageName=`cat Dockerfile | head -n 1 | sed "s/# //g"`
  docker build . -t $imageName $1
  cd $basedir
done
