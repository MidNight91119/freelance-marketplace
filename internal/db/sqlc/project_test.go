package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func createRandomProject(t *testing.T) Project {
	user := createRandomUser(t)

	budgetMin := util.RandomInt(1000, 50000)
	arg := CreateProjectParams{
		ClientID:    user.ID,
		Title:       util.RandomString(10),
		Description: util.RandomString(10),
		Category:    util.RandomString(10),
		BudgetMin:   budgetMin,
		BudgetMax:   budgetMin + util.RandomInt(1000, 50000),
		Deadline:    time.Now().AddDate(0, 0, 7),
	}

	project, err := testStore.CreateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	return project
}

func TestCreateProject(t *testing.T) {
	createRandomProject(t)
}

func TestCreateProjectConstraints(t *testing.T) {
	user := createRandomUser(t)
	future := time.Now().AddDate(0, 0, 1)
	past := time.Now().AddDate(0, 0, -1)

	tests := []struct {
		name       string
		budgetMin  int64
		budgetMax  int64
		deadline   time.Time
		constraint string
	}{
		{
			name:       "max below min",
			budgetMin:  5000,
			budgetMax:  1000,
			deadline:   future,
			constraint: "budget_range",
		},
		{
			name:       "min is 0",
			budgetMin:  0,
			budgetMax:  1000,
			deadline:   future,
			constraint: "budget_positive",
		},
		{
			name:       "deadline is in past",
			budgetMin:  1000,
			budgetMax:  5000,
			deadline:   past,
			constraint: "deadline_after_creation",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arg := CreateProjectParams{
				ClientID:    user.ID,
				Title:       util.RandomString(10),
				Description: util.RandomString(20),
				Category:    util.RandomString(8),
				BudgetMin:   tc.budgetMin,
				BudgetMax:   tc.budgetMax,
				Deadline:    tc.deadline,
			}

			_, err := testStore.CreateProject(context.Background(), arg)
			require.Error(t, err)

			pgErr, ok := errors.AsType[*pgconn.PgError](err)
			require.True(t, ok)
			require.Equal(t, CheckViolation, pgErr.Code)
			require.Equal(t, tc.constraint, pgErr.ConstraintName)
		})
	}
}
