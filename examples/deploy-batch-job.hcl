job "hello-batch" {
  datacenters = ["dc1"]
  type        = "batch"
  
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