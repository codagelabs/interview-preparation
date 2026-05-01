package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	// You need to create an OAuth app on GitHub to get these:
	// Go to: Settings -> Developer settings -> OAuth Apps -> New OAuth App
	// Let Authorization callback URL be: http://localhost:8080/callback
	oauthConfig = &oauth2.Config{
		ClientID:     "YOUR_GITHUB_CLIENT_ID",
		ClientSecret: "YOUR_GITHUB_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// A random string used to protect against CSRF attacks.
	// In production, this should be a cryptographically secure random string stored in a session.
	oauthStateString = "randomized-state-string"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	fmt.Println("Server successfully started on http://localhost:8080")
	fmt.Println("Make sure you replace CLIENT_ID and CLIENT_SECRET in main.go!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	htmlIndex := `<html><body>
		<h1>Simple Go OAuth 2.0 Example (GitHub)</h1>
		<a href="/login">Log in with GitHub</a>
	</body></html>`
	fmt.Fprintf(w, "%s", htmlIndex)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// 1. Generate the URL to redirect the user to GitHub's Authorization Server
	// We pass the State string to prevent CSRF.
	url := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)

	// 2. Redirect user to GitHub
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// 3. GitHub redirects back here. First verify the state string.
	if r.FormValue("state") != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// 4. Get the authorization code from the query string
	code := r.FormValue("code")

	// 5. Exchange the Authorization Code for an Access Token
	// This happens on the backend out-of-band (Server to Server communication)
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConfig.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// 6. Use the access token to get user data from the Resource Server (GitHub API)
	client := oauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://api.github.com/user")
	if err != nil {
		fmt.Printf("client.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	// 7. Read and print the user data back to the browser
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("io.ReadAll() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Display the JSON response from GitHub API
	w.Header().Set("Content-Type", "application/json")
	w.Write(contents)
}
