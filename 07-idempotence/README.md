# Установка
## Добавить аддон ingress minikube
```shell
minikube addons enable ingress
```

## Развернуть БД и сервисы
```shell
helm install \
  --set "image.tag=8.0" \
  --set "auth.database=order" \
  --set "auth.username=user" \
  --set "auth.password=test1234" \
  --namespace=order \
  --create-namespace \
  order-db bitnami/mysql

kubectl apply -f ./k8s
```

# Запуск тестов
```shell
newman run --env-var="baseUrl=arch.homework" ./test.postman_collection.json
```