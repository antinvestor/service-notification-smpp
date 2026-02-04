package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"buf.build/gen/go/antinvestor/notification/connectrpc/go/notification/v1/notificationv1connect"
	"buf.build/gen/go/antinvestor/partition/connectrpc/go/partition/v1/partitionv1connect"
	"buf.build/gen/go/antinvestor/profile/connectrpc/go/profile/v1/profilev1connect"
	"github.com/antinvestor/apis/go/notification"
	"github.com/antinvestor/apis/go/partition"
	"github.com/antinvestor/apis/go/profile"
	"github.com/antinvestor/service-notification-smpp/config"
	"github.com/antinvestor/service-notification-smpp/service/handlers"
	"github.com/gorilla/mux"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
)

type templateServer struct {
	Service         *frame.Service
	Config          *config.TemplateConfig
	ProfileCli      profilev1connect.ProfileServiceClient
	PartitionCli    partitionv1connect.PartitionServiceClient
	NotificationCli notificationv1connect.NotificationServiceClient
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *templateServer) writeError(w http.ResponseWriter, err error, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")

	util.Log(context.Background()).
		With("code", code).
		With("message", msg).WithError(err).Error("internal service error")
	w.WriteHeader(code)

	encErr := json.NewEncoder(w).Encode(&ErrorResponse{
		Code:    code,
		Message: fmt.Sprintf(" internal processing err message: %s %s", msg, err),
	})
	if encErr != nil {
		util.Log(context.Background()).WithError(encErr).Error("could not write error to response")
	}
}

func (h *templateServer) addHandler(router *mux.Router,
	f func(w http.ResponseWriter, r *http.Request) error, path string, name string, method string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(frame.ToContext(r.Context(), h.Service))
		r = r.WithContext(profile.ToContext(r.Context(), h.ProfileCli))
		r = r.WithContext(partition.ToContext(r.Context(), h.PartitionCli))
		r = r.WithContext(notification.ToContext(r.Context(), h.NotificationCli))

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

// NewAuthRouterV1 creates the HTTP router with injected dependencies.
func NewAuthRouterV1(service *frame.Service,
	templateConfig *config.TemplateConfig,
	profileCli profilev1connect.ProfileServiceClient,
	partitionCli partitionv1connect.PartitionServiceClient,
	notificationCli notificationv1connect.NotificationServiceClient) *mux.Router {
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
