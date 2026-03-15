resource "liff_app" "example" {
  description            = "My LIFF App"
  permanent_link_pattern = "concat"
  bot_prompt             = "normal"
  scope                  = ["openid", "profile"]

  view {
    type        = "full"
    url         = "https://example.com"
    module_mode = false
  }

  features {
    qr_code = true
  }
}
