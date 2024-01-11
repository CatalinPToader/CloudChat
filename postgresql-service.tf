resource "kubernetes_service" "postgresql_service" {
  metadata {
    name      = "postgresql"
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
