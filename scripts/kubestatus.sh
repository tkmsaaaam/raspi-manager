#!/bin/bash
cd `dirname $0`
source ./env.sh
export TEXT="pods"$'\n'
export RESULT=`kubectl get pods -A | grep -v "Running"`
curl -H "Authorization: Bearer $SLACK_BOT_TOKEN" -d "text=$TEXT$RESULT" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
exit 0
