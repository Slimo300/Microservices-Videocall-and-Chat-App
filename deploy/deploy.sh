eval $(minikube docker-env)
docker build -f backend/group-service/Dockerfile -t chat/groupservice .
docker build -f backend/token-service/Dockerfile -t chat/tokenservice .
docker build -f backend/search-service/Dockerfile -t chat/searchservice .
docker build -f backend/user-service/Dockerfile -t chat/userservice .
docker build -f backend/ws-service/Dockerfile -t chat/wsservice .
docker build -f backend/message-service/Dockerfile -t chat/messageservice .
docker build -f frontend/Dockerfile -t chat/frontend .