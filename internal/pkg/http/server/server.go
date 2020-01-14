package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func ListenAndServe(conf Config, router http.Handler) {
	port := fmt.Sprint(conf.Port)
	if conf.Port == 0 {
		port = os.Getenv("PORT")
		if port == "" {
			port = "80"
		}
	}
	address := fmt.Sprintf("%v:%v", conf.Address, conf.Port)

	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  conf.ReadTimeOut,
		WriteTimeout: conf.WriteTimeOut,
	}
	logrus.Infof("HTTP SERVER is listening on: %s", address)

	if err := srv.ListenAndServe(); err != nil {
		log.Panicf("listen: %s\n", err)
	}
}
