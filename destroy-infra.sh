#!/bin/bash

terraform destroy -auto-approve
kind delete cluster --name goodcluster

rm -rf .terraform
rm .terraform.lock.hcl
rm mykubeconfig
rm terraform.tfstate
rm terraform.tfstate.backup