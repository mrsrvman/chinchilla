package queue

import (
	"fmt"
	"testing"

	"github.com/benschw/chinchilla/config"
	"github.com/benschw/chinchilla/example/ex"
	"github.com/stretchr/testify/assert"
)

func TestTopicConsume(t *testing.T) {
	// given
	epCfg := config.EndpointConfig{
		QueueConfig: map[interface{}]interface{}{
			"queuename":    "foos",
			"topicname":    "foo.update",
			"exchangename": "demo",
		},
	}

	publisher := &ex.Publisher{
		Conn:   conn,
		Config: &epCfg,
	}

	topic := &Topic{}

	ch, _ := conn.Channel()
	defer ch.Close()

	// when
	msgs, err := topic.Consume(ch, epCfg)

	for i := 0; i < 10; i++ {
		publisher.PublishTopic(fmt.Sprintf("test topic: #%d", i), "text/plain")
	}

	// then
	assert.Nil(t, err)
	cnt := countMessages(msgs)

	assert.Equal(t, 10, cnt, "wrong number of msgs")
}

func TestTopicConsumeGlob(t *testing.T) {
	// given
	epCfg := config.EndpointConfig{
		QueueConfig: map[interface{}]interface{}{
			"queuename":    "foos",
			"topicname":    "foo.update",
			"exchangename": "demo",
		},
	}

	publisher := &ex.Publisher{
		Conn:   conn,
		Config: &epCfg,
	}

	topic := &Topic{}

	ch, _ := conn.Channel()
	defer ch.Close()

	// when
	epCfg = config.EndpointConfig{
		QueueConfig: map[interface{}]interface{}{
			"queuename":    "foos",
			"topicname":    "foo.*",
			"exchangename": "demo",
		},
	}
	msgs, err := topic.Consume(ch, epCfg)

	for i := 0; i < 10; i++ {
		publisher.PublishTopic(fmt.Sprintf("test topic: #%d", i), "text/plain")
	}

	// then
	assert.Nil(t, err)
	cnt := countMessages(msgs)

	assert.Equal(t, 10, cnt, "wrong number of msgs")
}
