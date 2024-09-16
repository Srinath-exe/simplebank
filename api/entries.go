package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	db "github.com/Srinath-exe/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type searchEntriesRequest struct {
	SearchQuery string `json:"search_query" binding:"required"`
	Limit       int32  `json:"limit,omitempty"`
	Offset      int32  `json:"offset,omitempty"`
	OrderBy     string `json:"order_by,omitempty"`
	Column1     string `json:"column_1,omitempty"`
	MaxAmount   int64  `json:"max_amount,omitempty"`
	MinAmount   int64  `json:"min_amount,omitempty"`
	MaxDate     string `json:"max_date,omitempty"`
	MinDate     string `json:"min_date,omitempty"`
}

func (server *Server) searchEntries(ctx *gin.Context) {
	var req searchEntriesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println(req)

	// validiate the search query
	Column1 := "id"
	Limit := int32(10)
	Offset := int32(0)
	OrderBy := "DESC"
	minAmount := int64(-10000)
	maxAmount := int64(100000)
	startDate := time.Now().AddDate(-1, 0, 0)
	endDate := time.Now()
	var err error

	if req.Limit > 0 {
		Limit = req.Limit
	}

	if req.Offset > -0 {
		Offset = req.Offset
	}

	if req.OrderBy != "" {
		OrderBy = req.OrderBy
	}

	if req.Column1 != "" {
		Column1 = req.Column1
	}

	if req.MaxAmount != 0 {
		maxAmount = req.MaxAmount
	}

	if req.MinAmount != 0 {
		minAmount = req.MinAmount
	}
	layout := "2006-01-02T15:04:05Z"

	if req.MaxDate != "" {
		endDate, err = time.Parse(layout, req.MaxDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	if req.MinDate != "" {
		startDate, err = time.Parse(layout, req.MinDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	arg := db.SeachEntriesByAccountOwnerParams{
		Field:       Column1,
		Limit:       Limit,
		Offset:      Offset,
		SearchQuery: sql.NullString{String: req.SearchQuery, Valid: true},
		MinAmount:   minAmount,
		MaxAmount:   maxAmount,
		StartDate:   startDate,
		EndDate:     endDate,
		OrderBy:     OrderBy,
	}

	entries, err := server.store.SeachEntriesByAccountOwner(ctx, arg)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type ListEntryFromAccountIdRequest struct {
	ID     int64  `json:"id" binding:"required,min=1"`
	Offset *int32 `json:"offset,omitempty"`
	Limit  *int32 `json:"limit,omitempty"`
}

func (server *Server) listEntriesFromAccountId(ctx *gin.Context) {

	var req ListEntryFromAccountIdRequest

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

	arg := db.ListEntryFromAccountIdParams{
		AccountID: req.ID,
		Limit:     limit,
		Offset:    offset,
	}

	entries, err := server.store.ListEntryFromAccountId(ctx, arg)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)

}

type getEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getEntry(ctx *gin.Context) {
	var req getEntryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	entry, err := server.store.GetEntry(ctx, req.ID)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entry)
}
