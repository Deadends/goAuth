package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth/gothic"
)

// RegisterRoutes initializes all API routes
func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// CORS settings (explicit origins)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"}, // Explicitly allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/", s.indexHandler)
	r.Get("/health", s.healthHandler)
	r.Get("/auth/{provider}", s.authHandler)
	r.Get("/auth/callback/{provider}", s.getAuthCallbackFunction) // Google OAuth callback
	r.Get("/logout/{provider}", s.logoutHandler)

	return r
}

// IndexHandler - Renders the index page with login options
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	providerIndex := ProviderIndex{
		Providers:    []string{"google"},
		ProvidersMap: map[string]string{"google": "Google"},
	}
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(w, providerIndex)
}

// HelloWorldHandler - Basic test endpoint
func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// LogoutHandler - Handles user logout
func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// AuthHandler - Handles user authentication
func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	// Try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, gothUser)
	} else {
		// Ensure correct provider is set in context
		req := r.WithContext(context.WithValue(r.Context(), "provider", provider))
		// Start OAuth authentication
		gothic.BeginAuthHandler(w, req)
	}
}

// HealthHandler - API Health Check
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// Google OAuth Callback Handler
func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider)) // Preserve request context

	// Authenticate user via Gothic
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("Error during authentication: %v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Log user details
	log.Printf("User details: %+v", user)
	log.Printf("User email: %s", user.Email)
	log.Printf("User name: %s", user.Name)
	log.Printf("User nickname: %s", user.NickName)
	log.Printf("User location: %s", user.Location)
	log.Printf("User avatar URL: %s", user.AvatarURL)
	log.Printf("User description: %s", user.Description)
	log.Printf("User ID: %s", user.UserID)
	log.Printf("User access token: %s", user.AccessToken)
	log.Printf("User expires at: %s", user.ExpiresAt)
	log.Printf("User refresh token: %s", user.RefreshToken)

	// Redirect user to frontend after successful login
	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
}

// Templates
var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>`

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}
