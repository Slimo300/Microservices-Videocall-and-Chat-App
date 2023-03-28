eval $(minikube docker-env)
docker build -t chat/emailservice ./backend/email-service
docker build -t chat/groupservice ./backend/group-service
docker build -t chat/tokenservice ./backend/token-service
docker build -t chat/searchservice ./backend/search-service
docker build -t chat/userservice ./backend/user-service
docker build -t chat/wsservice ./backend/ws-service
docker build -t chat/messageservice ./backend/message-service
docker build -t chat/frontend ./frontend
docker image rm $(docker images -f "dangling=true" -q)

helm install kafka bitnami/kafka

kubectl apply -f deploy/secrets
kubectl apply -f deploy/dbs

echo "Waiting for dbs to start"
sleep 10s

kubectl apply -f deploy/services
kubectl apply -f deploy/frontend.yaml
kubectl apply -f deploy/ingress.yaml