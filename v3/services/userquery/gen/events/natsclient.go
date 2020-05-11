// Code generated by goa v3.0.9, DO NOT EDIT.
//
// NATS client
//
// Command:
// $ goa gen gitlab.com/wiserskills/v3/services/userquery/design

package events

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env"
	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
	log "gitlab.com/wiserskills/v3/services/userquery/gen/log"
)

// DefaultWriteTimeout defines the default timeout for write operations
const DefaultWriteTimeout = time.Second * 5

// DefaultRetryCount defines the default number of connection retry
const DefaultRetryCount = 5

// DefaultRetryDelay defines the delay in seconds between two connection retries
const DefaultRetryDelay = 10

// BusConfig defines the variables used to connect to the event bus
type BusConfig struct {

	// EventBusCluster defines the name of the eventbus cluster
	EventBusCluster string `env:"EVENTBUS_CLUSTER" json:"eventbus_cluster,omitempty"`

	// EventBusURL defines the url of the eventbus
	EventBusURL string `env:"EVENTBUS_URL" json:"eventbus_url,omitempty"`

	// Topic defines the default topic
	EventBusDefaultTopic string `env:"EVENTBUS_DEFAULT_TOPIC" json:"eventbus_default_topic,omitempty"`

	// WriteTimeout defines the timeout to write to the event bus in microseconds
	// Default: 5 seconds
	EventBusWriteTimeout time.Duration `env:"EVENTBUS_WRITE_TIMEOUT" json:"eventbus_write_timeout,omitempty"`

	// EventBusRetryCount defines the number of retry to connect to the event bus in case of connection failure
	// Default: 5
	EventBusRetryCount int `env:"EVENTBUS_RETRY_COUNT" json:"eventbus_retry_count,omitempty"`

	// EventBusRetryDelay defines the delay in seconds before retrying to connect to the event bus in case of failure
	// Default: 10
	EventBusRetryDelay int `env:"EVENTBUS_RETRY_DELAY" json:"eventbus_retry_delay,omitempty"`

	// EventBusClientCertPath defines the path to the client certificate file in PEM format used to connect to NATS
	EventBusClientCertPath string `env:"EVENTBUS_CLIENT_CERT_PATH" json:"eventbus_client_cert_path,omitempty"`

	// EventBusClientKeyPath defines the path to the client private key file in PEM format used to connect to NATS
	EventBusClientKeyPath string `env:"EVENTBUS_CLIENT_KEY_PATH" json:"eventbus_client_key_path,omitempty"`

	// Subscriptions represents the topics that are subscribed by the service
	Subscriptions []*TopicSubscription
}

// Valid defines if the Event Bus configuration is valid
func (c *BusConfig) Valid() bool {
	if c.EventBusCluster == "" || c.EventBusURL == "" {
		return false
	}

	return true
}

// NATSClient represents a NATS client
type NATSClient struct {
	logger        *log.Logger
	client        stan.Conn
	id            string
	cfg           BusConfig
	subscriptions []*stan.Subscription
}

// NewNATSClient creates a new NATS client with the specified parameters
func NewNATSClient(logger *log.Logger, cfg *BusConfig) (EventbusClient, error) {

	clientid, err := newUUID()

	if err != nil {
		return nil, err
	}

	if !cfg.Valid() {
		logger.Warn("Event Bus configuration is not valid. Service will not handle events.")
	}

	result := &NATSClient{logger: logger, id: clientid, cfg: *cfg, subscriptions: make([]*stan.Subscription, 0)}

	return result, nil
}

// GetBusConfig returns the event bus configuration
func GetBusConfig() (*BusConfig, error) {
	result := BusConfig{}
	serviceConfigPath := getEnv("USERQUERY_CONFIG_PATH", "userquery.cfg")
	_, err := parseFromFile(&result, serviceConfigPath)

	if err != nil {
		return nil, err
	}

	err = env.Parse(&result)

	if err != nil {
		return nil, err
	}

	if result.EventBusRetryCount == 0 {
		result.EventBusRetryCount = DefaultRetryCount
	}

	if result.EventBusRetryDelay == 0 {
		result.EventBusRetryDelay = DefaultRetryDelay
	}

	if result.EventBusWriteTimeout == 0 {
		result.EventBusWriteTimeout = DefaultWriteTimeout
	}

	return &result, nil
}

