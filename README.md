В данном API реализованые следующие возможности:

- Создание заказа:

        path:"/create_order"
        POST запрос, в качестве данных принимает JSON c двумя параметрами "dishes", "count"

        Пример запроса:
        response = requests.post(
            'http://localhost:8080/create_order',
            json={
                "dishes": ['1', '2', '50', '45'],
                "count": ["51", "10", "145", "20"]
            }
        )
        
        где "dishes" - список позиций из меню (список из номеров блюд),
            "count" - споисок с количеством порций.
        Данным dishes[0], coint[0] - соответствует блюдо и количестно первой позиции заказа.

        Результат: возвращает номер созданного заказа в формате {"order":"10"} или ошибка.
        
- Изменение статуса заказа( перевод заказа в статус "Готов" или "Выполнен"):

        path:"/change_order_status"
        POST запрос, в качестве данных принимает JSON c двумя параметрами "number", "status"

        Примеры запроса:
            response = requests.post(
                'http://localhost:8080/change_order_status',
                json={'number': '2', 'status': 'ready'}
            )  # перевод в новый статус "готов"
            
            response = requests.post(
                'http://localhost:8080/change_order_status',
                json={'number': '2', 'status': 'close'}
            )  # перевод в новый статус "выполнен"
        
        где "number" - номер заказа, 
            "status" - новый статус (принимает значения "ready" - готов, "close" - выполнен)

        Результат: ответ в формате {"id":2,"number":2,"status":"ready"} или ошибка.

- Получение списка заказов (ожидающих обработки, ожидающих выдачи):

        path:"/get_orders"
        GET запрос, в качестве параметров принимает "status"

        Примеры запроса:
            response = requests.get(
                'http://localhost:8080/get_orders',
                params={'status': "new"},
            )  # получение заказов ожидающих обработки
            
            response = requests.get(
                'http://localhost:8080/get_orders',
                params={'status': "ready"},
            )   # получение заказов ожидающих выдачи
        где "status" - статус (принимает значения "new" - новый, "ready" - готов)
        Параметры: список позиций из меню и количество
        Результат: ответ в формате [{"id":1,"number":1,"status":"new"},{"id":3,"number":3,"status":"new"}] или ошибка.

- Создание позиции меню (доп.):

        path:"/new_dish"
        POST запрос, в качестве данных принимает JSON c параметрами "name", "number"

        Пример запроса:
            response = requests.post(
                'http://localhost:8080/new_dish',
                json={
                    'name': ' weed cake', 
                    'number': '50'
                }
            )
        где "name" - название блюда, 
            "number" - номер блюда в меню
            
        результат: данные в формате {"id":7,"number":"50","name":" weed cake"} или ошибка