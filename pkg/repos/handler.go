package repos

import (
	"net/http"
	"scanner/pkg/internal/net"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type handler struct {
	logger *logrus.Entry
	s      Service
}

func NewHandler(l *logrus.Entry, s Service) *handler {
	return &handler{l, s}
}

func (h *handler) AddRoutes(r *mux.Router) {
	r.HandleFunc("/repos/add", h.addRepos).Methods(http.MethodPost)
	r.HandleFunc("/repos/get", h.getRepos).Methods(http.MethodPost)
	r.HandleFunc("/repos/update", h.updateRepos).Methods(http.MethodPost)
	r.HandleFunc("/repos/delete", h.deleteRepos).Methods(http.MethodPost)
	r.HandleFunc("/repos/scan", h.scanRepos).Methods(http.MethodPost)
	r.HandleFunc("/repos/scan/results", h.scanResult).Methods(http.MethodPost)
}

func (h *handler) addRepos(w http.ResponseWriter, r *http.Request) {
	var (
		request  = &ReqAddRepos{}
		response = &RespAddRepos{}
		err      error
	)

	if err = net.Bind(r, request); err != nil {
		net.WriteJSON(w, nil, err)
		return
	}
	response, err = h.s.AddRepos(r.Context(), request)
	net.WriteJSON(w, response, err)
}

func (h *handler) scanRepos(w http.ResponseWriter, r *http.Request) {
	var (
		request  = &ReqScan{}
		response = &RespScan{}
		err      error
	)

	if err = net.Bind(r, request); err != nil {
		net.WriteJSON(w, nil, err)
		return
	}
	response, err = h.s.StartScanRepos(r.Context(), request)
	net.WriteJSON(w, response, err)
}

func (h *handler) updateRepos(w http.ResponseWriter, r *http.Request) {
	var (
		request  = &ReqUpdateRepos{}
		response = &RespUpdateRepos{}
		err      error
	)

	if err = net.Bind(r, request); err != nil {
		net.WriteJSON(w, nil, err)
		return
	}
	response, err = h.s.UpdateRepos(r.Context(), request)
	net.WriteJSON(w, response, err)
}

func (h *handler) deleteRepos(w http.ResponseWriter, r *http.Request) {
	var (
		request  = &ReqDeleteRepos{}
		response = &RespDeleteRepos{}
		err      error
	)

	if err = net.Bind(r, request); err != nil {
		net.WriteJSON(w, nil, err)
		return
	}
	response, err = h.s.ArchiveRepos(r.Context(), request)
	net.WriteJSON(w, response, err)
}

func (h *handler) getRepos(w http.ResponseWriter, r *http.Request) {
	var (
		request  = &ReqGetRepos{}
		response = &RespGetRepos{}
		err      error
	)

	if err = net.Bind(r, request); err != nil {
		net.WriteJSON(w, nil, err)
		return
	}
	response, err = h.s.GetRepos(r.Context(), request)
	net.WriteJSON(w, response, err)
}

func (h *handler) scanResult(w http.ResponseWriter, r *http.Request) {
	var (
		request  = &ReqGetResult{}
		response = &RespGetResult{}
		err      error
	)

	if err = net.Bind(r, request); err != nil {
		net.WriteJSON(w, nil, err)
		return
	}
	response, err = h.s.GetResult(r.Context(), request)
	net.WriteJSON(w, response, err)
}
