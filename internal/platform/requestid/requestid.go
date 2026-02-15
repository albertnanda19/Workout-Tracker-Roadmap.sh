package requestid

import "context"

type contextKey string

const key contextKey = "request_id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

func Get(ctx context.Context) (string, bool) {
	v := ctx.Value(key)
	s, ok := v.(string)
	return s, ok && s != ""
}
