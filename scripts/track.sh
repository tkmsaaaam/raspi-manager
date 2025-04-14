#!/bin/bash
cd $(dirname $0)
cd ..
git pull
kubectl apply -f k8s -R
