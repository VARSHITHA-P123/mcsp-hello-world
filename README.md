# MultiCloud SaaS Platform — Hello World Application

## Overview

This repository contains the implementation of a MultiCloud SaaS Platform (MCSP) built on Red Hat OpenShift. It covers two scenarios — deploying a SaaS application with a complete CI/CD pipeline and automating multi-tenant customer onboarding.

---

## Scenarios Implemented

### Scenario 1: SaaS Application Deployment
A complete CI/CD and GitOps pipeline for deploying a Node.js application on OpenShift.

- Automated pipeline using Tekton with 4 stages — test, build, image tag update, and ArgoCD sync
- GitOps-based deployment using ArgoCD watching the GitHub repository
- RHACM governance policy ensuring application compliance
- TLS certificate provisioned automatically using Cert Manager

---

### Scenario 2: Multi-Tenant Customer Onboarding
Automated customer environment provisioning triggered by a single file pushed to GitHub.

- RHACM automatically creates isolated customer namespace with resource quotas
- ArgoCD ApplicationSet detects new customer file and deploys the application
- External Secrets operator syncs customer secrets automatically
- Cert Manager provisions a TLS certificate per customer

**To onboard a new customer**, simply add a file under the customers/ folder and push to GitHub. Everything else is automated.


---

## Repository Structure
```
mcsp-hello-world/
|-- app/                        # Node.js application source code
|-- k8s/                        # Kubernetes deployment manifests
|-- tekton/                     # Tekton pipeline and tasks
|-- gitops/                     # ArgoCD application configuration
|-- cert-manager/               # TLS certificate configuration
|-- rhacm/                      # RHACM governance policies
|-- customers/                  # Customer onboarding trigger files
|-- tenant-config/              # Customer environment templates
```

---

## Operators Used

| Operator | Purpose |
|---|---|
| RHACM | Cluster governance and placement policies |
| OpenShift GitOps | GitOps-based automated deployments |
| OpenShift Pipelines | CI/CD pipeline automation |
| Cert Manager | Automatic TLS certificate provisioning |
| External Secrets | Secure secret injection per customer |

---

