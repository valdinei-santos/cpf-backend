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

// AccessCounterMiddleware -- Middleware para Contar Acessos
func AccessCounterMiddleware(c *gin.Context) {
	// Usa o path da rota para contagem. Ex: /api/v1/cpfs
	stats.GlobalStats.Increment(c.Request.URL.Path)
	c.Next() // Processa o restante da requisição
}

// GetStatsHandler -- Função Handler para exibir as estatísticas
func GetStatsHandler(c *gin.Context) {
	statsData := stats.GlobalStats.GetStats()
	uptime := time.Since(stats.GlobalStats.StartTime).Round(time.Second).String()

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"uptime":        uptime,
		"access_counts": statsData,
	})
}

func InitRoutes(router *gin.RouterGroup, log logger.ILogger, db *mongo.Database) {

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8888", "http://127.0.0.1:8888"}, // Para liberar o Swagger
		// OU use AllowAllOrigins: true para permitir TUDO (apenas para DEV/TESTE)

		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Limita tempo navegador pode armazenar em cache as informações do preflight
	}))
	// ---------------------------

	router.Use(AccessCounterMiddleware) // Adicionar o Middleware de Contagem antes de todas as rotas
	router.GET("/stats", GetStatsHandler)
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := router.Group("/api/v1")
	prod := v1.Group("/cliente")

	prod.POST("/", func(c *gin.Context) {
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

	prod.GET("/", func(c *gin.Context) {
		log.Info("### Start endpoint " + c.Request.Method + " " + c.Request.URL.Path)
		cliente.StartGetAll(log, c, db)
	})

	prod.PUT("/:id", func(c *gin.Context) {
		log.Info("### Start endpoint " + c.Request.Method + " " + c.Request.URL.Path)
		cliente.StartUpdate(log, c, db)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

}
