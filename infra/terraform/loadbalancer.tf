resource "digitalocean_loadbalancer" "relrel_lb" {
  name   = "relrel-lb"
  region = "fra1"

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"

    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 22
    protocol = "tcp"
  }

  lifecycle {
    ignore_changes = [forwarding_rule]
  }
}