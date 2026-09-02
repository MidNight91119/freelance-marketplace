package api

import (
	"encoding/json"
	"net/http"
	"time"

	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
)

type createProjectRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"required"`
	Category    string `json:"category" validate:"required,max=255"`
	BudgetMin   int64  `json:"budgetMin" validate:"required,gt=0"`
	BudgetMax   int64  `json:"budgetMax" validate:"required,gtefield=BudgetMin"`
	Deadline    string `json:"deadline" validate:"required,datetime=2006-01-02"`
}

type projectResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	BudgetMin   int64     `json:"budgetMin"`
	BudgetMax   int64     `json:"budgetMax"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newProjectResponse(project db.Project) projectResponse {
	return projectResponse{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description,
		Category:    project.Category,
		BudgetMin:   project.BudgetMin,
		BudgetMax:   project.BudgetMax,
		Deadline:    project.Deadline,
		Status:      string(project.Status),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

func (server *Server) createProject(w http.ResponseWriter, r *http.Request) {
	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "request body is not valid JSON")
		return
	}

	err := server.validate.Struct(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	payload, ok := payloadFrom(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired access token")
		return
	}

	deadline, err := time.Parse("2006-01-02", req.Deadline)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "deadline must be in YYYY-MM-DD format")
		return
	}
	if !deadline.After(time.Now()) {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "deadline must be in the future")
		return
	}

	arg := db.CreateProjectParams{
		ClientID:    payload.UserID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		BudgetMin:   req.BudgetMin,
		BudgetMax:   req.BudgetMax,
		Deadline:    deadline,
	}

	project, err := server.store.CreateProject(r.Context(), arg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	writeJSON(w, http.StatusCreated, newProjectResponse(project))
}
