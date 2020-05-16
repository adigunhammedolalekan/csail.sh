### csail.sh
csail is a platform built on Kubernetes to provide an easy way to deploy a docker container, with few commands, without writing `yaml specs`. It is made for monolithic apps and for companies or an individuals that builds or run multiple projects, csail will make it easy to run projects in a K8s cluster and enjoy all the advantages of k8s without the complexities or even any knowledge of it at all.

## Features
* Support for all programming languages, you only need a valid `Dockerfile` to deploy your code.
* Automatic `https` on default sub-domain. `https://{app-name}.hostname.com`.
* Support scaling up and down - by modifying `deployment.replicas` value of app's kubernetes deployment.
* Environment variables management. Add, delete, update and list deployment's environment with a single command.
* Multiple custom domain. csail uses `K8s Ingress` internally to facilitates addition and removal of custom domains.
* `https` for custom domains. You can add and update `https` certs for custom domains. (Automated certs management still in the work).
* View application or container's logs.
* Supports deployment rollbacks. Release is created on every deployment or changes in environment variable, you can rollback or update to any of application's releases.
* Supports database and storage systems provisioning. Currently supports
	* MySQL
	* PostgreSQL
	* MongoDB
	* Redis
	* Minio
* Support database dump for MySQL, PostgresSQL and MongoDB
* Database restore (In progress)
* Automatic database backup (In progress)

## Installation

### Requirements
The following is required to run `csail`
* A K8s cluster
* Docker
* Mac/Linux/Window

1. #### Install `csail-cli`
```
$ 
```
2. #### Setup K8s cluster
csail deploys all application in a kubernetes cluster. It is more like an advanced `kubectl for web apps`, so a K8s cluster is needed.
```
$ export CSAIL_CLUSTER_CONFIG=${HOME}/.kube/config && csail init k8s
```
The above command will setup will point `csail` to your k8s cluster. If you run `csail k8s init` without defining `CSAIL_CLUSTER_CONFIG`, csail will default to `${HOME}/.kube/config`

3. #### Install `nginx ingress`.  [Installation instructions](https://kubernetes.github.io/ingress-nginx/deploy/)
4. #### Install csail-server and all the necessary k8s applications.
```
$ kubectl apply -f https://github.com/adigunhammedolalekan/csail.sh/k8s/install/spec.yml
```
5. Setup domain and access
6. Deploy your first app
