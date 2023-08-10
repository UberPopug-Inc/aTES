package main

import (
	"context"
	"crypto/tls"

	"github.com/Nerzal/gocloak/v13"
)

func main() {
	client := gocloak.NewClient("http://localhost:8080/auth")
	restyClient := client.RestyClient()
	restyClient.SetDebug(true)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		panic("Something wrong with the credentials or url")
	}

	gr := []string{"Managers"}
	user := gocloak.User{
		Username:  gocloak.StringP("CoolGuy manager"),
		Enabled:   gocloak.BoolP(true),
		FirstName: gocloak.StringP("Bob"),
		LastName:  gocloak.StringP("Uncle"),
		Email:     gocloak.StringP("manager@really.wrong"),
		Groups:    &gr,
	}

	_, err = client.CreateUser(ctx, token.AccessToken, "master", user)
	if err != nil {
		panic("Oh no!, failed to create user :(")
	}
}
