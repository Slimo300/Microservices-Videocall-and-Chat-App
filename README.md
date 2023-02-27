# MicroservicesChatApp

## This is an example Golang microservice application build in educational purposes.

The goal of this project is to explore microservice architecture and patterns associated with it, 
such as Cloud Native Development, Event Sourcing, etc. This application consists of: 

1. user-service - REST service handling user profiles (profile pictures), registration, login, logout and password reset
2. group-service - REST service handling chat groups, invites to them and user rights in context of a group
3. message-service - REST service for obtaining and deleting messages, connected with message broker to ws-service
4. ws-service - WS service working with websocket connections (uses REST to establish them)
5. token-service - gRPC service working with redis to store and handle tokens, it also distributes public key to other services for validation
6. search-service - REST service working with elasticsearch to query user index to implement search-as-you-type functionality
7. frontend - React frontend for application (needs testing and better error handling :/)
8. File storage - application uses AWS S3 service to store user and group pictures. It also stores files sent by users  (for now only jpeg and png)
9. Message broker - application uses Kafka for asynchronous communication
10. Application is dockerized 
11. Application is ready to be deployed to Kubernetes if needed config files and environment variables are prepared

## How to setup?

In order to setup this project you'll need a couple of things first: 

1. configuration file - every microservice starts with loading config file from file $CHAT_CONFIG/config.yaml. CHAT_CONFIG is an environment describing config file location
The way how config file should look can be seen in sample-conf.yaml

Because I wanted to test this app in different environments (local, docker, kubernetes) they are all looking for different config file
- when started locally services will look for config.yaml (in container runtime other names are also changed to this)
- when using docker (or docker-compose) it will look for docker-conf.yaml (Dockerfile default). Can be changed when building 
```sh
    docker build -f backend/user-service/Dockerfile --build-arg configFile=otherFile.yaml .
```
- when starting kubernetes it looks for kube-conf.yaml
It can be changed by changing deploy/deploy/sh file

2. certificate and private key - since all REST services start both http and https servers its gonna need this project requires certificate and private key in pem format to enable TLS. certificate and private key files should be specified in config file and in Dockerfiles

To generate such cert (it won't be signed by CA so it will be useless in browser but will be good enough to start project) you can use:
```sh
go run %GOROOT%/src/crypto/tls/generate_cert.go
```
You can also check your GOROOT by typing
```sh
go env | grep GOROOT
```

3. private key for token service with which all of access tokens will be signed
There is a tool for generating one in backend/cmd folder. Just run this and move result where you want it to be
```sh 
go run generate_key.go
``` 

4. to run it in kubernetes environment the app needs secrets defined for MYSQL databases, AWS Credentials and Redis password

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: awskey
type: Opaque
data:
  AWS_ACCESS_KEY_ID: <BASE64_ENCODED_KEY>
  AWS_SECRET_ACCESS_KEY: <BASE64_ENCODED_SECRET>
```
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysqlcreds
  labels:
    chat/tier: database
data:
  MYSQL_DATABASE: <MYSQL_DB>
  MYSQL_USER: <MYSQL_USER>
  MYSQL_PASSWORD: <MYSQL_PASS>
  MYSQL_ROOT_PASSWORD: <MYSQL_ROOT_PASS>
```
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: redispass
type: Opaque
data:
  REDIS_PASSWORD: <REDIS_PASS>
```
Secret (only AWS and redis need to be base64 encoded)
to base64 encode on Linux distributions you can use: 

```sh
echo "YOUR_SECRET" -n | base64
```

5. Run
```sh
docker-compose up -d --build
```

or in Minikube

```sh
minikube start
minikube enable addons ingress
source deploy/deploy.sh
```

Yeah, this setup needs to change