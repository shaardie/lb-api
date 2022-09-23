#!/bin/bash

set -eu

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
ip: $(hostname -I | cut -f 1 -d ' ')
EOF

systemctl enable --now lb-api
