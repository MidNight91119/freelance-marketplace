package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
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


type listProjectResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	BudgetMin   int64     `json:"budgetMin"`
	BudgetMax   int64     `json:"budgetMax"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	ClientName  string    `json:"clientName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func newlistProjectResponse(projects []db.ListProjectsRow) []listProjectResponse {
	var newProjects []listProjectResponse
	for _, project := range projects {
		newProjects = append(newProjects, listProjectResponse{
			ID: project.ID,
			Title: project.Title,
			Description: project.Description,
			Category: project.Category,
			BudgetMin: project.BudgetMin,
			BudgetMax: project.BudgetMax,
			Deadline: project.Deadline,
			Status: string(project.Status),
			ClientName: project.ClientName,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
		})
	}
	return newProjects
}

func (server *Server) listProjects(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	category := q.Get("category")
	arg := db.ListProjectsParams{
		Category: pgtype.Text{
			String: category,
			Valid:  category != "",
		},
	}
	if s := q.Get("minBudget"); s != "" {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "minBudget must be a number")
			return
		}
		arg.MinBudget = pgtype.Int8{
			Int64: n,
			Valid: true,
		}
	}
	if s := q.Get("maxBudget"); s != "" {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "maxBudget must be a number")
			return
		}
		arg.MaxBudget = pgtype.Int8{
			Int64: n,
			Valid: true,
		}
	}

	projects, err := server.store.ListProjects(r.Context(), arg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	writeJSON(w, http.StatusOK, newlistProjectResponse(projects))
}
