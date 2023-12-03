#!/bin/bash
TEXT=`cat -n /var/log/auth.log | grep systemd-logind | tail -n 50`
curl -H "Authorization: Bearer $SLACK_TOKEN" -d "text=$TEXT" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
