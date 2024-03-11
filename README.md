# Continuous Deployment of Kubernetes using FluxCD, Flagger, and CDK

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

## Steps
### 1. Kubernetes Deployment Flow using CDK
After following the documentation on installing CDK and running the cdk bootstrap command, the next step would be to deploy the template using the
`cdk deploy` command. After running the command it should prompt you on whether you are sure about deploying the template to AWS CloudFormation.
After typing `y` the command is going to automatically create a change set in AWS CloudFormation and then deploy the stack.

### 2. Inside the Instance/Cluster
After SSHing inside the EC2 Instance, check if all of the components that are part of the `userData.sh`
have been successfully installed/created; if they have, the next part would be to set up our FluxCD. Since we have
already installed FluxCD with the userdata, all we have to do is run the bootstrap command. 
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
After applying this command, your terminal is going to ask you for your GitHub PAT(Personal Access Token) which you can create in your GitHub Settings.

This bootstrap command is going to push the Flux manifests to your Git repository and deploy Flux to your cluster.

### 3. Flagger Progressive Delivery