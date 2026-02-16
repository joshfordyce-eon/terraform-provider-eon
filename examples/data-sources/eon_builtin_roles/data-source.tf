# Use flat attributes instead of hardcoding built-in role UUIDs
data "eon_builtin_roles" "builtin" {}

# Example: Use in eon_idp_group.role_ids
# resource "eon_idp_group" "admins" {
#   idp_id            = "your-idp-id"
#   provider_group_id = "your-idp-group-id"
#
#   role_ids = [
#     data.eon_builtin_roles.builtin.global_admin,
#     data.eon_builtin_roles.builtin.admin,
#   ]
# }

output "builtin_roles" {
  description = "All built-in role UUIDs (for use in eon_idp_group.role_ids)"
  value = {
    global_admin  = data.eon_builtin_roles.builtin.global_admin
    global_viewer = data.eon_builtin_roles.builtin.global_viewer
    viewer        = data.eon_builtin_roles.builtin.viewer
    admin         = data.eon_builtin_roles.builtin.admin
    operator      = data.eon_builtin_roles.builtin.operator
  }
}

output "global_admin_id" {
  description = "Global Admin built-in role UUID"
  value       = data.eon_builtin_roles.builtin.global_admin
}
