package resolver

import "context"

func (r *resolver) OnMessage(ctx context.Context) <-chan *messageEventResolver {
	c := make(chan *messageEventResolver)
	// NOTE: this could take a while
	r.onMessageSubscriber <- &onMessageSubscriber{events: c, stop: ctx.Done()}

	return c
}
