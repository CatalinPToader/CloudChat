resource "kubernetes_service" "backend_go_service" {
  metadata {
    name      = "backend"
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
