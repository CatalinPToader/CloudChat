resource "kubernetes_service" "portainer_service" {
  metadata {
    name = "portainer"
    namespace = kubernetes_namespace.portainer_namespace.metadata[0].name
  }

  spec {
    selector = {
      app = "portainer"
    }

    port {
      port        = 9000
      target_port = 9000
      node_port   = 31000 # this can be changed to any suitable port
    }

    type = "NodePort"
  }
}