package server

import (
	"authentication/service"
	"encoding/json"
	"net/http"
)

type Handler struct {
	authService service.Service
}

func NewHandler(authService service.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	session, code, err := h.authService.CreateAuthSession(guid)
	if err != nil {
		w.WriteHeader(code)
		return
	}

	w.WriteHeader(code)

	response, err := json.Marshal(&TokenResponse{
		AccessToken:  session.AccessToken,
		TokenType:    "Bearer",
		ExpiresAt:    session.ExpiresAt,
		RefreshToken: session.RefreshToken,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	guid, code, err := h.authService.DeleteAuthSession(request.AccessToken, request.RefreshToken)

	if err != nil {
		w.WriteHeader(code)
		return
	}

	session, code, err := h.authService.CreateAuthSession(guid)
	if err != nil {
		w.WriteHeader(code)
		return
	}

	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(&TokenResponse{
		AccessToken:  session.AccessToken,
		TokenType:    "Bearer",
		ExpiresAt:    session.ExpiresAt,
		RefreshToken: session.RefreshToken,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(response)
}
