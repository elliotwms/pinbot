variable "stack" {
  type = string
}

variable "aws_access_key" {
  type      = string
  sensitive = true
}

variable "aws_secret_key" {
  type      = string
  sensitive = true
}

variable "discord_application_id" {
  type        = string
  description = "Discord Application ID"
}

variable "discord_token" {
  type        = string
  sensitive   = true
  description = "Discord bot token"
}

variable "discord_public_key" {
  type        = string
  sensitive   = true
  description = "Discord bot public key"
}