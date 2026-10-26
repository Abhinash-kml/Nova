package channels

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	config  *config.ChannelsConfig
	logger  *zap.Logger
}

func NewController(s Service, c *config.ChannelsConfig, l *zap.Logger) *Controller {
	return &Controller{
		config:  c,
		service: s,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.GetAll")
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

	decodedCursor, err := utils.DecodeCursorUUID(req.Cursor)
	if err != nil {
		ctx.Error(err)
		return
	}

	channels, err := c.service.GetAll(sctx, decodedCursor, req.Limit)
	if err != nil {
		c.logger.Error("Failed to get all channels", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(channels))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.Get")
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

	parsedId, _ := uuid.Parse(req.Id)

	channel, err := c.service.GetById(sctx, parsedId)
	if err != nil {
		c.logger.Error("Failed to get channel",
			zap.String("channel_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, channel)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req CreateRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	channel, err := c.service.Add(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create channel", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, channel)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.Modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UpdateRequest

	err = ctx.ShouldBindUri(&req.ChannelId)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ctx.ShouldBindBodyWithJSON(&req.ChannelModifications)
	if err != nil {
		ctx.Error(err)
		return
	}

	updatedChannel, err := c.service.Modify(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update channel",
			zap.String("channel_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, updatedChannel)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteRequest

	err = ctx.ShouldBindUri(&req.ChannelId)
	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete channel",
			zap.String("channel_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) BulkAdd(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.BulkAdd")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req BulkCreateRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(req.Channels)))

	results, err := c.service.BulkAdd(sctx, req)
	if err != nil {
		c.logger.Error("Failed to bulk add channels", zap.Error(err))
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
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.BulkModify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req BulkModifyRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(req.Updates)))

	results, err := c.service.BulkModify(sctx, req)
	if err != nil {
		c.logger.Error("Failed to to bulk modify channels", zap.Error(err))
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
	sctx, span := tracer.Start(ctx.Request.Context(), "channels.controller.BulkDelete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req BulkDeleteRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("count", len(req.Channels)))

	results, err := c.service.BulkDelete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to bulk delete channels", zap.Error(err))
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
