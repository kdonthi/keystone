package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type RewindTimeRequest struct {
	TroopId       int `json:"troopID"`
	SecondsBefore int `json:"secondsBefore"`
}

func RewindTime(ctx *EngineCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := RewindTimeRequest{}
		DecodeRequestBody(c, &req)

		id := uuid.New().String()

		tickReq := RewindTimeRequest{
			TroopId:       req.TroopId,
			SecondsBefore: req.SecondsBefore,
		}

		jsonBytes, _ := json.Marshal(tickReq)
		jsonString := string(jsonBytes)

		tick := ctx.Ticker.TickNumber + 1 // TODO what if tick passes by?

		AddTickJob(ctx.World, tick, MoveCalculationTickID, jsonString, strconv.Itoa(req.TileId))
		ctx.AddTickTransaction(MoveCalculationTickID, tick, jsonString)

		c.JSON(http.StatusOK, CreateBasicResponseObject(id))
	}
}
