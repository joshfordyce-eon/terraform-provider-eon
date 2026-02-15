# Example: List all IDP groups
data "eon_idp_groups" "all" {}

# Example: Filter groups by IDP
locals {
  # Groups keyed by Identity Provider ID
  groups_by_idp = {
    for idp_id in distinct([for g in data.eon_idp_groups.all.groups : g.idp_id]) :
    idp_id => [for g in data.eon_idp_groups.all.groups : g if g.idp_id == idp_id]
  }
}

# Example: Groups that have at least one role assigned
locals {
  groups_with_roles = [
    for g in data.eon_idp_groups.all.groups :
    g if length(g.role_ids) > 0
  ]
}

# Example: Output IDP group information
output "idp_groups_count" {
  description = "Total number of IDP groups"
  value       = length(data.eon_idp_groups.all.groups)
}

output "idp_groups_with_roles_count" {
  description = "Number of IDP groups that have roles assigned"
  value       = length(local.groups_with_roles)
}

output "idp_groups_summary" {
  description = "Summary of all IDP groups"
  value = {
    for g in data.eon_idp_groups.all.groups :
    g.id => {
      idp_id            = g.idp_id
      provider_group_id = g.provider_group_id
      role_count        = length(g.role_ids)
    }
  }
}

output "idp_group_ids_by_idp" {
  description = "IDP group IDs grouped by Identity Provider"
  value = {
    for idp_id, groups in local.groups_by_idp :
    idp_id => [for g in groups : g.id]
  }
}
