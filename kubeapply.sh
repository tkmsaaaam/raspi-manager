#!/bin/bash
ls k8s | grep -e "yaml" -e "yml" | xargs -I@ sh -c "kubectl apply -f k8s/@"
