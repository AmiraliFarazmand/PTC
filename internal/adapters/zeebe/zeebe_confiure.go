package zeebe

import (
	// "log"

	"log"

	"github.com/camunda-community-hub/zeebe-client-go/v8/pkg/zbc"
)

func NewZeebeClient() zbc.Client {
    client, err := zbc.NewClient(&zbc.ClientConfig{
        GatewayAddress: "localhost:26500", 
        UsePlaintextConnection: true,
    })
    if err != nil {
        // log.Fatalf("###Failed to create Zeebe client: %v", err)
    }
    log.Printf("###Zeebe client created successfully\n")
    return client
}

func MustCloseClient(client zbc.Client) {
	// log.Println("###closing client")
	_ = client.Close()
}
