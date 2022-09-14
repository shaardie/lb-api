# Loadbalancer API

This is a small tool to make it possible to configure your loadbalancer via HTTP API calls.

## Installation

You can install the tool with the following step.
For production services a proper package or configuration management tool should be used.

```bash
# Install dependencies
apt-get update
apt-get install -y haproxy golang curl make

# Build binary
make build

# Install
mkdir -p /var/lib/lb-api/ /etc/lb-api/
cp lb-api /usr/bin/
cat << EOF > /etc/systemd/system/lb-api.service
[Unit]
Description=Loadbalancer API Service
After=network-online.target

[Service]
ExecStart=/usr/bin/lb-api -config /etc/lb-api/lb-api.conf

[Install]
WantedBy=multi-user.target
EOF
cat << EOF > /etc/lb-api/lb-api.conf
db_filename: /var/lib/lb-api/db.json
configurator_filename: /etc/haproxy/haproxy.cfg
configurator_command: ["systemctl", "reload", "haproxy"]
ip: $(ip -4 -o a | grep -v "scope host" | awk '{print $4}' | awk -F '/' '{print $1}')
EOF
```

## Development

For development you need a configuration file to test this tool.
The following one puts all database and configuration files in the current directory and do not actually reload any proxy:

```yaml
db_filename: ./db.json
configurator_filename: ./haproxy.conf
configurator_command: ["logger", "-t", "lb-api", "reload"]
hostname: lb-api.example.com
```
