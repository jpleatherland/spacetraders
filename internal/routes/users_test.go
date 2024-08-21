package routes

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

func TestUserLogin(t *testing.T) {
	err := createTestUser(&resources, "testLoginUser1", "1234")
	if err != nil {
		t.Fatalf("unable to create test user: %v", err)
	}
	cookie, err := testUserLogin(&resources, "testLoginUser1", "1234")
	if err != nil {
		t.Fatalf("unable to login user: %v", err)
	}
	if cookie.Value == "" {
		t.Errorf("no token in cookie value")
	}
}

func TestGetSessionByToken(t *testing.T) {
	err := createTestUser(&resources, "testSessionUser1", "1234")
	if err != nil {
		t.Fatalf("unable to create test user: %v", err)
	}
	cookie, err := testUserLogin(&resources, "testSessionUser1", "1234")
	if err != nil {
		t.Fatalf("unable to login user: %v", err)
	}
	if cookie.Value == "" {
		t.Fatal("no token in cookie value")
	}

	session, err := testGetSessionById(&resources, cookie)
	if err != nil {
		t.Fatalf("unable to get session: %v", err.Error())
	}

	if session.ID != cookie.Value {
		t.Errorf("session ID does not match token")
	}

	log.Printf("session id: %v\nuser id:%v\nagent id:%v\nexpires at: %v\n",
		session.ID, session.UserID, session.AgentID, session.ExpiresAt)
}
