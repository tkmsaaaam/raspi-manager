#!/bin/bash
cd `dirname $-0`
source ./env.sh
TEXT='ログイン履歴\n'
RESULT=`cat -n /var/log/auth.log | grep systemd-logind | tail -n 50`
curl -H "Authorization: Bearer $SLACK_TOKEN" -d "text=$TEXT$RESULT" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
exit 0
