package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Lobby(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"html": strings.TrimSpace(lobbyHTML)})
}

const lobbyHTML = `
<div class="lobby">
  <h2>ggg</h2>
  <p>ggg.</p>
  <div class="actions">
    <button onclick="CreateGame()">gg</button>
    <button onclick="ConnectToGame()">gg</button>
  </div>
</div>
`
