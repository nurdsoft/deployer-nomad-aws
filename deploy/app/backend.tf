terraform {
  backend "s3" {
    bucket = "aws-usw2-nurdsoft-terraform-state"
    key    = "usw2/deployer-nomad/terraform.tfstate"
    region = "us-west-2"
  }
}
