package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/containous/traefik/log"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/prometheus/prompb"
)

func main() {
	log.Info("creating pulsar producer")
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/receive", receiveHandler())
	r.Run()
}

func receiveHandler() func(c *gin.Context) {
	return func(c *gin.Context) {

		compressed, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logrus.WithError(err).Error("couldn't read body")
			return
		}

		reqBuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			logrus.WithError(err).Error("couldn't decompress body")
			return
		}

		var req prompb.WriteRequest
		if err := proto.Unmarshal(reqBuf, &req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			logrus.WithError(err).Error("couldn't unmarshal body")
			return
		}

		metrics, err := processWriteRequest(&req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			logrus.WithError(err).Error("couldn't process write request")
			return
		}

		producer, err := client.CreateProducer(pulsar.ProducerOptions{
			Topic: pulsarTopic,
		})

		defer producer.Close()

		for _, metric := range metrics {
			_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
				Payload: metric,
			})

			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				logrus.WithError(err).Error("couldn't produce message in pulsar")
				return
			}
		}

	}
}

func removeSpecialCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}

	trimmedInput := strings.TrimLeft(input, "_")
	return strings.Map(filter, trimmedInput)
}
