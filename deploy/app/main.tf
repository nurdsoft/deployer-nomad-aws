locals {
  binary_name      = "bootstrap"
  lambda_name      = "nomad-deploy"
  lambda_path      = "${path.module}/../../build/lambda.zip"
  runtime          = "provided.al2023"
  api_gateway_name = "deployer-nomad-api"
  lambda_sg_name   = "deployer-nomad-lambda-sg"
}

# Lambda security group
module "lambda_sg" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "5.1.0"

  name   = local.lambda_sg_name
  vpc_id = data.aws_vpc.vpc.id

  ingress_cidr_blocks = [data.aws_vpc.vpc.cidr_block]
  ingress_rules       = ["all-all"]
  egress_rules        = ["nomad-http-tcp"]
}

# Lambda function (Go binary packaged externally)
module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.5.0"

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
    API_KEYS   = join(",", var.api_keys)
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api-gateway.api_execution_arn}/*/*"
    }
  }
}

# API Gateway (HTTP)
module "api-gateway" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "5.0.0"

  name                  = local.api_gateway_name
  protocol_type         = "HTTP"
  stage_name            = "v1"
  create_domain_name    = false
  create_domain_records = false
  create_certificate    = false

  cors_configuration = {
    allow_headers = ["x-api-key"]
    allow_methods = ["POST"]
    allow_origins = ["*"]
  }

  routes = {
    "POST /deploy" = {
      integration = {
        uri                    = module.lambda_function.lambda_function_arn
        payload_format_version = "2.0"
        timeout_milliseconds   = 12000
      }
    }
  }

  stage_access_log_settings = {
    create_log_group            = true
    log_group_retention_in_days = 7

    format = jsonencode({
      requestId      = "$context.requestId"
      sourceIP       = "$context.identity.sourceIp"
      protocol       = "$context.protocol"
      routeKey       = "$context.routeKey"
      status         = "$context.status"
      responseLength = "$context.responseLength"
    })
  }
}
