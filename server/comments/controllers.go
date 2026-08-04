package comments

import (
	"net/http"

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
	config  *config.Config
	logger  *zap.Logger
}

func NewController(s Service, c *config.Config, l *zap.Logger) *Controller {
	return &Controller{
		config:  c,
		service: s,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.getall")
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

	comments, err := c.service.GetAll(sctx, decodedCursor, dto.Limit)
	if err != nil {
		c.logger.Error("Failed to get all comments",
			zap.String("cursor", dto.Cursor),
			zap.Int("limit", dto.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(comments))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.get")
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

	span.SetAttributes(attribute.String("comment_id", dto.Id))

	parsedUuid, _ := uuid.Parse(dto.Id)

	comment, err := c.service.GetById(sctx, parsedUuid)
	if err != nil {
		c.logger.Error("Failed to get comment by id",
			zap.String("comment_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.create")
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

	comment, err := c.service.Add(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create comment", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UpdateDTO

	err = ctx.ShouldBindUri(&dto.CommentId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("comment_id", dto.Id))

	err = ctx.ShouldBindWith(&dto.Body, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	modifiedComment, err := c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update comment",
			zap.String("comment_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, modifiedComment)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto DeleteDTO

	err = ctx.ShouldBindUri(&dto.CommentId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("comment_id", dto.Id))

	err = ctx.ShouldBindQuery(&dto.DeleteOptions)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("type", dto.Type))

	_, err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete comment",
			zap.String("comment_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.replace")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto ReplaceDTO

	err = ctx.ShouldBindUri(&dto.CommentId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("comment_id", dto.Id))

	err = ctx.ShouldBindWith(&dto.Body, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	replacedComment, err := c.service.Replace(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to replace comment",
			zap.String("comment_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, replacedComment)
}

func (c *Controller) BulkAdd(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.bulkadd")
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

	span.SetAttributes(attribute.Int("count", len(dto.Comments)))

	results, err := c.service.BulkAdd(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to perform bulk add comments", zap.Error(err))
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
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.bulkmodify")
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
		c.logger.Error("Failed to perform bulk update comments", zap.Error(err))
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
	sctx, span := tracer.Start(ctx.Request.Context(), "comments.controller.bulkdelete")
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

	span.SetAttributes(attribute.Int("count", len(dto.Comments)))

	results, err := c.service.BulkDelete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to perform bulk delete comments", zap.Error(err))
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
