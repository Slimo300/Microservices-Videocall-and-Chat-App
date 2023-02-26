eval $(minikube docker-env)
docker build -f backend/group-service/Dockerfile --build-arg configFile=config/kube-conf.yaml -t chat/groupservice .
docker build -f backend/token-service/Dockerfile --build-arg configFile=config/kube-conf.yaml -t chat/tokenservice .
docker build -f backend/search-service/Dockerfile --build-arg configFile=config/kube-conf.yaml -t chat/searchservice .
docker build -f backend/user-service/Dockerfile --build-arg configFile=config/kube-conf.yaml -t chat/userservice .
docker build -f backend/ws-service/Dockerfile --build-arg configFile=config/kube-conf.yaml -t chat/wsservice .
docker build -f backend/message-service/Dockerfile --build-arg configFile=config/kube-conf.yaml -t chat/messageservice .
docker build -f frontend/Dockerfile -t chat/frontend .
docker image rm $(docker images -f "dangling=true" -q)

helm install kafka bitnami/kafka

kubectl apply -f deploy/secrets
kubectl apply -f deploy/dbs
kubectl apply -f deploy/services
kubectl apply -f deploy/frontend.yaml
kubectl apply -f deploy/ingress.yaml