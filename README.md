# Go Chat Server

Что я думаю про язык GO?
- https://www.youtube.com/watch?v=eIiyTq4VHK4

Изучение Go в одном видео уроке за 30 минут!
- https://www.youtube.com/watch?v=pfmxPtLIW34
- https://itproger.com/course/one-lesson/14

Алексей Акулович — Плюсы и минусы Go, а также его применения в ВКонтакте
- https://www.youtube.com/watch?v=2fxNbhy2gt0

Работа с сетью в Go. Алексей Акулович, Вконтакте.
- https://www.youtube.com/watch?v=p1ILhiq5Clw

Продуктовая разработка на Go: история одного проекта. Максим Рындин, Gett.
- https://www.youtube.com/watch?v=ppnnuDotxZM

Чат на Go (часть 1)
- https://habr.com/ru/post/306840/

Чат на iOS: используем сокеты
- https://habr.com/ru/post/467909

GoLang Terminal Messenger
- https://github.com/Shivam010/GoLang-Terminal-Messenger

Golang - Xây dựng ứng dụng chat Server-Client với TCP
- https://www.youtube.com/watch?v=ibvtI3PSyno
- https://github.com/code4func/golang-chat-tcp

За что ругают Golang и как с этим бороться?
- https://habr.com/ru/post/282588/

6 рекомендаций по разработке безопасных Go-приложений
- https://habr.com/ru/company/ruvds/blog/484614/

Привіт. Ласкаво просимо до туру з мови програмування Go.
- https://go-tour-ua.appspot.com
- https://go-tour-ua.appspot.com/list

- https://golang.org
- https://play.golang.org
- https://play.golang.org/p/jeKYspgp8wM
- https://ru.wikipedia.org/wiki/Go
- https://github.com/golang
- https://github.com/djherbis/nio

Отборный список фреймворков, библиотек и программного обеспечения на языке Go с открытым исходным кодом
- https://github.com/avelino/awesome-go
- https://github.com/avelino/awesome-go#server-applications

Разбираемся с новым sync.Map в Go 1.9
- https://habr.com/ru/post/338718/

Олег Облеухов — Как сделать высоконагруженный сервис, не зная количество нагрузки
- https://www.youtube.com/watch?v=hy5OSruLqvU

Миллион WebSocket и Go
- https://habr.com/ru/company/mailru/blog/331784/

Room Chat - packagemain #20: Building a TCP Chat in Go
- https://www.youtube.com/watch?v=Sphme0BqJiY
- https://github.com/plutov/packagemain/tree/master/20-tcp-chat

# Linux Build
- ~# cd /home/go        `- Перейти в директорию /home/go`
- /home/go# go build .  `- Скомпилировать исходник server.go исполняемый файл *go`
- /home/go# ./go        `- Запустить сервер (исполняемый файл *go)`
- 2021/03/28 10:22:03 TCP Server Started on port: 3000
- [Ctrl+C]              `- Остановить сервер`

# Linux Auto Start
- cd /usr/bin 		    `– Перейти в директорию /usr/bin`
- nano run-server.sh 	`– Создаем файл run-server.sh с текстом:`
- #!/bin/bash
- cd /home/go
- ./go
- [Ctrl+s] 		`– Сохранить изменения в тексте`
- [Ctrl+x] 		`– Выход из редактора`

- chmod ugo+x run-server.sh 	`– Сменить права файла`
- run_server.sh			          `– Зауск сервера по созданному скрипту (проверка)`

- systemctl			              `– Показывает список запущенных служб`
- cd /usr/lib/systemd/system	`– Переходим в директорию system`
- nano run-server.service	    `– Создаем файл run-server.service с текстом:`
```php
[Unit]
Description=Run Server
After=multi-user.target
[Service]
Type=forking
ExecStart=/usr/bin/run-server.sh
Restart=always
[Install]
WantedBy=multi-user.target
```
- [Ctrl+s] 			`– Сохранить изменения в тексте`
- [Ctrl+x] 			`– Выход из редактора`

- systemctl daemon-reload	    `– После изменений необходимо перечитать изменения`
- systemctl enable run-server	`– Разрешить автозапуск сервиса run-server`
- systemctl -q is-active run-server	`– В crontab добавить команду стобы сервер сделал перезапуск сервиса в случае его остановки`
- systemctl start run-server		`– Запустить сервис run-server`
- systemctl status run-server		`– Cмотрим статус сервиса run-server`
- systemctl stop run-server		  `– Остановить сервис run-server`
- reboot naw				            `– Перезапуск Linux`
