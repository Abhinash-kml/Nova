package users

import (
	"fmt"
	"net/http"

	"github.com/abhinash-kml/nova/server/auth"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	logger  *zap.Logger
	config  *config.Config
}

func NewController(s Service, c *config.Config, l *zap.Logger) *Controller {
	return &Controller{
		service: s,
		config:  c,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.getall")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetAllRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	fmt.Println("Passed cursor:", req.Cursor)
	decodedCursor, err := utils.DecodeCursorUUID(req.Cursor)
	if err != nil {
		ctx.Error(err)
		return
	}
	fmt.Println("decoded cursor:", decodedCursor.String())

	users, err := c.service.GetAll(sctx, decodedCursor, req.Limit, req.SortBy, req.SortDirection)
	if err != nil {
		c.logger.Error("Failed to get all users",
			zap.String("cursor", decodedCursor.String()),
			zap.Int("limit", req.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(users))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", req.Id))

	parsedId, _ := uuid.Parse(req.Id)

	user, err := c.service.GetById(sctx, parsedId)
	if err != nil {
		c.logger.Error("Failed to get user",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req CreateRequest

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	user, err := c.service.Add(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create user", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UpdateRequest

	err = ctx.ShouldBindUri(&req.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", req.Id))

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	modifiedUser, err := c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update user",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, modifiedUser)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteRequest

	err = ctx.ShouldBindUri(&req.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", req.Id))

	err = ctx.ShouldBindQuery(&req.DeleteOptions)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("type", req.Type))

	_, err = c.service.Delete(sctx, req)
	if err != nil {
		ctx.Error(err)
		c.logger.Error("Failed to delete user",
			zap.String("user_id", req.Id),
			zap.Error(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.replace")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req ReplaceRequest

	err = ctx.ShouldBindUri(&req.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", req.Id))

	err = ctx.ShouldBindWith(&req.ReplacementData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	replacedUser, err := c.service.Replace(sctx, req)
	if err != nil {
		c.logger.Error("Failed to replace user",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, replacedUser)
}

func (c *Controller) BulkAdd(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.bulkadd")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req BulkCreateRequest

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(req.Users)))

	results, err := c.service.BulkAdd(sctx, req)
	if err != nil {
		c.logger.Error("Bulk add operation failed", zap.Error(err))
		ctx.Error(err)
		return
	}

	var errCount int
	var addResults []common.AddResult
	for index := range results {
		if !results[index].Success {
			errCount++
		}

		addResults = append(addResults, common.AddResult{
			Id:      results[index].Id,
			Status:  results[index].Status,
			Message: results[index].Message,
		})
	}

	response := common.BulkResponse{
		Success: errCount == 0,
		Summary: common.BulkSummary{
			TotalProcessed: len(results),
			TotalSuccess:   len(results) - errCount,
			TotalErrors:    errCount,
		},
		Added: addResults,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) BulkModify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.bulkmodify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req BulkModifyRequest

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(req.Updates)))

	results, err := c.service.BulkModify(sctx, req)
	if err != nil {
		c.logger.Error("Bulk modify operation failed", zap.Error(err))
		ctx.Error(err)
		return
	}

	var errCount int
	var replaceResults []common.ReplaceResult
	for index := range results {
		if !results[index].Success {
			errCount++
		}

		replaceResults = append(replaceResults, common.ReplaceResult{
			Id:      results[index].Id,
			Status:  results[index].Status,
			Message: results[index].Message,
		})
	}

	response := common.BulkResponse{
		Success: errCount == 0,
		Summary: common.BulkSummary{
			TotalProcessed: len(results),
			TotalSuccess:   len(results) - errCount,
			TotalErrors:    errCount,
		},
		Deleted: replaceResults,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) BulkDelete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.bulkdelete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req BulkDeleteRequest

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(req.Users)))

	results, err := c.service.BulkDelete(sctx, req)
	if err != nil {
		c.logger.Error("Bulk delete operation failed", zap.Error(err))
		ctx.Error(err)
		return
	}

	var errCount int
	var deleteResults []common.DeleteResult
	for index := range results {
		if !results[index].Success {
			errCount++
		}

		deleteResults = append(deleteResults, common.DeleteResult{
			Id:      results[index].Id,
			Status:  results[index].Status,
			Message: results[index].Message,
		})
	}

	response := common.BulkResponse{
		Success: errCount == 0,
		Summary: common.BulkSummary{
			TotalProcessed: len(results),
			TotalSuccess:   len(results) - errCount,
			TotalErrors:    errCount,
		},
		Deleted: deleteResults,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Login(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.login")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var loginRequest auth.LoginRequest

	err = ctx.ShouldBindBodyWith(&loginRequest, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	loginResponse, err := c.service.Login(sctx, loginRequest)
	if err != nil {
		c.logger.Error("Failed to login user", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)
}

func (c *Controller) Refresh(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.refresh")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req auth.TokenRefreshRequest

	err = ctx.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Refresh(sctx, req)
	if err != nil {
		c.logger.Error("Failed to refresh user token", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Ban(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.ban")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to ban user", zap.Error(err))
		}
		span.End()
	}()

	var req BanRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	bannedUser, err := c.service.Ban(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, bannedUser)
}

func (c *Controller) Unban(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.unban")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to unban user", zap.Error(err))
		}
		span.End()
	}()

	var req UnbanRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.UnBan(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) CheckUsername(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.checkusername")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to check username", zap.Error(err))
		}
		span.End()
	}()

	var req CheckUsernameRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	valid, err := c.service.CheckUsername(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, valid)
}

func (c *Controller) GetAllUsernames(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "users.controller.getallusernames")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to get all usernames", zap.Error(err))
		}
		span.End()
	}()

	var req UserId

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	usernames, err := c.service.GetAllUsernames(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, usernames)
}

func (c *Controller) Test(ctx *gin.Context) {
	request := struct {
		Head  string `header:"head" json:"head"`
		Uri   string `uri:"uri" json:"uri"`
		Query string `form:"query" json:"query"`
		BodyA string `json:"body_a"`
		BodyB int    `json:"body_b"`
	}{}

	err := ctx.ShouldBindAllSource(&request)
	if err != nil {
		ctx.String(200, err.Error())
	}

	fmt.Println("Header:", request.Head)
	fmt.Println("Uri:", request.Uri)
	fmt.Println("Query:", request.Query)

	ctx.JSON(200, request)
}
