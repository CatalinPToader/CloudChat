resource "kubernetes_service" "postgresql_service" {
  metadata {
    name = "postgresql"
    namespace = kubernetes_namespace.postgresql_namespace.metadata[0].name
  }

  spec {
    selector = {
      app = "postgresql"
    }

    port {
      port        = 5432
      target_port = 5432
    }
  }
}
