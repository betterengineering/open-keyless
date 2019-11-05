// Copyright 2019 Mark Spicer
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package server

import (
	"github.com/lodge93/open-keyless/api/shopmanager"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*shopmanager.User).UserName,
		"text":     "Hello World.",
	})
}

type Server struct {
	appName    string
	signingKey string
	router     *gin.Engine
}

func NewServer(appName string, signingKey string) (*Server, error) {
	router := gin.New()
	srv := &Server{
		appName: appName,
		router:  router,
	}

	// Configure router.
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(ginprom.PromMiddleware(nil))

	// Configure authentication.
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           appName,
		Key:             []byte(signingKey),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     srv.payloadFunc,
		IdentityHandler: srv.identityHandler,
		Authenticator:   srv.authenticator,
		Authorizator:    srv.authorizator,
		Unauthorized:    srv.unauthorized,
	})
	if err != nil {
		return nil, err
	}

	// Configure admin routes.
	router.GET("/healthz", srv.healthz)
	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	// Configure login.
	router.POST("/api/shop-manager/login", authMiddleware.LoginHandler)

	// Configure api endpoints.
	api := router.Group("/api/shop-manager/v1")
	api.Use(authMiddleware.MiddlewareFunc())
	api.GET("/refresh-token", authMiddleware.RefreshHandler)
	api.GET("/hello", helloHandler)

	// Configure no route.
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return srv, err
}

func (srv *Server) Run() {
	srv.router.Run()
}
