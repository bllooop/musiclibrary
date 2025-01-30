package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/gin-gonic/gin"
)

type getSongsResponse struct {
	Data map[string]interface{} `json:"data"`
}

func (h *Handler) getSongs(c *gin.Context) {
	sort := c.DefaultQuery("sort", "artist")
	order := strings.ToUpper(c.DefaultQuery("order", "asc"))
	name := c.Query("name")
	group := c.Query("artist")
	text := c.Query("text")
	releasedate := c.Query("releasedate")
	link := c.Query("link")

	page := c.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	lists, err := h.usecases.GetSongsLibrary(order, sort, pageInt, name, group, text, releasedate, link)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getSongsResponse{
		Data: lists,
	})
}

func (h *Handler) createSong(c *gin.Context) {
	var input domain.UpdateSong
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.Group == nil || input.Name == nil {
		newErrorResponse(c, http.StatusBadRequest, "empty fields for adding a song")
		return
	}
	group := ""
	name := ""
	if input.Group != nil {
		group = *input.Group
	}
	if input.Name != nil {
		name = *input.Name
	}
	url := fmt.Sprintf("https://api.example.com/info?group=%s&song=%s", group, name)
	resp, err := http.Get(url)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error when making external get request")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error reading response body")
		return
	}
	var song domain.UpdateSong
	err = json.Unmarshal(body, &song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error unmarshaling JSON")
		return
	}
	id, err := h.usecases.CreateSong(input, song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
func (h *Handler) getSongById(c *gin.Context) {
}
func (h *Handler) deleteSong(c *gin.Context) {
	songid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id value")
		return
	}
	err = h.usecases.DeleteSong(songid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}
func (h *Handler) updateSong(c *gin.Context) {

	songid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id value")
		return
	}

	var input domain.UpdateSong
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecases.Update(songid, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
