resource "helm_release" "kafka" {
  name       = "kafka"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "kafka"
  version    = "v20.0.6"
}

data "kubectl_path_documents" "dbs" {
    pattern = "${var.manifests_path}/dbs/*.yaml"
}

resource "kubectl_manifest" "dbs" {
    for_each  = toset(data.kubectl_path_documents.dbs.documents)
    yaml_body = each.value

    depends_on = [ kubectl_manifest.redis-creds, kubectl_manifest.mysql-creds ]
}

data "kubectl_path_documents" "services" {
    pattern = "${var.manifests_path}/services/*.yaml"
}

resource "kubectl_manifest" "services" {
    for_each  = toset(data.kubectl_path_documents.services.documents)
    yaml_body = each.value

    depends_on = [ 
        kubectl_manifest.dbs, 
        kubectl_manifest.brevo-creds, 
        kubectl_manifest.digitalocean_spaces,
        kubectl_manifest.private-key,
        kubectl_manifest.public-key,
        kubectl_manifest.refresh-secret,
        kubectl_manifest.turn-credentials,
    ]
}