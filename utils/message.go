package utils

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type DefaultResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JsonToStruct(jsonStr string, dynamic any) error {
	if jsonStr == "" {
		return nil
	}

	err := json.Unmarshal([]byte(jsonStr), dynamic)
	if err != nil {
		print("\nError JsonToStruct: " + err.Error())
		return err
	}
	return nil
}
func JsonToArrayStruct(jsonStr string, dynamic any) error {
	if jsonStr == "" {
		return nil
	}
	err := json.Unmarshal([]byte(jsonStr), dynamic)
	if err != nil {
		return err
	}
	return nil
}
func StructToJson(v any) string {
	res, err := json.Marshal(v)
	if err != nil {
		print("\nError StructToJson : " + err.Error())
	}
	return string(res)
}
func ShowResponseDefault(ctx *fasthttp.RequestCtx, statuscode int, msg string, datastring string) {
	ctx.Response.SetStatusCode(statuscode)
	fmt.Fprintf(ctx, StructToJson(DefaultResponse{Status: statuscode, Message: msg, Data: datastring}))
	return
}

func ShowResponseJson(ctx *fasthttp.RequestCtx, statuscode int, msg string, structdata any) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(statuscode)

	resp := DefaultResponse{
		Status:  statuscode,
		Message: msg,
		Data:    structdata, // <- biarkan tetap object
	}

	if jsonBytes, err := json.Marshal(resp); err == nil {
		ctx.Write(jsonBytes)
	} else {
		// fallback kalau gagal marshal
		ctx.WriteString(`{"status":500,"message":"Internal Server Error"}`)
	}
}