// SetDefaultTopic sets the default topic
func (c *NATSClient) SetDefaultTopic(topic string) {
	c.cfg.EventBusDefaultTopic = topic
}

// SetWriteTimeout sets the default topic
func (c *NATSClient) SetWriteTimeout(timeout time.Duration) {
	c.cfg.EventBusWriteTimeout = timeout
}

// SetRetryCount sets the retry count value
func (c *NATSClient) SetRetryCount(count int) {
	c.cfg.EventBusRetryCount = count
}

// SetRetryDelay sets the retry delay value
func (c *NATSClient) SetRetryDelay(delay int) {
	c.cfg.EventBusRetryDelay = delay
}

// Connect opens a connection with a remote event bus
func (c *NATSClient) Connect() error {

	if c.client == nil && c.cfg.Valid() {

		var nc *nats.Conn
		var err error

		if c.cfg.EventBusClientCertPath == "" || c.cfg.EventBusClientKeyPath == "" {
			c.logger.Warn("no client certificate or key provided to connect to the event bus.")
		} else {
			nc, err = nats.Connect(c.cfg.EventBusURL, nats.ClientCert(c.cfg.EventBusClientCertPath, c.cfg.EventBusClientKeyPath))

			if err != nil {
				return err
			}
		}

		var client stan.Conn

		if nc != nil {
			client, err = stan.Connect(c.cfg.EventBusCluster, c.id, stan.NatsConn(nc))
		} else {
			client, err = stan.Connect(c.cfg.EventBusCluster, c.id, stan.NatsURL(c.cfg.EventBusURL))
		}

		if err != nil {
			return err
		}

		c.client = client
		c.loadSubscriptions()
	}

	return nil
}

// Close closes the NATS client connection
func (c *NATSClient) Close() error {

	if c.client != nil {
		err := c.client.Close()

		if err != nil {
			return err
		}

		c.client = nil

		return nil
	}

	return nil
}

