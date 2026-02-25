package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/intyouss/AI-Task-Hub/ent/schema/hooks"
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
		hooks.TimeMixin{},
	}
}
