package rbac

const (
	SsoServiceType = "sso-service"

	ReadUsers        = "read:users"
	WriteUsers       = "write:users"
	ReadTokens       = "read:tokens"
	WriteTokens      = "write:tokens"
	ReadPermissions  = "read:permissions"
	WritePermissions = "write:permissions"
	ReadRoles        = "read:roles"
	WriteRoles       = "write:roles"
	ReadScopes       = "read:scopes"
	WriteScopes      = "write:scopes"
)

func HasPermission(claimPermissions []string, requiredPermission string) bool {
	return has(claimPermissions, requiredPermission)
}

func HasScope(claimScope []string) bool {
	return has(claimScope, SsoServiceType)
}

func has(collection []string, item string) bool {
	for _, r := range collection {
		if r == item {
			return true
		}
	}

	return false
}
