# gRPC сервис для загрузки thumbnail'ов с видеороликов YouTube
## Инструкция по использованию
1. Перед запуском поднять инстанс `Redis`, который будет использоваться в качестве кэша для thumbnail'ов командой из `Makefile` - `make redis`.
2. Запустить сервер, который находится в папке `server` командой `go run server.go`.
3. Запустить клиента из папки `client` командой `go run client.go`. Варианты использования:
- Флаг `--urls`. Передача ссылок, по которым нужно скачать thumbnail'ы. Пример: `go run client.go --urls=https://youtu.be/Ilg3gGewQ5U,https://www.youtube.com/watch?v=ueAoUtacdNw`.
- Флаг `--async`. Позволяет скачивать большое количество файлов асинхронно.
- Флаг `-f`. Позволяет скачивать thumbnail'ы из всех видео, ссылки которых будут указаны в файле. В качестве примера можно использовать файл `videos.txt` в этом репозитории. Пример: `go run client.go --async -f videos.txt`

## Результат
При запуске сервера (!) будет создана папка `images`, которая будет содержать в себе все скаченные thubmnail'ы.
