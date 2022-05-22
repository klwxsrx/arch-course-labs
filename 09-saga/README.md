# Процесс проведения платежа
![docs/services.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/09-saga/docs/purchase.puml)

Сервис `Order` является оркестратором процесса, реализует паттерн Saga

# Установка
## Добавить аддон ingress minikube
```shell
minikube addons enable ingress
```

## Развернуть БД и сервисы
```shell
helm install \
  --set "image.tag=8.0" \
  --set "auth.database=archcourse" \
  --set "auth.username=user" \
  --set "auth.password=test1234" \
  --namespace=arch-course \
  --create-namespace \
  arch-course-db bitnami/mysql

kubectl apply -f ./k8s
```

# Запуск тестов
```shell
newman run --env-var="baseUrl=arch.homework" ./test.postman_collection.json
```