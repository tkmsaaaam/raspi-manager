#!/bin/bash
cd `dirname $0`
source ./env.sh
export TEXT="スキャン結果"$'\n'
export RESULT=`clamdscan -i $1`
echo "$RESULT" >> /var/log/clamdscan.log
export BODY=`echo "$RESULT" | grep -v "Access denied. ERROR"`
export MESSAGE=$TEXT$BODY
./slackPublisher
exit 0
