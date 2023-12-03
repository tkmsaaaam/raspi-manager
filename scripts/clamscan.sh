#!/bin/bash
TEXT=`clamscan / -r -i`
curl -H "Authorization: Bearer $SLACK_TOKEN" -d "text=$TEXT" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
