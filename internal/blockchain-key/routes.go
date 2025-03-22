package blockchainkey

import (
	"soft-hsm/internal/blockchain-key/handlers"

	"github.com/go-chi/chi/v5"
)

func BlockchainKeyRoutes(blockchainKeyHandler *handlers.BlockchainKeyHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", blockchainKeyHandler.GenerateKey)

	return r
}
