package api

import "github.com/gin-gonic/gin"

func (h *Handler) getSongs(c *gin.Context) {
	sort := c.DefaultQuery("sort", "date")
	page := c.DefaultQuery("page", "1")
}

func (h *Handler) createSong(c *gin.Context) {
	sort := c.DefaultQuery("sort", "group")
	order := c.DefaultQuery("order", "asc")
	page := c.DefaultQuery("page", "1")
}
func (h *Handler) getSongById(c *gin.Context) {
}
func (h *Handler) deleteSong(c *gin.Context) {
}
func (h *Handler) updateSong(c *gin.Context) {
}
