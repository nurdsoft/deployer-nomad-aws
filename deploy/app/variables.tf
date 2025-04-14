variable "vpc_name" {
  type = string
}

variable "region" {
  type    = string
  default = "us-west-2"
}

variable "nomad_addr" {
  type = string
}

variable "api_keys" {
  type        = list(string)
  description = "List of api keys"
  default     = []
}

# variable "default_tags" {
#   description = "The default set of tags to pass to all resources"
#   type        = map(string)
#   default = {
#     terraform = "true"
#   }
# }

# variable "custom_tags" {
#   description = "The custom set of tags to pass to all resources. This is merged with default_tags and passed via tfvars"
#   type        = map(string)
#   default     = {}
# }
