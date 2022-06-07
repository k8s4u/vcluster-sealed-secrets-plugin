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
vcluster connect vcluster --namespace vcluster -- kubectl get secret my-secret -o yaml
```
