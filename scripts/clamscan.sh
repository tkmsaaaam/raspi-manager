#!/bin/bash
TEXT=`clamscan / -r -i`
curl -H "Authorization: Bearer $TOKEN" -d "text=$TEXT" "https://slack.com/api/chat.postMessage?channel=$CHANNEL"
