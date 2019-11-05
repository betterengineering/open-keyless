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
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/lodge93/open-keyless/api/shopmanager"
)

const (
	identityKey = "id"
)

func (srv *Server) payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*shopmanager.User); ok {
		return jwt.MapClaims{
			identityKey: v.UserName,
		}
	}
	return jwt.MapClaims{}
}

func (srv *Server) identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &shopmanager.User{
		UserName: claims[identityKey].(string),
	}
}

func (srv *Server) authenticator(c *gin.Context) (interface{}, error) {
	var loginVals shopmanager.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &shopmanager.User{
			UserName:  userID,
			LastName:  "Bo-Yi",
			FirstName: "Wu",
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

func (srv *Server) authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*shopmanager.User); ok && v.UserName == "admin" {
		return true
	}

	return false
}

func (srv *Server) unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
