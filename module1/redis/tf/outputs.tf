output "public_ip" {
  description = "vm public ip"
  value       = tencentcloud_instance.web[0].public_ip
}

output "kube_config" {
  description = "kubeconfig"
  value       = "${path.module}/config.yaml"
}

output "password" {
  description = "vm ubuntu password"
  value       = var.password
}