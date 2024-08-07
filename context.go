package daisysms

import "context"

type ContextKey string

const DAISYSMS_KEY ContextKey = "daisysms"

func WithContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, DAISYSMS_KEY, client)
}

func FromContext(ctx context.Context) *Client {
	client, ok := ctx.Value(DAISYSMS_KEY).(*Client)
	if !ok {
		return nil
	}
	return client
}
