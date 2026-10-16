package handler

import (
	"github.com/azuvicenna/selaras-services/bff-web/internal/response"
)

var (
	WriteSuccess   = response.WriteSuccess
	WriteCreated   = response.WriteCreated
	WritePaginated = response.WritePaginated
	WriteError     = response.WriteError
	WriteGRPCError = response.WriteGRPCError
	WriteJSON      = response.WriteJSON
	ReadJSON       = response.ReadJSON
)
