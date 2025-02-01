package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bllooop/musiclibrary/internal/domain"
	logger "github.com/bllooop/musiclibrary/pkg"
	"github.com/gin-gonic/gin"
)

type getSongsResponse struct {
	Data map[string]interface{} `json:"data"`
}

func (h *Handler) getSongById(c *gin.Context) {

	logger.Log.Info().Msg("Received request for getting song text")

	songName := c.Param("name")
	if songName == "" {
		newErrorResponse(c, http.StatusBadRequest, "name can't be empty")
	}

	begin, err := strconv.Atoi(c.DefaultQuery("begin", "1"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	end, err := strconv.Atoi(c.DefaultQuery("end", "1"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Log.Debug().Msgf("Name: %s, verse_number: %s-%s", songName, begin, end)
	list, err := h.usecases.GetSongsById(songName, begin, end)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for getting song text")

	c.JSON(http.StatusOK, gin.H{
		"song_name": songName,
		"verse":     list,
	})

}
func (h *Handler) getSongs(c *gin.Context) {
	logger.Log.Info().Msg("Received request for get songs")

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
	logger.Log.Debug().Msgf("Successfully read Sort: %s, Order: %s, Page: %v", sort, order, page)
	logger.Log.Debug().Msgf("Successfully read Name: %s, Group: %s, Text: %s, Release Date: %s, Link: %s ", name, group, text, releasedate, link)

	lists, err := h.usecases.GetSongsLibrary(order, sort, pageInt, name, group, text, releasedate, link)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for get songs")

	c.JSON(http.StatusOK, getSongsResponse{
		Data: lists,
	})
}

func (h *Handler) createSong(c *gin.Context) {
	logger.Log.Info().Msg("Received request for create song")

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
	logger.Log.Debug().Msgf("Successfully read Name: %s, Group: %s", name, group)
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
	logger.Log.Debug().Str("response_body", string(body)).Msg("Successfully read response body")

	var song domain.UpdateSong
	err = json.Unmarshal(body, &song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error unmarshaling JSON")
		return
	}
	logger.Log.Debug().
		Interface("parsed_song", song).
		Msg("Successfully unmarshaled JSON into UpdateSong struct")

	id, err := h.usecases.CreateSong(input, song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	logger.Log.Info().Msg("Received response for creating song")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) deleteSong(c *gin.Context) {
	logger.Log.Info().Msg("Received request for delete song")

	songid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id value")
		return
	}
	logger.Log.Debug().Int("id parameter", songid).Msg("Successfully read song id")
	err = h.usecases.DeleteSong(songid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for deleting song")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}
func (h *Handler) updateSong(c *gin.Context) {
	logger.Log.Info().Msg("Received request for updating song")

	songid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id value")
		return
	}
	logger.Log.Debug().Int("id parameter", songid).Msg("Successfully read song id")

	var input domain.UpdateSong
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Log.Debug().
		Interface("binded_song", input).
		Msg("Successfully binded JSON input to struct")

	if err := h.usecases.Update(songid, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for updating song")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
