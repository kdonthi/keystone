package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type RewindTimeRequest struct {
	TroopId     int `json:"troopID"`
	TicksBefore int `json:"ticksBefore"`
}

func RewindTime(ctx *EngineCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := RewindTimeRequest{}
		DecodeRequestBody(c, &req)

		tickReq := RewindTimeRequest{
			TroopId:     req.TroopId,
			TicksBefore: req.TicksBefore,
		}

		jsonBytes, _ := json.Marshal(tickReq)
		jsonString := string(jsonBytes)

		tick := ctx.Ticker.TickNumber + 1 // TODO what if tick passes by? should we also have a method to calculate this?

		AddTickJob(ctx.World, tick, MoveCalculationTickID, jsonString, strconv.Itoa(req.TroopId)) // TODO what is the point of the tickID?
		ctx.AddTickTransaction(MoveCalculationTickID, tick, jsonString)

		id := uuid.New().String()
		c.JSON(http.StatusOK, CreateBasicResponseObject(id))
	}
}
