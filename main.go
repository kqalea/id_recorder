package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func getExternal() string {
	resp, err := http.Get("http://myexternalip.com/raw")

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
		return "0.0.0.0"
	}
	return string(body)
}

func main() {
	var ip = getExternal()
	currentTime := time.Now()
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		hostname = "error"
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile("<Your access.json>")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app %v\n", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	_, _, err = client.Collection("ip_record").Add(ctx, map[string]interface{}{
		"HostName": hostname,
		"DateTime": currentTime.Format("2006-01-02 15:04:05"),
		"IPv4":     ip,
	})

	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("error closing app... %v", err)
		}
	}(client)

}
