#!/bin/bash

kind create cluster --config kind-config.yaml --name goodcluster
kind get kubeconfig --name goodcluster > mykubeconfig
terraform init
terraform apply -auto-approve

worker_nodes=$(kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | grep 'worker')

echo "Portainer can be accessed by clicking one of the following links:"
for node in $worker_nodes; do
    internal_ip=$(kubectl get node $node -o=jsonpath='{.status.addresses[?(@.type=="InternalIP")].address}')
    echo "http://$internal_ip:31000"
done