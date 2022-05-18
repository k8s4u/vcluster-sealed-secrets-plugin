## Sealed Secret Plugin for vcluster
k8s4u (Kubernetes for you) is open source project/organization which targets for share well tested Kubernetes management code between organizations so that everyone does not need to solve same issues over and over again.

This repository contains [vcluster](https://vcluster.com/) [plugin](https://www.vcluster.com/docs/plugins/overview) which provides transparent version of [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets).


*Transparent* in this context means that you can use standard methods like `kubectl create secret ...` to create secrets (both plain text and already encrypted are supported) and `kubectl get secret ...` will always return only encrypted version of secret so it can be stored to GIT repo for GitOps type of use cases.

On other hand secrets inside of Kubernetes are still non-encrypted so pods can use those same way than normally.


**NOTE!!!** Currently hardcoded certificates from folder *tls* is used.
Do **NOT** use those in production.
