package main

import (
	"context"

	"RTLS_API/config"
	"RTLS_API/pkg/auth"
	"RTLS_API/pkg/barang"
	"RTLS_API/pkg/firebase"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	dbClient := firebase.NewDatabase(
		ctx,
		"service-account-key-test.json",
		cfg.DBURL,
	)

	jwtService := auth.NewJWTService(
		cfg.JWTSecret,
		cfg.JWTAlgo,
		cfg.JWTExpire,
	)

	barangService := barang.NewService(ctx, dbClient)
	barangHandler := barang.NewHandler(barangService)

	r := gin.Default()

	r.POST("/auth/token", func(c *gin.Context) {
		token, _ := jwtService.GenerateToken()
		c.JSON(200, gin.H{"access_token": token})
	})

	authMW := jwtService.Middleware()

	r.GET("/barang", barangHandler.Get)
	r.POST("/barang", authMW, barangHandler.Create)
	r.PATCH("/barang/:device_id", authMW, barangHandler.Update)
	r.DELETE("/barang/:device_id", authMW, barangHandler.Delete)
	r.DELETE("/meta", authMW, barangHandler.ResetSystem)

	r.Run(":8000")
}
