# Техническое решение
![docs/services.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/05-auth/docs/services.puml)
В качестве APIGateway используется Traefik.

Сервис `Auth` отвечает за регистрацию и авторизацию пользователей.
Запросы по префиксу `/auth/*` форвардятся Traefik'ом в него.

При успешной логинации в `Auth` пользователю проставляется сессионная кука `sid` (рандомная строка).

Запросы в сервис `UserProfile` осуществляются только авторизованным пользователем.

На каждый запрос в `UserProfile` Traefik дергает `GET /auth` в `Auth`, где происходит сопоставление пользователя по сессионной куке.
Возвращаются заголовки с данными залогиненного пользователя `X-Auth-User-ID`, `X-Auth-User-Login`.

`UserProfile` использует заголовок с UserID в запросе при получении и изменении данных профилей пользователей.

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