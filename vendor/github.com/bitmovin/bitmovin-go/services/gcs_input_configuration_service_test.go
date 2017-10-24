package services

import (
	"encoding/json"
	"testing"

	"github.com/bitmovin/bitmovin-go/bitmovin"
	"github.com/bitmovin/bitmovin-go/models"
)

const apiKey = "INSERT_API_KEY"

func createClient() *RestService {
	bitmovin := bitmovin.NewBitmovinDefaultTimeout(apiKey, "https://api.bitmovin.com/v1/")
	return NewRestService(bitmovin)
}

func TestCreate(t *testing.T) {
	svc := createClient()
	gcsInput := &models.GCSInput{
		AccessKey:  stringToPtr(""),
		SecretKey:  stringToPtr(""),
		BucketName: stringToPtr(""),
	}
	json, _ := json.Marshal(*gcsInput)
	_, err := svc.Create(`encoding/inputs/gcs`, json)

	if err == nil {
		t.Fatal("Expected to receive error")
	}
	if err.Error() != "ERROR 1000: One or more fields are not present or invalid" {
		t.Fatalf("Expected error message - got %s", err.Error())
	}
}

func TestRetrieve(t *testing.T) {
	svc := createClient()
	_, err := svc.Retrieve(`encoding/inputs/gcs/invalid-id`)
	if err == nil {
		t.Fatal("Expected to receive error - got nil")
	}
	if err.Error() != "ERROR 1001: Input with the given id was not found in our system" {
		t.Fatalf("Expected error message - got %s", err.Error())
	}
}

func TestDelete(t *testing.T) {
	svc := createClient()
	_, err := svc.Delete(`encoding/inputs/gcs/invalid-id`)
	if err == nil {
		t.Fatal("Expected to receive error - got nil")
	}
	if err.Error() != "ERROR 1001: Input with the given id was not found in our system" {
		t.Fatalf("Expected error message - got %s", err.Error())
	}
}

func stringToPtr(s string) *string {
	return &s
}
