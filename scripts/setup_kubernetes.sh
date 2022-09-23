#!/bin/bash

set -eu

snap install microk8s --classic --channel=1.25
microk8s status --wait-ready
microk8s config > /etc/kubeconfig

cat << EOF > /etc/systemd/system/cloud-provider-manager.service
[Unit]
After=network-online.target

[Service]
ExecStart=/src/cloud-provider-manager --kubeconfig /etc/kubeconfig --cloud-provider lb-api

[Install]
WantedBy=multi-user.target
EOF

systemctl enable --now cloud-provider-manager

curl https://baltocdn.com/helm/signing.asc | gpg --dearmor > /usr/share/keyrings/helm.gpg
apt-get install apt-transport-https --yes
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
apt-get update
apt-get install helm --yes
