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
}

data "kubectl_path_documents" "services" {
    pattern = "${var.manifests_path}/services/*.yaml"
}

resource "kubectl_manifest" "services" {
    for_each  = toset(data.kubectl_path_documents.services.documents)
    yaml_body = each.value
}