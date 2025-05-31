output "instance_ip" {
  value = google_compute_instance.go_app_instance.network_interface[0].access_config[0].nat_ip
}


output "cloud_sql_ip" {
  value = data.google_sql_database_instance.existing.ip_address[0].ip_address
}