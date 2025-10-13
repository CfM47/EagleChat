package usecases

import "context"

type UseCase[T, V any] interface {
    Execute(ctx context.Context, input T) (V, error)
}
