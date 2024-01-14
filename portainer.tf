provider "kubernetes" {
  config_path = "./mykubeconfig" # Path to your kubeconfig file
}

resource "kubernetes_deployment" "portainer_deployment" {
  metadata {
    name = "portainer"
    namespace = kubernetes_namespace.portainer_namespace.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "portainer"
      }
    }

    template {
      metadata {
        labels = {
          app = "portainer"
        }
      }

      spec {
        container {
          name  = "portainer"
          image = "portainer/portainer-ce:latest"
          port {
            container_port = 9000
          }
          command = [
            "/portainer",
            // "--no-auth" - deprecated 
            "--admin-password",
            "$2y$05$w5wsvlEDXxPjh2GGfkoe9.At0zj8r7DeafAkXXeubs0JnmxLjyw/a",
          ]
        }
      }
    }
  }
}

resource "kubernetes_cluster_role" "portainer_cluster_role" {
  metadata {
    name = "portainer-cluster-role"
  }

  rule {
    api_groups = ["*"]
    resources  = ["*"]
    verbs      = ["*"]
  }
}

resource "kubernetes_cluster_role_binding" "portainer_cluster_role_binding" {
  metadata {
    name = "portainer-cluster-role-binding"
  }

  subject {
    kind = "ServiceAccount"
    name = "default"
    namespace = kubernetes_namespace.portainer_namespace.metadata[0].name
  }

  role_ref {
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.portainer_cluster_role.metadata[0].name
    api_group = "rbac.authorization.k8s.io"
  }
}