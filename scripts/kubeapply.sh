#!/bin/bash
cd $(dirname $0)
cd ..
find k8s | grep -e "yaml" -e "yml" | grep -v "node_modules" | xargs -I@ sh -c "kubectl apply -f @"
