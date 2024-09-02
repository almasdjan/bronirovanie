# bronirovanie




1. **Запустите контейнеры:**

   Выполните следующую команду, чтобы создать и запустить контейнеры в фоновом режиме:
   docker-compose up --build -d

   создаются и запустятся три контейнера, go-приложение и две базы данных
   
   ссылка на swagger:
   http://localhost:8443/swagger/index.html

   для бронирования пример body:
    {
        "room_id": 1,
        "start_time": "2024-09-01T10:00:00Z",
        "end_time": "2024-09-01T12:00:00Z"
    }

2. **Запустите тесты:**

    Команда: make test

    ВЫВОД: 
    TestCreateReservation_Success - добавляет одну запись в таблицу в тестовой бд и выводит 

    TestCreateReservationWithConflict_Success - пытается добавить две пересекающие записи, добавляется только первый и выводит все записи в таблице 

    TestCreateReservation_ConcurrentRequests - попытается добавить несколько одинаковых записей в одно время 

    код теста в pkg/service/test 


