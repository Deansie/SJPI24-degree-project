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

### SRE Rules and Best Practices
- [Google SRE Book](https://sre.google/sre-book/table-of-contents/)  
  Reference for principles and methodologies applied to operational reliability, observability, and automation.

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

## Author
**Dennis Andersen**  
Higher Vocational Education Diploma in System Development (Java & JavaScript)  
SJPI24 – Lernia Yrkeshögskola
