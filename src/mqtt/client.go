package client

import (
    "fmt"

    log "github.com/cloudwego/hertz/pkg/common/hlog"
    mqtt "github.com/eclipse/paho.mqtt.golang"
)

// var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
//     log.Infof("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
// }

var ikunHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    log.Infof("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    log.Infof("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    log.Infof("Connect lost: %v", err)
}

func Run() {
    var broker = "192.168.31.24"
    var port = 1883
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("admin")
    opts.SetPassword("public")
    // opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    token := client.Subscribe("ikun", 1, ikunHandler)
    token.Wait()
    log.Infof("sub end")
}
