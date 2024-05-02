# STORE-SERVICE

## Install docker
If use wsl : https://docs.docker.com/engine/install/ubuntu/

## Configure visual studio code to use WSL
https://code.visualstudio.com/docs/remote/wsl

## How to run
run the apps using command : 
go mod tidy
go mod vendor
sudo docker compose up --build

## Generate pb file from proto file
### Install buf
https://buf.build/docs/installation
If failed, run : brew install buf

### generate protobuf
Run : buf generate

## Reset database structures (Dont run this! Only if needed)
run : sudo docker compose down --volumes


## Naming Convention
### Error Response
code : grpc code (int32)
code_detail : following these format (string).
    generic : "ERR_UNKNOWN"
    database error : "ERR_TABLENAME_ERRORDETAIL"
    specific error : "ERR_MODULE_FEATURE_ERRORDETAIL"
message : error message (string)