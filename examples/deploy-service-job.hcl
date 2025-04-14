job "hello-service" {
  datacenters = ["us-west-2"]
  type        = "service"

# Deploy the job on the "shared" cluster
  constraint {
    attribute = "${node.class}"
    operator  = "="
    value     = "shd"
  }

  group "web-group" {
    count = 1

    task "web" {
      driver = "docker"

      config {
        image = "hashicorp/http-echo"
        args  = [
          "-text=Hello from Nomad service job!"
        ]
      }

      resources {
        cpu    = 100
        memory = 128
      }

      service {
        name = "hello-service"
        port = "http"

        check {
          type     = "http"
          path     = "/"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }

    network {
      port "http" {
        static = 8080
      }
    }
  }
}