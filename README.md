# SJPI24-Degree-Project

This repository contains the final degree project for my **Higher Vocational Education Diploma in System Development with Java and JavaScript**.

The project focuses on **automation of infrastructure for Kubernetes-based environments** within an **SRE (Site Reliability Engineering)** context, emphasizing reliability, security, and automation in CI/CD workflows.

---

## Project Overview

**Title:**  
Go-based CLI Tool with Cobra and Fang for YAML Validation, DNS Management, and Ingress Automation using Cert-Manager integrated via Jenkins CI/CD Pipelines

**Summary:**  
The project addresses common challenges in Kubernetes infrastructure management, such as faulty YAML manifests, inefficient manual DNS configuration, and limited TLS coverage with Cloudflare’s proxy.  

To solve these issues, a **Go-based CLI tool** was developed using **Cobra** and **Fang**. The tool:
- Validates Kubernetes YAML manifests against SRE-inspired rules  
- Automates DNS record creation in Cloudflare  
- Configures Kubernetes Ingress with Cert-Manager for **end-to-end (E2E) TLS encryption**  
- Integrates with **Jenkins** for continuous deployment automation  

The solution aims to eliminate deployment errors, ensure consistent and secure configurations, and streamline CI/CD processes through automation and adherence to SRE best practices such as **logging, metrics, and idempotency**.

---

## Knowledge References

- [https://sre.google/sre-book/table-of-contents/](https://sre.google/sre-book/table-of-contents/) - Reference for principles and methodologies applied to operational reliability, observability, and automation.

- [https://cobra.dev](https://cobra.dev) - Official site for Cobra CLI framework.

- [https://github.com/charmbracelet/fang](https://github.com/charmbracelet/fang) - Fang library for Go CLI flag and argument handling.  

- [https://go.dev/doc/](https://go.dev/doc/) - Official Go documentation and guides.  

- [https://cert-manager.io/docs/](https://cert-manager.io/docs/) - Documentation for Cert-Manager (TLS automation) in Kubernetes.  

- [https://github.com/cloudflare/cloudflare-go](https://github.com/cloudflare/cloudflare-go) - Cloudflare Go API library for DNS automation.  

- [https://github.com/kubernetes/client-go](https://github.com/kubernetes/client-go) – Official Kubernetes Go client for working with clusters and API objects.  

- [https://pkg.go.dev/sigs.k8s.io/yaml](https://pkg.go.dev/sigs.k8s.io/yaml) – YAML/JSON parsing utilities for Go, used to unmarshal and validate Kubernetes manifests.  

- [https://pkg.go.dev/k8s.io/apimachinery](https://pkg.go.dev/k8s.io/apimachinery) – Kubernetes API machinery for runtime objects, decoding, and schema handling.  

- [https://pkg.go.dev/k8s.io/api](https://pkg.go.dev/k8s.io/api) – Core Kubernetes API object definitions (Deployments, Services, Ingress, etc.).  

- [https://github.com/kubernetes/kube-openapi](https://github.com/kubernetes/kube-openapi) – OpenAPI schema validation used by Kubernetes for API consistency and validation.  

- [https://github.com/instrumenta/kubeval](https://github.com/instrumenta/kubeval) – YAML schema validation tool for Kubernetes manifests (reference for building custom validation).  

---

## Technologies Used
- Go (Cobra, Fang)
- Kubernetes
- Jenkins (CI/CD)
- Cloudflare (DNS)
- Cert-Manager / Let’s Encrypt
- Docker
- Proxmox

---

## UML-Diagram

![](assets/images/use-case-diagram.png)

---

## Author
**Dennis Andersen**  
Higher Vocational Education Diploma in System Development (Java & JavaScript)  
SJPI24 – Lernia Yrkeshögskola
