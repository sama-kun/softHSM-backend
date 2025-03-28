package handlers

import (
	"context"
	"fmt"
	"net/http"
	"soft-hsm/internal/blockchain-key/dto"
	"soft-hsm/internal/blockchain-key/models"
	"soft-hsm/internal/blockchain-key/services"
	"soft-hsm/internal/common/validators"
	"soft-hsm/internal/middleware"
	"strings"

	"github.com/google/uuid"
)

type BlockchainKeyHandler struct {
	blockchainKeyService services.BlockchainKeyServiceInterface
}

func NewBlockchainKeyHandler(blockchainKeyService services.BlockchainKeyServiceInterface) *BlockchainKeyHandler {
	return &BlockchainKeyHandler{blockchainKeyService: blockchainKeyService}
}

func (h *BlockchainKeyHandler) KeyDetail(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[2] == "" {
		http.Error(w, "Key ID is required", http.StatusBadRequest)
		return
	}

	keyIDStr := parts[3]
	fmt.Println("Extracted UUID string:", keyIDStr)

	// Парсим UUID
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid key ID")
		return
	}

	resp, err := h.blockchainKeyService.KeyDetail(context.Background(), keyID, int64(user.Id))

	if err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "key failed")
		return
	}

	middleware.JSONResponse(w, http.StatusCreated, resp)
}

func (h *BlockchainKeyHandler) GenerateKey(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)

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

	resp, err := h.blockchainKeyService.GenerateEthereumKey(context.Background(), int64(user.Id), req)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Generate key failed")
		return
	}

	middleware.JSONResponse(w, http.StatusCreated, resp)
}

func (h *BlockchainKeyHandler) ImportKey(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.GetUserFromContext(r)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
		return
	}
	var req dto.ImportKeyDTO

	if err := middleware.DecodeJSON(r, &req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := validators.ValidateStruct(req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "invalid input")
		return
	}

	sessionUser, err := middleware.ExtractAndDecryptSessionToken(r)
	if err != nil {
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized Session")
		return
	}

	resp, err := h.blockchainKeyService.ImportEthereumKey(context.Background(), int64(sessionUser.Id), req)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Generate key failed")
		return
	}

	middleware.JSONResponse(w, http.StatusCreated, resp)
}

func (h *BlockchainKeyHandler) GetKeysByUserId(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	resp, err := h.blockchainKeyService.FindKeysByUserID(context.Background(), int64(user.Id))

	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Generate key failed")
		return
	}

	middleware.JSONResponse(w, http.StatusCreated, resp)
}
