#!/bin/bash
cd `dirname $0`
source ./env.sh
TEXT="スキャン結果"$'\n'
RESULT=`clamdscan -i $1`
echo "$RESULT" > ./clamscan.log
curl -H "Authorization: Bearer $SLACK_BOT_TOKEN" -d "text=$TEXT$RESULT" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
exit 0
