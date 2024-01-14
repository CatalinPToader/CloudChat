resource "kubernetes_deployment" "auth_deployment" {
  metadata {
    name = "auth"
    namespace = kubernetes_namespace.auth_namespace.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "auth"
      }
    }

    template {
      metadata {
        labels = {
          app = "auth"
        }
      }

      spec {
        container {
          image = "localhost:5001/auth:latest"
          name  = "auth"
          port {
            container_port = 9000
          }
        }
      }
    }
  }
}