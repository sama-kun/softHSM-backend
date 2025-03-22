package handlers

import (
	"context"
	"net/http"
	"soft-hsm/internal/blockchain-key/dto"
	"soft-hsm/internal/blockchain-key/models"
	"soft-hsm/internal/blockchain-key/services"
	"soft-hsm/internal/common/validators"
	"soft-hsm/internal/middleware"
)

type BlockchainKeyHandler struct {
	blockchainKeyService services.BlockchainKeyServiceInterface
}

func NewBlockchainKeyHandler(blockchainKeyService services.BlockchainKeyServiceInterface) *BlockchainKeyHandler {
	return &BlockchainKeyHandler{blockchainKeyService: blockchainKeyService}
}

func (h *BlockchainKeyHandler) GenerateKey(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.GetUserFromContext(r)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
		return
	}
	var req dto.GenerateKeyDTO

	if err := middleware.DecodeJSON(r, &req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := validators.ValidateStruct(req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "invalid input")
		return
	}

	if !models.IsValidBlockchain(req.Blockchain) {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid blockchain type")
		return
	}

	sessionUser, err := middleware.ExtractAndDecryptSessionToken(r)
	if err != nil {
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized Session")
		return
	}

	password := "Test"

	resp, err := h.blockchainKeyService.GenerateEthereumKey(context.Background(), int64(sessionUser.Id), password, req)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Generate key failed")
		return
	}

	middleware.JSONResponse(w, http.StatusCreated, resp)
}

// func (h *BlockchainKeyHandler) ImportPrivateKey(w http.ResponseWriter, r *http.Request) {
// 	user, err := middleware.GetUserFromContext(r)

// 	if err != nil {
// 		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
// 		return
// 	}

// 	var req dto.ImportKeyDTO

// 	if err := middleware.DecodeJSON(r, &req); err != nil {
// 		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
// 		return
// 	}

// 	if !models.IsValidBlockchain(req.Blockchain) {
// 		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid blockchain type")
// 		return
// 	}

// 	resp, err := h.blockchainKeyService.ImportPrivateKey(context.Background(), int64(user.Id), req)

// 	if err != nil {
// 		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Generate key failed")
// 		return
// 	}

// 	middleware.JSONResponse(w, http.StatusCreated, resp)
// }
