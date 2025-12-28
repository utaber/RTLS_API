package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

func NewDatabase(ctx context.Context, credFile, dbURL string) *db.Client {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL: dbURL,
	}, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
