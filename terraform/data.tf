data "aws_ssm_parameter" "rest_api_id" {
  name = "rest_api_id"
}

data "aws_iam_role" "existing_role" {
    name = "lambda-execution-role"
}

data "aws_ssm_parameter" "db_user" {
  name = "db_user"
}

data "aws_ssm_parameter" "db_password" {
  name = "db_password"
}

data "aws_ssm_parameter" "mongo_host" {
  name = "mongo_host"
}

output "existing_role_arn" {
    value = data.aws_iam_role.existing_role.arn
}
