helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.12.0 \
  --set installCRDs=true
kubectl apply -f deploy/secrets/ca-secret.yaml
kubectl apply -f deploy/cert-manager/cluster-issuer.yaml
kubectl apply -f deploy/cert-manager/certificate.yaml
sleep 10

helm install \
  kafka bitnami/kafka \
  --namespace default \
  --version v20.0.6

kubectl apply -f deploy/secrets/*
kubectl apply -f deploy/dbs/*
kubectl apply -f deploy/services/*
kubectl apply -f deploy/ingress-prod.yaml