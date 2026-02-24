package schema

import "entgo.io/ent"

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return nil
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return nil
}

func (Task) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
