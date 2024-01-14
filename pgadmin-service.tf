resource "kubernetes_service" "pgadmin_service" {
  metadata {
    name = "pgadmin"
    namespace = kubernetes_namespace.pgadmin_namespace.metadata[0].name
  }

  spec {
    selector = {
      app = "pgadmin"
    }

    port {
      port        = 80
      target_port = 80
      node_port   = 32000 # You can choose a different port if needed
    }

    type = "NodePort"
  }
}
