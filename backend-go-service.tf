resource "kubernetes_service" "backend_go_service" {
  metadata {
    name = "backend"
    namespace = kubernetes_namespace.backend_namespace.metadata[0].name
  }

  spec {
    selector = {
      app = "backend"
    }

    port {
      port        = 8080
      target_port = 8080
      node_port   = 32080
    }

    type = "NodePort"
  }
}
