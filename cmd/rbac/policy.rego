package rbac.authz

# role_permissions maps a role to the set of scopes (permissions) it may use.
# move to a data document/DB table if roles need to be editable at runtime.
role_permissions := {
	"admin": {"clients:read", "clients:write", "clients:delete"},
	"editor": {"clients:read", "clients:write"},
	"viewer": {"clients:read"},
}

default allow := false

allow if input.scope in role_permissions[input.role]
