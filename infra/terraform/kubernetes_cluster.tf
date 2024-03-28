resource "digitalocean_kubernetes_cluster" "relrel_cluster" {
  name    = "relrel"
  region  = "fra1"
  version = "1.29.1-do.0"

  node_pool {
    name       = "relrel-pool"
    size       = "s-2vcpu-4gb"
    auto_scale = false
    node_count = 2
  }
}