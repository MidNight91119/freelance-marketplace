package api

import "net/http"

func (server *Server) listContracts(w http.ResponseWriter, r *http.Request) {
	payload, ok := payloadFrom(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired access token")
		return
	}
	
	contracts, err := server.store.ListContractsByUser(r.Context(), payload.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong")
		return
	}

	writeJSON(w, http.StatusOK, contracts)
}