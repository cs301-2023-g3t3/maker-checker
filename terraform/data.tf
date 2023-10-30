data "aws_ssm_parameter" "rest_api_id" {
  name = "rest_api_id"
}

data "aws_iam_role" "existing_role" {
    name = "lambda-execution-role"
}

data "aws_ssm_parameter" "mongo_user" {
  name = "mongo_user"
}

data "aws_ssm_parameter" "mongo_password" {
  name = "mongo_password"
}

data "aws_ssm_parameter" "mongo_host" {
  name = "mongo_host"
}

data "aws_ssm_parameter" "queue_name" {
  name = "queue_name"
}


output "existing_role_arn" {
    value = data.aws_iam_role.existing_role.arn
}
