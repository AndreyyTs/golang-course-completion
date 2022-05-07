В этом задании вы научитесь строить микросервис на базе фреймворка grpc

*В это задании нельзя использовать глобальные переменные, нужное вам храните в полях структуры, которая живёт в замыкании*

Вам потребуется реализовать:

* Сгенерировать необходимый код из proto-файла
* Базу микросервиса в возможностью остановки сервера
* ACL - контроль доступа от разных клиентов
* Систему логирования вызываемых методов
* Систему сбора сборки статистики ( просто счетчики ) по вызываемым методам

Микросервис будет состоять из 2-х частей:
* Какая-то бизнес-логика. В нашем примере она ничего не делает, её достаточно просто вызывать
* Модуль администрирования, где находится логирование и статистика

С первым всё просто, там логики нету.

Со вторым интереснее. Как правило в настоящих микросервисах и логирование, и статистика работают в единственном экземпляре, у нас же они будут доступны по потоковому ( streaming ) интерфейсу тому, кто подключится к сервису. Это значит, что к сервису может подключиться 2 клиента логирования и оба будут получать поток логов. Так же к сервису может подключиться 2 ( и более ) модуля статистики с разными интервалами получения статистики ( например, каждые 2, 3 и 5 секунд ) и она будет асинхронно отправляться по каждому интерфейсу.

Раз уж был упомянут асинхрон - в задании будут горутины, таймеры, мютексы, контекст с таймаутами/завершением.

Особенности задания:

Содержимое файлов service.pb.go и service_grpc.pb.go (которые получилось у вас при генерации proto-файла) вам необходимо поместить в service.go для загрузки 1 файлом
В этом задании нельзя использовать глобальные переменные. Всё что нам необходимо - храните в полях структуры.

Запускать тесты с go test -v -race

Используемые версии:
* libprotoc 3.19.3
* protoc-gen-go v1.27.1
* protoc-gen-go-grpc 1.2.0

Все гошные зависимости есть в папке vendor, команды для установки плагинов и генерации есть в Makefile