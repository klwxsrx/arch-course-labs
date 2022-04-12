# Варианты решения
## RESTful-подход
![docs/services.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/06-order-service/docs/restful.puml)

[OpenAPI-схема](https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/06-order-service/docs/restful.openapi.yaml)



## Event Notifications
![docs/services.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/06-order-service/docs/event_notifications.puml)

[OpenAPI-схема](https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/06-order-service/docs/event_notifications.openapi.yaml)

[AsyncAPI-схема](https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/06-order-service/docs/event_notifications.asyncapi.yaml)

В событиях передаются только минимальные необходимые данные, например `userID`, `orderID` и т.п. Если сервису понадобится дополнительная информация, то он выполняет GET-запрос в сервис.

Таким образом для отправки email-пользователю необходимо получить информацию по заказу из `Order` (например чтобы в письме прописать наименования товаров), email-пользователя из `User`.

## Event Collaboration
![docs/services.puml](http://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/06-order-service/docs/event_collaboration.puml)

[AsyncAPI-схема](https://raw.githubusercontent.com/klwxsrx/arch-course-labs/master/06-order-service/docs/event_collaboration.asyncapi.yaml)

Для осуществления полностью событийного взаимодействия необходимо иметь локальные копии данных модели из других сервисов.

Например `Notification` должен слушать события добавления/изменения/удаления пользователей в `User` для того чтобы иметь актуальные данные email пользователей.

Аналогично `Notification` необходимо слушать события `OrderCreated` для того чтобы иметь актуальное содержимое заказа для последующей отправки письма по заказу.

Для поддержания согласованности данных между сервисами необходимо гарантировать доставку событий до всех потребителей и их обработку.

## Приемлемое решение
Зависит от потребности.

Полностью синхронный подход быстрее и проще в разработке, но хуже всего справится с высокой нагрузкой.

Второй и третий подходы лучше справятся с нагрузкой, однако сложнее в разработке. Третий подход сложнее и более критичен к согласованнооси данных между сервисами.

В рамках лабораторных работ можно выбрать любой вариант из трех, например первый.
