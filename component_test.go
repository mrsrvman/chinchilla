package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"testing"

	"github.com/benschw/chinchilla/config"
	"github.com/benschw/chinchilla/ep"
	"github.com/benschw/chinchilla/ep/queue"
	"github.com/benschw/chinchilla/example/ex"
	"github.com/benschw/opin-go/rando"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

var conn *amqp.Connection

type Msg struct {
	Message string
}

func init() {
	c, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	conn = c
}
func GetPublisher(cfg *config.EndpointConfig) *ex.Publisher {
	p := &ex.Publisher{
		Conn:   conn,
		Config: cfg,
	}
	return p
}

func GetServices() (*ep.EndpointApp, *ex.Server, *ex.Publisher, *ex.Publisher) {
	port := uint16(rando.Port())

	server := ex.NewServer(fmt.Sprintf(":%d", port))

	epCfg := config.EndpointConfig{
		Name:        "Foo",
		QueueName:   "test.foo",
		ServiceHost: fmt.Sprintf("http://localhost:%d", port),
		Uri:         "/foo",
		Method:      "POST",
	}
	epCfg2 := config.EndpointConfig{
		Name:        "Bar",
		QueueName:   "test.bar",
		ServiceHost: fmt.Sprintf("http://localhost:%d", port),
		Uri:         "/bar",
		Method:      "POST",
	}

	p := GetPublisher(&epCfg)
	p2 := GetPublisher(&epCfg2)

	repo := &config.StaticRepo{
		Address: config.RabbitAddress{
			User:     "guest",
			Password: "guest",
			Host:     "localhost",
			Port:     5672,
		},
		Endpoints: []config.EndpointConfig{epCfg, epCfg2},
	}

	qReg := ep.NewQueueRegistry()
	qReg.Add(qReg.DefaultKey, &queue.Queue{C: &queue.DefaultWorker{}, D: &queue.DefaultDeliverer{}})

	mgr := ep.NewApp(repo, repo, qReg)
	return mgr, server, p, p2
}

func TestPublish(t *testing.T) {
	// setup
	mgr, server, p, _ := GetServices()
	go server.Start()
	go mgr.Run()

	// wait for queue creation to prevent race condition... do this better
	time.Sleep(200 * time.Millisecond)

	api := &Msg{Message: "Hello World"}
	apiB, _ := json.Marshal(api)
	apiStr := string(apiB)

	// when
	err := p.Publish(apiStr, "application/json")
	assert.Nil(t, err)

	time.Sleep(200 * time.Millisecond)

	mgr.Stop()
	server.Stop()

	// then

	statLen := len(server.H.Stats["Foo"])
	assert.Equal(t, 1, statLen, "wrong number of stats")
	if statLen > 0 {
		foundApi := &Msg{}
		err := json.Unmarshal([]byte(server.H.Stats["Foo"][0]), foundApi)
		assert.Nil(t, err, "err should be nil")

		assert.True(t, reflect.DeepEqual(api, foundApi), fmt.Sprintf("\n   %+v\n!= %+v", api, foundApi))
	}
}
func TestPublishLotsAndLots(t *testing.T) {
	// setup
	mgr, server, p, p2 := GetServices()
	go server.Start()
	go mgr.Run()

	body := "Hello World"

	// when
	for i := 0; i < 100; i++ {
		err := p.Publish(body, "text/plain")
		assert.Nil(t, err)

		err = p2.Publish(body, "text/plain")
		assert.Nil(t, err)

	}
	server.Stop()
	mgr.Stop()

	// then
	assert.Equal(t, 100, len(server.H.Stats["Foo"]), "wrong number of stats")
	assert.Equal(t, 100, len(server.H.Stats["Bar"]), "wrong number of stats")
}
