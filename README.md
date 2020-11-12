## KubeXCloud

KubeXCloud (KXC) is a minimalist self service cloud platform built on top of Kubernetes. Built with Go.

**THIS SOFTWARE IS WORK IN PROGRESS / ALPHA RELEASE AND IS NOT MEANT FOR USAGE IN PRODUCTION SYSTEMS**

### Features
- Deploy anywhere you can host a Kubernetes cluster
- Users management
- Launch apps
- Groups apps within projects
- Isolate projects from each other
- Expose apps via http/https  

### Architecture
- KXC API server: receives REST requests and interacts with the Kubernetes API server to create Custom Resources
- KXC Operator/Controllers: monitors Custom Resources created by the KXC API server and reconciliates the internal Kubernetes resources (deployments/services/etc)
- KXC CLI: command line tool to interact with the KXC API server 

