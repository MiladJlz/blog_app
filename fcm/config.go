package fcm

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

const filePath = "./service_account_key.json"

type FirebaseMessagingClient struct {
	client *messaging.Client
}

func NewFirebaseMessagingClient(ctx context.Context) (*FirebaseMessagingClient, error) {
	opt := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &FirebaseMessagingClient{client: fcmClient}, nil
}

func (c *FirebaseMessagingClient) SendNotification(ctx context.Context, tokens []string, message string) error {
	_, err := c.client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Data: map[string]string{
			message: message,
		},
		Tokens: tokens,
	})
	return err
}
