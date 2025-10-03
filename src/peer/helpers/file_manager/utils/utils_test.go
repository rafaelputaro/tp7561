package utils

import (
	"sync"
	"testing"
	"time"
	"tp/peer/helpers/file_manager/config_fm"
)

func TestKeys(t *testing.T) {
	ret := PathExists("/tmp")
	print("Retorno:", ret)
}

func TestFolderConfig(t *testing.T) {
	inputDataFolder := "/tmp/data"
	storeIpfsFolder := "/tmp/peer/store"
	config_fm.LocalStorageConfig = *config_fm.NewStorageConfig(inputDataFolder, storeIpfsFolder)
	config_fm.LocalStorageConfig.LogConfig()
	t.Logf("Store: %v", GenerateIpfsStorePath("file1"))
	t.Logf("Download: %v", GenerateIpfsDownloadPath("file1"))
	t.Logf("Download Part: %v", GenertaIpfsDownloadPartPath("file1", 1))
	t.Logf("Restore: %v", GenerateIpfsRestorePath("file1"))
}

func TestGoFunc(t *testing.T) {
	max := 10
	wg := new(sync.WaitGroup)
	for n := range max {
		wg.Add(1)
		go func() {
			for range 10 {
				println("Echo ", n)
				t := time.Duration(1) * time.Second
				time.Sleep(t)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
