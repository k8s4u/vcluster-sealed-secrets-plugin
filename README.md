## Sealed Secret Plugin for vcluster
k8s4u (Kubernetes for you) is open source project/organization which targets for share well tested Kubernetes management code between organizations so that everyone does not need to solve same issues over and over again.

This repository contains [vcluster](https://vcluster.com/) [plugin](https://www.vcluster.com/docs/plugins/overview) which provides transparent version of [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets).


*Transparent* in this context means that you can use standard methods like `kubectl create secret ...` to create secrets (both plain text and already encrypted are supported) and `kubectl get secret ...` will always return only encrypted version of secret so it can be stored to GIT repo for GitOps type of use cases.

On other hand secrets inside of Kubernetes are still non-encrypted so pods can use those same way than normally.


**NOTE!!!** Currently hardcoded certificates from folder *tls* is used.
Do **NOT** use those in production.


## Using the Plugin

At least version [0.9.1-beta.0](https://github.com/loft-sh/vcluster/releases/tag/v0.9.1-beta.0) of vcluster is needed.

To use the plugin, create a new vcluster with the `plugin.yaml`:

```
# Use public plugin.yaml
vcluster create vcluster -n vcluster -f https://raw.githubusercontent.com/k8s4u/vcluster-sealed-secrets-plugin/main/plugin.yaml
```

After vcluster is started create test secret and immediately read its content:
```
vcluster connect vcluster --namespace vcluster -- kubectl create secret generic my-secret --from-literal=key1=supersecret
vcluster connect vcluster --namespace vcluster -- kubectl annotate secret my-secret vcluster.loft.sh/force-sync=true
vcluster connect vcluster --namespace vcluster -- kubectl get secret my-secret -o yaml
```

## Building the Plugin
To just build the plugin image and push it to the registry, run:
```
# Build
docker build . -t k8s4u/vcluster-sealed-secrets-plugin:dev

# Push
docker push k8s4u/vcluster-sealed-secrets-plugin:dev
```

Then exchange the image in the `plugin.yaml`.

## Development

General vcluster plugin project structure:
```
.
├── go.mod              # Go module definition
├── go.sum
├── devspace.yaml       # Development environment definition
├── devspace_start.sh   # Development entrypoint script
├── Dockerfile          # Production Dockerfile
├── Dockerfile.dev      # Development Dockerfile
├── main.go             # Go Entrypoint
├── plugin.yaml         # Plugin Helm Values
├── syncers/            # Plugin Syncers
└── manifests/          # Additional plugin resources
```

Before starting to develop, make sure you have installed the following tools on your computer:
- [docker](https://docs.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) with a valid kube context configured
- [helm](https://helm.sh/docs/intro/install/), which is used to deploy vcluster and the plugin
- [vcluster CLI](https://www.vcluster.com/docs/getting-started/setup) v0.6.0 or higher
- [DevSpace](https://devspace.sh/cli/docs/quickstart), which is used to spin up a development environment

If you want to develop within a remote Kubernetes cluster (as opposed to docker-desktop or minikube), make sure to exchange `PLUGIN_IMAGE` in the `devspace.yaml` with a valid registry path you can push to.

After successfully setting up the tools, start the development environment with:
```
devspace dev -n vcluster
```

After a while a terminal should show up with additional instructions. Enter the following command to start the plugin:
```
go run -mod vendor ./cmd/main.go
```

The output should look something like this:
```
I0124 11:20:14.702799    4185 logr.go:249] plugin: Try creating context...
I0124 11:20:14.730044    4185 logr.go:249] plugin: Waiting for vcluster to become leader...
I0124 11:20:14.731097    4185 logr.go:249] plugin: Starting syncers...
[...]
I0124 11:20:15.957331    4185 logr.go:249] plugin: Successfully started plugin.
```

You can now change a file locally in your IDE and then restart the command in the terminal to apply the changes to the plugin.

Delete the development environment with:
```
devspace purge -n vcluster
```


### TODO
Fix this. Needs learning about sealed secrets logic..
```
MutateGetVirtual called
MutateCreatePhysical called
MutateCreatePhysical: Secret looks to be already non-encrypted, error: no key could decrypt secret (key1)
MutateGetVirtual called
```