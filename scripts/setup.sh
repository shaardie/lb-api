set -eu

ip="$(hostname -I | cut -d ' ' -f 1)"

# Install lb-api
apt-get install --yes haproxy
mkdir -p /var/lib/lb-api/ /etc/lb-api/
cat << EOF > /etc/systemd/system/lb-api.service
[Unit]
After=network-online.target

[Service]
ExecStart=/src/lb-api -config /etc/lb-api/lb-api.conf

[Install]
WantedBy=multi-user.target
EOF
cat << EOF > /etc/lb-api/lb-api.conf
db_filename: /var/lib/lb-api/db.json
configurator_filename: /etc/haproxy/haproxy.cfg
configurator_command: ["systemctl", "reload", "haproxy"]
ip: $ip
EOF
systemctl enable --now lb-api

# Install docker
mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list
apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Install kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.15.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# create kind cluster
cat << EOF > /etc/kind.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: $ip
nodes:
- role: control-plane
- role: worker
- role: worker
EOF
kind create cluster --config /etc/kind.yaml

# install cloud-provider-manager
mkdir -p /etc/cloud-provider-manager/
cp /root/.kube/config /etc/cloud-provider-manager/kubeconfig
cat << EOF > /etc/cloud-provider-manager/cloud.yaml
loadbalancer:
  url: http://$ip:8080
EOF
cat << EOF > /etc/systemd/system/cloud-provider-manager.service
[Unit]
After=network-online.target

[Service]
ExecStart=/src/cloud-provider-manager \
    --cloud-config /etc/cloud-provider-manager/cloud.yaml \
    --kubeconfig /etc/cloud-provider-manager/kubeconfig \
    --cloud-provider lb-api
    -v 4

[Install]
WantedBy=multi-user.target
EOF
systemctl enable --now cloud-provider-manager
