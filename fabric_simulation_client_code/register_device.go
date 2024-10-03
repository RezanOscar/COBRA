/////////////////////////////////////////////////////////////////////////////////////////////////
//
// Objet : GO that Registers 100 devices
//
// version : 6.1
//
// Author : Rêzan OSCAR
// Infos :
//      - This script registers 100 devices with 10% Edge Servers (EC) and 90% UAVs
//
/////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
    "fmt"
    "log"
    "math/rand"
    "sync"
    "time"

    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
)

// Function to initialize the SDK and channel client
func initSDKAndClient(configPath, channelID, user, org string) (*fabsdk.FabricSDK, *channel.Client, error) {
    sdk, err := fabsdk.New(config.FromFile(configPath))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create SDK: %w", err)
    }

    channelClient, err := channel.New(sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org)))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create channel client: %w", err)
    }

    return sdk, channelClient, nil
}

func main() {
    startTime := time.Now()

    // Initialize the SDK and channel client
    sdk, channelClient, err := initSDKAndClient("cobra-config.yaml", "channelcoop", "Admin", "Provider1MSP")
    if err != nil {
        log.Fatalf("Initialization failed: %s", err)
    }
    defer sdk.Close()

    rand.Seed(time.Now().UnixNano())
    var wg sync.WaitGroup
    sem := make(chan struct{}, 10) // Limit to 10 concurrent goroutines

    totalDevices := 30
    totalEC := int(float64(totalDevices) * 0.1) // 10% EC
    totalUAV := totalDevices - totalEC           // 90% UAV

    deviceIDSet := make(map[string]bool) // Tracking generated IDs for duplicates
    var mu sync.Mutex

    // Function to generate a unique device ID
    generateDeviceID := func() string {
        for {
            id := fmt.Sprintf("%04d", rand.Intn(100)) // DeviceID random between 001 and 099
            mu.Lock()
            if !deviceIDSet[id] {
                deviceIDSet[id] = true
                mu.Unlock()
                return id
            }
            mu.Unlock()
        }
    }

    // Function to register a device with the ledger
    registerDevice := func(deviceID, deviceType string, batteryLife, initialBattery float64, computeResources float64, initialResources float64, tasksCompleted int, totalTasks int, timeTasks int, computeCostDevice float64, taskLimit int, reputation float64, previousreputation float64) {
        defer wg.Done()
        sem <- struct{}{} // Slot

        _, err := channelClient.Execute(channel.Request{
            ChaincodeID: "cobra_algo", // Update the Chaincode ID to match the latest version
            Fcn:         "RegisterDevice",
            Args: [][]byte{
                []byte(deviceID),
                []byte(deviceType),
                []byte("Available"),
                []byte(fmt.Sprintf("%.2f", batteryLife)),
                []byte(fmt.Sprintf("%.2f", initialBattery)),
                []byte(fmt.Sprintf("%.2f", computeResources)),
                []byte(fmt.Sprintf("%.2f", initialResources)),
                []byte(fmt.Sprintf("%d", tasksCompleted)), 
                []byte(fmt.Sprintf("%d", totalTasks)),    
                []byte(fmt.Sprintf("%d", timeTasks)),      
                []byte(fmt.Sprintf("%.2f", computeCostDevice)),
                []byte(fmt.Sprintf("%d", taskLimit)),      
                []byte(fmt.Sprintf("%.2f", reputation)),
                []byte(fmt.Sprintf("%.2f", previousreputation)),   
            },
        })

        if err != nil {
            fmt.Printf("Failed to register %s: %s\n", deviceID, err)
        } else {
            fmt.Printf("Registered device %s type %s with battery life %.2f and compute resources %.2f\n", deviceID, deviceType, batteryLife, computeResources)
        }

        <-sem // Release slot
    }


    // Register EC devices
    for i := 0; i < totalEC; i++ {
        wg.Add(1)
        go registerDevice(generateDeviceID(), "EC", 50.0, 50.0, 100.0, 100, 0, 0, 0, 0, 0, 1.0, 0.0)
    }

    // Register UAV devices
    for i := 0; i < totalUAV; i++ {
        wg.Add(1)
        go registerDevice(generateDeviceID(), "UAV", 50.0, 50.0, 10.0, 10, 0, 0, 0, 0, 0, 1.0, 0.0)
    }

    wg.Wait()

    endTime := time.Now()
    duration := endTime.Sub(startTime)
    fmt.Printf("Total time to register devices: %s\n", duration)
}

