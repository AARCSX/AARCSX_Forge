package contextx

import "context"

type key string

const (
	requestIDKey   key = "request_id"
	traceIDKey     key = "trace_id"
	tenantIDKey    key = "tenant_id"
	actorIDKey     key = "actor_id"
	rolesKey       key = "roles"
	permissionsKey key = "permissions"
)

type IdentityContext struct {
	ActorID     string
	Roles       []string
	Permissions []string
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func WithIdentity(ctx context.Context, identity IdentityContext) context.Context {
	ctx = context.WithValue(ctx, actorIDKey, identity.ActorID)
	ctx = context.WithValue(ctx, rolesKey, identity.Roles)
	ctx = context.WithValue(ctx, permissionsKey, identity.Permissions)
	return ctx
}

func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

func TraceID(ctx context.Context) string {
	v, _ := ctx.Value(traceIDKey).(string)
	return v
}

func TenantID(ctx context.Context) string {
	v, _ := ctx.Value(tenantIDKey).(string)
	return v
}

func Identity(ctx context.Context) IdentityContext {
	actorID, _ := ctx.Value(actorIDKey).(string)
	roles, _ := ctx.Value(rolesKey).([]string)
	permissions, _ := ctx.Value(permissionsKey).([]string)
	return IdentityContext{
		ActorID:     actorID,
		Roles:       append([]string{}, roles...),
		Permissions: append([]string{}, permissions...),
	}
}
