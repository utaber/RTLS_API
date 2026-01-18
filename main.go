package main

import (
	"context"

	"RTLS_API/auth"
	"RTLS_API/config"
	"RTLS_API/pkg/barang"
	"RTLS_API/pkg/firebase"
	"RTLS_API/pkg/user"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	dbClient := firebase.NewDatabase(
		ctx,
		cfg.DBKEY,
		cfg.DBURL,
	)

	firebaseAdapter := firebase.NewAdapter(dbClient)

	userService := user.NewService(ctx, firebaseAdapter)

	jwtService := auth.NewJWTService(
		cfg.JWTSecret,
		cfg.JWTAlgo,
		cfg.JWTExpire,
	)

	authHandler := auth.NewHandler(jwtService, userService)
	barangService := barang.NewService(ctx, dbClient)
	barangHandler := barang.NewHandler(barangService)

	r := gin.Default()

	r.POST("/auth/login", authHandler.Login)

	authMW := jwtService.Middleware()

	r.GET("/barang", barangHandler.Get)
	r.POST("/barang", authMW, barangHandler.Create)
	r.PATCH("/barang/:device_id", authMW, barangHandler.Update)
	r.DELETE("/barang/:device_id", authMW, barangHandler.Delete)
	r.DELETE("/meta", authMW, barangHandler.ResetSystem)

	r.Run(":8000")
}
