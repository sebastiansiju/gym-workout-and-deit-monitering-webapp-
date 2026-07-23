package controllers

import (
	"github.com/Cawlumm/sebu-backend/config"
	"github.com/Cawlumm/sebu-backend/utils"
	"github.com/gin-gonic/gin"
)

// ServerInfo is a public, unauthenticated endpoint clients use to verify that a
// configured server URL is reachable and is actually a Sebu backend — the
// "test connection" behind the in-app server selector. It runs under the same
// CORS policy as the rest of the API, so a successful probe honestly predicts
// that authenticated requests from the same origin will be allowed too.
func (h *Handler) ServerInfo(c *gin.Context) {
	utils.OK(c, gin.H{
		"name":    "sebu",
		"version": config.C.Version,
	})
}
