/////////////////////////////////////////////////////////////////////////////////////////////////
//
// Objet : Go code to simulate send of task by an device for the blockchain
//
// version : 7.2
//
// Author : Rêzan OSCAR
// Infos :
//      - This script simulate task send with an proportion of task type and for each task an  
//      energyCost with computeCost for this task
//
/////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
    cryptoRand "crypto/rand"
    "encoding/csv"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "math"
    "math/rand"
    "os"
    "sync"
    "time"

    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
    numTasks       = 2000 // Number of tasks to be sent
    maxRetries     = 5    // Max retries for a failed task
    maxGoroutines  = 10   // Number of goroutines to simulate concurrent users
    reportInterval = 10   // Intervals by task to show stats information in csv
    reportIntervalScreen = 100   // Intervals by task to show stats information in the screen
    CLevel = 0.95 // 95% CI
    epsilon         = 0.7  // Weight for the energy priority in TaskOffloadCobra
    lambda          = 0.3  // Weight for reputation and previous reputation in TaskOffloadCobra
)

var (
    functionUsed = "TaskOffloadCobra" // Name of the smart contract function
    networkUsed  = "cobra_algo"                 // Name of the blockchain network
)

type TaskType struct {
    Name        string
    EnergyCost  float64
    ComputeCost float64
}

var taskTypes = []TaskType{
    {"IC", 2.2, 2.7},       // Immersive Communication
    {"HRLLC", 1.1, 1.9},    // Hyper-Reliable and Low-Latency Communication
    {"UC", 0.5, 0.9},       // Ubiquitous Connectivity
    {"MC", 0.9, 1.4},       // Massive Communication
    {"AIC", 2.7, 3.0},      // AI and Communication
    {"ISC", 1.2, 2.0},      // Integrated Sensing and Communication
}

// Proportion of task by %
var taskDistributionPercentage = map[string]int{
    "IC":    10,
    "AIC":   15,
    "HRLLC": 10,
    "ISC":   25,
    "UC":    20,
    "MC":    15,
}

type Device struct {
    DeviceID         string  `json:"deviceID"`
    DeviceType       string  `json:"deviceType"`
    Status           string  `json:"status"`
    BatteryLife      float64 `json:"batteryLife"`
    ComputeCostDevice float64 `json:"computeCostDevice"`
    TotalTasks       int     `json:"totalTasks"`
}

// Initialize SDK and create a channel client
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

// Generates a random string for task data
func generateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := cryptoRand.Read(bytes)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes)[:length], nil
}

//  Interval Calculation
func CIC(mean, stddev, size float64, CLevel float64) (float64, float64) {
    z := 1.96 // For 96% 
    marginOfError := z * (stddev / math.Sqrt(size))
    return mean - marginOfError, mean + marginOfError
}

// Queries all devices from the ledger
func queryAllDevices(client *channel.Client) ([]Device, error) {
    response, err := client.Query(channel.Request{
        ChaincodeID: networkUsed,
        Fcn:         "QueryAllDevices",
    })
    if err != nil {
        return nil, err
    }

    var devices []Device
    err = json.Unmarshal(response.Payload, &devices)
    return devices, err
}

// Calculate average battery and count available UAVs
func calculateUAVStats(devices []Device) (float64, int) {
    var totalBattery float64
    var availableUAVs int
    var totalUAVs int

    for _, device := range devices {
        if device.DeviceType == "UAV" {
            // If the battery level is negative, treat it as 0
            batteryLevel := device.BatteryLife
            if batteryLevel < 0 {
                batteryLevel = 0
            }

            totalBattery += batteryLevel
            totalUAVs++ // Count all UAVs, regardless of their status

            // Only count UAVs that are "Available" and have a battery life > 0
            if device.Status == "Available" && batteryLevel > 0 {
                availableUAVs++
            }
        }
    }

    if totalUAVs == 0 {
        totalUAVs = 1 // Avoid division by zero
    }

    averageBattery := totalBattery / float64(totalUAVs) // Calculate the average battery level
    return averageBattery, availableUAVs
}

