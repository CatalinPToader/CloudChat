resource "kubernetes_deployment" "backend_deployment" {
  metadata {
    name = "backend"
    namespace = kubernetes_namespace.backend_namespace.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "backend"
      }
    }

    template {
      metadata {
        labels = {
          app = "backend"
        }
      }

      spec {
        container {
          image = "localhost:5001/cloudchat-backend:latest"
          name  = "backend"
          port {
            container_port = 8080
          }
        }
      }
    }
  }
}
