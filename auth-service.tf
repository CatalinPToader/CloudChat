resource "kubernetes_service" "auth_service" {
  metadata {
    name = "auth"
    namespace = kubernetes_namespace.auth_namespace.metadata[0].name
  }

  spec {
    selector = {
      app = "auth"
    }

    port {
      port        = 9000
      target_port = 9000
      node_port   = 32700
    }

    type = "NodePort"
  }
}
