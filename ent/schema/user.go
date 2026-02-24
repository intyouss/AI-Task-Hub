package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/intyouss/AI-Task-Hub/ent/hook"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New),
		field.String("nickname"),
		field.String("api_key").Unique(),
		field.String("phone").Default(""),
		field.String("email").Default(""),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

type TimeMixin struct {
	ent.Schema
}

func (TimeMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(mutator ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
					err := mutation.SetField("updated_at", time.Now())
					if err != nil {
						return nil, err
					}

					err = mutation.SetField("created_at", time.Now())
					if err != nil {
						return nil, err
					}
					return mutator.Mutate(ctx, mutation)
				})
			}, ent.OpCreate),
		hook.On(
			func(mutator ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
					err := mutation.SetField("updated_at", time.Now())
					if err != nil {
						return nil, err
					}
					return mutator.Mutate(ctx, mutation)
				})
			}, ent.OpUpdate|ent.OpUpdateOne),
	}
}
