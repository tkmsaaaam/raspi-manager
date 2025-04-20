#!/bin/bash

cd $(dirname $0)
cd ..
rootdir=`pwd`

# Set up promtail configuration
sudo ln -s $rootdir/symlinks/etc/ptomtail/config.yml /etc/promtail/config.yml

# Set up prometheus configuration
sudo mkdir /etc/prometheus
sudo cp $rootdir/settings/etc/prometheus/prometheus.yml /etc/prometheus/prometheus.yml

# Set up Nginx configuration
sudo rm -rf /etc/nginx/conf.d
sudo ln -s $rootdir/symlinks/etc/nginx/conf.d/ /etc/nginx/conf.d
sudo systemctl restart nginx
