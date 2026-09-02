package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type createProposalRequest struct {
	CoverLetter           string `json:"coverLetter" validate:"required"`
	ProposedPrice         int64  `json:"proposedPrice" validate:"required,gt=0"`
	EstimatedDurationDays int64  `json:"estimatedDurationDays" validate:"required,gt=0"`
}

func (server *Server) createProposal(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(r.PathValue("projectId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid project id")
		return
	}

	var req createProposalRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "request body is not valid JSON")
		return
	}
	if err = server.validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	payload, ok := payloadFrom(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired access token")
		return
	}

	project, err := server.store.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "PROJECT_NOT_FOUND", "does not exist")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}
	if project.Status != db.ProjectStatusOpen {
		writeError(w, http.StatusConflict, "PROJECT_NOT_OPEN", "project is not open")
		return
	}

	arg := db.CreateProposalParams{
		ProjectID:             projectID,
		FreelancerID:          payload.UserID,
		CoverLetter:           req.CoverLetter,
		ProposedPrice:         req.ProposedPrice,
		EstimatedDurationDays: req.EstimatedDurationDays,
	}

	proposal, err := server.store.CreateProposal(r.Context(), arg)
	if err != nil {
		pgErr, ok := errors.AsType[*pgconn.PgError](err)
		if ok && pgErr.Code == db.UniqueViolation {
			writeError(w, http.StatusConflict, "PROPOSAL_ALREADY_EXISTS", "Freelancer already submitted a proposal")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	writeJSON(w, http.StatusCreated, proposal)
}
