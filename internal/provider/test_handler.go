// Copyright (c) HashiCorp, Inc.

package provider

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func newTestHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/token", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("grant_type") == "urn:ietf:params:oauth:grant-type:jwt-bearer" {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"access_token": "validToken"}`))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		logrus.WithField("grant_type", r.FormValue("grant_type")).Error("Invalid grant_type")
		w.WriteHeader(http.StatusUnauthorized)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		// Get device
		r.Get("/v2/projects/{projectID}/devices/{deviceID}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			name := strings.TrimPrefix(r.URL.Path, "/v2/")
			body, err := json.Marshal(dt.Device{
				Name: name,
				Type: "temperature",
				Labels: map[string]string{
					"key": "value",
				},
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = w.Write(body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})

		r.Get("/v2/projects/{projectID}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			name := strings.TrimPrefix(r.URL.Path, "/v2/")
			body, err := json.Marshal(dt.Project{
				Name:                    name,
				DisplayName:             "Test Project",
				Inventory:               true,
				Organization:            "organizations/your-organization-id",
				OrganizationDisplayName: "Test Organization",
				SensorCount:             10,
				CloudConnectorCount:     1,
				Location: dt.Location{
					Latitude:     63.44539,
					Longitude:    10.910202,
					TimeLocation: "Europe/Oslo",
				},
			})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = w.Write(body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	})

	return r
}

// authMiddleware is a middleware that checks for a valid token in the Authorization header.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer validToken" {
			logrus.WithField("Authorization", r.Header.Get("Authorization")).Debug("Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
