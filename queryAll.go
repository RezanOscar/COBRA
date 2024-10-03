/////////////////////////////////////////////////////////////////////////////////////////////////
// 
// Objet : GO Script to show data in the ledger for tasks and devices
// 
// version : 4.1
//
// Author : Rêzan OSCAR
// Infos :
//      - Shows the ledger data; can filter by task or device with optional parameters.
//          - Task => DeviceID - TaskType - Status
//          - Device => DeviceID - Status - DeviceType - BatteryLife  
//      ex : ./QueryAll task URLLC   ./QueryAll device EC    
//
/////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
)

// Task structure as per your smart contract
type Task struct {
    TaskID     string  `json:"taskID"`
    DeviceID   string  `json:"deviceID"`
    TaskData   string  `json:"taskData"`
    TaskType   string  `json:"taskType"`
    EnergyCost float64 `json:"energyCost"`
    ComputeCost float64 `json:"computeCost"`
    Status     string  `json:"status"`
}

// Device structure as per your smart contract
type Device struct {
    DeviceID          string  `json:"deviceID"`
    DeviceType        string  `json:"deviceType"`
    Status            string  `json:"status"`
    BatteryLife       float64 `json:"batteryLife"`
    InitialBattery    float64 `json:"initialBattery"`   // Initial battery when the device was registered
    ComputeResources  float64 `json:"computeResources"`
    InitialResources  float64 `json:"initialResources"`
    TasksCompleted    int     `json:"tasksCompleted"`
    TotalTasks        int     `json:"totalTasks"`
    TimeTasks         int     `json:"timeTasks"`        // Tasks completed within the time threshold
    ComputeCostDevice float64 `json:"computeCostDevice"`
    TaskLimit         int     `json:"taskLimit"`
    Reputation        float64 `json:"reputation"`
    PreviousReputation        float64 `json:"previousreputation"`  
}

// Function to initialize the SDK and channel client
func initSDKAndClient(configPath, channelID, user, org string) (*fabsdk.FabricSDK, *channel.Client, error) {
    sdk, err := fabsdk.New(config.FromFile(configPath))
    if err != nil {
        return nil, nil, fmt.Errorf("Failed to create SDK: %w", err)
    }
    channelClient, err := channel.New(sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org)))
    if err != nil {
        return nil, nil, fmt.Errorf("Failed to create channel client: %w", err)
    }
    return sdk, channelClient, nil
}

// Function to query tasks from the ledger with an optional filter
func queryTasks(channelClient *channel.Client, filter string) {
    response, err := channelClient.Query(channel.Request{ChaincodeID: "cobra_algo", Fcn: "QueryAllTasks"})
    if err != nil {
        log.Fatalf("Failed to query tasks: %s", err)
    }
    var tasks []Task
    json.Unmarshal(response.Payload, &tasks)
    for _, task := range tasks {
        if filter == "" || task.DeviceID == filter || task.TaskType == filter || task.Status == filter {
            fmt.Printf("TaskID: %s, DeviceID: %s, TaskData: %s, TaskType: %s, EnergyCost: %.2f, ComputeCost: %.2f, Status: %s\n",
                task.TaskID, task.DeviceID, task.TaskData, task.TaskType, task.EnergyCost, task.ComputeCost, task.Status)
        }
    }
}

// Function to query devices from the ledger with an optional filter
func queryDevices(channelClient *channel.Client, filter string) {
    response, err := channelClient.Query(channel.Request{ChaincodeID: "cobra_algo", Fcn: "QueryAllDevices"})
    if err != nil {
        log.Fatalf("Failed to query devices: %s", err)
    }
    var devices []Device
    json.Unmarshal(response.Payload, &devices)
    for _, device := range devices {
        if filter == "" || device.DeviceID == filter || device.DeviceType == filter || device.Status == filter || 
           (filter == "battery" && device.BatteryLife > 11.0) {
            fmt.Printf("DeviceID: %s, Type: %s, Status: %s, Battery: %.2f, Init Battery: %.2f, ComputeResources: %.2f, TaskCompleted: %d, TotalTask: %d, TimeTask: %d, ComputeCost: %.2f, TaskLimit: %d, Reputation: %.2f, PreviousReputation: %.2f\n",
                device.DeviceID, device.DeviceType, device.Status, device.BatteryLife, device.InitialBattery, device.ComputeResources, device.TasksCompleted, device.TotalTasks, device.TimeTasks, device.ComputeCostDevice, device.TaskLimit, device.Reputation, device.PreviousReputation)
        }
    }
}

func main() {
    if len(os.Args) < 2 || len(os.Args) > 3 {
        log.Fatalf("Usage: ./query <task|device> [optional_filter] <DeviceID|TaskType|DeviceType|Status|battery>")
    }

    sdk, channelClient, err := initSDKAndClient("cobra-config.yaml", "channelcoop", "Admin", "Provider1MSP")
    if err != nil {
        log.Fatalf("Failed to initialize: %s", err)
    }
    defer sdk.Close()

    queryType := os.Args[1]
    var filter string
    if len(os.Args) == 3 {
        filter = os.Args[2]
    }

    switch queryType {
    case "task":
        queryTasks(channelClient, filter)
    case "device":
        queryDevices(channelClient, filter)
    default:
        log.Fatalf("Use 'task' (DeviceID - TaskType - Status) or 'device' (DeviceID - DeviceType - Status - battery)")
    }
}

