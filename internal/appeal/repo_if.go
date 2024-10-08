// Code generated by ifacemaker; DO NOT EDIT.

package appeal

import (
	"context"
)

// Repository ...
type Repository interface {
	SaveAppeal(ctx context.Context, p *Appeal) (*NewAppeal, error)
	AnswerAppeal(ctx context.Context, id int64) error
	GetAppeal(ctx context.Context, id int64) (*Appeal, error)
	AdminExists(ctx context.Context, division string) (bool, error)
}
