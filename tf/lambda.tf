data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "pinbot" {
  name               = "pinbot"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}


// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "function_archive" {
  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

// create the lambda function from zip file
resource "aws_lambda_function" "pinbot" {
  function_name = "pinbot"
  description   = "Pinbot"
  role          = aws_iam_role.pinbot.arn
  handler       = "provided"
  memory_size   = 128

  filename         = local.archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  runtime = "provided.al2023"
  architectures = ["arm64"]
}

resource "aws_lambda_function_url" "endpoint" {
  authorization_type = "NONE"
  function_name      = aws_lambda_function.pinbot.function_name
  cors {
    allow_origins = ["*"]
    allow_methods = ["POST"]
  }
}

data "aws_iam_policy_document" "allow_lambda_logging" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }
}

// create a policy to allow writing into logs and create logs stream
resource "aws_iam_policy" "function_logging_policy" {
  name        = "AllowLambdaLoggingPolicy"
  description = "Policy for lambda cloudwatch logging"
  policy      = data.aws_iam_policy_document.allow_lambda_logging.json
}

// attach policy to out created lambda role
resource "aws_iam_role_policy_attachment" "lambda_logging_policy_attachment" {
  role       = aws_iam_role.pinbot.id
  policy_arn = aws_iam_policy.function_logging_policy.arn
}

// create log group in cloudwatch to gather logs of our lambda function
resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.pinbot.function_name}"
  retention_in_days = 7
}

output "url" {
  value = aws_lambda_function_url.endpoint.function_url
}