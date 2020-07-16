package mqttclient

import (
	"errors"
	"fmt"
	"time"

	"github.com/austinjan/AtopIOTServer/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttClientOption  settings of MqttClient
type MqttClientOption struct {
	ClientID  string
	Port      int32
	Reconnect bool
	Host      string
	KeepAlive int32
	Key       string
	UserName  string
	Password  string
}

// MqttClient interface defined for a client of mqtt. Provieds methods to
// connect broker, subscribe , publish etc..
type MqttClient interface {
	Connect() error
	Subscribe(topics map[string]byte) error
}

// MqttClient mqtt client type
type mqttClient struct {
	client  mqtt.Client
	connect bool
}

func defaultMqttClient() *mqttClient {
	return &mqttClient{client: nil, connect: false}
}

// func NewClient(opt MqttClientOption) MqttClient {
// 	c := defaultMqttClient()
// 	ops := mqtt.NewClientOptions()
// 	ops.SetClientID(opt.ClientID)
// 	ops.SetAutoReconnect(opt.Reconnect)
// 	ops.Set
// }

// NewLocalClient create mqtt client which connected to localhost broker
func NewLocalClient() MqttClient {
	c := defaultMqttClient()
	ops := mqtt.NewClientOptions()
	c.setDefaultMqttClientOptions(ops)
	clientID := utils.GetUniqueID()
	fmt.Println("start mqtt local client with clientID: ", clientID)
	ops.SetClientID(clientID)
	ops.AddBroker("127.0.0.1:1883")

	c.client = mqtt.NewClient(ops)

	return c
}

// Connect connect to host
func (c *mqttClient) Connect() error {
	connToken := c.client.Connect()
	// waitting connect
	for !connToken.WaitTimeout(time.Second * 3) {
		return connToken.Error()
	}

	return connToken.Error()
}

// Subscribe  subsrib topics,  ie: c.Subscribe(map[string]byte{"test/topic":0,"test/qos2/topic":2})
func (c *mqttClient) Subscribe(topics map[string]byte) error {
	token := c.client.SubscribeMultiple(topics, handleSubscribe)
	for !token.WaitTimeout(time.Second * 3) {
		return token.Error()
	}

	return token.Error()
}

func (c *mqttClient) handleConnect(cli mqtt.Client) {
	fmt.Println("mqtt client connected")
}

func (c *mqttClient) handleDisconnect(cli mqtt.Client, e error) {
	fmt.Println("mqtt client disconnected with err: ", e.Error())
}

// process received mqtt message
func handleSubscribe(c mqtt.Client, msg mqtt.Message) {
	fmt.Printf("*mqttserver [%s] %s\n", msg.Topic(), string(msg.Payload()))

}

func (c *mqttClient) setDefaultMqttClientOptions(ops *mqtt.ClientOptions) *mqtt.ClientOptions {
	ops.SetOnConnectHandler(c.handleConnect)
	ops.SetConnectionLostHandler(c.handleDisconnect)
	return ops
}

// Connect connet mqtt server
func connect(client mqtt.Client) error {

	if client == nil {
		return errors.New("mqtt connect error: null client")
	}

	if client.IsConnected() {
		return nil
	}

	connToken := client.Connect()
	// waitting connect
	for !connToken.WaitTimeout(time.Second * 3) {
	}

	return connToken.Error()
}

func ConnectLocal() {
	ops := mqtt.NewClientOptions()
	clientID := utils.GetUniqueID()
	fmt.Println("mqtt client clientID ", clientID)
	ops.SetClientID(clientID)
	ops.AddBroker("127.0.0.1:1883")

	client := mqtt.NewClient(ops)

	if client == nil {
		fmt.Print("ops is nil")
	}

	if err := connect(client); err != nil {
		fmt.Println("Mqtt broker connect fail.")
	}

	client.Subscribe("test/+", byte(0), handleSubscribe)
}
