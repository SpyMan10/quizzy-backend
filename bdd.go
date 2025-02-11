package main

import (
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

func ConnectToFirebase() {
	// Replace with your actual path to the service account key file
	opt := option.WithCredentialsFile("firebase.conf.json")

	log.Println("Connecting to Firebase...")
	// Initialize the Firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Get a Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
	}
	// Add a document to a collection
	_, err = client.Collection("users").Doc("alice").Set(context.Background(), map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
	})
	if err != nil {
		log.Fatalf("error adding document: %v", err)
	}

	// Read a document
	doc, err := client.Collection("users").Doc("alice").Get(context.Background())
	if err != nil {
		log.Fatalf("error getting document: %v", err)
	}
	log.Println("Document data:", doc.Data())

	defer client.Close()
}
