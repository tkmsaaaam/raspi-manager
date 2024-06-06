#!/bin/bash
cd `dirname $0`
source ./env.sh
export TEXT="pods"$'\n'
export RESULT=`kubectl get pods -A | grep -v "Running"`
export MESSAGE=$TEXT$RESULT
./slackPublisher
exit 0
