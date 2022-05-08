terraform {
    required_providers {
        indigo = {
            version = "0.0.1"
            source = "undefx.com/dev/indigo"
        }
    }
}

variable "indigo_api_key" {
    type = string
}

variable "indigo_api_secret" {
    type = string
}

provider "indigo" {
    host = "https://api.customer.jp"
    //api_key = var.indigo_api_key
    //api_secret = var.indigo_api_secret
}

resource "indigo_ssh_key" "key00" {
    name = "key00"
    key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC0DQty7U5izXgnzhIcgegby4EuV/BsAdb8BJZCCTFxBv5JttTV8+hd9v6XVXt+HKs2LEmRv1Bj2hw5VKV8JKVO2HBqFFRVqw4oTWPJhifboXO+WfrOy49/19nkBRVoTmK+vcRu+MaSd40vC2x8CYF0IizhOGNkJ5keKpCbllzO+nbWb7wIpr9lOevXsnAQ7fg2tihhAr3Y+CLnAJrnxHgYj9DNzB2GVbWKXeHhaPMmXIl5D6kKjdVCR7f47OXbNMp+cxUsCaT7P4dCWtyTwg2K3KFHH/Kr5oqRxJQa+SikhP0CylYTpX0fWOjLN+TjNwnvY+tAW5LXZ/h2HCZoiVkY81nda8raElV/rjBSEbpmpB0D5I7Ddaei3+4QA6BUucIxTlaKV06M+bCGroAwjfPjYt+XADm/ZHVIU7mHc0AIP2YJDB1AyRT8VXYag/xjDsbVYY/qOeYv6EHSie+h4glUdj9LjRzNZPrjIxT3CIcivle4B6QbX/CiJVy+y+aEAm0= user@example.com"
}


resource "indigo_instance" "inst00" {
    name = "inst00"
    ssh_key_id = indigo_ssh_key.key00.id
    region_id = 1
    os_id = 1
    plan_id = 1
}

resource "indigo_snapshot" "snapshot00" {
    name = "snapshot00"
    instance_id = indigo_instance.inst00.id
}

output "inst00_output" {
    value = indigo_instance.inst00
}

output "snapshot00_output" {
    value = indigo_snapshot.snapshot00
}
