#!/bin/bash
cd `dirname $0`
source ./env.sh
export TEXT="スキャン結果"$'\n'
export RESULT=`clamdscan -i $1`
echo "$RESULT" >> /var/log/clamdscan.log
export MESSAGE=`echo "$RESULT" | grep -v "Access denied. ERROR"`
curl -H "Authorization: Bearer $SLACK_BOT_TOKEN" -d "text=$TEXT$MESSAGE" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
exit 0
