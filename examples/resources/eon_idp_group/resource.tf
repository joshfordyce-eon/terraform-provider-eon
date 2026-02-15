terraform {
  required_providers {
    eon = {
      source = "eon-io/eon"
    }
  }
}

# Example: IDP group that assigns custom roles (reference eon_role resources)
resource "eon_role" "viewer" {
  name = "Viewer (IDP Example)"

  permission_grants = [
    { permission = "dashboard.view" },
    { permission = "inventory.view" },
  ]
}

resource "eon_role" "operator" {
  name = "Operator (IDP Example)"

  permission_grants = [
    { permission = "dashboard.view" },
    { permission = "inventory.view" },
    { permission = "snapshots.take_or_convert" },
    { permission = "jobs.view" },
  ]
}

resource "eon_idp_group" "viewers" {
  idp_id            = "your-idp-id"
  provider_group_id = "your-idp-group-id-for-viewers"

  role_ids = [
    eon_role.viewer.id,
  ]
}

# Example: IDP group with multiple roles
resource "eon_idp_group" "operators" {
  idp_id            = "your-idp-id"
  provider_group_id = "your-idp-group-id-for-operators"

  role_ids = [
    eon_role.viewer.id,
    eon_role.operator.id,
  ]
}

# Optional: Use the eon_roles data source to assign roles by name instead of creating them
# data "eon_roles" "all" {}
#
# resource "eon_idp_group" "admins" {
#   idp_id            = "your-idp-id"
#   provider_group_id = "your-idp-group-id-for-admins"
#
#   role_ids = [for r in data.eon_roles.all.roles : r.id if r.name == "Admin"]
# }
