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
