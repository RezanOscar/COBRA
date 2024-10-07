# COBRA Framework : A Cooperative Blockchain Resource Allocation Framework for 6G UAV Networks

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
- The ***"fabric_simulation_client_code"*** folder contains some code in go to interecact with the Hyperledger blockchain, the ***"clean"*** code allow to delete all data in the blockchain, the ***"queryAll"*** code allow to show the data of an ledger in this case the device ledger or the task ledger, the ***"register_device"*** allow to register massively device in the blockchain you can chosse the number of device and the proportion beetween EC or UAV, finnaly the ***"cobra-config"*** yaml file is the most important is allow the communication beetween the client and the blockcahin he containes parameter and credential to acces on the blockchain.
-  The ***"result"*** folder contains different csv result files of the simulation and also a python code to generate graphes.
-  For the 2 files, there are the ***"Cobra_Algo_SC"*** go file is the smart contract inplement in my Blockchain and the ***"simulation"*** go file to simulate the task send and have the result, a more detailed explanation is available below.

> [!NOTE]
> Depending on your usage, you will need to adapt the distribution of the proportion of tasks sent and the number of UVA and ES. To do this, you just need to modify the simulation file, and finally you can modify Lambda and Epsilon to change the weight of the energy importance and reputation of the devices.

> [!IMPORTANT]
> All these files can simply work with my Hyperledger blockchain to have the same results you will have to follow the installation of my architecture or adapt the configuration files to your blockchain

# Scenario Use Case

The simulation is carried out following a use case of the framework in the future indeed with the rapid evolution of 6G technologies and the rise of NTNs, such as UAVs and Edge servers, it becomes possible to deploy temporary and powerful communication networks in regions lacking infrastructure. The combination of these technologies with blockchain offers increased security and transparency, essential for the management of resources and data in decentralized environments. As the requirements for intelligent transport and automated industrial processes increase, this platform will provide an integrated and scalable solution, optimizing operations while ensuring a fair distribution of benefits for each actor in the network.

This scenario therefore proposes a cooperative computing platform using UAVs (drones) and Edge servers to create an integrated solution that simultaneously meets the needs of connectivity, data processing and optimization of transport systems. In rural and remote areas, the challenges of connectivity, automation of industrial processes and management of intelligent transport systems are compounded by the lack of suitable infrastructure.

UAVs from different suppliers play a vital role in providing aerial computing power and establishing a temporary communication network or strengthening existing networks in these isolated areas. This infrastructure not only supports industrial operations by automating production systems through AI, but also optimizes the movement and management of autonomous vehicle fleets, making transportation safer and more efficient.

Blockchain is integrated into this platform to ensure increased security and full transparency in interactions between service providers, ensuring a fair distribution of gains based on the contributions of each actor in the network.

# Problem Statement
In Non-Terrestrial Networks (NTNs) composed of UAVs and Edge Computing (EC) servers, the challenge is to efficiently allocate computational tasks to devices with limited resources. The goal is to select the optimal device for each task while minimizing energy consumption, especially for UAVs with limited battery life, and maintaining high service quality. An incentive mechanism is also needed to encourage active participation from infrastructure providers.

# CoBRA Algorithm Overview
The CoBRA framework uses three key metrics to guide its decision-making process:
- Task Cost Index (TCI): Evaluates devices based on energy and computational costs to find the most efficient option.
- Reputation Index (Rp): Measures the reliability of devices based on their past performance, ensuring consistent task execution.
- Reliability Index (RI): Assesses UAVs and nodes based on their battery level and available resources, prioritizing devices that can reliably complete tasks.
The algorithm begins by calculating these indices for each eligible device. It then selects the device with the optimal combination of low TCI and high RI to execute the task. Throughout the process, the blockchain is used to monitor performance and update reputation scores, ensuring a fair and efficient task allocation while promoting collaboration within the network.

# Smart Contract: COBRA Framework
Task Offloading Algorithms: Implements multiple task offloading strategies, including:
- Round Robin: Tasks are distributed to devices in a round-robin fashion.
- Random Selection: Tasks are assigned randomly to available devices.
- ECP : Task Offload based on a choice that priotorize Edge Server  
- Energy-Aware Task Scheduling: Optimizes task scheduling based on energy levels and compute resources, extending the operational time of the network based on Energy-Aware Task Scheduling proposed by Ningning Wang
- COBRA Algorithm: Uses a combination of Task Cost Index (TCI) and Reliability Index (RI) to assign tasks based on device reputation, available resources, and energy levels.

