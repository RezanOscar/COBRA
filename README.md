# COBRA
COBRA Framework : A Cooperative Blockchain Resource Allocation Framework for 6G UAV Networks

# Overview
This repository contains the COBRA Framework with is simulation, a smart contract designed to manage task offloading in a decentralized edge computing environment. The framework leverages Non-Terrestrial Networks (NTNs), UAVs (drones), and Edge Servers (ECs) to distribute computational tasks efficiently, focusing on reliability, reputation, cost, and energy-awareness.

The repository also includes a Go-based simulation program that generates tasks and interacts with the blockchain to simulate real-world usage. The simulation tracks task completion, UAV battery life, EC compute resources, and system performance metrics like bandwidth, task duration, and consensus time.

All the UAV and EC are simulate, an EC is 10 time more powerful than a UAV and for the task we simulate 6 different type of task :
      - **Immersive communication**: This will provide users with a rich and interactive video experience, including interactions with machine interfaces. Typical use cases include communication for immersive XR, remote multi-sensory telepresence, and holographic communications. While 5G is theoretically also capable of supporting these use cases, the numbers supported by a single cell are very limited, so this is an extension of those capabilities.
      - **Hyper-reliable and low-latency communication**: This will enhance communications in an industrial environment for full automation, control, and operation. These types of communications can help various applications such as machine interactions, emergency services, telemedicine, and monitoring for electrical power transmission and distribution. This is an extension of 5G’s URLLC.
      - **Ubiquitous connectivity**: This will bridge the digital divide, especially for rural and remote areas. Typical use cases include IoT and mobile broadband communication but can also include access for hikers, farmers and others.
      - **Massive communication**: This will connect many devices or sensors for various use cases and applications, including in smart cities, transportation, logistics, health, energy, environmental monitoring, and agriculture. This, too, is an extension of what has been defined in 5G. 
      - **AI and communication**: This will support automated driving, autonomous collaboration between devices for medical assistance applications, offloading heavy computation operations across devices and networks, creation of and prediction with digital twins, and more.
      - **Integrated sensing and communication**: This will improve applications and services requiring sensing capabilities, such as assisted navigation, activity detection, movement tracking, environmental monitoring, and sensing data on surroundings for AI and XR.

 We simulate the execution of an task with a time-sleep method based on five paper mainly :
      - Ultra-reliable and low-latency communications: applications, opportunities and challenges by FENG et al.
      - Towards an Evolved Immersive Experience: Exploring 5G- and Beyond-Enabled Ultra-Low-Latency Communications for Augmented and Virtual Reality Hazarika et al.
      - High-Reliability and Low-Latency Wireless Communication for Internet of Things: Challenges, Fundamentals, and Enabling Technologies by Ma et al.
      - Sensing and Communication Integrated System for Autonomous Driving Vehicles by Zhang et al.
      - Realizing XR Applications Using 5G-Based 3D Holographic Communication and Mobile Edge Computing Yuan et al.

You can find in this repository, 2 Folder and 2 Files :
          - The "fabric_simulation_client_code" folder contains some code in go to interecact with the Hyperledger blockchain, the "clean" code allow to delete all data in the blockchain, the "queryAll" code allow to show the data of an ledger in this case the device ledger or the task ledger, the "register_device" allow to register massively device in the blockchain you can chosse the number of device and the proportion beetween EC or UAV, finnaly the "cobra-config" yaml file is the most important is allow the communication beetween the client and the blockcahin he containes parameter and credential to acces on the blockchain.
          -  The "result" folder contains different csv result files of the simulation and also a python code to generate graphes.
          -  For the 2 files, there are the "Cobra_Algo_SC" go file is the smart contract inplement in my Blockchain and the "simulation" go file to simulate the task send and have the result, a more detailed explanation is available below.

