package g

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	f, err := os.OpenFile(setting.Server.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Open log file failed: ", err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	if gin.Mode() == gin.DebugMode {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(f)
	}
}
