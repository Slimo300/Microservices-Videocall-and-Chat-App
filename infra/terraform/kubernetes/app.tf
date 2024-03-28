resource "helm_release" "nginx_ingress_chart" {
  name       = "nginx-ingress-controller"
  namespace  = "ingress-nginx"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx-ingress-controller"
  create_namespace = true

  set {
    name  = "service.type"
    value = "LoadBalancer"
  }
  set {
    name  = "service.annotations.kubernetes\\.digitalocean\\.com/load-balancer-id"
    value = var.loadbalancer_id
  }
}

resource "helm_release" "cert-manager-release" {
  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  version    = "v1.12.0"
  namespace  = "cert-manager"
  timeout    = 120

  set {
    name  = "createCustomResource"
    value = "true"
  }
  set {
    name  = "installCRDs"
    value = "true"
  }
}

resource "helm_release" "argo_release" {
  name = "argo-cd"
  namespace = "argo-cd"
  repository = "https://charts.bitnami.com/bitnami"
  chart = "argo-cd"
  create_namespace = true
}

resource "kubectl_manifest" "argo-application" {

  yaml_body = file("${var.manifests_path}/argo-application.yaml")

  depends_on = [ helm_release.argo_release, helm_release.cert-manager-release ]
}