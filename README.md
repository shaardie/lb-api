# Loadbalancer API

This is a small tool to make it possible to configure your loadbalancer via HTTP API calls.

## Development

You can setup a development environment using Vagrant.
With a proper Vagrant installation simply run

```
make setup_init
```

and a VM is created and properly set up with lb-api, haproxy, a kind cluster with 3 Nodes and the Cloud Provider Manager.

Evething properly configured to talk to each other.

This will also create a `kubeconfig` which you can use to access the cluster.

A simple test for the loadbalancer would, e.g. to install an ingress

```
helm repo add nginx-stable https://helm.nginx.com/stable
helm repo update
helm --kubeconfig kubeconfig upgrade --install --namespace ingress --create-namespace ingress nginx-stable/nginx-ingress
```

You can use the vagrant command from within [scripts](./scripts) to take a look at them.

To compile and update the binaries in the setup during development simply run:

```
make setup_update
```

And to destroy the setup again run:

```
make setup_destroy
```

to enter the VM run

```
make setup_ssh
```
