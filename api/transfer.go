package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Srinath-exe/simplebank/db/sqlc"
	"github.com/Srinath-exe/simplebank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)

	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccID: req.FromAccountID,
		ToAccID:   req.ToAccountID,
		Amount:    req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

type searchTransferRequest struct {
	SearchQuery string `json:"search_query" binding:"required"`
	Offset      *int32 `json:"offset,omitempty"`
	Limit       *int32 `json:"limit,omitempty"`
}

func (server *Server) searchTransfers(ctx *gin.Context) {
	var searchRequest searchTransferRequest

	if err := ctx.ShouldBindJSON(&searchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//  := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	limit := int32(2)  // Set default limit
	offset := int32(0) // Set default offset

	if searchRequest.Limit != nil && (*searchRequest.Limit > 5 && *searchRequest.Limit < 25) {
		limit = *searchRequest.Limit
	}
	if searchRequest.Offset != nil && (*searchRequest.Offset > 0) {
		offset = *searchRequest.Offset
	}

	req := db.SeachTransfersByAccountOwnerParams{
		SearchQuery: sql.NullString{String: searchRequest.SearchQuery, Valid: true},
		Limit:       limit,
		Offset:      offset,
	}

	transfers, err := server.store.SeachTransfersByAccountOwner(ctx, req)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)

}

type listTransferRequest struct {
	ID     int64  `json:"id" binding:"required,min=1"`
	Offset *int32 `json:"offset,omitempty"`
	Limit  *int32 `json:"limit,omitempty"`
}

func (server *Server) listTransfersFromAccountId(ctx *gin.Context) {

	var req listTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(2)  // Set default limit
	offset := int32(0) // Set default offset

	if req.Limit != nil && (*req.Limit > 5 && *req.Limit < 25) {
		limit = *req.Limit
	}
	if req.Offset != nil && (*req.Offset > 0) {
		offset = *req.Offset
	}

	arg := db.ListTransfersFromAccountIdParams{
		FromAccountID: req.ID,
		Limit:         limit,
		Offset:        offset,
	}

	transfers, err := server.store.ListTransfersFromAccountId(ctx, arg)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)

}
