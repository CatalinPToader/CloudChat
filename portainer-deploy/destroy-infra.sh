#!/bin/bash

terraform destroy -auto-approve
kind delete cluster --name goodcluster