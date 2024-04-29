terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}

resource "docker_image" "golang_chatbot" {
  name = "golang-chatbot:v1.0"
  build {
    context    = "../"
    dockerfile = "docker/chatbot.Dockerfile"
  }
}

resource "docker_image" "redis" {
  name = "redis:custom-v1.0"
  build {
    context    = "../docker"
    dockerfile = "redis.Dockerfile"
  }
}

resource "docker_network" "app_network" {
  name = "app_network"
}

resource "docker_container" "webserver" {
  image = docker_image.golang_chatbot.name
  name  = "webserver"
  ports {
    internal = 8080
    external = 8080 // Port exposed on the host
  }
  networks_advanced {
    name = docker_network.app_network.name
  }
}

resource "docker_container" "redis-server" {
  image = docker_image.redis.name
  name  = "redis-server"
  ports {
    internal = 6379
    external = 6379
  }
  networks_advanced {
    name = docker_network.app_network.name
  }
}
