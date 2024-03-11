#!/bin/bash

# Update the package repository and install required dependencies
sudo apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo useradd ec2-user
sudo usermod -aG docker ec2-user

# Install Kubernetes components (kubeadm, kubectl, kind)
sudo curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.29/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.29/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo snap install kubeadm --classic
apt-mark hold kubeadm kubectl
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo cp ./kind /usr/local/bin/kind
rm -rf kind

# Install Helm
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

# Install FluxCD
curl -s https://fluxcd.io/install.sh | sudo bash

# Install K9s
wget https://github.com/derailed/k9s/releases/download/v0.24.7/k9s_Linux_x86_64.tar.gz
tar -xzvf k9s_Linux_x86_64.tar.gz
sudo mv k9s /usr/local/bin

# Add kubectl bash completion
echo "source <(kubectl completion bash)" >> ~/.bashrc

# Create Kubernetes Cluster fluxcd
sudo kind create cluster --name fluxcd

# Add repo and install Istio + create istio namespace
sudo helm repo add istio https://istio-release.storage.googleapis.com/charts
sudo helm repo update
sudo kubectl create namespace istio-system
sudo helm install istio-base istio/base -n istio-system --set defaultRevision=default

# Add repo and install Flagger
sudo helm repo add flagger https://flagger.app
sudo kubectl apply -f https://raw.githubusercontent.com/fluxcd/flagger/main/artifacts/flagger/crd.yaml
sudo helm upgrade -i flagger flagger/flagger \
--namespace=istio-system \
--set crd.create=false \
--set meshProvider=istio \
--set metricsServer=http://prometheus:9090