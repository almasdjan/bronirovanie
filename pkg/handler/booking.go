package handler

import (
	"bronirovanie/models"
	"strconv"

	"encoding/json"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// @Summary make reservation
// @Tags reservations
// @Param input body models.CreateReservation true "reservation"
// @Success 200 {integer} string
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /reservations [post]
func (h *Handler) CreateReservationHandler(w http.ResponseWriter, r *http.Request) {

	var input models.CreateReservation

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.Booking.Create(r.Context(), &input)
	if err != nil {
		logrus.Print(err.Error())
		if err.Error() == "Conflict" {
			respondWithError(w, 409, "Conflict")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithStatus(w, http.StatusCreated, "Created")

}

// @Summary Get all reservations for a room
// @Tags reservations
// @Param room_id path int true "Room ID"
// @Success 200 {array} models.Reservation
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /reservations/{room_id} [get]
func (h *Handler) GetReservationHandler(w http.ResponseWriter, r *http.Request) {

	roomID, err := strconv.Atoi(chi.URLParam(r, "room_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	reservations, err := h.services.Booking.GetAll(r.Context(), roomID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(reservations) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	formattedJSON, err := json.MarshalIndent(reservations, "", "  ")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to format JSON response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)

}
