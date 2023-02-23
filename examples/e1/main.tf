terraform {
  required_providers {
    education = {
      source = "raygecao.cn/test/education-demo"
    }
  }
}

provider "education" {
  user = "admin"
  passwd = "admin123"
}

resource "education_teacher" "bob" {
  id = 10001
  name = "bob"
  subject = "english"
  salary = 16000
  organ  = data.education_school.organ.name
}

output "bob" {
  value = education_teacher.bob
}

data "education_school" "organ" {}
