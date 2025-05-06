#!/bin/bash

sudo apt-get -y update
sudo apt-get -y upgrade

sudo apt-get -y install ufw nginx avahi-daemon

# https://kubernetes.io/ja/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#kubeadm-kubelet-kubectl%E3%81%AE%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB

sudo ufw enable
sudo ufw deny ssh
sudo ufw allow 22222
sudo ufw allow 22223/tcp
sudo ufw allow 22224/tcp
sudo ufw allow 22225/tcp
sudo ufw reload

curl https://github.com/tkmsaaaam.keys >> ~/.ssh/authorized_keys
echo "Port 22222" >> /etc/ssh/sshd_config
sudo systemctl restart ssh
