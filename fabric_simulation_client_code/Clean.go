/////////////////////////////////////////////////////////////////////////////////////////////////
// 
// Objet : GO Script to clean the ledger data
// 
// version : 2
//
// Author : Rêzan OSCAR
// Infos :
//      - Clean the ledger data use parameter task or device or all
//
/////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
    "fmt"
    "log"
    "os"

    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
)

func main() {
    if len(os.Args) != 2 {
        log.Fatalf("Usage: ./delete <tasks|devices|all>")
    }

    deleteType := os.Args[1]
    if deleteType != "tasks" && deleteType != "devices" && deleteType != "all" {
        log.Fatalf("Invalid argument: %s. Must be 'tasks', 'devices', or 'all'", deleteType)
    }

    // Init SDK + Channel 
    sdk, err := fabsdk.New(config.FromFile("cobra-config.yaml"))
    if err != nil {
        log.Fatalf("Failed to create SDK: %s", err)
    }
    defer sdk.Close()

    channelClient, err := channel.New(sdk.ChannelContext("channelcoop", fabsdk.WithUser("Admin"), fabsdk.WithOrg("Provider1MSP")))
    if err != nil {
        log.Fatalf("Failed to create new channel client: %s", err)
    }

    // Invoke delete function if error see chaincode ID
    _, err = channelClient.Execute(channel.Request{
        ChaincodeID: "cobra_algo",
        Fcn:         "DeleteAll",
        Args:        [][]byte{[]byte(deleteType)},
    })

    if err != nil {
        log.Fatalf("Failed to delete %s: %s", deleteType, err)
    } else {
        fmt.Printf("Successfully deleted all %s.\n", deleteType)
    }
}
