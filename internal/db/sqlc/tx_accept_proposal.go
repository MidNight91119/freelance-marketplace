package db

import "context"

type AcceptProposalTxParams struct {
	ProposalID int64
}

type AcceptProposalTxResult struct {
	Proposal Proposal
	Project  Project
	Contract Contract
}

func (store *SQLStore) AcceptProposalTx(ctx context.Context, arg AcceptProposalTxParams) (AcceptProposalTxResult, error) {
	var result AcceptProposalTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Proposal, err = q.AcceptProposal(ctx, arg.ProposalID)
		if err != nil {
			return err
		}

		rejectArg := RejectOtherProposalsParams{
			ProjectID: result.Proposal.ProjectID,
			ID:        result.Proposal.ID,
		}
		err = q.RejectOtherProposals(ctx, rejectArg)
		if err != nil {
			return err
		}

		projectArg := UpdateProjectStatusParams{
			ID:     result.Proposal.ProjectID,
			Status: ProjectStatusInProgress,
		}
		result.Project, err = q.UpdateProjectStatus(ctx, projectArg)
		if err != nil {
			return err
		}

		contractArg := CreateContractParams{
			ProjectID:    result.Project.ID,
			ProposalID:   result.Proposal.ID,
			ClientID:     result.Project.ClientID,
			FreelancerID: result.Proposal.FreelancerID,
			Amount:       result.Proposal.ProposedPrice,
		}
		result.Contract, err = q.CreateContract(ctx, contractArg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
