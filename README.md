### csail.sh
csail is a platform built on Kubernetes to provide an easy way to deploy docker containers. It provides an easy way to run monolithic apps in a kubernetes cluster without the usual complexities of running a microservice. If you're a company or an individual running multiple projects, csail is an easy way to run these projects in a K8s cluster and enjoy all the advantages of k8s without the complexities or even any knowledge of it at all.

## Features
* Support for all programming languages, you only need a valid `Dockerfile` to deploy your code
* Automatic `https` on default sub-domain. `https://{app-name}.hostname.com`
* Support scaling up and down - supports scaling by modifying `deployment.replica` value of app's kubernetes deployment
* Environment variables management. Add, delete, update deployment's environment with a single command
* Multiple custom domain. csail uses `K8s Ingress` internally to facilitates addition and removal of custom domains
* `https` for custom domains. You can add and update `https` certs for custom domains. (Automated certs management still in the work)
* View application or container's logs
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
