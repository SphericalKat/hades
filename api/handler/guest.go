package handler

import (
	"github.com/ATechnoHazard/hades-2/pkg/guest"
	"net/http"
)

func CreateGuest(guestService guest.Service) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

	}
}