terraform {
  required_providers {
    kickflow = {
      source  = "hashicorp.com/edu/kickflow"
      version = "1.0.0"
    }
  }
}

provider "kickflow" {
}

data "kickflow_user" "example" {
  user_id = "1702fa8d-86b7-44c0-9205-b40c62e142cf"
}

output "user_email" {
  value = data.kickflow_user.example.email
}
