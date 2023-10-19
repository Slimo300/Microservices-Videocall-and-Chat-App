resource "kubernetes_namespace" "cert-manager-namespace" {
  metadata {
    name = "cert-manager"
  }
}

resource "helm_release" "cert-manager-release" {
  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  version    = "v1.12.0"
  namespace  = "cert-manager"
  timeout    = 120
  depends_on = [
    kubernetes_namespace.cert-manager-namespace
  ]
  set {
    name  = "createCustomResource"
    value = "true"
  }
  set {
    name  = "installCRDs"
    value = "true"
  }

}

resource "kubectl_manifest" "cluster_issuer" {
  yaml_body = file("${var.manifests_path}/cert-manager/cluster-issuer.yaml")

  depends_on = [kubectl_manifest.digitalocean_token, helm_release.cert-manager-release]
}

resource "kubectl_manifest" "certificate" {
  yaml_body = file("${var.manifests_path}/cert-manager/certificate.yaml")

  depends_on = [kubectl_manifest.cluster_issuer]
}