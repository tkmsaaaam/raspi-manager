#!/bin/bash
TEXT='スキャン結果\n'
RESULT=`clamscan / -r -i`
curl -H "Authorization: Bearer $SLACK_TOKEN" -d "text=$TEXT$RESULT" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
