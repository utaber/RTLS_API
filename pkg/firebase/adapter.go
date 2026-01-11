package firebase

import (
	"context"

	"firebase.google.com/go/v4/db"
)

type Adapter struct {
	client *db.Client
}

func NewAdapter(client *db.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) Get(ctx context.Context, path string, dest interface{}) error {
	return a.client.NewRef(path).Get(ctx, dest)
}
