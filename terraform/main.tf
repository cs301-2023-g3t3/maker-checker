provider "aws" {
  region = "ap-southeast-1"
}

terraform {
  backend "s3" {
    bucket         = "itsa-terraform-states"
    key            = "makerchecker-terraform.tfstate"
    region         = "ap-southeast-1"
    dynamodb_table = "makerchecker-terraform-state-lock"
  }
}

resource "aws_lambda_function" "this" {
  function_name    = "makerchecker-api"
  runtime          = "go1.x"
  handler          = "main"
  role             = data.aws_iam_role.existing_role.arn
  filename         = "./main.zip"
  source_code_hash = filebase64sha256("./main.zip")
  timeout          = 10

  vpc_config {
    subnet_ids         = ["subnet-067ead2dfed1c40f8", "subnet-07722d2cb666693fd"]
    security_group_ids = ["sg-05a445100fc469383"]
  }

  environment {
    variables = {
      MONGO_USERNAME = data.aws_ssm_parameter.mongo_user.value
      MONGO_PASSWORD = data.aws_ssm_parameter.mongo_password.value
      MONGO_HOST     = data.aws_ssm_parameter.mongo_host.value
      QUEUE_NAME     = data.aws_ssm_parameter.queue_name.value
      JWT_SECRET     = data.aws_ssm_parameter.jwt_secret.value
      ENV            = "lambda"
      # REDIS_NODE_1   = data.aws_ssm_parameter.redis_node_1.value
      # REDIS_NODE_2   = data.aws_ssm_parameter.redis_node_2.value
    }
  }
}

resource "aws_api_gateway_resource" "root" {
  rest_api_id = data.aws_ssm_parameter.rest_api_id.value
  parent_id   = "9gy5jtm4yf"
  path_part   = "makerchecker"
}

resource "aws_api_gateway_resource" "this" {
  rest_api_id = data.aws_ssm_parameter.rest_api_id.value
  parent_id   = aws_api_gateway_resource.root.id
  path_part   = "{proxy+}"
}

resource "aws_lambda_permission" "this" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:ap-southeast-1:345215350058:${data.aws_ssm_parameter.rest_api_id.value}/*/*/makerchecker/*"
}

resource "aws_api_gateway_integration" "this" {
  rest_api_id             = data.aws_ssm_parameter.rest_api_id.value
  resource_id             = aws_api_gateway_resource.this.id
  http_method             = aws_api_gateway_method.this.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.this.invoke_arn
}

resource "aws_api_gateway_method" "this" {
  rest_api_id      = data.aws_ssm_parameter.rest_api_id.value
  resource_id      = aws_api_gateway_resource.this.id
  http_method      = "ANY"
  authorization    = "CUSTOM"
  authorizer_id    = "kjkxid"
  api_key_required = false
}