## Smart Contract: COBRA Framework
   Task Offloading Algorithms: Implements multiple task offloading strategies, including:
          - Round Robin: Tasks are distributed to devices in a round-robin fashion.
          - Random Selection: Tasks are assigned randomly to available devices.
          - Energy-Aware Task Scheduling: Optimizes task scheduling based on energy levels and compute resources, extending the operational time of the network.
          - COBRA Algorithm: Uses a combination of Task Cost Index (TCI) and Reliability Index (RI) to assign tasks based on device reputation, available resources, and energy levels.
   Device Management: Supports registration, status updates, and tracking of UAVs and ECs. Devices are monitored for their compute resources, energy consumption, and overall reputation.

   Metrics: Tracks key metrics such as:

   Task completion time
   Device reputation (with periodic recalculation based on performance)
   Energy usage and compute costs per task

## Simulation Program
The simulation program is a Go-based script that:
          - Generates 2,000 tasks with different types (e.g., Immersive Communication, AI, Ubiquitous Connectivity).
          - Sends tasks to the blockchain and evaluates system performance with real-time stats.
          - Tracks metrics such as:
                - Task success/failure rates
                - Average task duration and standard deviation
                - Bandwidth (tasks per second)
                - Task confirmation and consensus time
                - Device stats: average UAV battery life, available UAVs, EC compute costs, and task distributions across devices.

# Prerequisites
Before starting, make sure you have the following installed:
      - Go (v1.18 or Lesser)
      - Hyperledger Fabric (v2.x or higher)
      - Docker (for running Hyperledger Fabric containers)

Fabric Network Setup
Ensure you have set up a Hyperledger Fabric network with multiple organizations and peers. You can follow the official Hyperledger Fabric documentation to set up the environment if not already done "https://hyperledger-fabric.readthedocs.io/en/release-2.2/prereqs.html" Or my tuto to set up a rapid network configure with you preference "https://github.com/RezanOscar/Hyperledger-Blockchain-Network-Network.git"

OR

You can use my VM already configure with a Hyperledger Blockchain with 5 Peer / 1 Orderer and based on Raft Consensus "https://github.com/RezanOscar/COBRA-VM-Hyperledger-Blockchain.git"

# Implementation
## Smart Contract Installation on you Blockchain :

To deploy the COBRA Framework on your Hyperledger Fabric network:
! All this command bellow work with my VM but you can addapt for your blockchain 

Prepare the environment:
```export ORDERER_CA=/opt/gopath/fabric-samples/research-network/crypto-config/ordererOrganizations/research-network.com/orderers/orderer.research-network.com/msp/tlscacerts/tlsca.research-network.com-cert.pem```

Edit and Save the Smart Contract:
```vim /opt/gopath/src/chain/bto_chaincode/go/test_cobra/test_cobra.go```

Package the Chaincode:
      `
      peer lifecycle chaincode package cobra_algo.tar.gz --path opt/gopath/src/chain/bto_chaincode/go/test_cobra/ --lang golang --label cobra_algo
      `

Install the Chaincode (in all peer of you blockchain in my case 5 peer):
      `
      cd /
      peer lifecycle chaincode install cobra_algo.tar.gz
      peer lifecycle chaincode queryinstalledè
      `

Set the Package ID (in all organizations):
    `PACKAGE_ID=cobra_algo:<your-package-id>`

Approve the Chaincode (in all organizations):
    `
    peer lifecycle chaincode approveformyorg -o orderer.research-network.com:7050 --tls true --cafile $ORDERER_CA --channelID channelcoop --name cobra_algo --version 1 --init-required --package-id $PACKAGE_ID --sequence 1 --signature-policy "OR('Provider1MSP.peer', 'Provider2MSP.peer', 'Provider3MSP.peer', 'Provider4MSP.peer', 'Provider5MSP.peer')"
    `

