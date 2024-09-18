#!/bin/bash
cd $(dirname $0)
cd ..
kubectl apply -f k8s -R
