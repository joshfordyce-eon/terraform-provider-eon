terraform {
  required_providers {
    eon = {
      source = "eon-io/eon"
    }
  }
}

# Example: IDP group that assigns custom roles (reference eon_role resources)
resource "eon_role" "read_only" {
  name = "Read-only (IDP Example)"

  permission_grants = [
    { permission = "dashboard.view" },
    { permission = "inventory.view" },
  ]
}

resource "eon_role" "ops" {
  name = "Ops (IDP Example)"

  permission_grants = [
    { permission = "dashboard.view" },
    { permission = "inventory.view" },
    { permission = "snapshots.take_or_convert" },
    { permission = "jobs.view" },
  ]
}

resource "eon_idp_group" "read_only" {
  idp_id            = "your-idp-id"
  provider_group_id = "your-idp-group-id-for-read-only"

  role_ids = [
    eon_role.read_only.id,
  ]
}

# Example: IDP group with multiple roles
resource "eon_idp_group" "ops" {
  idp_id            = "your-idp-id"
  provider_group_id = "your-idp-group-id-for-ops"

  role_ids = [
    eon_role.read_only.id,
    eon_role.ops.id,
  ]
}

# ---------------------------------------------------------------------------
# Assigning EON built-in roles to IDP groups
# ---------------------------------------------------------------------------
# Built-in roles (Global Admin, Global Viewer, Viewer, Admin, Operator)
# can be assigned by raw UUID or by stable key via the eon_builtin_roles data source.

# Example: Assign built-in roles by raw role ID (UUID)
resource "eon_idp_group" "admins_by_id" {
  idp_id            = "your-idp-id"
  provider_group_id = "your-idp-group-id-for-admins-by-uuid"

  role_ids = [
    "379a1104-838a-4bf3-af96-da3af27c5712", # Global Admin
    "a675e456-8602-4550-9c65-66583404e0d6", # Admin
  ]
}

# Example: Assign built-in roles by stable key (recommended; keys are stable if display names change)
data "eon_builtin_roles" "builtin" {}

resource "eon_idp_group" "viewers_by_key" {
  idp_id            = "your-idp-id"
  provider_group_id = "your-idp-group-id-for-viewers-by-key"

  role_ids = [
    data.eon_builtin_roles.builtin.global_viewer,
  ]
}

resource "eon_idp_group" "operators_by_key" {
  idp_id            = "your-idp-id"
  provider_group_id = "your-idp-group-id-for-operators-by-key"

  role_ids = [
    data.eon_builtin_roles.builtin.viewer,
    data.eon_builtin_roles.builtin.operator,
  ]
}

# For custom roles, create eon_role resources and reference them in role_ids.
# For built-in roles, use data.eon_builtin_roles.builtin.<attribute> (e.g. .global_admin, .viewer; see above).
