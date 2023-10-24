resource "helm_release" "nginx_ingress_chart" {
  name       = "nginx-ingress-controller"
  namespace  = "default"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx-ingress-controller"
  set {
    name  = "service.type"
    value = "LoadBalancer"
  }
  set {
    name  = "service.annotations.kubernetes\\.digitalocean\\.com/load-balancer-id"
    value = var.loadbalancer_id
  }
}

resource "kubectl_manifest" "ingress" {
  depends_on = [ helm_release.nginx_ingress_chart ]

  yaml_body = file("${var.manifests_path}/ingress.yaml")
}