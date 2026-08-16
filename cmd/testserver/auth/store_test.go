package auth_test

import (
	"testing"
	"time"

	"github.com/Hari-Kiri/virest/cmd/testserver/auth"
)

func TestAuthenticateAndAllowlist(t *testing.T) {
	cfg := auth.Config{Users: []auth.User{{
		Username:            "john",
		Password:            "password",
		Role:                "operator",
		HypervisorAllowlist: []string{"qemu:///system"},
	}}}
	store, err := auth.NewStore(cfg, "app", []byte("secret-key-for-tests"), time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	token, err := store.AuthenticateBasic("john", "password")
	if err != nil {
		t.Fatalf("auth: %v", err)
	}
	claims, err := store.ParseBearer("Bearer " + token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.Subject != "john" || claims.Role != "operator" {
		t.Fatalf("claims=%+v", claims)
	}
	if !store.AllowHypervisor("john", "qemu:///system") {
		t.Fatal("expected allow")
	}
	if store.AllowHypervisor("john", "qemu:///session") {
		t.Fatal("expected deny")
	}
	if _, err := store.AuthenticateBasic("john", "wrong"); err == nil {
		t.Fatal("expected bad password")
	}
}
