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

# resource "kubernetes_ingress" "default_cluster_ingress" {
#   depends_on = [
#     helm_release.nginx_ingress_chart,
#   ]
#   metadata {
#     name      = "relrel-ingress"
#     namespace = "default"
#     annotations = {
#       "kubernetes.io/ingress.class"          = "nginx"
#       "ingress.kubernetes.io/rewrite-target" = "/"
#       "cert-manager.io/cluster-issuer"       = "letsencrypt-production"
#       "nginx.ingress.kubernetes.io/use-regex" : "true"
#     }
#   }
#   spec {
#     ingress_class_name = "nginx"
#     rule {
#       host = "api.relrel.org"
#       http {
#         path {
#           path = "/groups/?(.*)"
#           backend {
#             service_name = "group-service"
#             service_port = 8080
#           }
#         }
#         path {
#           path = "/messages/?(.*)"
#           backend {
#             service_name = "message-service"
#             service_port = 8080
#           }
#         }
#         # path {
#         #   path = "/search/?(.*)"
#         #   backend {
#         #     service_name = "search-service"
#         #     service_port = 8080
#         #   }
#         # }
#         path {
#           path = "/users/?(.*)"
#           backend {
#             service_name = "user-service"
#             service_port = 8080
#           }
#         }
#         path {
#           path = "/ws/?(.*)"
#           backend {
#             service_name = "ws-service"
#             service_port = 8080
#           }
#         }
#         path {
#           path = "/video-call/?(.*)/websocket"
#           backend {
#             service_name = "webrtc-gateway-service"
#             service_port = 8080
#           }
#         }
#         path {
#           path = "/video-call/?(.*)/accessCode"
#           backend {
#             service_name = "webrtc-service"
#             service_port = 8080
#           }
#         }
#       }
#     }
#     rule {
#       host = "www.relrel.org"
#       http {
#         path {
#           path = "/"
#           backend {
#             service_name = "frontend"
#             service_port = 80
#           }
#         }
#       }
#     }
#     tls {
#       secret_name = "relrel-tls"
#       hosts = [
#         "www.relrel.org",
#         "api.relrel.org"
#       ]
#     }
#   }
# }