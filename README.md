# Loadbalancer API

This is a small tool to make it possible to configure your loadbalancer via HTTP API calls.

## Development

You can setup a development environment using Vagrant.
With a proper Vagrant installation simply run

```
make init_setup
```

and two VMs are created and properly set up.
One with lb-api and haproxy installed and running as a loadbalancer.
And one with MicroK8s and the Cloud Provider Manager running as the Cluster.

Both should be properly configured.

You can use the vagrant command from within [scripts](./scripts) to take a look at them.
