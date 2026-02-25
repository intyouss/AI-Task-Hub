package schema

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	gen "github.com/intyouss/AI-Task-Hub/ent"
	"github.com/intyouss/AI-Task-Hub/ent/hook"
	"github.com/intyouss/AI-Task-Hub/ent/intercept"
)

type SoftDeleteMixin struct {
	ent.Schema
}

func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").Optional().Nillable().Comment("删除时间"),
	}
}

func (d SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(mutator ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
					if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
						return mutator.Mutate(ctx, mutation)
					}

					mx, ok := mutation.(interface {
						SetOp(ent.Op)
						SetDeletedAt(time.Time)
						Client() *gen.Client
						WhereP(...func(*sql.Selector))
					})
					if !ok {
						return nil, errors.New("unexpected mutation type for soft delete")
					}
					d.P(mx)
					mx.SetOp(ent.OpUpdate)
					mx.SetDeletedAt(time.Now().UTC())
					return mx.Client().Mutate(ctx, mutation)
				})
			}, ent.OpDelete|ent.OpDeleteOne),
	}
}

func (d SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, query intercept.Query) error {
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			d.P(query)
			return nil
		}),
	}
}

// P 添加软删除查询条件
func (d SoftDeleteMixin) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(
		sql.FieldIsNull(d.Fields()[0].Descriptor().Name),
	)
}

type softDeleteKey struct{}

// SkipSoftDelete 跳过软删除
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}
