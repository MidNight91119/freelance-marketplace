package db

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Name:           util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: util.RandomString(8),
		Role:           randomRole(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Role, user.Role)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	user1 := createRandomUser(t)
	var pgErr *pgconn.PgError

	arg := CreateUserParams{
		Name:           util.RandomName(),
		Email:          user1.Email,
		HashedPassword: util.RandomString(8),
		Role:           randomRole(),
	}

	_, err := testQueries.CreateUser(context.Background(), arg)
	require.ErrorAs(t, err, &pgErr)
	require.Equal(t, UniqueViolation, pgErr.Code)
}

func randomRole() Roles {
	if rand.IntN(2) == 0 {
		return RolesClient
	}
	return RolesFreelancer
}
