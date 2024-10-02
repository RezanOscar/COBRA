/////////////////////////////////////////////////////////////////////////////////////////////////
//
// Objet : Smart Contract COBRA framework
//
// version : 6
//
// Author : Rêzan OSCAR
// Infos :
//      - Implementation of RI (ReliabilityIndex) and TCI (TaskCostIndex) in task offloading
//      - Addition of TaskOffloadFirstRoundRobin, TaskOffloadRandom, TaskOffloadSemiRandom, and TaskOffloadCobra functions
//      - SC completed with all model 
//
/////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "time"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/////////////////////////////////////////////////////////////////////////////////////////////////
// Section 1 : Initial structure with the different basic funtion and structure                // 
/////////////////////////////////////////////////////////////////////////////////////////////////

type SmartContract struct {
    contractapi.Contract
    lastUsedEC     string // To track the last EC assigned a task
    currentUAVTasks int   // To track how many UAV tasks are left in the current UAV phase
}

// Device represents a UAV or Edge Server properties 
type Device struct {
    DeviceID          string  `json:"deviceID"`
    DeviceType        string  `json:"deviceType"`
    Status            string  `json:"status"`
    BatteryLife       float64 `json:"batteryLife"`
    InitialBattery    float64 `json:"initialBattery"`   
    ComputeResources  float64 `json:"computeResources"`
    InitialResources  float64 `json:"initialResources"`
    TasksCompleted    int     `json:"tasksCompleted"`   // Total Task Completed
    TotalTasks        int     `json:"totalTasks"`       // Total Task affect for an Device
    TimeTasks         int     `json:"timeTasks"`        // Total Tasks realize in the time affect for the task
    ComputeCostDevice float64 `json:"computeCostDevice"`// Total Compute Cost of task realize by an device
    TaskLimit         int     `json:"taskLimit"`        // Limit size of task to simulate update 3 for an UAV and 30 for an EC
    Reputation        float64 `json:"reputation"`       // Repuation
    PreviousReputation        float64 `json:"previousreputation"`       // Repuation        
}

// Task represents the task details to be offloaded
type Task struct {
    TaskID       string  `json:"taskID"`
    DeviceID     string  `json:"deviceID"`
    TaskData     string  `json:"taskData"`
    TaskType     string  `json:"taskType"`
    EnergyCost   float64 `json:"energyCost"`
    ComputeCost  float64 `json:"computeCost"`
    Status       string  `json:"status"`
}

// InitLedger initializes the ledger with some sample devices
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    devices := []Device{
        {DeviceID: "0001", DeviceType: "EC",  Status: "Available", BatteryLife: 100, InitialBattery: 100, ComputeResources: 100, InitialResources: 100, TasksCompleted: 0, TotalTasks: 0, TimeTasks: 0, ComputeCostDevice: 0, TaskLimit: 0, Reputation: 0, PreviousReputation: 0},
        {DeviceID: "0002", DeviceType: "UAV", Status: "Available", BatteryLife: 100, InitialBattery: 100, ComputeResources: 10, InitialResources: 10, TasksCompleted: 0, TotalTasks: 0, TimeTasks: 0, ComputeCostDevice: 0, TaskLimit: 0, Reputation: 0, PreviousReputation: 0},
        {DeviceID: "0003", DeviceType: "EC",  Status: "Available", BatteryLife: 100, InitialBattery: 100, ComputeResources: 100, InitialResources: 100, TasksCompleted: 0, TotalTasks: 0, TimeTasks: 0, ComputeCostDevice: 0, TaskLimit: 0, Reputation: 0, PreviousReputation: 0},
        {DeviceID: "0004", DeviceType: "UAV", Status: "Available", BatteryLife: 100, InitialBattery: 100, ComputeResources: 10, InitialResources: 10, TasksCompleted: 0, TotalTasks: 0, TimeTasks: 0, ComputeCostDevice: 0, TaskLimit: 0, Reputation: 0, PreviousReputation: 0},
    }

    for _, device := range devices {
        deviceAsBytes, _ := json.Marshal(device)
        err := ctx.GetStub().PutState("D"+device.DeviceID, deviceAsBytes)
        if err != nil {
            return err
        }
    }

    return nil
}


