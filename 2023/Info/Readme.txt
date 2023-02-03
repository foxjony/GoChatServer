
	=== Go TCP Chat Server (178.238.225.218 port 50013) ===

https://www.youtube.com/watch?v=Sphme0BqJiY
https://github.com/plutov/packagemain/tree/master/20-tcp-chat


			=== Команди Сервера ===

/user <name> 		- Set a Name, otherwise user will stay Guest.
			  Задати ім'я, інакше користувач залишиться гостем.
/order <name> 		- Join a Order or Room.
				(if doesn't exist, the will be created)
				(user can be only in one Order or Room at the same time)
			  Приєднатися до Замовлення чи Кімнати.
				(якщо не існує, буде створено)
				(користувач одночасно може знаходитися тільки в одному Замовленні чи Кімнаті)
/list 			- Show List of available Orders and Rooms.
			  Показати Список доступних Замовлень чи Кімнат.
/msg <msg> 		- Broadcast message to everyone in a Order or Room.
			  Відправити повідомлення всім в Замовленні чи Кімнаті.
/quit 			- Disconnect from the Chat Server.
			  Відключитися від Чат Сервера.
/ping <0|0|0|0|0> 	- Ping to maintain communication with the server and transfer gps coordinates
			  Ping для підтримки зв'язку з сервером та передача gps координат 
			  "Lat|Lon|Altitude|Accuracy|Speed"


			=== Linux Build ===
~# cd ../home/go 			- Перейти в директорію
/home/go# go build . 			- Зкомпілювати ісходник server.go в файл *go
/home/go# ./go 				- Запустити сервер (файл *go)
2021/03/28 10:22:03 TCP Server Started on port: 50013
[Ctrl+C] 				- Зупинити сервер


			=== Linux Auto Start ===
cd /usr/bin 				– Перейти в директорію /usr/bin
nano run-server.sh 			– Створити файл run-server.sh з текстом:
#!/bin/bash
cd /home/go
./go
[Ctrl+s] 				– Зберегти зміни в тексті
[Ctrl+x] 				– Вихід з редактора

chmod ugo+x run-server.sh 		– Змінити права файла
run_server.sh				– Запуск сервера по створеному скрипту (перевірка)

systemctl				– Показує список запущених служб
cd /usr/lib/systemd/system		– Переходимо в директорію system
nano run-server.service			– Створюєм файл run-server.service з текстом:
[Unit]
Description=Run Server
After=multi-user.target
[Service]
Type=forking
ExecStart=/usr/bin/run-server.sh
Restart=always
[Install]
WantedBy=multi-user.target
[Ctrl+s] 				– Зберегти зміни в тексті
[Ctrl+x] 				– Виход з редактора

systemctl daemon-reload			– Після змін необхідно перезапустити службу
systemctl enable run-server		– Дозволити автозапуск сервісу run-server
systemctl -q is-active run-server	– В crontab добавити команду щоб сервер зробив 
					  перезапуск сервіса в випадку його зупинки
systemctl start run-server		– Запустити сервіс run-server
systemctl status run-server		– Дивимось статус сервіса run-server
(systemctl stop run-server		– Зупинити сервіс run-server)
reboot naw				– Перезапуск Linux


			=== Windows Setup ===
pkgmgr /iu:"TelnetClient"		(ввімкнути telnet в Windows 10) 		telnet 178.238.225.218 50013
Run Windows => ./chat


			=== Додатково ===
Миллион WebSocket и Go
https://habr.com/ru/company/mailru/blog/331784/
