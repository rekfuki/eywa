# Description
Eywa is an end-to-end open source FaaS platform that was developed as an Honours project at the University of Dundee. 

### Dissertation: https://drive.google.com/file/d/1DEb2UZUNqcSHW7OGN9pZ1QDMBoU6JH9f/view?usp=sharing

### Demo: https://youtu.be/QFQ6kmHnY74

## Eywa Layout Guide
9
The repository contains the following folder structure:

```bash
$tree -d -L 1
    .
    ├── app
    ├── charts
    ├── envoy
    ├── execution-tracker
    ├── eywa-cluster
    ├── gateway
    ├── go-libs
    ├── idler
    ├── registry
    ├── tugrik
    ├── warden
    └── watchdog
```

`./app` folder contains the front-end app.

`./charts` folder contains all the Helm charts that are used to manage Kubernetes objects (custom and third party).

`./envoy` folder contains envoy source code.

`./execution-tracker` folder contains both the API (`./execution-tracker/api`) and the Consumer (`./execution-tracker/consumer`) source code.

`./eywa-cluster` folder contains all the declarative Kubernetes objects that are being observed by the Flux toolkit. Reflects the state of the cluster.

`./gateway` folder contains both the API (`./gateway/api`) and the Consumer (`./gateway/consumer`) source code.

`./go-libs` folder contains any common libraries that are being shared between different components.

`./idler` folder contains the idler source code.

`./registry` folder contains the registry source code.

`./tugrik` folder contains the tugrik source code.

`./warden` folder contains the warden source code.

`./watchdog` folder contains the watchdog source code.

## Eywa Deployment Guide

#### Warning: Requires advanced knowledge of:

##### * Kubernetes ([docs](https://kubernetes.io/docs/home/))

##### * Kubernetes Container Storage Interface ([docs](https://kubernetes-csi.github.io/docs/))

##### * Kubernetes Secrets Store CSI Driver ([docs](https://github.com/kubernetes-sigs/secrets-store-csi-driver))

##### * HashiCorp Vault Provider for Secrets Store CSI Driver ([docs](https://github.com/hashicorp/vault-csi-provider))

##### * Flux GitOps Toolkit ([docs](https://toolkit.fluxcd.io/))

##### Building Containers:

