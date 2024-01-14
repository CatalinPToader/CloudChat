resource "kubernetes_namespace" "postgresql_namespace" {
  metadata {
    name = "postgresql-namespace"
  }
}

resource "kubernetes_namespace" "pgadmin_namespace" {
  metadata {
    name = "pgadmin-namespace"
  }
}

resource "kubernetes_namespace" "backend_namespace" {
  metadata {
    name = "backend-namespace"
  }
}

resource "kubernetes_namespace" "portainer_namespace" {
  metadata {
    name = "portainer-namespace"
  }
}

resource "kubernetes_namespace" "auth_namespace" {
  metadata {
    name = "auth-namespace"
  }
}
