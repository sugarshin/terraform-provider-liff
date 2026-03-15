data "liff_apps" "all" {}

output "liff_app_ids" {
  value = data.liff_apps.all.apps[*].liff_id
}
