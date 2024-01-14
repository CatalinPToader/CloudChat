resource "kubernetes_config_map" "pgadmin_configmap" {
  metadata {
    name = "pgadmin-config"
    namespace = kubernetes_namespace.pgadmin_namespace.metadata[0].name
  }

  data = {
    servers = file("./pgadmin-conf.json")
  }
}

resource "kubernetes_deployment" "pgadmin_deployment" {
  metadata {
    name = "pgadmin"
    namespace = kubernetes_namespace.pgadmin_namespace.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "pgadmin"
      }
    }

    template {
      metadata {
        labels = {
          app = "pgadmin"
        }
      }

      spec {
        container {
          name  = "pgadmin"
          image = "dpage/pgadmin4:latest"
          port {
            container_port = 80
          }
          env {
            name  = "PGADMIN_DEFAULT_EMAIL"
            value = "admin@test.com"
          }
          env {
            name  = "PGADMIN_DEFAULT_PASSWORD"
            value = "admin"
          }
          env {
            name  = "PGADMIN_CONFIG_SERVER_MODE"
            value = "False"
          }
          env {
            name  = "PGADMIN_CONFIG_MASTER_PASSWORD_REQUIRED"
            value = "False"
          }
          env {
            name  = "PGADMIN_SERVER_JSON_FILE"
            value = "/pgadmin4/config/servers"
          }
          volume_mount {
            name       = "pgadmin-config-volume"
            mount_path = "/pgadmin4/config"
          }
          command = [
            "/bin/sh",
            "-c",
            "/bin/echo 'postgresql.postgresql-namespace:5432:custom_db:custom_user:custom_passwd' > /tmp/pgpassfile && chmod 600 /tmp/pgpassfile && /entrypoint.sh",
          ]
        }
        volume {
          name = "pgadmin-config-volume"
          config_map {
            name = kubernetes_config_map.pgadmin_configmap.metadata[0].name
          }
        }
      }
    }
  }
}