Before you can deploy the cluster, you will need to build all the containers. Make sure you have installed Go programming language (documentation can be found [here](https://golang.org/doc/install)) and have setup a Docker Registry with basic auth credentials (documentation can be found [here](https://docs.docker.com/registry/deploying/)). Remember to save Docker Registry auth credentials for later, required when deploying the cluster.

There are total of 11 container images that need to be built:

| Repo Path                    | Image Name                 | Chart Path                          |
| ---------------------------- | -------------------------- | ----------------------------------- |
| ./envoy                      | envoy                      | ./charts/envoy                      |
| ./execution-tracker/api      | execution-tracker-api      | ./charts/execution-tracker-api      |
| ./execution-tracker/consumer | execution-tracker-consumer | ./charts/execution-tracker-consumer |
| ./gateway/api                | gateway-api                | ./charts/gateway-api                |
| ./gateway/consumer           | gateway-consumer           | ./charts/gateway-conusmer           |
| ./idler                      | idler                      | ./charts/idler                      |
| ./registry                   | registry                   | ./charts/registry                   |
| ./tugrik                     | tugrik                     | ./charts/tugrik                     |
| ./warden                     | warden                     | ./charts/warden                     |
| ./app                        | app                        | ./charts/app                        |
| ./watchdog                   | of-watchgod                | -                                   |

In order to build a container image of each you will need to first `git tag` each individual component. 

The tag should be of the format `{image_name}-{version}`. The `image_name` can be found in the table above or alternatively, inside every components directory (`Repo Path` in the table) there is `Makefile` that contains `IMAGE` variable at the top of the file. 

The `version` should be set to the same one as the `appVersion` in side `Chart.yaml` of every single component. Every component's  `Chart.yaml` file can be found in the directory indicated by the `Chart Path` column from the table above.

The only exception is the `watchdog`  component. While it still needs to be tagged, there is no chart for it and therefore any git tag can be used.

You will also need to update every component's `Makefile` and set the `DOCKER_REGISTRY` variable to point to your deployed docker registry.

**An example using `gateway-api`:**

First check the `appVersion` inside the `./charts/gateway-api/Chart.yaml` file. It should be set to `appVersion: 1.0.3` (the version might be inaccurate as development is still occurring). 

Then navigate to `./gateway/api` folder and create a new git tag: 

```bash
# Generate new tag
git tag gateway-api-1.0.3
# Push tag
git push --tags
```

After the tag is successfully created and pushed, run make command for building, tagging and pushing:

```bash
make build tag push
```



##### Cluster Deployment Steps:

1. Deploy a Kubernetes cluster. Recommended to use Minikube due to how easy it is to set it up. Documentation on how to setup a Minikube cluster can be found [here](https://minikube.sigs.k8s.io/docs/start/) (any other solution will of course work too).

2. Deploy HashiCorp Vault. Documentation on how to set it up can be found [here](https://learn.hashicorp.com/tutorials/vault/deployment-guide).

3. Patch `./charts/key-vault-csi/values.yaml` and set `vaultAddress` to point to your Vault instance.

4. Bootstrap WeaveWorks Flux to your Kubernetes cluster. Documentation on how to do it can be found [here](https://toolkit.fluxcd.io/get-started/). When specifying `--path` parameter to the bootstrap command, make sure to set it to `./eywa-cluster`.

5. Setup Vault to accept Kubernetes authentication.

   1. Enable Kubernetes authentication (documentation on it can be found [here](https://www.vaultproject.io/docs/auth/kubernetes)) by running the following command:
      `vault auth enable kubernetes`

   2. Write Kubernetes auth config:

      ```bash
      export KUBE_CA_CERT=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.certificate-authority-data}' | base64 --decode);
      export KUBE_HOST=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.server}');
      export VAULT_SA_NAME=$(kubectl get sa secrets-store-csi-driver -n csi -o jsonpath="{.secrets[*]['name']}");
      export SA_JWT_TOKEN=$(kubectl get secret $VAULT_SA_NAME -n csi -o jsonpath="{.data.token}" | base64 -d) &&
      vault write auth/kubernetes/config \
        kubernetes_host="$KUBE_HOST" \
        kubernetes_ca_cert="$KUBE_CA_CERT" \
        token_reviewer_jwt="$SA_JWT_TOKEN"
      ```

   3. Write Vault Kubernetes role:

      ```bash
      vault write auth/kubernetes/role/secrets-loader \
        bound_service_account_names=secrets-store-csi-driver \
        bound_service_account_namespaces=csi \
        policies=default,secrets-loader \
        ttl=20m
      ```

      

6. Create GitHub OAuth2.0 application and save the credentials for later. The registration page can be found [here](https://github.com/settings/applications/new).

7. Write both Docker Registry basic auth credentials as well as GitHub OAuth2.0 application credentials to Vault:

   ```bash
   # Write GitHub OAuth2.0 application credentials
   vault kv put secret/github-oauth \
   client_id={YOUR_GITHUB_APP_CLIENT_ID} \
   callback_url={YOUR_GITHUB_APP_CALLBACK_URL} \
   client_secret={YOUR_GITHUB_APP_CLIENT_SECRET}
   ```

   ```bash
   # Write Docker Registry credentials
   vault kv put secret/docker-registry-auth \
   docker-server={DOCKER_REGISTRY_SERVER} \
   docker-username={DOCKER_REGISTRY_USERNAME} \
   docker-password={DOCKER_REGISTRY_PASSWORD} \
   docker-email={DOCKER_REGISTRY_EMAIL}
   ```

8. Write Vault policies:

   ```bash
   vault policy write secrets-loader - <<EOF
   path "sys/renew/*" {
     capabilities = ["update"]
   }
   path "sys/mounts" {
     capabilities = ["read"]
   }
   path "secret/data/docker-registry-auth" {
     capabilities = ["read", "list"]
   }
   path "secret/data/github-oauth" {
     capabilities = ["read", "list"]
   }
   EOF
   ```

9. Restart the following component's pods:

   | Namespace   | Pod/Deployment Name                              |
   | ----------- | ------------------------------------------------ |
   | csi         | csi-secrets-store-provider-vault (pod)           |
   | csi         | csi-secrets-store-secrets-store-csi-driver (pod) |
   | envoy       | kube-pull-loader (deployment)                    |
   | faas        | kube-pull-loader (deployment)                    |
   | faas-system | kube-pull-loader (deployment)                    |
   | frontend    | kube-pull-loader (deployment)                    |

10. Everything should start rolling out. 
