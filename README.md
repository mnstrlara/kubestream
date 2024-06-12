# GitOps-Enabled Kubernetes Deployment with FluxCD, Flagger, and CDK

## Objective:
The objective of this project is to set up a continuous deployment pipeline for a sample application using FluxCD for GitOps-style continuous delivery, Flagger for progressive delivery (canary deployments), and CDK for infrastructure as code (IaC) to provision and manage AWS resources.

## Components:
### 1. Sample Application
We need a sample application to deploy. I am going to be using an AWS EC2 Instance that has a Kubernetes cluster deployed inside of it using AWS CloudFormation.
### 2. Flux CD
FluxCD will be used to automate the deployment of the application. It will watch a Git repository where the application manifests are stored and ensure that the cluster stays in sync with the desired state defined in the Git repository.
### 3. Flagger
Flagger will be used for progressive delivery, enabling canary deployments and automated analysis of metrics (e.g., latency, error rates) to determine whether to promote or rollback a new version of the application.
### 4. CDK (AWS Cloud Development Kit)
CDK will be used to provision and manage AWS resources required by the application.

## Pre-requisites:
- AWS Root account created and available
- AWS CLI installed on local machine
- AWS CDK installed on local machine

## .env Configuration
Since there are some secret credentials that need to stay hidden they have been added inside the `.env` file. This is an example of the contents inside the file:
```shell
ACCOUNT_ID=[string]
ACCOUNT_REGION=[string]
VPC_ID=[string]
AMI_ID=[string]
SUBNET_ID=[string]
EC2_KEYPAIR=[string]
```

## Steps:
### 1. Kubernetes Deployment Flow using CDK
After following the documentation on installing CDK and running the cdk synth and bootstrap command, the next step would be to deploy the template using the
`cdk deploy` command. After running the command it should prompt you on whether you are sure about deploying the template to AWS CloudFormation.
After typing `y` the command is going to automatically create a change set in AWS CloudFormation and then deploy the stack.

### 2. Inside the Instance/Cluster
After SSHing inside the EC2 Instance, check if all of the components that are part of the `userData.sh`
have been successfully installed/created; if they have, the next part would be to set up our FluxCD. Since we have
already installed FluxCD with the userdata script, all we have to do is run the bootstrap command. 
But before we get to do that we should check our prerequisites with the `flux check --pre` command. If all the checks passed, we can now run our bootstrap command.

Example:
```bash
flux bootstrap github \
  --owner=$GITHUB_USER \
  --repository=example-name \
  --branch=main \
  --path=./example-path \
  --personal
```
After applying this command, your terminal is going to ask you for your GitHub PAT(Personal Access Token) which you can create in your GitHub Settings>Developer Settings>Personal Access Tokens>Tokens(classic)>Generate New Token>Generate New Token(classic); and then choose the Expiration date and the Scopes that fit your project best and finish off with "Generate Token".

This bootstrap command is going to push the Flux manifests to your Git repository and deploy Flux to your cluster. And as it becomes a part of your Git repository, since it is located in the path where your Kubernetes manifests are located, it will automatically deploy your manifests to your cluster without any extra commands.

### 3. Flagger Progressive Delivery
As Istio and Flagger installations are part of the `userData.sh` script, there is no further commands needed. The `flagger.yaml` canary deployment is located in the `k8s-deployments` directory with the rest of the Kubernetes manifests, which means that FluxCD does the deployment of the canary as well. To check if the canary is available, run the following command:
```shell
kubectl get canaries --all-namespaces
```

### Documentation
#### [AWS CLI install and update instructions](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) - This topic describes how to install or update the latest release of the AWS Command Line Interface (AWS CLI) on supported operating systems
#### [Getting started with AWS CDK](https://docs.aws.amazon.com/cdk/v2/guide/hello_world.html) - Quick Tutorial on starting with AWS CDK
#### [Flux Bootstrap Command](https://fluxcd.io/flux/cmd/flux_bootstrap/) - Several options on which bootstrap command to use and the proper way to use it
#### [Istio Canary Deployments](https://docs.flagger.app/tutorials/istio-progressive-delivery) - This guide shows you how to use Istio and Flagger to automate canary deployments