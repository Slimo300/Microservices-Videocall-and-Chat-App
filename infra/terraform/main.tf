terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

variable "do_token" {}

provider "digitalocean" {
  token = var.do_token
}

module "kubernetes" {
  source = "./kubernetes"

  cluster_endpoint       = digitalocean_kubernetes_cluster.relrel_cluster.endpoint
  cluster_token          = digitalocean_kubernetes_cluster.relrel_cluster.kube_config[0].token
  cluster_ca_certificate = digitalocean_kubernetes_cluster.relrel_cluster.kube_config[0].cluster_ca_certificate
  loadbalancer_id       = digitalocean_loadbalancer.relrel_lb.id

  manifests_path = "../k8s/prod"
}