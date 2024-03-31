package slices

import "context"

func Map[T, U any](in []T, fn func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}

func MapWithCtx[T, U any](
	ctx context.Context,
	in []T,
	fn func(context.Context, T) U,
) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = fn(ctx, v)
	}
	return out
}
