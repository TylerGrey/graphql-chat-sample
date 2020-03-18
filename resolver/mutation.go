package resolver

import (
	"encoding/json"
	"time"
)

func (r *resolver) SendMessage(args struct{ Msg string }) *messageEventResolver {
	e := messageEvent{
		ID:        randomID(),
		Msg:       args.Msg,
		CreatedAt: time.Now().String(),
	}

	go func() {
		payload, _ := json.Marshal(e)
		r.client.Publish("mychannel1", string(payload))
	}()

	return &messageEventResolver{e: e}
}
