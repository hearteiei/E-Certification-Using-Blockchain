

# Verification System of E-Certificate using Hyperledger Fabric Blockchain back-end api

### Build Status

You can use Fabric samples to get started working with Hyperledger Fabric, explore important Fabric features, and learn how to build applications that can interact with blockchain networks using the Fabric SDKs. To learn more about Hyperledger Fabric, visit the [Fabric documentation](https://hyperledger-fabric.readthedocs.io/en/latest).

## Getting started with the Fabric samples

To use the Fabric samples, you need to download the Fabric Docker images and the Fabric CLI tools. First, make sure that you have installed all of the [Fabric prerequisites](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html). You can then follow the instructions to [Install the Fabric Samples, Binaries, and Docker Images](https://hyperledger-fabric.readthedocs.io/en/latest/install.html) in the Fabric documentation. In addition to downloading the Fabric images and tool binaries, the Fabric samples will also be cloned to your local machine.

## how to start

in wsl teminal

1. cd test-network
2. ./network.sh down
3. ./network.sh up createChannel -c mychannel -ca
4. ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go

cd ..
cd ..

go to rest-api-go

1. cd asset-transfer-basic
2. cd rest-api-go
3. go mod download     
4. go run main.go

ถ้าอยากแก้ไขchain code ต้องไปแก้ใน chaincode-go สำหรับคนที่ต้องการพัฒนา
ซึ่งอยู่ใน asset-transfer-basic\chaincode-go\chaincode\smartcontract.go