// Calculate the task proportions for UAVs and ECs as a percentage of total tasks
func calculateTaskProportions(devices []Device) (float64, float64) {
    var totalTasksUAV, totalTasksEC int

    for _, device := range devices {
        if device.DeviceType == "UAV" {
            totalTasksUAV += device.TotalTasks
        } else if device.DeviceType == "EC" {
            totalTasksEC += device.TotalTasks
        }
    }

    totalTasks := totalTasksUAV + totalTasksEC
    if totalTasks == 0 {
        totalTasks = 1 // Avoid division by zero
    }

    totalTaskUAVPercentage := (float64(totalTasksUAV) / float64(totalTasks)) * 100
    totalTaskECPercentage := (float64(totalTasksEC) / float64(totalTasks)) * 100

    return totalTaskUAVPercentage, totalTaskECPercentage
}

// Calculate compute cost stats and average tasks per device type
func calculateDeviceStats(devices []Device) (float64, float64, float64, float64, float64) {
    var totalComputeCostAll, totalComputeCostUAV, totalComputeCostEC float64
    var totalTasksUAV, totalTasksEC int
    var totalUAV, totalEC int

    for _, device := range devices {
        totalComputeCostAll += device.ComputeCostDevice
        if device.DeviceType == "UAV" {
            totalComputeCostUAV += device.ComputeCostDevice
            totalTasksUAV += device.TotalTasks
            totalUAV++
        } else if device.DeviceType == "EC" {
            totalComputeCostEC += device.ComputeCostDevice
            totalTasksEC += device.TotalTasks
            totalEC++
        }
    }

    if totalUAV == 0 {
        totalUAV = 1 // Prevent division by zero
    }
    if totalEC == 0 {
        totalEC = 1 // Prevent division by zero
    }

    avgComputeCostAll := totalComputeCostAll / float64(len(devices))
    avgComputeCostUAV := totalComputeCostUAV / float64(totalUAV)
    avgComputeCostEC := totalComputeCostEC / float64(totalEC)

    avgTasksUAV := float64(totalTasksUAV) / float64(totalUAV)
    avgTasksEC := float64(totalTasksEC) / float64(totalEC)

    return avgComputeCostAll, avgComputeCostUAV, avgComputeCostEC, avgTasksUAV, avgTasksEC
}

// Send a task to the blockchain with retry mechanism
func sendTask(client *channel.Client, taskData string, taskType TaskType, wg *sync.WaitGroup, mu *sync.Mutex, results chan<- map[string]interface{}) {
    defer wg.Done()

    args := [][]byte{
        []byte(taskData),
        []byte(taskType.Name),
        []byte(fmt.Sprintf("%.2f", taskType.EnergyCost)),
        []byte(fmt.Sprintf("%.2f", taskType.ComputeCost)),
        []byte(fmt.Sprintf("%.2f", lambda)),  // Pass lambda to the blockchain
        []byte(fmt.Sprintf("%.2f", epsilon)), // Pass epsilon to the blockchain
    }

    var success bool
    var attempts int
    var start time.Time
    var end time.Time

    for attempts = 1; attempts <= maxRetries; attempts++ {
        mu.Lock()
        start = time.Now()
        _, err := client.Execute(channel.Request{
            ChaincodeID: networkUsed,
            Fcn:         functionUsed,
            Args:        args,
        })
        end = time.Now()
        mu.Unlock()

        if err == nil {
            success = true
            break
        }

        time.Sleep(100 * time.Millisecond * time.Duration(attempts)) // Linear backoff
    }

    duration := end.Sub(start)

    results <- map[string]interface{}{
        "taskData":  taskData,
        "success":   success,
        "attempts":  attempts,
        "startTime": start,
        "endTime":   end,
        "duration":  duration.Seconds(),
    }
}


// Generate task distribution based on percentages
func generateTaskDistribution(taskTypes []TaskType, numTasks int) []TaskType {
    var tasks []TaskType
    totalTasks := 0

    // Distribute tasks based on percentages
    for taskName, percentage := range taskDistributionPercentage {
        count := (percentage * numTasks) / 100
        for _, taskType := range taskTypes {
            if taskType.Name == taskName {
                for i := 0; i < count; i++ {
                    tasks = append(tasks, taskType)
                }
                totalTasks += count
            }
        }
    }

    // Handle leftover tasks due to rounding issues
    leftoverTasks := numTasks - totalTasks
    for i := 0; i < leftoverTasks; i++ {
        tasks = append(tasks, taskTypes[0])
    }

    // Shuffle tasks for randomness
    rand.Shuffle(len(tasks), func(i, j int) { tasks[i], tasks[j] = tasks[j], tasks[i] })

    return tasks
}

