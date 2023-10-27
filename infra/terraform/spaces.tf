resource "digitalocean_spaces_bucket" "relrelspaces" {
  name = "relrelspaces"
  
  region = "fra1"
  acl = "public-read"
}

resource "digitalocean_spaces_bucket_cors_configuration" "relrelspaces-cors" {
  bucket = "relrelspaces"
  region = "fra1"

  cors_rule {
    allowed_methods = ["PUT", "GET", "DELETE"]
    allowed_origins = [ digitalocean_domain.relrel_domain.name ]
    max_age_seconds = 3000
    allowed_headers = [ "Authorization", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "accept", "origin", "Cache-Control", " X-Requested-With", "X-AMZ-ACL" ]
  }
}