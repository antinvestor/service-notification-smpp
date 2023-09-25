package service

import (
	"encoding/json"
	"fmt"
	notificationV1 "github.com/antinvestor/service-notification-api"
	partitionV1 "github.com/antinvestor/service-partition-api"
	profileV1 "github.com/antinvestor/service-profile-api"
	"github.com/antinvestor/template-service/config"
	"github.com/antinvestor/template-service/service/handlers"
	"github.com/gorilla/mux"
	"github.com/pitabwire/frame"
	"net/http"
)

type templateServer struct {
	Service         *frame.Service
	Config          *config.TemplateConfig
	ProfileCli      *profileV1.ProfileClient
	PartitionCli    *partitionV1.PartitionClient
	NotificationCli *notificationV1.NotificationClient
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *templateServer) writeError(w http.ResponseWriter, err error, code int, msg string) {

	w.Header().Set("Content-Type", "application/json")

	h.Service.L().
		WithField("code", code).
		WithField("message", msg).WithError(err).Error("internal service error")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(&ErrorResponse{
		Code:    code,
		Message: fmt.Sprintf(" internal processing err message: %s %s", msg, err),
	})
	if err != nil {
		h.Service.L().WithError(err).Error("could not write error to response")
	}
}

func (h *templateServer) addHandler(router *mux.Router,
	f func(w http.ResponseWriter, r *http.Request) error, path string, name string, method string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(frame.ToContext(r.Context(), h.Service))
		r = r.WithContext(profileV1.ToContext(r.Context(), h.ProfileCli))
		r = r.WithContext(partitionV1.ToContext(r.Context(), h.PartitionCli))
		r = r.WithContext(notificationV1.ToContext(r.Context(), h.NotificationCli))

		err := f(w, r)
		if err != nil {
			h.writeError(w, err, 500, "could not process request")
		}
	})

	router.Path(path).
		Name(name).
		Handler(handler).
		Methods(method)
}

// NewAuthRouterV1 NewRouterV1 -
func NewAuthRouterV1(service *frame.Service,
	templateConfig *config.TemplateConfig,
	profileCli *profileV1.ProfileClient,
	partitionCli *partitionV1.PartitionClient,
	notificationCli *notificationV1.NotificationClient) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	holder := &templateServer{
		Service:         service,
		Config:          templateConfig,
		ProfileCli:      profileCli,
		PartitionCli:    partitionCli,
		NotificationCli: notificationCli,
	}

	holder.addHandler(router, handlers.IndexEndpoint, "/", "IndexEndpoint", "GET")

	return router
}