// Calculate mean and stddev for task durations
func calculateMeanAndStdDev(durations []float64) (float64, float64) {
    var sum, mean, variance, stddev float64
    n := len(durations)

    // Calculate mean
    for _, duration := range durations {
        sum += duration
    }
    mean = sum / float64(n)

    // Calculate variance
    for _, duration := range durations {
        variance += math.Pow(duration-mean, 2)
    }
    variance /= float64(n)

    // Calculate standard deviation
    stddev = math.Sqrt(variance)

    return mean, stddev
}

// Write results into the CSV with additional stats, including TotalTaskUAV and TotalTaskEC in percentage
func writeResultsToCSV(filename string, uavBatteryAvg []float64, uavAvailable []int, avgComputeCostAll []float64, avgComputeCostUAV []float64, avgComputeCostEC []float64, avgTasksUAV []float64, avgTasksEC []float64, timeDelay []float64, totalTaskUAVPercentage []float64, totalTaskECPercentage []float64) error {
    // Convert the uavAvailable slice from []int to []float64
    uavAvailableFloat := intSliceToFloat64Slice(uavAvailable)

    // Find the minimum length to avoid index out-of-range errors
    minLength := minSliceLength(uavBatteryAvg, uavAvailableFloat, avgComputeCostAll, avgComputeCostUAV, avgComputeCostEC, avgTasksUAV, avgTasksEC, timeDelay)

    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create CSV file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{
        "Total Tasks", "UAV Battery Avg", "UAV Available", "Avg ComputeCost (All)", "Avg ComputeCost (UAV)", "Avg ComputeCost (EC)",
        "Avg Tasks (UAV)", "Avg Tasks (EC)", "Time Delay", "TotalTaskUAV (%)", "TotalTaskEC (%)",
    }
    err = writer.Write(header)
    if err != nil {
        return fmt.Errorf("failed to write CSV header: %w", err)
    }

    for i := 0; i < minLength; i++ {
        row := []string{
            fmt.Sprintf("%d", i*reportInterval), // Total tasks up to this point
            fmt.Sprintf("%.2f", uavBatteryAvg[i]*2), // Battery Avg multiplied by 2
            fmt.Sprintf("%.0f", uavAvailableFloat[i]), // No decimals for UAV Available
            fmt.Sprintf("%.2f", avgComputeCostAll[i]),
            fmt.Sprintf("%.2f", avgComputeCostUAV[i]),
            fmt.Sprintf("%.2f", avgComputeCostEC[i]),
            fmt.Sprintf("%.2f", avgTasksUAV[i]),
            fmt.Sprintf("%.2f", avgTasksEC[i]),
            fmt.Sprintf("%.2f", timeDelay[i]),
            fmt.Sprintf("%.2f", totalTaskUAVPercentage[i]), // TotalTaskUAV% 
            fmt.Sprintf("%.2f", totalTaskECPercentage[i]),  // TotalTaskEC%
        }
        err := writer.Write(row)
        if err != nil {
            return fmt.Errorf("failed to write CSV row: %w", err)
        }
    }

    return nil
}


// Helper function to convert []int to []float64
func intSliceToFloat64Slice(intSlice []int) []float64 {
    float64Slice := make([]float64, len(intSlice))
    for i, val := range intSlice {
        float64Slice[i] = float64(val)
    }
    return float64Slice
}

// Find the minimum slice length among the slices
func minSliceLength(slices ...interface{}) int {
    minLength := len(slices[0].([]float64))
    for _, slice := range slices[1:] {
        length := len(slice.([]float64))
        if length < minLength {
            minLength = length
        }
    }
    return minLength
}