Check Commit Readiness:
      `
      peer lifecycle chaincode checkcommitreadiness --channelID channelcoop --name cobra_algo --version 1 --sequence 1 --output json --init-required --signature-policy "OR('Provider1MSP.peer', 'Provider2MSP.peer', 'Provider3MSP.peer', 'Provider4MSP.peer', 'Provider5MSP.peer')"
      `

Commit the Chaincode (on one peer):
      `
      PRO1_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro1.research-network.com/peers/peer0.pro1.research-network.com/tls/ca.crt
      PRO2_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro2.research-network.com/peers/peer0.pro2.research-network.com/tls/ca.crt
      PRO3_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro3.research-network.com/peers/peer0.pro3.research-network.com/tls/ca.crt
      PRO4_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro4.research-network.com/peers/peer0.pro4.research-network.com/tls/ca.crt
      PRO5_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro5.research-network.com/peers/peer0.pro5.research-network.com/tls/ca.crt
      peer lifecycle chaincode commit -o orderer.research-network.com:7050 --tls true --cafile $ORDERER_CA --channelID channelcoop --name cobra_algo --peerAddresses peer0.pro1.research-network.com:7051 --tlsRootCertFiles $PRO1_CERTFILES_PEER --peerAddresses peer0.pro2.research-network.com:7051 --tlsRootCertFiles $PRO2_CERTFILES_PEER --peerAddresses peer0.pro3.research-network.com:7051 --tlsRootCertFiles $PRO3_CERTFILES_PEER --peerAddresses peer0.pro4.research-network.com:7051 --tlsRootCertFiles $PRO4_CERTFILES_PEER --peerAddresses peer0.pro5.research-network.com:7051 --tlsRootCertFiles $PRO5_CERTFILES_PEER --version 1 --sequence 1 --init-required --signature-policy "OR('Provider1MSP.peer', 'Provider2MSP.peer', 'Provider3MSP.peer', 'Provider4MSP.peer', 'Provider5MSP.peer')"
      `

Query Committed Chaincode:
      `
      peer lifecycle chaincode querycommitted --channelID channelcoop --name cobra_algo
      `

    
Initialize the Ledger (on one peer):
      `
      peer chaincode invoke -o orderer.research-network.com:7050 --tls true --cafile $ORDERER_CA -C channelcoop -n cobra_algo --peerAddresses peer0.pro1.research-network.com:7051 --tlsRootCertFiles $PRO1_CERTFILES_PEER --peerAddresses peer0.pro2.research-network.com:7051 --tlsRootCertFiles $PRO2_CERTFILES_PEER --peerAddresses peer0.pro3.research-network.com:7051 --tlsRootCertFiles $PRO3_CERTFILES_PEER --peerAddresses peer0.pro4.research-network.com:7051 --tlsRootCertFiles $PRO4_CERTFILES_PEER --peerAddresses peer0.pro5.research-network.com:7051 --tlsRootCertFiles $PRO5_CERTFILES_PEER --isInit -c '{"function":"initLedger","Args":[]}'
      `


## Simulation Program

If you don't have a go version than 1.18 follow this step :
     `
      sudo rm -rf /usr/local/go
      wget https://golang.org/dl/go1.17.5.linux-amd64.tar.gz
      sudo tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz
      export PATH=$PATH:/usr/local/go/bin
      source ~/.bashrc  # OR source ~/.zshrc if you use zsh
      go version
      rm go1.17.5.linux-amd64.tar.gz
      `

Build the Simulation:
      ```mkdir fabric-client & cd fabric-client```
      
  And copy inside all the contenu of the fabric_simulation_client_code
      
   ```go mod init fabric-client
      go get github.com/hyperledger/fabric-sdk-go
      go mod tidy
      go mod vendor```

  After do juste the go build name_of the file

     `go build Simulation.go
      ./Simulation`

The simulation will generate 2,000 tasks, send them to the blockchain, and provide performance metrics about 2 Hours for each simulation with a model you can choose one of the 5 different model and modify the number of task and the proportion of task type.







