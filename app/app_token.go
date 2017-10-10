package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"gopkg.in/guregu/null.v3"

	"github.com/danfixeads/livepush/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

type appToken struct {
	Aud string // Audience
	Exp int    // Expiration Time
	Iat int    // Issued at
	Iss string // Issuer
	Nbf int    // Not before
}

func (a *App) returnClientID(r *http.Request) string {
	if rv := context.Get(r, "decoded"); rv != nil {
		var token appToken
		mapstructure.Decode(rv, &token)
		return token.Iss
	}
	return ""
}

func (a *App) returnClientIDNullString(r *http.Request) null.String {

	clientID := a.returnClientID(r)

	return null.String{NullString: sql.NullString{
		String: clientID,
		Valid:  len(clientID) > 0,
	}}
}

// Middleware to protect private pages
func (a *App) tokenValidate(protectedPage http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {

			config := models.ReturnConfig()

			verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(config.TokenKey))

			t, err := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return verifyKey, nil
			})

			fmt.Printf("%v - %v - %v", t, err, t.Claims)

			/*
				var token appToken
				mapstructure.Decode(t.Claims, &token)
				fmt.Printf("TOKEN - %v", token)
			*/

			if err != nil {
				a.respondWithError(w, req, http.StatusBadRequest, err.Error())
				return
			}
			if t.Valid {
				context.Set(req, "decoded", t.Claims)
				protectedPage(w, req)
			} else {
				a.respondWithError(w, req, http.StatusForbidden, "Invalid authorization token")
			}
		} else {
			a.respondWithError(w, req, http.StatusForbidden, "An authorization header is required")
		}
	})
}
