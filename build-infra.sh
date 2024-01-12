#!/bin/bash

reg_name='kind-registry'
reg_port='5001'

if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --network bridge --name "${reg_name}" \
    registry:2
fi

kind create cluster --config kind-config.yaml --name goodcluster

kind get kubeconfig --name goodcluster > mykubeconfig

REGISTRY_DIR="/etc/containerd/certs.d/localhost:${reg_port}"
for node in $(kind get nodes -n goodcluster); do
  docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://${reg_name}:5000"]
EOF
done

if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect "kind" "${reg_name}"
fi

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

docker build -t localhost:${reg_port}/cloudchat-backend:latest -f backend-go/backend_dockerfile ./backend-go/
docker push localhost:${reg_port}/cloudchat-backend:latest

terraform init
terraform apply -auto-approve

worker_nodes=$(kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | grep 'worker')

echo "Portainer can be accessed by clicking one of the following links:"
for node in $worker_nodes; do
    internal_ip=$(kubectl get node $node -o=jsonpath='{.status.addresses[?(@.type=="InternalIP")].address}')
    echo "http://$internal_ip:31000"
done