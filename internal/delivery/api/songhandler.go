package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bllooop/musiclibrary/internal/domain"
	logger "github.com/bllooop/musiclibrary/pkg/logging"
	"github.com/gin-gonic/gin"
)

type getSongsResponse struct {
	Data map[string]interface{} `json:"data"`
}

// @Summary Find song text
// @Tags songList
// @Description получение текста песни
// @ID find-songtext
// @Produce  json
// @Param       name    query     string  false  "text search by id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {string} message
// @Failure 500 {string} message
// @Failure default {string} message
// @Router /api/songs/song-text [get]
func (h *Handler) getSongById(c *gin.Context) {
	logger.Log.Info().Msg("Received request for getting song text / Получен запрос на получение текста песни")
	if c.Request.Method != http.MethodGet {
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос GET")
		return
	}
	songName := c.Query("name")
	if songName == "" {
		newErrorResponse(c, http.StatusBadRequest, "name can't be empty / имя не может быть пустым")
		return
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

	//logger.Log.Debug().Msgf("Name: %s, verse_number: %s-%s", songName, begin, end)
	list, err := h.usecases.GetSongsById(songName, begin, end)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for getting song text / Получен ответ для получения текста песни")

	c.JSON(http.StatusOK, gin.H{
		"song_name": songName,
		"verse":     list,
	})

}

// @Summary Get all songs
// @Tags songList
// @Description получение списка песен
// @ID get-songs
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {string} message
// @Failure 500 {string} message
// @Failure default {string} message
// @Router /api/songs [get]
func (h *Handler) getSongs(c *gin.Context) {
	logger.Log.Info().Msg("Received request for get songs / Получили запрос на получение песен")
	if c.Request.Method != http.MethodGet {
		newErrorResponse(c, http.StatusBadRequest, "GET request required / Требуется запрос GET")
		return
	}
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
	logger.Log.Info().Msg("Received response for get songs / Получен ответ на получение песен")

	c.JSON(http.StatusOK, getSongsResponse{
		Data: lists,
	})
}

// @Summary Create song
// @Tags songList
// @Description добавление песни в базу данных
// @ID create-song
// @Accept  json
// @Produce  json
// @Param input body domain.UpdateSong true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {string} message
// @Failure 500 {string} message
// @Failure default {string} message
// @Router /api/songs [post]
func (h *Handler) createSong(c *gin.Context) {
	logger.Log.Info().Msg("Received request for create song / Получен запрос на создание песни")
	if c.Request.Method != http.MethodPost {
		newErrorResponse(c, http.StatusBadRequest, "POST request required / Требуется запрос POST")
		return
	}
	var input domain.UpdateSong
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.Group == nil || input.Name == nil {
		newErrorResponse(c, http.StatusBadRequest, "empty fields for adding a song / поля для добавления песни пустые")
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
	url := fmt.Sprintf("https://api.example.com/info?group=%s&song=%s", url.QueryEscape(group), url.QueryEscape(name))
	resp, err := http.Get(url)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error when making external get request / ошибка при выполнении внешнего запроса на получение")
		return
	}
	if resp.StatusCode != http.StatusOK {
		newErrorResponse(c, resp.StatusCode, "error: received non-200 status code / ошибка: получен код состояния не-200")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error reading response body / ошибка чтения тела ответа")
		return
	}
	if len(body) == 0 {
		newErrorResponse(c, http.StatusInternalServerError, "error: empty response body / ошибка: пустое тело ответа")
		return
	}
	logger.Log.Debug().Str("response_body", string(body)).Msg("Successfully read response body / Успешно прочитано тело ответа")

	var song domain.UpdateSong
	err = json.Unmarshal(body, &song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error unmarshaling JSON / ошибка при обработке JSON")
		return
	}
	logger.Log.Debug().
		Interface("parsed_song", song).
		Msg("Successfully unmarshaled JSON into UpdateSong struct / Успешно обработали JSON в структуре UpdateSong")
	id, err := h.usecases.CreateSong(input, song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for creating song / Получен ответ на создание песни")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Delete song from list
// @Security ApiKeyAuth
// @Tags songList
// @Description delete song from list by id
// @ID delete-list
// @Produce  json
// @Param       id    query     int  false  "song delete by id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {string} message
// @Failure 500 {string} message
// @Failure default {string} message
// @Router /api/songs [delete]
func (h *Handler) deleteSong(c *gin.Context) {
	logger.Log.Info().Msg("Received request for delete song / Получен запрос на удаление песни")
	if c.Request.Method != http.MethodDelete {
		newErrorResponse(c, http.StatusBadRequest, "DELETE request required / Требуется запрос DELETE")
		return
	}
	songid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id value / недопустимое значение идентификатора")
		return
	}
	logger.Log.Debug().Int("id parameter", songid).Msg("Successfully read song id / Успешно прочитан идентификатор песни")
	err = h.usecases.DeleteSong(songid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for deleting song / Получен ответ на удаление песни")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Update song
// @Tags songList
// @Description обновление данных песни
// @ID update-list
// @Accept  json
// @Produce  json
// @Param input body domain.UpdateSong true "list info"
// @Param       id    query     int  false  "song update by id"
// @Success 200 {integer} integer 1
// @Failure 400,404 {string} message
// @Failure 500 {string} message
// @Failure default {string} message
// @Router /api/songs [put]
func (h *Handler) updateSong(c *gin.Context) {
	logger.Log.Info().Msg("Received request for updating song / Получен запрос на обновление песни")
	if c.Request.Method != http.MethodPut {
		newErrorResponse(c, http.StatusBadRequest, "PUT request required / Требуется запрос PUT")
		return
	}
	songid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id value / недопустимое значение идентификатора")
		return
	}
	logger.Log.Debug().Int("id parameter", songid).Msg("Successfully read song id / Успешно прочитан идентификатор песни")

	var input domain.UpdateSong
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Log.Debug().
		Interface("binded_song", input).
		Msg("Successfully binded JSON input to struct / Успешное обработка JSON со структурой")

	if err := h.usecases.Update(songid, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Received response for updating song / Получен ответ на обновление песни")

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
