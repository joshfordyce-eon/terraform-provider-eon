# Example: List all roles (built-in and custom)
data "eon_roles" "all" {}

# Example: Filter built-in vs custom roles
locals {
  built_in_roles = [
    for r in data.eon_roles.all.roles :
    r if r.is_built_in_role
  ]

  custom_roles = [
    for r in data.eon_roles.all.roles :
    r if !r.is_built_in_role
  ]
}

# Example: Map of role display name to role ID (for auditing/display only; display names can change)
# For assigning built-in roles to IDP groups, use the eon_builtin_roles data source (stable keys) instead.
locals {
  role_ids_by_name = {
    for r in data.eon_roles.all.roles :
    r.name => r.id
  }
}

# Example: Roles that have access conditions (custom roles with scoped permissions)
locals {
  roles_with_access_conditions = [
    for r in data.eon_roles.all.roles :
    r if length(r.access_conditions) > 0
  ]
}

# Example: Output role information
output "roles_count" {
  description = "Total number of roles"
  value       = length(data.eon_roles.all.roles)
}

output "built_in_roles_count" {
  description = "Number of built-in Eon roles"
  value       = length(local.built_in_roles)
}

output "custom_roles_count" {
  description = "Number of custom roles"
  value       = length(local.custom_roles)
}

output "roles_summary" {
  description = "Summary of all roles (includes access_conditions when present)"
  value = {
    for r in data.eon_roles.all.roles :
    r.name => {
      id                = r.id
      is_built_in_role  = r.is_built_in_role
      permission_count  = length(r.permission_grants)
      access_conditions = r.access_conditions
    }
  }
}

# Example: Output roles that use access conditions (e.g. to audit scoped permissions)
output "roles_with_access_conditions" {
  description = "Roles that define access conditions (id, effect, expression)"
  value = {
    for r in local.roles_with_access_conditions :
    r.name => {
      id                = r.id
      access_conditions = r.access_conditions
    }
  }
}

output "roles_with_access_conditions_count" {
  description = "Number of roles that have at least one access condition"
  value       = length(local.roles_with_access_conditions)
}

output "role_ids_by_name" {
  description = "Map of role display name to role ID (for auditing/display only; for built-in roles in eon_idp_group use eon_builtin_roles instead)"
  value       = local.role_ids_by_name
}

# For assigning built-in roles (e.g. Admin) to IDP groups, use the eon_builtin_roles data source
# (e.g. role_ids = [data.eon_builtin_roles.builtin.admin]), not lookup by display name.
