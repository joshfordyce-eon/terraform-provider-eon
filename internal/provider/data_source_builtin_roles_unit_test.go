package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinRolesDataSource_Unit(t *testing.T) {
	t.Parallel()

	ds := NewBuiltinRolesDataSource()
	assert.NotNil(t, ds, "Data source should not be nil")
}

// Ensure builtinRoleIDs is unchanged (stable keys and UUIDs).
func TestBuiltinRoleIDs_StableMapping(t *testing.T) {
	t.Parallel()

	expected := map[string]string{
		"global_admin":  "379a1104-838a-4bf3-af96-da3af27c5712",
		"global_viewer": "543bad56-e9b2-421f-8456-b43c53fcebfe",
		"viewer":        "d6afa067-d3a0-457e-923d-27cd26c9e5cb",
		"admin":         "a675e456-8602-4550-9c65-66583404e0d6",
		"operator":      "21d0ae2b-9bbc-4a41-bd5e-98011e9f10a5",
	}
	assert.Equal(t, expected, builtinRoleIDs)
}
