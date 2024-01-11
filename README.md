
# CloudChat

This branch contains a deployment of Portainer into a Kubernetes cluster with 1 control plane and 2 workers using Terraform.

## Pre-requisites

One needs to have installed on own machine the following:

 - Docker
 - Kind - for the creation of Kubernetes cluster
 - Terraform - for applying the configuration onto the cluster

## How to run everything

**_Note:_** You might need to use **sudo** before sourcing scripts

Go to portainer-deploy directory 

Here execute commands:
```bash
  chmod +x build-infra.sh
  . build-infra.sh
```

After the script is done you can access Portainer at the links provided in the output of the script. Google Chrome works better than Firefox in the Portainer's Applications tab.

**_Note:_** The IPs that are working can be verified using the following command:
```bash
  kubectl get nodes -o wide
```
They are actually the ones from INTERNAL IP column.
## How to shut down the whole deployment

Here execute commands:
```bash
  chmod +x destroy-infra.sh
  . destroy-infra.sh
```
