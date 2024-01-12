resource "kubernetes_deployment" "backend_deployment" {
  metadata {
    name = "backend"
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
        }
      }
    }
  }
}