Device Management: 
- Supports registration, status updates, and tracking of UAVs and ECs. Devices are monitored for their compute resources, energy consumption, and overall reputation.

Tracks key metrics such as:

- Task completion time
- Device reputation (with periodic recalculation based on performance)
- Energy usage and compute costs per task

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
## Before starting, make sure you have the following installed:
- Go (v1.18 or Lesser)
- Hyperledger Fabric (v2.x or higher)
- Docker (for running Hyperledger Fabric containers)

## Fabric Network Setup :
Ensure you have set up a Hyperledger Fabric network with multiple organizations and peers. You can follow the official Hyperledger Fabric documentation to set up the environment if not already done https://hyperledger-fabric.readthedocs.io/en/release-2.2/prereqs.html 


Or 

My tuto to set up a rapid network configure with already configure with a Hyperledger Blockchain with 5 Peer / 1 Orderer and based on Raft Consensus or with your preference [https://github.com/RezanOscar/Hyperledger-Blockchain-Network-Network.git](https://github.com/RezanOscar/Hyperledger-Blockchain-Fabric-Network-COBRA.git)


# Implementation
## Smart Contract Installation on you Blockchain:

To deploy the COBRA Framework on your Hyperledger Fabric network:
! All this command bellow work with my VM but you can addapt for your blockchain 

Prepare the environment:
```
export ORDERER_CA=/opt/gopath/fabric-samples/research-network/crypto-config/ordererOrganizations/research-network.com/orderers/orderer.research-network.com/msp/tlscacerts/tlsca.research-network.com-cert.pem
```

Edit and Save the Smart Contract:
```
vim /opt/gopath/src/chain/bto_chaincode/go/test_cobra/test_cobra.go
```

Package the Chaincode:
```
peer lifecycle chaincode package cobra_algo.tar.gz --path opt/gopath/src/chain/bto_chaincode/go/test_cobra/ --lang golang --label cobra_algo
```

Install the Chaincode (in all peer of you blockchain in my case 5 peer):
```
      cd /
      peer lifecycle chaincode install cobra_algo.tar.gz
      peer lifecycle chaincode queryinstalledè
```

Set the Package ID (in all organizations):
```
PACKAGE_ID=cobra_algo:<your-package-id>
```

Approve the Chaincode (in all organizations):
```
    peer lifecycle chaincode approveformyorg -o orderer.research-network.com:7050 --tls true --cafile $ORDERER_CA --channelID channelcoop --name cobra_algo --version 1 --init-required --package-id $PACKAGE_ID --sequence 1 --signature-policy "OR('Provider1MSP.peer', 'Provider2MSP.peer', 'Provider3MSP.peer', 'Provider4MSP.peer', 'Provider5MSP.peer')"
```

Check Commit Readiness:
```
      peer lifecycle chaincode checkcommitreadiness --channelID channelcoop --name cobra_algo --version 1 --sequence 1 --output json --init-required --signature-policy "OR('Provider1MSP.peer', 'Provider2MSP.peer', 'Provider3MSP.peer', 'Provider4MSP.peer', 'Provider5MSP.peer')"
```

Commit the Chaincode (on one peer):
```
      PRO1_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro1.research-network.com/peers/peer0.pro1.research-network.com/tls/ca.crt
      PRO2_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro2.research-network.com/peers/peer0.pro2.research-network.com/tls/ca.crt
      PRO3_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro3.research-network.com/peers/peer0.pro3.research-network.com/tls/ca.crt
      PRO4_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro4.research-network.com/peers/peer0.pro4.research-network.com/tls/ca.crt
      PRO5_CERTFILES_PEER=/opt/gopath/fabric-samples/research-network/crypto-config/peerOrganizations/pro5.research-network.com/peers/peer0.pro5.research-network.com/tls/ca.crt
      peer lifecycle chaincode commit -o orderer.research-network.com:7050 --tls true --cafile $ORDERER_CA --channelID channelcoop --name cobra_algo --peerAddresses peer0.pro1.research-network.com:7051 --tlsRootCertFiles $PRO1_CERTFILES_PEER --peerAddresses peer0.pro2.research-network.com:7051 --tlsRootCertFiles $PRO2_CERTFILES_PEER --peerAddresses peer0.pro3.research-network.com:7051 --tlsRootCertFiles $PRO3_CERTFILES_PEER --peerAddresses peer0.pro4.research-network.com:7051 --tlsRootCertFiles $PRO4_CERTFILES_PEER --peerAddresses peer0.pro5.research-network.com:7051 --tlsRootCertFiles $PRO5_CERTFILES_PEER --version 1 --sequence 1 --init-required --signature-policy "OR('Provider1MSP.peer', 'Provider2MSP.peer', 'Provider3MSP.peer', 'Provider4MSP.peer', 'Provider5MSP.peer')"
```

Query Committed Chaincode:
```
      peer lifecycle chaincode querycommitted --channelID channelcoop --name cobra_algo
```

    
Initialize the Ledger (on one peer):
```
      peer chaincode invoke -o orderer.research-network.com:7050 --tls true --cafile $ORDERER_CA -C channelcoop -n cobra_algo --peerAddresses peer0.pro1.research-network.com:7051 --tlsRootCertFiles $PRO1_CERTFILES_PEER --peerAddresses peer0.pro2.research-network.com:7051 --tlsRootCertFiles $PRO2_CERTFILES_PEER --peerAddresses peer0.pro3.research-network.com:7051 --tlsRootCertFiles $PRO3_CERTFILES_PEER --peerAddresses peer0.pro4.research-network.com:7051 --tlsRootCertFiles $PRO4_CERTFILES_PEER --peerAddresses peer0.pro5.research-network.com:7051 --tlsRootCertFiles $PRO5_CERTFILES_PEER --isInit -c '{"function":"initLedger","Args":[]}'
```

> [!TIP]
> In the case of modification of the smart contract, it will be necessary to reflect the modification on all the files, in the commands above the name is "cobra_algo" which is the one present in the sdk files, always be careful to use the correct SC and the sequence number for each modification the number is incremented

Now you have implement the SC in your Smart Contract, but now you have to use a SDK to communicate with the blockchain and so use your SC.

## Simulation Program and Creation of the Go SDK
The fist step in your user machine, is to install the good go version.
If you don't have a go version than 1.18 follow this step :
```
      sudo rm -rf /usr/local/go
      wget https://golang.org/dl/go1.17.5.linux-amd64.tar.gz
      sudo tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz
      export PATH=$PATH:/usr/local/go/bin
      source ~/.bashrc  # OR source ~/.zshrc if you use zsh
      go version
      rm go1.17.5.linux-amd64.tar.gz
```

Build the Simulation:
```
mkdir fabric-client & cd fabric-client
```
      
  And copy inside all the contenu of the fabric_simulation_client_code
  In the folder you will retrieve 5 files, the most important is the cobra-config.yaml this files allow to connect your machine with blokckchain, after copy all the files do the command bellow, i will create the go environement
      
```
      go mod init fabric-client
      go get github.com/hyperledger/fabric-sdk-go
      go mod tidy
      go mod vendor
```

  After do juste the go build name_of the file

```
  go build Simulation.go
  ./Simulation
```

The simulation will generate 2,000 tasks, send them to the blockchain, and provide performance metrics about 2 Hours for each simulation with a model you can choose one of the 5 different model and modify the number of task and the proportion of task type.

> [!TIP]
> If of course you want this to work with your blockchain, you will need to modify your main config file to allow the connection with your blockchain and modify all the files with the correct information, in particular the use / name part of the channel and the SC
>````
>  // Initialize the SDK and channel client
    sdk, channelClient, err := initSDKAndClient("cobra-config.yaml", "channelcoop", "Admin", "Provider1MSP")
> ````

# Some Example of Result that show the efficiency of my framework:

This graph the variation of the Battery after 1000 task send :

![UAV Battery Variation](https://github.com/user-attachments/assets/3ede3434-ebd9-4482-b958-791c3559bcda)

Here you can see the delay for task, this time delay includes the entire transaction time, also the time taken by the device to execute the task :

![Time Delay](https://github.com/user-attachments/assets/93a4f689-347b-46f2-9c88-fdf8f1ab84ee)

And here you can see the proportion of the offload on the UVA by the different model :

![Proportion](https://github.com/user-attachments/assets/effdd04b-6ee7-408d-86ef-fcf53ee6b150)







