package db

import (
	"context"
	"errors"
	"testing"

	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func createRandomProposal(t *testing.T, projectID int64) Proposal {
	freelancer := createRandomUserWithRole(t, RolesFreelancer)

	arg := CreateProposalParams{
		ProjectID:             projectID,
		FreelancerID:          freelancer.ID,
		CoverLetter:           util.RandomString(20),
		ProposedPrice:         util.RandomInt(1000, 100000),
		EstimatedDurationDays: util.RandomInt(1, 90),
	}
	proposal, err := testStore.CreateProposal(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, proposal)

	return proposal
}

func TestAcceptProposalTx(t *testing.T) {
	project := createRandomProject(t)

	proposal1 := createRandomProposal(t, project.ID)
	proposal2 := createRandomProposal(t, project.ID)
	proposal3 := createRandomProposal(t, project.ID)

	arg := AcceptProposalTxParams{
		ProposalID: proposal1.ID,
	}
	result, err := testStore.AcceptProposalTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, ProposalStatusAccepted, result.Proposal.Status)
	require.Equal(t, ProjectStatusInProgress, result.Project.Status)
	require.Equal(t, proposal1.ProposedPrice, result.Contract.Amount)

	// re-reading the proposals to get updated values from real database as tx doesn't update values in your local func vairables (ie, proposal2 and 3)
	p2, err := testStore.GetProposal(context.Background(), proposal2.ID)
	require.NoError(t, err)
	require.Equal(t, ProposalStatusRejected, p2.Status)

	p3, err := testStore.GetProposal(context.Background(), proposal3.ID)
	require.NoError(t, err)
	require.Equal(t, ProposalStatusRejected, p3.Status)
}

func TestAcceptProposalTxRollback(t *testing.T) {
	project := createRandomProject(t)

	proposal1 := createRandomProposal(t, project.ID)
	proposal2 := createRandomProposal(t, project.ID)

	arg1 := AcceptProposalTxParams{
		ProposalID: proposal1.ID,
	}
	result, err := testStore.AcceptProposalTx(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, ProposalStatusAccepted, result.Proposal.Status)
	require.Equal(t, ProjectStatusInProgress, result.Project.Status)
	require.Equal(t, proposal1.ProposedPrice, result.Contract.Amount)

	arg2 := AcceptProposalTxParams{
		ProposalID: proposal2.ID,
	}
	_, err = testStore.AcceptProposalTx(context.Background(), arg2)
	require.Error(t, err)

	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	require.True(t, ok)
	require.Equal(t, "one_contract_per_project", pgErr.ConstraintName)

	p2, err := testStore.GetProposal(context.Background(), proposal2.ID)
	require.NoError(t, err)
	require.Equal(t, ProposalStatusRejected, p2.Status)

	p1, err := testStore.GetProposal(context.Background(), proposal1.ID)
	require.NoError(t, err)
	require.Equal(t, ProposalStatusAccepted, p1.Status)
}
