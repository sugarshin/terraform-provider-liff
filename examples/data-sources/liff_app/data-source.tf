data "liff_app" "example" {
  liff_id = "1234567890-AbCdEfGh"
}

output "liff_app_url" {
  value = data.liff_app.example.view.url
}
