package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/dgrijalva/jwt-go"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestAuth(t *testing.T) {
	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single users.", testID)
		{
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a private key: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a private key.", success, testID)

			// The key id we are stating represents the public key in the
			// public key store.
			const keyID = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
			lookup := func(kid string) (*rsa.PublicKey, error) {
				switch kid {
				case keyID:
					return &privateKey.PublicKey, nil
				}
				return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
			}

			a, err := auth.New("RS256", keyID, lookup, auth.Keys{keyID: privateKey})
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create an authenticator: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create an authenticator.", success, testID)

			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				User: auth.User{
					Username: "test_name",
					ID:       "test_id",
				},
			}

			token, err := a.GenerateToken(keyID, claims)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to generate a JWT: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to generate a JWT.", success, testID)

			parsedClaims, err := a.ValidateToken(token)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to parse the claims: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to parse the claims.", success, testID)

			if exp, got := claims.User.Username, parsedClaims.User.Username; exp != got {
				t.Logf("\t\tTest %d:\texp: %s", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %s", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected username: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected username.", success, testID)

			if exp, got := claims.User.ID, parsedClaims.User.ID; exp != got {
				t.Logf("\t\tTest %d:\texp: %s", testID, exp)
				t.Logf("\t\tTest %d:\tgot: %s", testID, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected ID: %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected ID.", success, testID)
		}
	}
}
