
locals {
  binary_name = "bootstrap"
  lambda_name = "nomad-deploy"
  lambda_path = "${path.module}/../../build/lambda.zip"
  runtime     = "provided.al2023"
}

module "lambda_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.1.0"

  name                = "nomad-deploy-lambda-sg"
  vpc_id              = data.aws_vpc.vpc.id
  ingress_cidr_blocks = ["${data.aws_vpc.vpc.cidr_block}"]
  egress_rules        = ["nomad-http-tcp"] # Outgoing traffic only to Nomad API Endpoint
  ingress_rules       = ["all-all"]
}

module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.5.0"
  count   = 1

  function_name          = local.lambda_name
  handler                = local.binary_name
  runtime                = local.runtime
  publish                = true
  create_package         = false
  local_existing_package = local.lambda_path
  timeout                = 200
  attach_network_policy  = true
  vpc_subnet_ids         = toset(data.aws_subnets.private.ids)
  vpc_security_group_ids = [module.lambda_sg.security_group_id]

  environment_variables = {
    NOMAD_ADDR = var.nomad_addr
  }

  #   allowed_triggers = {
  #     AllowExecutionFromAPIGateway = {
  #       service    = "apigateway"
  #       source_arn = "${module.httpapi.api_execution_arn}/*/*"
  #     }
  #   }
  #   tags                   = merge(var.tags, { component = local.functions[count.index].component })
}
