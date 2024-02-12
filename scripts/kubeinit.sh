#!/bin/bash
sudo kubeadm reset
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
sudo cat /etc/kubernetes/admin.conf > $HOME/.kube/config
kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
# kubectl taint nodes --all node-role.kubernetes.io/control-plane-node/k8s-cplane.novalocal untainted
kubectl taint nodes --all node-role.kubernetes.io/control-plane-
# https://kubernetes.io/ja/docs/setup/production-environment/tools/kubeadm/
# https://tech.virtualtech.jp/entry/2022/08/24/172753
