package listener

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
)

type CallbackFn func(ctx context.Context, key, contentType string)

type Listener struct {
	client *minio.Client
	bucket string
}

func NewListener(bucket string, client *minio.Client) *Listener {
	return &Listener{
		client: client,
		bucket: bucket,
	}
}

func (l Listener) Listen(ctx context.Context, fn CallbackFn, events ...string) {
	log.Println(fmt.Sprintf("listening for %s in %s", events, l.bucket))
	for notificationInfo := range l.client.ListenBucketNotification(ctx, l.bucket, "", "", events) {
		if notificationInfo.Err != nil {
			log.Println(notificationInfo.Err)
		}

		for _, record := range notificationInfo.Records {
			key := record.S3.Object.Key
			log.Println(fmt.Sprintf("got %s in %s for %s", record.EventName, l.bucket, key))
			fn(context.Background(), key, record.S3.Object.ContentType)
		}
	}
}
