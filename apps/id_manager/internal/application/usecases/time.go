package usecases

import (
	"context"
	"eaglechat/common/clock"
	"time"
)

type GetTimeUseCase struct {
	clock clock.Clock
}

func NewGetTimeUseCase(clock clock.Clock) *GetTimeUseCase {
	return &GetTimeUseCase{
		clock: clock,
	}
}

type GetTimeOutput struct {
	CurrentTime time.Time
}

func NewGetTimeOutput(t time.Time) GetTimeOutput {
	return GetTimeOutput{
		CurrentTime: t,
	}
}

func (uc *GetTimeUseCase) Execute(ctx context.Context, req struct{}) GetTimeOutput {
	return NewGetTimeOutput(uc.clock.Now())
}
