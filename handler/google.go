package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func newGoogleOAuthConfig() *oauth2.Config {
  clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
  clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

	return &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(c *fiber.Ctx) error {
  conf := newGoogleOAuthConfig()

  url := conf.AuthCodeURL("state")
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
  conf := newGoogleOAuthConfig()
	token, err := conf.Exchange(c.Context(), code)
	if err != nil {
		return err
	}

  userInfo, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
  if err != nil {
    return err
  }

  userData, err := ioutil.ReadAll(userInfo.Body)
  if err != nil {
    return err
  }

	log.Printf("User Info: %#v\n", string(userData))
	return c.Status(http.StatusOK).JSON(userData)
}
