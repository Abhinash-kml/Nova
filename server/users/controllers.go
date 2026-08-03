package users

import (
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

	var dto GetAllDTO

	err = ctx.ShouldBindQuery(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	decodedCursor, err := utils.DecodeCursor(dto.Cursor)
	if err != nil {
		ctx.Error(err)
		return
	}

	users, err := c.service.GetAll(sctx, decodedCursor, dto.Limit)
	if err != nil {
		c.logger.Error("Failed to get all users",
			zap.Int("cursor", decodedCursor),
			zap.Int("limit", dto.Limit),
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

	var dto GetDTO

	err = ctx.ShouldBindUri(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", dto.Id))

	parsedId, _ := uuid.Parse(dto.Id)

	user, err := c.service.GetById(sctx, parsedId)
	if err != nil {
		c.logger.Error("Failed to get user",
			zap.String("user_id", dto.Id),
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

	var dto CreateDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	user, err := c.service.Add(sctx, dto)
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

	var dto UpdateDTO

	err = ctx.ShouldBindUri(&dto.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", dto.Id))

	if err := ctx.ShouldBindWith(&dto.FieldUpdates, binding.JSON); err != nil {
		ctx.Error(err)
		return
	}

	modifiedUser, err := c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update user",
			zap.String("user_id", dto.Id),
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

	var dto DeleteDTO

	err = ctx.ShouldBindUri(&dto.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", dto.Id))

	err = ctx.ShouldBindQuery(&dto.DeleteOptions)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("type", dto.Type))

	_, err = c.service.Delete(sctx, dto)
	if err != nil {
		ctx.Error(err)
		c.logger.Error("Failed to delete user",
			zap.String("user_id", dto.Id),
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

	var dto ReplaceDTO

	err = ctx.ShouldBindUri(&dto.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", dto.Id))

	err = ctx.ShouldBindWith(&dto.ReplacementData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	replacedUser, err := c.service.Replace(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to replace user",
			zap.String("user_id", dto.Id),
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

	var dto BulkCreateDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(dto.Users)))

	results, err := c.service.BulkAdd(sctx, dto)
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

	var dto BulkModifyDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(dto.Updates)))

	results, err := c.service.BulkModify(sctx, dto)
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

	var dto BulkDeleteDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(dto.Users)))

	results, err := c.service.BulkDelete(sctx, dto)
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

	var dto auth.TokenRefreshRequest

	err = ctx.ShouldBindBodyWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Refresh(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to refresh user token", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