/////////////////////////////////////////////////////////////////////////////////////////////////
// Section 3 : Task Offload algorithm                                                          // 
/////////////////////////////////////////////////////////////////////////////////////////////////

// TaskOffloadFirstAvailable assigns a task to the first available device (Test function)
func (s *SmartContract) TaskOffloadFirstAvailable(ctx contractapi.TransactionContextInterface, taskData string, taskType string, energyCost float64, computeCost float64) error {
    devices, err := s.getAvailableDevices(ctx, computeCost)
    if err != nil {
        return err
    }

    if len(devices) == 0 {
        return fmt.Errorf("No available devices with sufficient ressources")
    }

    selectedDevice := devices[0] // Select the first available device

    return s.assignTask(ctx, selectedDevice, taskData, taskType, energyCost, computeCost)
}


/////////////////////////////////////////////////////////////////////////////////////////////////
// Task Offload based on RoundRobin                                                            // 
/////////////////////////////////////////////////////////////////////////////////////////////////

// TaskOffloadingRoundRobin assigns tasks to devices in a round-robin fashion
func (s *SmartContract) TaskOffloadingRoundRobin(ctx contractapi.TransactionContextInterface, taskData string, taskType string, energyCost float64, computeCost float64) error {
    // Get all available devices
    devices, err := s.getAvailableDevices(ctx, computeCost)
    if err != nil {
        return err
    }

    if len(devices) == 0 {
        return fmt.Errorf("No available devices with sufficient ressources")
    }

    // Get the total number of devices
    totalDevices := len(devices)

    // Initialize or reset the index if needed
    if s.lastUsedEC == "" {
        s.lastUsedEC = devices[0].DeviceID // Start with the first device
    }

    // Find the index of the last used device
    var lastUsedIndex int
    for i, device := range devices {
        if device.DeviceID == s.lastUsedEC {
            lastUsedIndex = i
            break
        }
    }

    // Select the next device in a round-robin fashion
    nextDeviceIndex := (lastUsedIndex + 1) % totalDevices
    selectedDevice := devices[nextDeviceIndex]

    // Update the last used device
    s.lastUsedEC = selectedDevice.DeviceID

    // Assign the task to the selected device
    return s.assignTask(ctx, selectedDevice, taskData, taskType, energyCost, computeCost)
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Task Offload based on a total Random Choice                                                 // 
/////////////////////////////////////////////////////////////////////////////////////////////////

// TaskOffloadRandom assigns a task to a randomly chosen device
func (s *SmartContract) TaskOffloadRandom(ctx contractapi.TransactionContextInterface, taskData string, taskType string, energyCost float64, computeCost float64) error {
    devices, err := s.getAvailableDevices(ctx, computeCost)
    if err != nil {
        return err
    }

    if len(devices) == 0 {
        return fmt.Errorf("No available devices with sufficient ressources")
    }

    // Randomly select a device
    rand.Seed(time.Now().UnixNano())
    selectedDevice := devices[rand.Intn(len(devices))]

    return s.assignTask(ctx, selectedDevice, taskData, taskType, energyCost, computeCost)
}


/////////////////////////////////////////////////////////////////////////////////////////////////
// Task Offload based on a choice that priotorize Edge Server                                  //
/////////////////////////////////////////////////////////////////////////////////////////////////

// TaskOffloadECP (Edge Server Prioritize) assigns a task first to a random EC, but not consecutively, 
// and if all ECs have completed their tasks, UAVs will handle twice the number of tasks as ECs.
func (s *SmartContract) TaskOffloadECP(ctx contractapi.TransactionContextInterface, taskData string, taskType string, energyCost float64, computeCost float64) error {
    devices, err := s.getAvailableDevices(ctx, computeCost)
    if err != nil {
        return err
    }

    var ecs, uavs []Device
    for _, device := range devices {
        if device.DeviceType == "EC" {
            ecs = append(ecs, device)
        } else if device.DeviceType == "UAV" {
            uavs = append(uavs, device)
        }
    }

    // Total ECs and UAVs
    totalECs := len(ecs)
    totalUAVs := len(uavs)

    // If no ECs or UAVs are available
    if totalECs == 0 && totalUAVs == 0 {
        return fmt.Errorf("No available devices with sufficient ressources")
    }

    // Check if we're in the UAV phase: UAVs 
    if s.currentUAVTasks > 0 {
        // Decrease the UAV task count, meaning we're still in the UAV phase
        selectedUAV := s.selectRandomDevice(uavs)
        s.currentUAVTasks--
        return s.assignTask(ctx, selectedUAV, taskData, taskType, energyCost, computeCost)
    }

    // Ensure no consecutive EC selection and check EC phase logic
    if !s.areAllECsAssigned(ecs) {
        selectedEC := s.selectRandomEC(ctx, ecs)
        return s.assignTask(ctx, selectedEC, taskData, taskType, energyCost, computeCost)
    }

    // If all ECs have completed a task, switch to UAV phase f
    s.currentUAVTasks = totalECs * 3
    selectedUAV := s.selectRandomDevice(uavs)
    s.currentUAVTasks--
    return s.assignTask(ctx, selectedUAV, taskData, taskType, energyCost, computeCost)
}


// selectRandomDevice randomly selects a device from the available devices
func (s *SmartContract) selectRandomDevice(devices []Device) Device {
    rand.Seed(time.Now().UnixNano())
    return devices[rand.Intn(len(devices))]
}

// selectRandomEC ensures no consecutive EC selection and tracks the last assigned EC
func (s *SmartContract) selectRandomEC(ctx contractapi.TransactionContextInterface, ecs []Device) Device {
    var eligibleECs []Device
    for _, ec := range ecs {
        if ec.DeviceID != s.lastUsedEC { // Ensure we don't pick the last EC used
            eligibleECs = append(eligibleECs, ec)
        }
    }

    // If all ECs were previously used, reset the last used EC and pick a new one from the pool
    if len(eligibleECs) == 0 {
        s.lastUsedEC = "" // Reset, allowing any EC to be selected again
        eligibleECs = ecs
    }

    selectedEC := s.selectRandomDevice(eligibleECs)
    s.lastUsedEC = selectedEC.DeviceID // Update the last used EC
    return selectedEC
}

// Check if all ECs have been assigned at least once before switching to UAV phase
func (s *SmartContract) areAllECsAssigned(ecs []Device) bool {
    for _, ec := range ecs {
        if ec.TasksCompleted == 0 {
            return false // If any EC hasn't completed a task, return false
        }
    }
    return true // All ECs have completed a task
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Task Offload based on a Energy-Aware Task Scheduling proposed by Ningning Wang              //
/////////////////////////////////////////////////////////////////////////////////////////////////
//      The paper Energy-Efficient Task Scheduling in Mobile Edge Computing Systems, focus on  //
//      optimizing task scheduling based on the energy levels of devices. Their approach aims  //  
//      extend the network's operational time by ensuring that tasks are offloaded to devices  //  
//      remaining energy. The function TaskOffloadEnergyAware follows the same principle by    //
//      balancing with sufficient battery life and compute resources through an energy         //
//      efficiency score. This scheduling technique ensures that tasks are assigned not only   //
//      based on available computational power but also by considering the energy reserves     //
//      of UAVs.                                                                               //
/////////////////////////////////////////////////////////////////////////////////////////////////

// TaskOffloadEnergyAware assigns a task to ECs first. Once all ECs have completed tasks, it switches to UAVs based on energy efficiency.
func (s *SmartContract) TaskOffloadEnergyAware(ctx contractapi.TransactionContextInterface, taskData string, taskType string, energyCost float64, computeCost float64) error {
    // Get available devices (both ECs and UAVs)
    devices, err := s.getAvailableDevices(ctx, computeCost)
    if err != nil {
        return err
    }

    var ecs, uavs []Device
    for _, device := range devices {
        if device.DeviceType == "EC" {
            ecs = append(ecs, device)
        } else if device.DeviceType == "UAV" {
            uavs = append(uavs, device)
        }
    }

    // Ensure we have at least one device to process tasks
    if len(ecs) == 0 && len(uavs) == 0 {
        return fmt.Errorf("No available devices with sufficient ressources")
    }

    // Check if we are in the UAV phase, where UAVs must handle 
    if s.currentUAVTasks > 0 {
        // Assign tasks to UAVs based on energy efficiency score
        selectedUAV := s.selectBestUAVByEnergyScore(uavs, energyCost, computeCost)
        s.currentUAVTasks-- // Decrease UAV task count
        return s.assignTask(ctx, selectedUAV, taskData, taskType, energyCost, computeCost)
    }

    // Prioritize ECs if available and they haven't all completed tasks
    if !s.areAllECsAssigned(ecs) {
        selectedEC := s.selectRandomEC(ctx, ecs)
        return s.assignTask(ctx, selectedEC, taskData, taskType, energyCost, computeCost)
    }

    // Once all ECs have completed tasks, switch to UAVs 
    s.currentUAVTasks = len(ecs) * 3
    selectedUAV := s.selectBestUAVByEnergyScore(uavs, energyCost, computeCost)
    s.currentUAVTasks-- // Decrease UAV task count
    return s.assignTask(ctx, selectedUAV, taskData, taskType, energyCost, computeCost)
}

// selectBestUAVByEnergyScore selects the UAV with the highest energy efficiency score
func (s *SmartContract) selectBestUAVByEnergyScore(uavs []Device, energyCost float64, computeCost float64) Device {
    var selectedUAV Device
    highestScore := -1.0

    // Iterate through UAVs and calculate their energy efficiency score
    for _, uav := range uavs {
        if uav.ComputeResources >= computeCost {
            score := s.calculateEnergyEfficiencyScore(uav, energyCost, computeCost)
            if score > highestScore {
                highestScore = score
                selectedUAV = uav
            }
        }
    }

    return selectedUAV
}

// calculateEnergyEfficiencyScore computes a score based on the device's battery life and compute resources.
// Higher scores represent better candidates for task assignment.
func (s *SmartContract) calculateEnergyEfficiencyScore(device Device, energyCost float64, computeCost float64) float64 {
    batteryWeight := 0.75 // Prioritize battery life (since energy efficiency is key)
    computeWeight := 0.25 // Compute resources are less important

    // Calculate energy efficiency score as a weighted suma of battery life and compute resources
    score := (batteryWeight * device.BatteryLife) + (computeWeight * device.ComputeResources)

    // Penalize UAVs with low battery to avoid rapid depletion
    if device.DeviceType == "UAV" && device.BatteryLife < energyCost {
        score -= 10.0 // Arbitrary penalty for low battery UAVs
    }

    return score
}


// assignTask handles the task assignment and device updates for all normal model
func (s *SmartContract) assignTask(ctx contractapi.TransactionContextInterface, device Device, taskData string, taskType string, energyCost float64, computeCost float64) error {
    // Generate a unique TaskID
    taskID := ctx.GetStub().GetTxID()

    // Deduct ComputeCost and EnergyCost (if applicable)
    device.ComputeResources -= computeCost
    device.ComputeCostDevice += computeCost 
    if device.DeviceType == "UAV" {
        device.BatteryLife -= energyCost
        if device.BatteryLife < 3 {
            device.Status = "Unavailable"
        }
    }
    if device.ComputeResources < 3 {
        device.Status = "Busy"
    }

    // Simulate task execution delay based on task type (randomized sleep) based on the paper "Ultra-reliable and low-latency communications: 
    // applications, opportunities and challenges" by Daquan FENG 2021 and "Evolved Immersive Experience: Exploring 5G- and
    // Beyond-Enabled Ultra-Low-Latency Communications for Augmented and Virtual Reality" by Hazarika2023

    var minSleep, maxSleep int
    switch taskType {
    case "IC":
        minSleep, maxSleep = 850, 1100
    case "HRLLC":
        minSleep, maxSleep = 50, 150
    case "UC":
        minSleep, maxSleep = 700, 900
    case "MC":
        minSleep, maxSleep = 550, 700
    case "AIC":
        minSleep, maxSleep = 1400, 2100
    case "ISC":
        minSleep, maxSleep = 400, 650
    }
    
    randomSleepDuration := time.Duration(rand.Intn(maxSleep-minSleep)+minSleep) * time.Millisecond
    time.Sleep(randomSleepDuration)

    // Calculate avgSleep in milliseconds and convert to time.Duration
    avgSleep := (minSleep + maxSleep) / 2
    avgSleepDuration := time.Duration(avgSleep) * time.Millisecond

    // Convert tolerance to time.Duration (10% of avgSleep)
    tolerance := time.Duration(float64(avgSleep) * 0.1) * time.Millisecond

    // Compare randomSleepDuration to avgSleep + tolerance
    if randomSleepDuration <= avgSleepDuration+tolerance {
        device.TimeTasks++ // Increment the count of on-time tasks
    }

    // Update task status to Completed
    task := Task{
        TaskID:       taskID,
        DeviceID:     device.DeviceID,
        TaskData:     taskData,
        TaskType:     taskType,
        EnergyCost:   energyCost,
        ComputeCost:  computeCost,
        Status:       "Completed",
    }
    taskAsBytes, err := json.Marshal(task)
    if err != nil {
        return err
    }
    err = ctx.GetStub().PutState("T"+taskID, taskAsBytes)
    if err != nil {
        return err
    }

    // Update device state after task completion
    device.TasksCompleted++
    device.TaskLimit++
    device.TotalTasks++    

    // Check if the device has completed enough tasks to reset its ComputeResources
    if (device.DeviceType == "UAV" && device.TaskLimit >= 3) || (device.DeviceType == "EC" && device.TaskLimit >= 30) {
        device.ComputeResources = device.InitialResources // Reset compute resources to initial value
        device.TaskLimit = 0 // Reset task counter after resource reset
    }


    if device.ComputeResources >= 3 && device.Status == "Busy" {
        device.Status = "Available"
    }

    updatedDeviceAsBytes, err := json.Marshal(device)
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState("D"+device.DeviceID, updatedDeviceAsBytes)
}




/////////////////////////////////////////////////////////////////////////////////////////////////
// Task Offload based on COBRA Framework                                                       //
/////////////////////////////////////////////////////////////////////////////////////////////////
//      CoBRA framework use TCI and RI to choose for each task the best Devices based          //
//      on the porperties of the task                                                          //
/////////////////////////////////////////////////////////////////////////////////////////////////

// TaskOffloadCobra assigns tasks using the COBRA algorithm based on RI and TCI and Reputation
func (s *SmartContract) TaskOffloadCobra(ctx contractapi.TransactionContextInterface, taskData string, taskType string, energyCost float64, computeCost float64, lambda float64, epsilon float64) error {
    devices, err := s.getAvailableDevices(ctx, computeCost)
    if err != nil {
        return err
    }

    var ecs, uavs []Device
    for _, device := range devices {
        if device.DeviceType == "EC" {
            ecs = append(ecs, device)
        } else if device.DeviceType == "UAV" {
            uavs = append(uavs, device)
        }
    }

    // Calculate TCI for the task
    tci := s.calculateTaskCostIndex(energyCost, computeCost, epsilon)

    // If TCI is low, prefer UAVs, otherwise prefer ECs
    if tci < 0.55 && len(uavs) > 0 {
        bestUAV := s.selectBestDeviceByRI(uavs, lambda, epsilon)
        return s.assignTaskCobra(ctx, bestUAV, taskData, taskType, energyCost, computeCost, lambda)
    } else if len(ecs) > 0 {
        bestEC := s.selectBestECByRI(ecs, lambda, epsilon)
        return s.assignTaskCobra(ctx, bestEC, taskData, taskType, energyCost, computeCost, lambda)
    } else if len(uavs) > 0 {
        bestUAV := s.selectBestDeviceByRI(uavs, lambda, epsilon)
        return s.assignTaskCobra(ctx, bestUAV, taskData, taskType, energyCost, computeCost, lambda)
    }

    return fmt.Errorf("No available devices with sufficient ComputeResources")
}

// Utility Functions
// getAvailableDevices retrieves all available devices with sufficient compute resources
func (s *SmartContract) getAvailableDevices(ctx contractapi.TransactionContextInterface, computeCost float64) ([]Device, error) {
    startKey := "D0001"
    endKey := "D9999"

    resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var devices []Device
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var device Device
        err = json.Unmarshal(queryResponse.Value, &device)
        if err != nil {
            return nil, err
        }

        if device.Status == "Available" && device.ComputeResources >= computeCost {
            devices = append(devices, device)
        }
    }

    return devices, nil
}

// assignTaskCobra handles the task assignment and device updates
func (s *SmartContract) assignTaskCobra(ctx contractapi.TransactionContextInterface, device Device, taskData string, taskType string, energyCost float64, computeCost float64, lambda float64) error {
    // Generate a unique TaskID
    taskID := ctx.GetStub().GetTxID()

    // Deduct ComputeCost and EnergyCost (if applicable)
    device.ComputeResources -= computeCost
    device.ComputeCostDevice += computeCost 
    if device.DeviceType == "UAV" {
        device.BatteryLife -= energyCost
        if device.BatteryLife < 3 {
            device.Status = "Unavailable"
        }
    }
    if device.ComputeResources < 3 {
        device.Status = "Busy"
    }

    // Simulate task execution delay based on task type (randomized sleep) based on the paper "Ultra-reliable and low-latency communications: 
    // applications, opportunities and challenges" by Daquan FENG 2021 and "Evolved Immersive Experience: Exploring 5G- and
    // Beyond-Enabled Ultra-Low-Latency Communications for Augmented and Virtual Reality" by Hazarika2023

    var minSleep, maxSleep int
    switch taskType {
    case "IC":
        minSleep, maxSleep = 850, 1100
    case "HRLLC":
        minSleep, maxSleep = 50, 150
    case "UC":
        minSleep, maxSleep = 700, 900
    case "MC":
        minSleep, maxSleep = 550, 700
    case "AIC":
        minSleep, maxSleep = 1400, 2100
    case "ISC":
        minSleep, maxSleep = 400, 650
    }
    
    randomSleepDuration := time.Duration(rand.Intn(maxSleep-minSleep)+minSleep) * time.Millisecond
    time.Sleep(randomSleepDuration)

    // Calculate avgSleep in milliseconds and convert to time.Duration
    avgSleep := (minSleep + maxSleep) / 2
    avgSleepDuration := time.Duration(avgSleep) * time.Millisecond

    // Convert tolerance to time.Duration (10% of avgSleep)
    tolerance := time.Duration(float64(avgSleep) * 0.1) * time.Millisecond

    // Compare randomSleepDuration to avgSleep + tolerance
    if randomSleepDuration <= avgSleepDuration+tolerance {
        device.TimeTasks++ // Increment the count of on-time tasks
    }

    // Update task status to Completed
    task := Task{
        TaskID:       taskID,
        DeviceID:     device.DeviceID,
        TaskData:     taskData,
        TaskType:     taskType,
        EnergyCost:   energyCost,
        ComputeCost:  computeCost,
        Status:       "Completed",
    }
    taskAsBytes, err := json.Marshal(task)
    if err != nil {
        return err
    }
    err = ctx.GetStub().PutState("T"+taskID, taskAsBytes)
    if err != nil {
        return err
    }

    // Update device state after task completion
    device.TasksCompleted++
    device.TaskLimit++
    device.TotalTasks++    


    // Recalculate reputation every 5 tasks
    if device.TasksCompleted%5 == 0 {
        // Save current reputation in PreviousReputation
        device.PreviousReputation = device.Reputation

        // Ensure TotalTasks is non-zero to avoid division by zero
        totalTasks := device.TotalTasks
        if totalTasks == 0 {
        totalTasks = 1 // Default to 1 if no tasks have been recorded yet
        }

        // Recalculate and update Reputation
        // Calculate the success rate: Completed tasks / Total tasks
        successRate := float64(device.TasksCompleted) / float64(totalTasks)

        // Calculate the on-time task completion rate: On-time tasks / Total tasks
        timeRate := float64(device.TimeTasks) / float64(totalTasks)

        // Reputation calculation
        reputationScore := (lambda * (successRate + timeRate)) + ((1 - lambda) * float64(device.PreviousReputation))
        device.Reputation = reputationScore
    }

    // Check if the device has completed enough tasks to reset its ComputeResources
    if (device.DeviceType == "UAV" && device.TaskLimit >= 3) || (device.DeviceType == "EC" && device.TaskLimit >= 30) {
        device.ComputeResources = device.InitialResources // Reset compute resources to initial value
        device.TaskLimit = 0 // Reset task counter after resource reset
    }


    if device.ComputeResources >= 3 && device.Status == "Busy" {
        device.Status = "Available"
    }

    updatedDeviceAsBytes, err := json.Marshal(device)
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState("D"+device.DeviceID, updatedDeviceAsBytes)
}

// calculateReliabilityIndexAndReputation calculates RI based on the updated formula with the reputation
func (s *SmartContract) calculateReliabilityIndexAndReputation(device Device, lambda float64, epsilon float64) float64 {
    

    reputationScore := (lambda * float64(device.Reputation)) + ((1 - lambda) * float64(device.PreviousReputation))

    // Calculate the resource availability ratio: Current compute resources / Initial compute resources
    resourceRatio := device.ComputeResources / device.InitialResources

    // Calculate battery life ratio: Current battery life / Initial battery life (only for UAVs)
    var batteryRatio float64
    if device.DeviceType == "UAV" {
        if device.BatteryLife < 0 {
            device.BatteryLife = 0 // Ensure battery life is non-negative
        }
        batteryRatio = device.BatteryLife / device.InitialBattery
    } else {
        batteryRatio = 0 // Battery ratio is ignored for ECs (λ = 0)
    }

    // Return the calculated Reliability Index
    return epsilon*(batteryRatio+resourceRatio) + (1-epsilon)*reputationScore
}


// calculateTaskCostIndex calculates TCI for a given task
func (s *SmartContract) calculateTaskCostIndex(energyCost, computeCost float64, epsilon float64) float64 {
    maxEnergyCost := 3.0
    maxComputeCost := 3.0

    // Weight with epsilon for energy
    

    return (energyCost/maxEnergyCost)*(1 - epsilon) + (computeCost/maxComputeCost)*epsilon
}

// selectBestECByRI selects the best EC by RI, ensuring the last selected EC is not used consecutively
func (s *SmartContract) selectBestECByRI(ecs []Device, lambda float64, epsilon float64) Device {
    var bestEC Device
    highestRI := -1.0

    for _, ec := range ecs {
        if ec.DeviceID != s.lastUsedEC { // Ensure it's not the last used EC
            ri := s.calculateReliabilityIndexAndReputation(ec, lambda, epsilon)
            if ri > highestRI {
                highestRI = ri
                bestEC = ec
            }
        }
    }

    // If all ECs were excluded (e.g., lastUsedEC was the only one), fallback to the best available EC
    if bestEC.DeviceID == "" && len(ecs) > 0 {
        bestEC = ecs[0] // Fallback to the first EC in the list
    }

    s.lastUsedEC = bestEC.DeviceID // Update the last used EC
    return bestEC
}

// selectBestDeviceByRI selects the best device by RI (for UAVs)
func (s *SmartContract) selectBestDeviceByRI(devices []Device, lambda float64, epsilon float64) Device {
    var bestDevice Device
    highestRI := -1.0

    for _, device := range devices {
        ri := s.calculateReliabilityIndexAndReputation(device, lambda, epsilon)
        if ri > highestRI {
            highestRI = ri
            bestDevice = device
        }
    }

    return bestDevice
}

/////////////////////////////////////////////////////////////////////////////////////////////////
// Section 3 : Other function for manage of the ledger and result                              //
/////////////////////////////////////////////////////////////////////////////////////////////////

// QueryAllDevices gets all devices from the world state
func (s *SmartContract) QueryAllDevices(ctx contractapi.TransactionContextInterface) ([]Device, error) {
    startKey := "D0001"
    endKey := "D9999"

    resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var devices []Device
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var device Device
        err = json.Unmarshal(queryResponse.Value, &device)
        if err != nil {
            return nil, err
        }

        devices = append(devices, device)
    }

    return devices, nil
}

