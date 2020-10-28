## KubeXCloud

KubeXCloud (KXC) is a minimalist self service cloud platform built on top of Kubernetes

### Features
- Deploy anywhere you can host a Kubernetes cluster
- Users management
- Launch containers
- Launch in-cluster db instances (planned feature)
- Expose apps via http/https  

### Architecture
- KXC API server: receives REST requests and interacts with the Kubernetes API server to create Custom Resources
- KXC Operator/Controllers: monitors Custom Resources created by the KXC API server and reconciliates the internal Kubernetes resources (deployments/services/etc)
- KXC CLI: command line tool to interact with the KXC API server 

