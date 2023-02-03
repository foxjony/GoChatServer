
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
