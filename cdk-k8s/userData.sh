#!/bin/bash

# Update the package repository and install required dependencies
sudo apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker ec2-user

# Install Kubernetes components (kubeadm, kubectl, kind)
curl -LO "https://dl.k8s.io/release/v1.27.4/bin/linux/amd64/kubectl"
curl -LO "https://dl.k8s.io/release/v1.27.4/bin/linux/amd64/kubectl.sha256"
install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
apt-get update
apt-get install -y apt-transport-https ca-certificates curl
curl -fsSL "https://dl.k8s.io/apt/doc/apt-key.gpg" | gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list
apt-get update
apt-get install -y kubectl
sudo snap install kubeadm --classic
apt-mark hold kubelet kubeadm kubectl
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