func main() {
    startTime := time.Now()

    // Initialize SDK and channel client
    sdk, channelClient, err := initSDKAndClient("cobra-config.yaml", "channelcoop", "Admin", "Provider1MSP")
    if err != nil {
        log.Fatalf("Initialization failed: %s", err)
    }
    defer sdk.Close()

    var wg sync.WaitGroup
    var mu sync.Mutex

    rand.Seed(time.Now().UnixNano())

    results := make(chan map[string]interface{}, numTasks)
    successCount := 0
    failCount := 0
    totalDuration := 0.0
    durations := []float64{} // Track task durations

    sem := make(chan struct{}, maxGoroutines)

    // Initial stats
    devices, err := queryAllDevices(channelClient)
    if err != nil {
        log.Fatalf("Failed to query devices: %v", err)
    }

    uavBatteryAvg := make([]float64, 0, (numTasks/reportInterval)+1)
    uavAvailable := make([]int, 0, (numTasks/reportInterval)+1)
    timeDelay := make([]float64, 0, (numTasks/reportInterval)+1)

    avgBattery, availableUAVs := calculateUAVStats(devices)
    fmt.Printf("Initial state:\n - Average UAV Battery: %.2f%%\n - Available UAVs: %d\n", avgBattery*2, availableUAVs)
    uavBatteryAvg = append(uavBatteryAvg, avgBattery) 
    uavAvailable = append(uavAvailable, availableUAVs)
    timeDelay = append(timeDelay, 0.0)

    avgComputeCostAll := make([]float64, 0, (numTasks/reportInterval)+1)
    avgComputeCostUAV := make([]float64, 0, (numTasks/reportInterval)+1)
    avgComputeCostEC := make([]float64, 0, (numTasks/reportInterval)+1)
    avgTasksUAV := make([]float64, 0, (numTasks/reportInterval)+1)
    avgTasksEC := make([]float64, 0, (numTasks/reportInterval)+1)
    totalTaskUAVPercentage := make([]float64, 0, (numTasks/reportInterval)+1)
    totalTaskECPercentage := make([]float64, 0, (numTasks/reportInterval)+1)



    // Gather initial device stats before any tasks are processed
    avgComputeAll, avgComputeUAV, avgComputeECBatch, avgTasksUAVBatch, avgTasksECBatch := calculateDeviceStats(devices)
    avgComputeCostAll = append(avgComputeCostAll, avgComputeAll)
    avgComputeCostUAV = append(avgComputeCostUAV, avgComputeUAV)
    avgComputeCostEC = append(avgComputeCostEC, avgComputeECBatch)
    avgTasksUAV = append(avgTasksUAV, avgTasksUAVBatch)
    avgTasksEC = append(avgTasksEC, avgTasksECBatch)

    totalTaskUAVPercent, totalTaskECPercent := calculateTaskProportions(devices)
    totalTaskUAVPercentage = append(totalTaskUAVPercentage, totalTaskUAVPercent)
    totalTaskECPercentage = append(totalTaskECPercentage, totalTaskECPercent)


    // Generate task distribution
    taskDistribution := generateTaskDistribution(taskTypes, numTasks)

    // Variables to capture time at specific task intervals
    var timeAt10, timeAt20, timeAt30, timeAt40, timeAt50, timeAt60, timeAt70  time.Duration

    csvFilename := fmt.Sprintf("graphe_result_%s.csv", functionUsed)

    // Process all tasks
    for i := 0; i < numTasks; i++ {
        wg.Add(1)

        taskData, err := generateRandomString(8)
        if err != nil {
            log.Fatalf("Failed to generate task data: %v", err)
        }

        taskType := taskDistribution[i]

        sem <- struct{}{}
        go func(taskType TaskType) {
            defer func() { <-sem }()
            sendTask(channelClient, taskData, taskType, &wg, &mu, results)
        }(taskType)

        wg.Wait()

        // Capture times at specific intervals
        elapsed := time.Since(startTime).Seconds()
        if i+1 == 10 {
            timeAt10 = time.Since(startTime)
        }
        if i+1 == 20 {
            timeAt20 = time.Since(startTime)
        }
        if i+1 == 30 {
            timeAt30 = time.Since(startTime)
        }
        if i+1 == 40 {
            timeAt40 = time.Since(startTime)
        }
        if i+1 == 50 {
            timeAt50 = time.Since(startTime)
        }
        if i+1 == 60 {
            timeAt60 = time.Since(startTime)
        }
        if i+1 == 70 {
            timeAt70 = time.Since(startTime)
        }

        // Report interval for stats
        if (i+1)%reportInterval == 0 {
            devices, _ = queryAllDevices(channelClient)
            avgBattery, availableUAVs = calculateUAVStats(devices)
            avgComputeAll, avgComputeUAV, avgComputeECBatch, avgTasksUAVBatch, avgTasksECBatch := calculateDeviceStats(devices)

            uavBatteryAvg = append(uavBatteryAvg, avgBattery)
            uavAvailable = append(uavAvailable, availableUAVs)
            avgComputeCostAll = append(avgComputeCostAll, avgComputeAll)
            avgComputeCostUAV = append(avgComputeCostUAV, avgComputeUAV)
            avgComputeCostEC = append(avgComputeCostEC, avgComputeECBatch)
            avgTasksUAV = append(avgTasksUAV, avgTasksUAVBatch)
            avgTasksEC = append(avgTasksEC, avgTasksECBatch)
            timeDelay = append(timeDelay, elapsed)
            totalTaskUAVPercent, totalTaskECPercent = calculateTaskProportions(devices)
            totalTaskUAVPercentage = append(totalTaskUAVPercentage, totalTaskUAVPercent)
            totalTaskECPercentage = append(totalTaskECPercentage, totalTaskECPercent)


            // Write updated stats to CSV after every reportInterval
            err = writeResultsToCSV(csvFilename, uavBatteryAvg, uavAvailable, avgComputeCostAll, avgComputeCostUAV, avgComputeCostEC, avgTasksUAV, avgTasksEC, timeDelay, totalTaskUAVPercentage, totalTaskECPercentage)
            if err != nil {
                log.Fatalf("Failed to write results to CSV: %v", err)
            }
        }


        // Screen Stats Display Interval
        if (i+1)%reportIntervalScreen == 0  {
            devices, _ = queryAllDevices(channelClient)
            avgBattery, availableUAVs = calculateUAVStats(devices)
            elapsed := time.Since(startTime).Seconds()

            fmt.Printf("After %d tasks:\n - Average UAV Battery: %.2f%%\n - Available UAVs: %d\n - Time elapsed: %.2f seconds\n", i+1, avgBattery*2, availableUAVs, elapsed)
        }
    }

    close(results)

    for result := range results {
        if result["success"].(bool) {
            successCount++
        } else {
            failCount++
        }
        duration := result["duration"].(float64)
        durations = append(durations, duration)
        totalDuration += duration
    }

    endTime := time.Now()
    duration := endTime.Sub(startTime)

    // Calculate mean and standard deviation from task durations
    mean, stddev := calculateMeanAndStdDev(durations)

    bandwidth := float64(numTasks) / duration.Seconds()

    // Calculate confirmation and consensus time (assuming they are the same for this simulation)
    transactionConfirmationTime := totalDuration / float64(numTasks)
    consensusTime := transactionConfirmationTime

    // Display results 
    ciLow, ciHigh := CIC(mean, stddev, float64(numTasks), CLevel)

    fmt.Printf("=====================================\n")
    fmt.Printf("Simulation Complete\n")
    fmt.Printf("Total Time: %.2f seconds\n", duration.Seconds())
    fmt.Printf("Number of Tasks Sent: %d\n", numTasks)
    fmt.Printf("Successful Tasks: %d\n", successCount)
    fmt.Printf("Failed Tasks: %d\n", failCount)
    fmt.Printf("Total Duration of Successful Tasks: %.2f seconds\n", totalDuration)
    fmt.Printf("Average Task Duration: %.2f seconds (95%% CI: %.2f, %.2f)\n", mean, ciLow, ciHigh)
    fmt.Printf("Bandwidth (tasks per second): %.2f\n", bandwidth)
    fmt.Printf("Transaction Confirmation Time: %.2f seconds\n", transactionConfirmationTime)
    fmt.Printf("Consensus Time: %.2f seconds\n", consensusTime)
    fmt.Printf("Blockchain Function Used: %s\n", functionUsed)
    fmt.Printf("Blockchain Network: %s\n", networkUsed)
    fmt.Printf("\nAverage time for first 10 tasks: %.2f seconds\n", timeAt10.Seconds())
    fmt.Printf("Average time for first 20 tasks: %.2f seconds\n", timeAt20.Seconds())
    fmt.Printf("Average time for first 30 tasks: %.2f seconds\n", timeAt30.Seconds())
    fmt.Printf("Average time for first 40 tasks: %.2f seconds\n", timeAt40.Seconds())
    fmt.Printf("Average time for first 50 tasks: %.2f seconds\n", timeAt50.Seconds())
    fmt.Printf("Average time for first 60 tasks: %.2f seconds\n", timeAt60.Seconds())
    fmt.Printf("Average time for first 70 tasks: %.2f seconds\n", timeAt70.Seconds())
    fmt.Printf("=====================================\n")

}
