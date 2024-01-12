#!/bin/bash

terraform destroy -auto-approve
kind delete cluster --name goodcluster

docker stop kind-registry
docker rm kind-registry

rm -rf .terraform
rm .terraform.lock.hcl
rm mykubeconfig
rm terraform.tfstate
rm terraform.tfstate.backup