// QueryAllTasks gets all tasks from the world state
func (s *SmartContract) QueryAllTasks(ctx contractapi.TransactionContextInterface) ([]Task, error) {
    startKey := "T00000001"
    endKey := "T99999999"

    resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var tasks []Task
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var task Task
        err = json.Unmarshal(queryResponse.Value, &task)
        if err != nil {
            return nil, err
        }

        tasks = append(tasks, task)
    }

    return tasks, nil
}

// RegisterDevice registers UAVs or Edge Servers in the blockchain network
func (s *SmartContract) RegisterDevice(ctx contractapi.TransactionContextInterface, deviceID string, deviceType string, status string, batteryLife float64, initialBattery float64,  computeResources float64, initialResources float64, tasksCompleted int, totalTasks int, timeTasks int, computeCostDevice float64, taskLimit int, reputation float64, previousreputation float64  ) error {
    device := Device{
        DeviceID:         deviceID,
        DeviceType:       deviceType,
        Status:           status,
        BatteryLife:      batteryLife,
        InitialBattery:   initialBattery,
        ComputeResources: computeResources,
        InitialResources: initialResources,
        TasksCompleted:   tasksCompleted,
        TotalTasks:       totalTasks,
        TimeTasks:        timeTasks,
        ComputeCostDevice: computeCostDevice,
        TaskLimit:         taskLimit,
        Reputation:        reputation,
        PreviousReputation: previousreputation,  
    }

    deviceAsBytes, err := json.Marshal(device)
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState("D"+deviceID, deviceAsBytes)
    if err != nil {
        return err
    }

    return nil
}


// Delete to clean the all the ledger or just an selection of the ledger
func (s *SmartContract) DeleteAll(ctx contractapi.TransactionContextInterface, deleteType string) error {
    if deleteType == "tasks" || deleteType == "all" {
        startKey := "T00000001"
        endKey := "T99999999"

        resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
        if err != nil {
            return err
        }
        defer resultsIterator.Close()

        for resultsIterator.HasNext() {
            queryResponse, err := resultsIterator.Next()
            if err != nil {
                return err
            }
            err = ctx.GetStub().DelState(queryResponse.Key)
            if err != nil {
                return err
            }
        }
    }

    if deleteType == "devices" || deleteType == "all" {
        startKey := "D0001"
        endKey := "D9999"

        resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
        if err != nil {
            return err
        }
        defer resultsIterator.Close()

        for resultsIterator.HasNext() {
            queryResponse, err := resultsIterator.Next()
            if err != nil {
                return err
            }
            err = ctx.GetStub().DelState(queryResponse.Key)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        panic(err)
    }

    if err := chaincode.Start(); err != nil {
        panic(err)
    }
}
