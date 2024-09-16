package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/Srinath-exe/simplebank/db/sqlc"
	"github.com/Srinath-exe/simplebank/token"
	"github.com/Srinath-exe/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {

		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {

		// Check if the error is a not found error
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Return the user
	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPasswordHash(req.Password, user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)

}

type updatePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
	Username    string `json:"username" binding:"required,alphanum"`
}

func (server *Server) updatePassword(ctx *gin.Context) {
	var req updatePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if user.Username != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdatePasswordParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}

	err = server.store.UpdatePassword(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "password updated"})
}

type deleteUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// user, err := server.store.GetUser(ctx, req.Username)

	// if err != nil {

	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, errorResponse(err))
	// 		return
	// 	}

	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// if user.Username != authPayload.Username {
	// 	err := errors.New("account doesn't belong to the authenticated user")
	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	// 	return
	// }

	result, err := server.store.DeleteUserWithAccountsTx(ctx, req.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type searchUsersRequest struct {
	Username string `json:"username" binding:"required"`
	Limit    *int32 `json:"limit,omitempty"`
	Offset   *int32 `json:"offset,omitempty"`
}

func (server *Server) searchUsers(ctx *gin.Context) {
	var req searchUsersRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(10)
	offset := int32(0)

	if req.Limit != nil {
		limit = *req.Limit
	}

	if req.Offset != nil {
		offset = *req.Offset
	}

	arg := db.SearchUsersParams{
		Column1: sql.NullString{String: req.Username, Valid: true},
		Limit:   limit,
		Offset:  offset,
	}

	users, err := server.store.SearchUsers(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	usersList := make([]userResponse, len(users))

	for i, user := range users {
		usersList[i] = newUserResponse(user)
	}

	ctx.JSON(http.StatusOK, usersList)
}

type getUsersRequest struct {
	Usernames []string `json:"usernames" binding:"required"`
	Limit     *int32   `json:"limit,omitempty"`
	Offset    *int32   `json:"offset,omitempty"`
}

func (server *Server) getUsers(ctx *gin.Context) {
	var req getUsersRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := int32(10)
	offset := int32(0)

	if req.Limit != nil {
		limit = *req.Limit
	}

	if req.Offset != nil {
		offset = *req.Offset
	}

	arg := db.GetUsersParams{
		Usernames: req.Usernames,
		Limit:     limit,
		Offset:    offset,
	}

	accounts, err := server.store.GetUsers(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	usersList := make([]userResponse, len(accounts))

	for i, user := range accounts {
		usersList[i] = newUserResponse(user)
	}

	ctx.JSON(http.StatusOK, usersList)
}
