package dynamicipupdater

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oddmario/tunnel-manager/config"
	"github.com/oddmario/tunnel-manager/logger"
	"github.com/oddmario/tunnel-manager/utils"
)

func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New() // using `.New()` instead of `.Default()` to initialize Gin without the Logger middleware to improve the performance of our app.
	r.Use(gin.Recovery())

	r.UseRawPath = true
	r.UnescapePathValues = false

	initErrors(r)
	initRoutes(r)

	srv := &http.Server{
		Handler:           r,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 30 * time.Second, // https://ieftimov.com/posts/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/#server-timeouts---first-principles
	}
	srv.SetKeepAlivesEnabled(true)

	srv.Addr = config.Config.DynamicIPUpdaterAPIListenAddress + ":" + utils.IToStr(config.Config.DynamicIPUpdaterAPIListenPort)

	err := srv.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start the dynamic IP updater HTTP server. " + err.Error())
	}
}
