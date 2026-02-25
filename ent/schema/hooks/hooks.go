package hooks

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	gen "github.com/intyouss/AI-Task-Hub/ent"
	"github.com/intyouss/AI-Task-Hub/ent/hook"
)

type TimeMixin struct {
	ent.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Optional().Nillable().Comment("创建时间"),
		field.Time("updated_at").Optional().Nillable().Comment("更新时间"),
	}
}

func (TimeMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(mutator ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
					now := time.Now().UTC()
					err := mutation.SetField("updated_at", now)
					if err != nil {
						return nil, err
					}

					err = mutation.SetField("created_at", now)
					if err != nil {
						return nil, err
					}
					return mutator.Mutate(ctx, mutation)
				})
			}, ent.OpCreate),
		hook.On(
			func(mutator ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
					err := mutation.SetField("updated_at", time.Now().UTC())
					if err != nil {
						return nil, err
					}
					return mutator.Mutate(ctx, mutation)
				})
			}, ent.OpUpdate|ent.OpUpdateOne),
	}
}

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
		ent.TraverseFunc(func(ctx context.Context, query ent.Query) error {
			return nil
		}),
	}
}

// P adds a storage-level predicate to the queries and mutations.
func (d SoftDeleteMixin) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(
		sql.FieldIsNull(d.Fields()[0].Descriptor().Name),
	)
}
