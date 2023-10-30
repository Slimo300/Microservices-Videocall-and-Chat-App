resource "digitalocean_domain" "relrel_domain" {
  name = "relrel.org"
}

resource "digitalocean_record" "relrel_a_record" {
  domain = digitalocean_domain.relrel_domain.name
  type   = "A"
  ttl    = 60
  name   = "@"
  value  = digitalocean_loadbalancer.relrel_lb.ip

  # depends_on = [module.kubernetes]
}

resource "digitalocean_record" "CNAME_www" {
  domain = digitalocean_domain.relrel_domain.name
  type   = "CNAME"
  name   = "www"
  value  = "@"
  ttl    = 30
}

resource "digitalocean_record" "CNAME_api" {
  domain = digitalocean_domain.relrel_domain.name
  type   = "CNAME"
  name   = "api"
  value  = "@"
  ttl    = 30
}