package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/codestates/WBABEProject-08/commits/main/controller"
)

type Router struct {
	buyer *controller.BuyerController
	seller *controller.SellerController
}

func NewRouter(ctl *controller.Controller) *Router {
	r := &Router{}
	r.buyer = ctl.GetBuyer()
	r.seller = ctl.GetSeller()

	return r
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		//허용할 header 타입에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		//허용할 method에 대해 열거
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (p *Router) Idx() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORS())

	bRoute := r.Group("/buyer")
	{
		///
	}

	sRoute := r.Group("/seller")
	{
		///
	}
}