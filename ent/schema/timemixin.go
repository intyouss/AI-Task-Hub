package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
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
