package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"

	"go-gin-test/controllers"
	dbCon "go-gin-test/db/sqlc"
	routes "go-gin-test/route"
	"go-gin-test/util"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	_ "go-gin-test/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Go Gin Test Api
//	@version		1.0

var (
	server *gin.Engine
	db     *dbCon.Queries
	ctx    context.Context

	ContactController controllers.ContactController
	ContactRoutes     routes.ContactRoutes
)

func main() {
	ctx = context.TODO()
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("could not loadconfig: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), config.DbSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			fmt.Println("Error closing connection...")
		}
	}(conn, context.Background())

	db = dbCon.New(conn)

	fmt.Println("PostgreSql connected successfully...")

	ContactController = *controllers.NewContactController(db, ctx)
	ContactRoutes = routes.NewRouteContact(ContactController)
	server = gin.Default()

	router := server.Group("/api")

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "The contact APi is working fine"})
	})

	ContactRoutes.ContactRoute(router)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	log.Fatal(server.Run(":" + config.ServerAddress))
}
