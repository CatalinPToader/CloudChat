resource "kubernetes_config_map" "psql_configmap" {
  metadata {
    name = "psql-config"
  }

  data = {
    "custom-db-init.sql" = file("./custom-db-init.sql")
  }
}

resource "kubernetes_deployment" "postgresql_deployment" {
  metadata {
    name      = "postgresql"
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "postgresql"
      }
    }

    template {
      metadata {
        labels = {
          app = "postgresql"
        }
      }

      spec {
        container {
          name  = "postgresql"
          image = "postgres:latest"
          port {
            container_port = 5432
          }
          env {
            name  = "POSTGRES_PASSWORD"
            value = "custom_passwd"
          }
          env {
            name  = "POSTGRES_USER"
            value = "custom_user"
          }
          env {
            name  = "POSTGRES_DB"
            value = "custom_db"
          }
          volume_mount {
            name       = "psql-config-volume"
            mount_path = "/docker-entrypoint-initdb.d/initdb"
          }
          startup_probe {
            exec {
              command = [
                "bin/sh", 
                "-c", 
                "psql -h localhost -U custom_user -d custom_db -f /docker-entrypoint-initdb.d/initdb/custom-db-init.sql"
              ]
            }
          }
        }
        volume {
          name = "psql-config-volume"
          config_map {
            name = kubernetes_config_map.psql_configmap.metadata[0].name
          }
        }
      }
    }
  }
}
