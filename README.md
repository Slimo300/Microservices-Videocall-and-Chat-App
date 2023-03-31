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
11. Application is ready to be deployed to Kubernetes if needed config files are present and environment variables are set

## How to setup?

1. Using minikube 

In order to setup this project you'll first need to provide configuration to files awsSecrets.yaml and smtpCreds.yaml in folder deploy/secrets_templates: 

awsSecrets.yaml stores AWS keys and name of S3 bucket to store files pictures sent to app either as profile pictures or messages
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: awssecrets
type: Opaque
data:
  AWS_ACCESS_KEY_ID: <BASE64_AWS_ACCESS_KEY_ID>
  AWS_SECRET_ACCESS_KEY: <BASE64_AWS_SECRET_ACCESS_KEY>
  S3_BUCKET: <S3_BUCKET_NAME>
```

smtpCreds.yaml stores SMTP credentials of SMTP provider to send verification and reset password email by email-service
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: smtpcreds
type: Opaque
data:
  SMTP_HOST: <SMTP_HOST>
  SMTP_PORT: <SMTP_PORT>
  SMTP_USER: <SMTP_USER>
  SMTP_PASS: <SMTP_PASS>
```
Secret (only AWS and redis need to be base64 encoded)
to base64 encode on Linux distributions you can use: 

```sh
echo -n "YOUR_SECRET" | base64
```

After that rename secrets_templates directory to secrets and run

```sh
minikube start
minikube enable addons ingress
source deploy/deploy.sh
```

2. Using docker-compose 

Setup with docker-compose is similar to that of minikube but instead of setting values in secrets you have to set them 
as environment variables: 

```bash
export AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID>
export AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY>
export S3_BUCKET=<S3_BUCKET>

export SMTP_HOST=<SMTP_HOST>
export SMTP_PORT=<SMTP_PORT>
export SMTP_USER=<SMTP_USER>
export SMTP_PASS=<SMTP_PASS>
```

After that just run:

```bash
docker-compose up -d
```