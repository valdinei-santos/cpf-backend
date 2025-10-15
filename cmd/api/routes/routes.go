// Pacote routes configura as rotas da API.
package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/valdinei-santos/cpf-backend/cmd/api/stats"
	_ "github.com/valdinei-santos/cpf-backend/docs" // swagger docs
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoutes(router *gin.RouterGroup, log logger.ILogger, db *mongo.Database) {

	router.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins:     []string{"http://192.168.37.143:8888", "http://localhost:8888", "http://127.0.0.1:8888"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// ---------------------------

	router.Use(AccessCounterMiddleware) // Adicionar o Middleware de Contagem antes de todas as rotas
	router.GET("/status", GetStatusHandler)
	router.GET("/ping", GetPingHandler)

	v1 := router.Group("/api/v1")
	prod := v1.Group("/cliente")

	prod.OPTIONS("", func(c *gin.Context) {
		c.Status(204)
	})

	prod.POST("", func(c *gin.Context) {
		// Para saber algumas informações que vem no Header que o CORS se importa
		// origin := c.GetHeader("Origin")
		// host := c.GetHeader("X-Forwarded-Host")
		// realHost := c.Request.Host
		// fmt.Printf("CORS DEBUG: Origin Header: %s\n", origin)
		// fmt.Printf("CORS DEBUG: X-Forwarded-Host Header: %s\n", host)
		// fmt.Printf("CORS DEBUG: Request Host (Go): %s\n", realHost)

		log.Info("### Start endpoint " + c.Request.Method + " " + c.Request.URL.Path)
		cliente.StartCreate(log, c, db)
	})

	prod.DELETE("/:id", func(c *gin.Context) {
		log.Info("### Start endpoint " + c.Request.Method + " " + c.Request.URL.Path)
		cliente.StartDelete(log, c, db)
	})

	prod.GET("/:id", func(c *gin.Context) {
		log.Info("### Start endpoint " + c.Request.Method + " " + c.Request.URL.Path)
		cliente.StartGet(log, c, db)
	})

	prod.GET("", func(c *gin.Context) {
		log.Info("### Start endpoint " + c.Request.Method + " " + c.Request.URL.Path)
		cliente.StartGetAll(log, c, db)
	})

	prod.OPTIONS("/:id", func(c *gin.Context) {
		// Isso garante que o Preflight Request seja recebido e respondido com 204.
		c.Status(204)
	})

	prod.PUT("/:id", func(c *gin.Context) {
		log.Info("### Start endpoint " + c.Request.Method + " " + c.Request.URL.Path)
		cliente.StartUpdate(log, c, db)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

}

// AccessCounterMiddleware -- Middleware para Contar Acessos
func AccessCounterMiddleware(c *gin.Context) {
	// Usa o path da rota para contagem. Ex: /api/v1/cpfs
	stats.GlobalStats.Increment(c.Request.URL.Path)
	c.Next() // Processa o restante da requisição
}

// @Summary      Retorna o status da API
// @Description  Retorna o status da API
// @Tags         util
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} dto.OutputDefault
// @Router       /status [get]
func GetStatusHandler(c *gin.Context) {
	statsData := stats.GlobalStats.GetStats()
	uptime := time.Since(stats.GlobalStats.StartTime).Round(time.Second).String()

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"uptime":        uptime,
		"access_counts": statsData,
	})
}

// @Summary      Retorna pong
// @Description  Retorna pong se estiver tudo ok com a API
// @Tags         util
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      400 {object} dto.OutputDefault
// @Router       /ping [get]
func GetPingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
