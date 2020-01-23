package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"

	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

var (
	apiKey         string
	apiSecret      string
	apiAccessToken string

	kc *kiteconnect.Client
)

func init() {
	var ok bool

	if apiKey, ok = os.LookupEnv("KITE_API_KEY"); !ok {
		log.Fatal("KITE_API_KEY not found")
	}
	if apiSecret, ok = os.LookupEnv("KITE_SECRET_KEY"); !ok {
		log.Fatal("KITE_SECRET_KEY not found")
	}
	setupKiteConnectClient()
	if apiAccessToken, ok = os.LookupEnv("KITE_ACCESS_TOKEN"); !ok {
		getAndSetKiteAccessToken()
	}

	kc.SetAccessToken(apiAccessToken)
}

func setupKiteConnectClient() {
	kc = kiteconnect.New(apiKey)

}

func getAndSetKiteAccessToken() {
	var (
		requestToken string
	)

	// Login URL from which request token can be obtained
	fmt.Println("Open the following url in your browser:\n", kc.GetLoginURL())

	// Obtain request token after Kite Connect login flow
	// Run a temporary server to listen for callback
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/api/user/callback/kite/", func(w http.ResponseWriter, r *http.Request) {
		requestToken = r.URL.Query()["request_token"][0]
		log.Println("request token", requestToken)
		go srv.Shutdown(context.TODO())
		w.Write([]byte("login successful!"))
		return
	})
	srv.ListenAndServe()

	// Get user details and access token
	data, err := kc.GenerateSession(requestToken, apiSecret)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	apiAccessToken = data.AccessToken
	fmt.Println("please run the following to avoid creating new session next time:")
	fmt.Printf("export KITE_ACCESS_TOKEN=%s\n", apiAccessToken)
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			userCmd,
			portfolioCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
