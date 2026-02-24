package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
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
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Unique().Field("user_id"),
	}
}

func (Task) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
