package main

import (
	"log"
	"testing"
)

func TestCreateUser(t *testing.T) {
	err := createTestUser(&resources, "testCreateUser1", "1234")
	if err != nil {
		t.Fatalf("unable to create test user: %v", err)
	}
	log.Printf("created user")
}
