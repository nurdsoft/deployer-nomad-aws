package nomad

import (
	"log"
	"sync"

	"github.com/hashicorp/nomad/api"
)

type instance struct {
	nomadClient *api.Client
}

var singleton = &instance{}
var once sync.Once

func init() {
	once.Do(func() {
		log.Println("initializing nomad client ...")
		nomadClient, err := api.NewClient(api.DefaultConfig())
		if err != nil || nomadClient == nil {
			panic(err)
		}
		singleton.nomadClient = nomadClient
		log.Println("nomad client initialized")
	})
}

func GetInstance() *api.Client {
	return singleton.nomadClient
}

func Destroy() error {
	singleton.nomadClient.Close()
	return nil
}
