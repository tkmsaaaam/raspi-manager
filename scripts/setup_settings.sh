#!/bin/bash

cd $(dirname $0)
cd ..
rootdir=`pwd`

# Set up Nginx configuration
sudo rm -rf /etc/nginx/conf.d
sudo ln -s $rootdir/symlinks/etc/nginx/conf.d/ /etc/nginx/conf.d
sudo systemctl restart nginx
