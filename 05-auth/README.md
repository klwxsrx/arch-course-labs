# Установка
## Добавить ingress в minikube
```shell
minikube addons enable ingress
```

## Traefik
```shell
helm repo add traefik https://helm.traefik.io/traefik
helm repo update
helm install --version "10.14.1" \
  --set ports.websecure.expose=false \
  --set providers.kubernetesIngress.enabled=false \
  --set providers.kubernetesCRD.namespaces=user \
  traefik traefik/traefik --namespace traefik --create-namespace
```

## Применить k8s манифесты
```shell
kubectl apply -f ./k8s
```

# Запуск тестов
```shell
newman run ./Tests.postman_collection.json
```