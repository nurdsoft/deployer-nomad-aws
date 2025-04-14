job "hello-batch" {
  datacenters = ["us-west-2"]
  type        = "batch"

  # Deploy the job on the "shared" cluster
  constraint {
    attribute = "${node.class}"
    operator  = "="
    value     = "shd"
  }

  group "echo-group" {
    task "echo-task" {
      driver = "docker"

      config {
        image   = "alpine"
        command = "echo"
        args    = ["Hello from Nomad via Lambda!"]
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}