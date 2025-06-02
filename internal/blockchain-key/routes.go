package blockchainkey

import (
	"soft-hsm/internal/blockchain-key/handlers"

	"github.com/go-chi/chi/v5"
)

func BlockchainKeyRoutes(blockchainKeyHandler *handlers.BlockchainKeyHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", blockchainKeyHandler.GenerateKey)
	r.Post("/import", blockchainKeyHandler.ImportKey)
	r.Get("/", blockchainKeyHandler.GetKeysByUserId)
	r.Get("/{id}", blockchainKeyHandler.KeyDetail)
	r.Get("/export-and-delete/{id}", blockchainKeyHandler.ExportAndDeleteKey)
	r.Post("/transaction", blockchainKeyHandler.SendEthereumTransaction)

	return r
}
