#!/bin/bash
cd `dirname $0`
source ./env.sh
TEXT='ログイン履歴\n'
DATE=`date '+%b %e'`
RESULT=`cat -n /var/log/auth.log | grep -a -e systemd-logind -e sshd | grep -a "$DATE"`
curl -H "Authorization: Bearer $SLACK_TOKEN" -d "text=$TEXT$RESULT" "https://slack.com/api/chat.postMessage?channel=$SLACK_CHANNEL"
exit 0
