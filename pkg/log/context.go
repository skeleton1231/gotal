package log

import "context"

type key int

const (
	contextKey key = iota
)

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, l)
}
