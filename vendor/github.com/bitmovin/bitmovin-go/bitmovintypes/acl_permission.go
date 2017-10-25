package bitmovintypes

type ACLPermission string

const (
	ACLPermissionPublicRead ACLPermission = "PUBLIC_READ"
	ACLPermissionPrivate    ACLPermission = "PRIVATE"
)
