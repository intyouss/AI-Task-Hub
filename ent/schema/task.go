package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	ent2 "github.com/intyouss/AI-Task-Hub/ent"
	"github.com/intyouss/AI-Task-Hub/ent/hook"
	"github.com/intyouss/AI-Task-Hub/ent/task"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("model_name"),
		field.Text("prompt"),
		field.Text("output"),
		field.Enum("status").Values("pending", "processing", "completed", "failed"),
		field.UUID("user_id", uuid.UUID{}),
		field.Time("finished_at").Optional(),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Unique().Required().Ref("tasks").Field("user_id"),
	}
}

func (Task) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Task) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
					update, ok := m.(*ent2.TaskMutation)
					if !ok {
						return next.Mutate(ctx, m)
					}
					status, ok := update.Status()
					if !ok {
						return next.Mutate(ctx, m)
					}
					if status == task.StatusCompleted {
						err := m.SetField(task.FieldFinishedAt, time.Now().UTC())
						if err != nil {
							return nil, err
						}
					}
					return next.Mutate(ctx, m)
				})
			},
			ent.OpUpdateOne|ent.OpUpdate,
		),
	}
}
