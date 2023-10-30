terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    ssh = {
      source  = "loafoe/ssh"
      version = "2.6.0"
    }
  }
}

variable "do_token" {
  sensitive = true
}
variable "pvt_key" {
  sensitive = true
}
variable "user" {
  sensitive = true
}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "relrel_ssh" {
  name = "relrel_ssh"
}

resource "digitalocean_droplet" "turn_droplet" {

  image = "ubuntu-22-04-x64"
  name  = "turn"

  region = "fra1"
  size   = "s-1vcpu-1gb"
  ssh_keys = [
    data.digitalocean_ssh_key.relrel_ssh.id
  ]
}

resource "digitalocean_domain" "turn_around" {
  name       = "turn-around.pro"
  ip_address = digitalocean_droplet.turn_droplet.ipv4_address
}

resource "ssh_resource" "coturn_config" {
  host = digitalocean_droplet.turn_droplet.ipv4_address

  user        = "root"
  private_key = file(var.pvt_key)

  when = "create" # Default

  commands = [
    "sudo apt install docker.io -y",
    "sudo apt install coturn -y",
    "docker run slimo300/turn-conf-generator ${digitalocean_droplet.turn_droplet.ipv4_address} ${var.user} > /etc/turnserver.conf",
  ]
}

resource "ssh_resource" "certbot_config" {
  host = digitalocean_droplet.turn_droplet.ipv4_address

  user        = "root"
  private_key = file(var.pvt_key)

  when = "create" # Default

  commands = [
    # Issuing a certififcate
    "sudo NEEDRESTART_MODE=a apt install certbot -y",
    "sudo certbot certonly --noninteractive --agree-tos --register-unsafely-without-email --standalone --preferred-challenges http -d ${digitalocean_domain.turn_around.name}",

    # Making certificate visible to coturn
    "sudo mkdir -p /etc/coturn/certs",
    "sudo chown -R turnserver:turnserver /etc/coturn",
    "sudo chmod -R 700 /etc/coturn",
    "sudo cp /etc/letsencrypt/live/${digitalocean_domain.turn_around.name}/fullchain.pem /etc/coturn/certs/${digitalocean_domain.turn_around.name}.cert",
    "sudo cp /etc/letsencrypt/live/${digitalocean_domain.turn_around.name}/privkey.pem /etc/coturn/certs/${digitalocean_domain.turn_around.name}.key",
    "sudo chown turnserver /etc/coturn/certs/${digitalocean_domain.turn_around.name}.cert /etc/coturn/certs/${digitalocean_domain.turn_around.name}.key",
    "sudo chmod 400 /etc/coturn/certs/${digitalocean_domain.turn_around.name}.cert /etc/coturn/certs/${digitalocean_domain.turn_around.name}.key",

    # Restart coturn
    "sudo systemctl start coturn"
  ]

  depends_on = [digitalocean_domain.turn_around, ssh_resource.coturn_config]

}