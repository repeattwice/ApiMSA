(Маркет пласе)
Тут я опишу эндпоинты и query, и тело http/json запроса, которые будут приходить на API сервис
во превых host-localhost
порт должен быть обьявлен либо в переменной окружения(DB_PORT), либо при запуске программы с флагом --port
Подключение к БД будет также через переменные окружения(host - DB_HOST, port - DB_PORT, user - DB_USER, password - DB_PASSWORD)

Ицициализация:
POST/CreateAccount    /gRPC # болодя
GET/Avtorizacion      /gRPC # болодя
DELETE/DeleteAccount  /gRPC # вадик
(все в теле запроса, отдельная бд аккаунтов)


и так, у нас есть такие EndPoint для :
GET/ShowAllItemsInCort  /gRPC  queryes(в праметре передаем в какю бд будем записывать, то что записываем в теле) # Санчез
POST/AddToCort  /kafka  queryes(в праметре передаем в какю бд будем записывать, то что записываем в теле) # болодя
PATCH/ChangePrice  /kafka  queryes(что какое поле изменяем, в теле то на что изменяем) # Санчез
DELETE/DeleteBuyFromKorzina  /kafka  queryes(какую строку удаляем, тело пустое) # вадику


Схема БД
users:
user_id(INTEGER, SERIAL PRIMARY KEY)    user_name(VARCHAR, NOT NULL)    last_name(VARCHAR, NOT NULL)    email(VARCHAR)

cort:
user_id()    item_name()    item_price()     

item:
item_name(VARCHAR, NOT NULL, UNIQUE)    item_price(INTEGER, NOT NULL)