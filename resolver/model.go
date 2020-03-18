package resolver

import "github.com/graph-gophers/graphql-go"

type messageEvent struct {
	ID  string
	Msg string
	CreatedAt string
}

type messageEventResolver struct {
	e messageEvent
}

func (r *messageEventResolver) ID() graphql.ID {
	return graphql.ID(r.e.ID)
}

func (r *messageEventResolver) Msg() string {
	return r.e.Msg
}

func (r *messageEventResolver) CreatedAt() string {
	return r.e.CreatedAt
}
