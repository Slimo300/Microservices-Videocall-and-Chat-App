terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.0.1"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}

variable "cluster_endpoint" {
  type        = string
  description = "Cluster Endpoint"
}

variable "cluster_ca_certificate" {
  type        = string
  description = "Certificate for TLS authentication"
}

variable "cluster_token" {
  type        = string
  description = "Cluster Token"
}

variable "manifests_path" {
  type        = string
  description = "Path where k8s manifests are stored"
  default     = "../../k8s/prod"
}

variable "loadbalancer_id" {
  type        = string
  description = "ID for provisioned load balancer"
}

provider "kubernetes" {
  host  = var.cluster_endpoint
  token = var.cluster_token
  cluster_ca_certificate = base64decode(
    var.cluster_ca_certificate
  )
}

provider "helm" {
  kubernetes {
    host  = var.cluster_endpoint
    token = var.cluster_token
    cluster_ca_certificate = base64decode(
      var.cluster_ca_certificate
    )
  }
}

provider "kubectl" {
  host  = var.cluster_endpoint
  token = var.cluster_token
  cluster_ca_certificate = base64decode(
    var.cluster_ca_certificate
  )
}