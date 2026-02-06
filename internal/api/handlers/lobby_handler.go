package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Lobby(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"html": strings.TrimSpace(lobbyHTML)})
}

func LobbyData(ctx *gin.Context) {
	var req struct {
		Data string `json:"data"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	modified := strings.TrimSpace(req.Data)
	if modified == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "data is required"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"ok":       true,
		"modified": modified,
		"length":   len(modified),
	})
}

func LobbyFiles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
		return
	}

	type fileInfo struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	}

	resp := make([]fileInfo, 0, len(files))
	for _, f := range files {
		resp = append(resp, fileInfo{Name: f.Filename, Size: f.Size})
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true, "files": resp})
}

const lobbyHTML = `
<div class="lobby">
  <h2>Лобби</h2>
 
  <label>Данные</label>
  <input id="lobbyData" type="text" placeholder="Введите текст" />
  <div id="lobbyDataStatus"></div>
  <label>Файлы</label>
  <input id="lobbyFiles" type="file" multiple />
  <div id="lobbyFilesStatus"></div>
  <div class="actions">
    <button onclick="SendLobbyData()">Отправить</button>
    <button onclick="UploadLobbyFiles()">Загрузить файлы</button>
 
  </div>
</div>
`
