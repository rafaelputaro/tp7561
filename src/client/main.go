package main

import (
	"os"
	"tp/client/helpers"
	"tp/common"
	"tp/common/communication/url"
	"tp/common/files_common"
	filetransfer "tp/common/files_common/file_transfer"
	"tp/common/keys"
	rpc_ops_common "tp/common/rpc_ops"

	"github.com/prometheus/client_golang/prometheus"
)

const MESSAGE_START = "Starting client..."

type metrics struct {
	addFile prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		addFile: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "client",
			Name:      "add_file",
			Help:      "Number of add file requests",
		}),
	}
	reg.MustRegister(m.addFile)
	return m
}

func main() {
	/*	reg := prometheus.NewRegistry()
		m := NewMetrics(reg)
		m.addFile.Set(float64(1))
		promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
		http.Handle("/metrics", promHandler)
	*/

	// iniciar cliente
	common.Log.Info(MESSAGE_START)
	common.InitLogger()
	config := helpers.LoadConfig()
	helpers.InitStore(*config)
	// esperar a que la mayoría de los pares se inicialicen intercambiando contactos
	common.SleepOnStart(config.NumberOfPairs)
	// agregar archivos en peer-1
	keysAdded := [][]byte{}
	// to check
	check := map[string]bool{}
	urlPeer := url.GenerateURLPeer(1)
	files_common.OpOverDir(helpers.GenerateInputFilePath(*config, ""),
		func(fileName string) error {
			key, err := rpc_ops_common.AddFile(urlPeer, fileName, helpers.GenerateInputFilePath(*config, fileName))
			if err == nil {
				keyS := keys.KeyToLogFormatString(key)
				common.Log.Debugf("File added: %v | %v | %v", fileName, keyS, urlPeer)
				keysAdded = append(keysAdded, key)
				check[fileName] = false
			}
			common.SleepShort(config.NumberOfPairs)
			return err
		})
	// escuchar llegada de archivos
	filetransfer.NewReceiver(
		config.Url,
		func(fileName string) string {
			return helpers.GenerateDownloadPath(*config, fileName)
		},
		func([]byte, string) {},
	)
	common.SleepOnStart(config.NumberOfPairs)
	urlPeer = url.GenerateURLPeer(config.NumberOfPairs - 1)
	// solicitar archivos a el último peer
	for _, key := range keysAdded {
		_, errGet := rpc_ops_common.GetFile(config.Url, urlPeer, key)
		if errGet != nil {
			common.Log.Debugf("Error get file %v", errGet)
		}
		common.SleepShort(config.NumberOfPairs)
	}
	// chequear si llegaron todos los archivos
	for range 100 {
		for file := range check {
			_, err := os.Stat(helpers.GenerateDownloadPath(*config, file))
			if err == nil {
				delete(check, file)
			}
		}
		if len(check) == 0 {
			common.Log.Infof("Pass")
			return
		}
		common.SleepBetweenRetries()
	}
	common.Log.Infof("Fail")
}
