package db

import (
	"context"
	"testing"

	"github.com/MidNight91119/freelance-marketplace/internal/util"
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

func TestListProposalsByFreelancer(t *testing.T) {
	project := createRandomProject(t)
	mine := createRandomProposal(t, project.ID)
	other := createRandomProposal(t, project.ID)

	proposals, err := testStore.ListProposalsByFreelancer(context.Background(), mine.FreelancerID)
	require.NoError(t, err)
	require.Len(t, proposals, 1)
	require.Equal(t, mine.FreelancerID, proposals[0].FreelancerID)
	require.NotEqual(t, other.FreelancerID, proposals[0].FreelancerID)
	require.Equal(t, project.Title, proposals[0].ProjectTitle)
}