// Write sends the passed byte array to the default NATS topic asynchronously (required by zap)
func (c *NATSClient) Write(p []byte) (n int, err error) {

	if c.client != nil && p != nil {

		ch := make(chan bool)
		var glock sync.Mutex
		var guid string

		ack := func(lguid string, err error) {

			glock.Lock()
			c.logger.Debugf("received ACK for guid: %s\n", lguid)
			defer glock.Unlock()
			if err != nil {
				c.logger.Errorf("error in server ack for guid %s: %v\n", lguid, err)
			}
			if lguid != guid {
				c.logger.Errorf("expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
			}

			ch <- true
		}

		glock.Lock()
		guid, err = c.client.PublishAsync(c.cfg.EventBusDefaultTopic, p, ack)
		if err != nil {
			c.logger.Errorf("error during async publish: %v\n", err)
		}
		glock.Unlock()
		if guid == "" {
			c.logger.Errorf("expected non-empty guid to be returned.")
		}

		select {
		case <-ch:
			break
		case <-time.After(c.cfg.EventBusWriteTimeout):
			c.logger.Errorf("timeout while sending event to event bus")
		}

		if err == nil {
			n = len(p)
		}
	}

	return
}

// WriteWithTopic sends the passed byte array to a specified NATS topic asynchronously
func (c *NATSClient) WriteWithTopic(topic string, p []byte, async bool, success EventSuccessHandler, failure EventErrorHandler) (err error) {

	if c.client != nil && p != nil {

		if async == true {

			ch := make(chan bool)
			var glock sync.Mutex
			var guid string

			ack := func(lguid string, err error) {

				glock.Lock()
				defer glock.Unlock()
				if err != nil {
					c.logger.Errorf("error in server ack for guid %s: %v\n", lguid, err)
				}
				if lguid != guid {
					c.logger.Errorf("expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
				}
				ch <- true
			}

			glock.Lock()
			guid, err = c.client.PublishAsync(topic, p, ack)
			if err != nil {
				c.logger.Errorf("error during async publish: %v\n", err)
			}
			glock.Unlock()
			if guid == "" {
				c.logger.Errorf("expected non-empty guid to be returned.")
			}

			select {
			case <-ch:
				success(guid)
				break
			case <-time.After(c.cfg.EventBusWriteTimeout):
				err := errors.New("timeout while sending event to event bus")
				failure(guid, err)
				c.logger.Errorf("timeout while sending event to event bus")
			}

		} else {

			err = c.client.Publish(topic, p)
			if err != nil {
				c.logger.Errorf("error during publish: %v\n", err)
			}
		}
	}

	return
}

// Sync is required to be compatible with zap
func (c *NATSClient) Sync() error {
	return nil
}

// Reload reloads the registered config/subscriptions
func (c *NATSClient) Reload(subscriptions []*TopicSubscription) error {

	if c.client != nil {

		for _, passedSub := range subscriptions {

			found := c.findTopicSubscription(passedSub.Topic)

			if found != nil {

				for _, eventSubs := range passedSub.Subscriptions {
					for _, eventSub := range eventSubs {
						if found.Subscriptions[eventSub.Name] == nil {
							found.Subscriptions[eventSub.Name][eventSub.ID] = eventSub
						}
					}
				}

			} else {
				c.cfg.Subscriptions = append(c.cfg.Subscriptions, passedSub)
			}
		}

		return c.loadSubscriptions()
	}

	return nil
}

func (c *NATSClient) loadSubscriptions() error {

	for _, sub := range c.cfg.Subscriptions {

		if sub.Registered {
			continue
		}

		h := func(msg *stan.Msg) {
			err := sub.Handler(msg.Subject, msg.Data)

			if err == nil {
				msg.Ack()
			}
		}

		var s stan.Subscription
		var err error

		options := getSubOptions(sub)

		if sub.Options != nil && sub.Options.GroupName != "" {
			s, err = c.client.QueueSubscribe(sub.Topic, sub.Options.GroupName, h, options...)
		} else {
			s, err = c.client.Subscribe(sub.Topic, h, options...)
		}

		if err != nil {
			return err
		}

		c.subscriptions = append(c.subscriptions, &s)
		sub.Registered = true
	}

	return nil
}

func (c *NATSClient) findTopicSubscription(topic string) *TopicSubscription {

	for _, sub := range c.cfg.Subscriptions {
		if sub.Topic == topic {
			return sub
		}
	}

	return nil
}

func parseFromFile(config *BusConfig, path string) (*BusConfig, error) {

	if _, err := os.Stat(path); !os.IsNotExist(err) {

		file, err := os.Open(path)

		if err != nil {
			return config, err
		}

		decoder := json.NewDecoder(file)
		err = decoder.Decode(config)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}

func getSubOptions(sub *TopicSubscription) []stan.SubscriptionOption {

	options := []stan.SubscriptionOption{}
	options = append(options, stan.MaxInflight(1))
	options = append(options, stan.SetManualAckMode())

	if sub.Options == nil {
		return options
	}

	if sub.Options.Durable {
		options = append(options, stan.DurableName(sub.Topic))
	}

	if sub.Options.Replay == 0 {
		options = append(options, stan.DeliverAllAvailable())
	}

	if sub.Options.Replay == 1 {
		options = append(options, stan.StartWithLastReceived())
	}

	if sub.Options.Replay == 2 {
		options = append(options, stan.StartAtSequence(sub.Options.ReplayParameter.(uint64)))
	}

	if sub.Options.Replay == 3 {
		options = append(options, stan.StartAtTime(sub.Options.ReplayParameter.(time.Time)))
	}

	if sub.Options.Replay == 4 {
		options = append(options, stan.StartAtTimeDelta(sub.Options.ReplayParameter.(time.Duration)))
	}

	return options
}
