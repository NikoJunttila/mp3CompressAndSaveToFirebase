package main

import (
	"fmt"
	"log"
  "strings"
	"os"
	"io"
	"context"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	_"cloud.google.com/go/firestore"
)

func main() {
  cleanAndCompressFiles()
  genre := "Rock"
	// Replace this with the path to your Firebase service account JSON file
	serviceAccountKeyFilePath := "./service-account.json"

	// Initialize the Firebase app
	ctx := context.Background()
	config := &firebase.Config{
		StorageBucket: "firesvelte-1.appspot.com", // Replace with your storage bucket name
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Create a storage client from the Firebase app
	client, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("error getting Storage client: %v\n", err)
	}

	// Get a reference to the default bucket
	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalf("error getting bucket handle: %v\n", err)
	}
 	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Read all files in the current directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}
    if !strings.HasSuffix(file.Name(), ".mp3"){
      continue
    }
	filenameWithoutExtension := strings.TrimSuffix(file.Name(), ".mp3")
  localFilePath := dir + "/" + file.Name()	
  
  fmt.Println("Current file "+  localFilePath)

  f, err := os.Open(localFilePath)
	if err != nil {
		log.Fatalf("error opening local file: %v\n", err)
	}
	defer f.Close()

	// Create a writer to the file location in Firebase Storage
	objectName := "music/" + genre + "/" + file.Name()
	wc := bucket.Object(objectName).NewWriter(ctx)

	// Copy the local file to the Firebase Storage bucket
	if _, err = io.Copy(wc, f); err != nil {
		log.Fatalf("error copying file to Firebase Storage: %v\n", err)
	}

	// Close the writer
	if err = wc.Close(); err != nil {
		log.Fatalf("error closing writer: %v\n", err)
	}
	// Make the file publicly accessible
	acl := bucket.Object(objectName).ACL()
	if err = acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Fatalf("error making file publicly accessible: %v\n", err)
	}

	// Generate the public URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.StorageBucket, objectName)
	fmt.Printf("File uploaded successfully. Public URL: %s\n", url)
	// Initialize Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v\n", err)
	}
	defer firestoreClient.Close()
	collectionName := "music"
  artist,name := SeparateArtistAndTitle(file.Name())
	doc := firestoreClient.Collection(collectionName).Doc(filenameWithoutExtension)
	_, err = doc.Set(ctx, map[string]interface{}{
		"name":       name,
    "artist":     artist,
    "genre":      genre,
		"url":        url,
	})
	if err != nil {
		log.Fatalf("error writing to Firestore: %v\n", err)
	}

	fmt.Printf("File details saved to Firestore in collection '%s'.\n", collectionName)

  }
	}
