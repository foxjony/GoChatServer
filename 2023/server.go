// packagemain #20: Building a TCP Chat in Go
// https://www.youtube.com/watch?v=Sphme0BqJiY
// https://github.com/plutov/packagemain/tree/master/20-tcp-chat
// pkgmgr /iu:"TelnetClient"		(включить telnet в Windows 10) 		telnet 178.238.225.218 50013
// Run Windows => ./chat
// (c) foxjony, 02.02.2023

// Миллион WebSocket и Go
// https://habr.com/ru/company/mailru/blog/331784/

/*			=== Linux Build ===
~# cd ../home/go 			- Перейти в директорию
/home/go# go build . 		- Скомпилировать исходник server.go исполняемый файл *go
/home/go# ./go 				- Запустить сервер (исполняемый файл *go)
2021/03/28 10:22:03 TCP Server Started on port: 50013
[Ctrl+C] 					- Остановить сервер
*/

/*			=== Linux Auto Start ===
cd /usr/bin 							– Перейти в директорию /usr/bin
nano run-server.sh 						– Создаем файл run-server.sh с текстом:
#!/bin/bash
cd /home/go
./go
[Ctrl+s] 								– Сохранить изменения в тексте
[Ctrl+x] 								– Выход из редактора

chmod ugo+x run-server.sh 				– Сменить права файла
run_server.sh							– Зауск сервера по созданному скрипту (проверка)

systemctl								– Показывает список запущенных служб
cd /usr/lib/systemd/system				– Переходим в директорию system
nano run-server.service					– Создаем файл run-server.service с текстом:
[Unit]
Description=Run Server
After=multi-user.target
[Service]
Type=forking
ExecStart=/usr/bin/run-server.sh
Restart=always
[Install]
WantedBy=multi-user.target
[Ctrl+s] 								– Сохранить изменения в тексте
[Ctrl+x] 								– Выход из редактора

systemctl daemon-reload					– После изменений необходимо перечитать изменения
systemctl enable run-server				– Разрешить автозапуск сервиса run-server
systemctl -q is-active run-server		– В crontab добавить команду стобы сервер сделал 
										  перезапуск сервиса в случае его остановки
systemctl start run-server				– Запустить сервис run-server
systemctl status run-server				– Cмотрим статус сервиса run-server
(systemctl stop run-server				– Остановить сервис run-server)
reboot naw								– Перезапуск Linux
*/

/*
/user <name> 	- Set a Name, otherwise user will stay Guest.
				  Задайте имя, иначе пользователь останется гостем.
/order <name> 	- Join a Order or Room.
					(if doesn't exist, the will be created)
					(user can be only in one Order or Room at the same time)
				  Присоедиться к Заказу или Комнате.
					(если не существует, будет создан)
					(пользователь одновременно может находиться только в одном Заказе или Комнате)
/list 			- Show List of available Orders and Rooms.
				  Показать Список доступных Заказов и Комнат.
/msg <msg> 		- Broadcast message to everyone in a Order or Room.
				  Отправить сообщение всем в Заказе или Комнате.
/quit 			- Disconnect from the Chat Server.
				  Отключиться от Чат Сервера.
*/

package main

import (
	"net"
	"bufio"
	"strings"
	"log"
	"fmt"
	"time"
	"runtime"
	"runtime/debug"
)

/* =========================================================
type struct: 
	command
	client
	order
	server
========================================================= */

var big = make([]byte, 1<<20)		// Allocate 1 MB

type commandID int

const (
	CMD_USER commandID = iota		// 0
	CMD_ORDER	                 	// 1
	CMD_LIST						// 2
	CMD_MSG							// 3
	CMD_QUIT						// 4
	CMD_DEL							// 5
	CMD_MEM							// 6
	CMD_PING						// 7
)

type command struct {
	id     commandID				// 0-6
	client *client
	args   []string
}

// conn.RemoteAddr().String() = "192.168.0.193:52626"
// conn.RemoteAddr().String() = "192.168.0.193:52628"
// conn: &{{%!s(*net.netFD=&{{{0 0 0} 4 {281472834653752} <nil> 0 0 true true false} 10 1 false tcp 0x400008a630 0x400008a660})}}
// conn: &{{%!s(*net.netFD=&{{{0 0 0} 6 {281472834653544} <nil> 0 0 true true false} 10 1 false tcp 0x40000e4030 0x40000e4060})}}
type client struct {
	conn     net.Conn
	user     string
	order    *order
	commands chan <- command
}

type order struct {
	name    string
	members map[net.Addr]*client
}

type server struct {
	orders   map[string]*order
	commands chan command
}

/* =========================================================
func:
	broadcast
	readInput
	err
	msgC
	newServer
	run
	newClient
	user
	order
	list
	msgS
	quit
	quitCurrentOrder
	del
	mem
	ping
	main
========================================================= */

func (r *order) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {m.msgC(msg)}
	}
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {return}
		msg = strings.Trim(msg, "\n")
		//log.Printf("MESSAGE %s", msg)
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		switch cmd {
			case "/user":  	c.commands <- command{id: CMD_USER,  client: c, args: args}
			case "/order":  c.commands <- command{id: CMD_ORDER, client: c, args: args}
			case "/list": 	c.commands <- command{id: CMD_LIST,  client: c,}
			case "/msg":   	c.commands <- command{id: CMD_MSG,   client: c, args: args}
			case "/quit":  	c.commands <- command{id: CMD_QUIT,	 client: c,}
			case "/del":  	c.commands <- command{id: CMD_DEL,   client: c,}
			case "/mem":  	c.commands <- command{id: CMD_MEM,   client: c,}
			case "/ping":  	c.commands <- command{id: CMD_PING,  client: c, args: args}
			default:		c.err(fmt.Errorf("Unknown Command %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("Error: " + err.Error() + "\n"))
}

func (c *client) msgC(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func newServer() *server {
	return &server{
		orders:   make(map[string]*order),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
			case CMD_USER:	s.user(cmd.client, cmd.args)
			case CMD_ORDER:	s.order(cmd.client, cmd.args)
			case CMD_LIST:	s.list(cmd.client)
			case CMD_MSG:	s.msgS(cmd.client, cmd.args)
			case CMD_QUIT:	s.quit(cmd.client)
			case CMD_DEL:	s.del(cmd.client)
			case CMD_MEM:	s.mem(cmd.client)
			case CMD_PING:	s.ping(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("Connected New Client: %s (Guest)", conn.RemoteAddr().String())
	c := &client{conn: conn, user: "Guest", commands: s.commands,}
	c.readInput()
}

func (s *server) user(c *client, args []string) {
	if len(args) < 2 {
		c.msgC("Error, usage: /user NAME")
		return
	}
	c.user = args[1]
	c.msgC(fmt.Sprintf("Hi, %s", c.user))
	log.Printf("%s (Guest) = %s", c.conn.RemoteAddr().String(), c.user)
}

func (s *server) order(c *client, args []string) {
	if len(args) < 2 {
		c.msgC("Error, usage: /order ORDER_NAME")
		return
	}
	orderName := args[1]
	o, ok := s.orders[orderName]
	if !ok {
		o = &order{name: orderName, members: make(map[net.Addr]*client),}
		s.orders[orderName] = o
	}
	o.members[c.conn.RemoteAddr()] = c
	s.quitCurrentOrder(c)
	c.order = o
	o.broadcast(c, fmt.Sprintf("%s Connected", c.user))
	c.msgC(fmt.Sprintf("Connected to Order %s", orderName))
	log.Printf("Connected to Order %s", orderName)
}

func (s *server) list(c *client) {
	var orders []string
	for name := range s.orders {
		orders = append(orders, name)
	}
	c.msgC(fmt.Sprintf("Orders: %s", strings.Join(orders, ", ")))
}

func (s *server) msgS(c *client, args []string) {
	if len(args) < 2 {
		c.msgC("Error, usage: /msg MSG")
		return
	}
	msg := strings.Join(args[1:], " ")
	log.Printf("ORDER %s, IDUSER %s, MSG %s", c.order.name, c.user, msg)
	c.order.broadcast(c, c.user+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("Disconnected Client: %s", c.conn.RemoteAddr().String())
	s.quitCurrentOrder(c)
	c.msgC("Quit Ok")
	c.conn.Close()
	runtime.GC()			// Go Clear Memory
}

func (s *server) quitCurrentOrder(c *client) {
	if c.order != nil {
		oldOrder := s.orders[c.order.name]
		delete(s.orders[c.order.name].members, c.conn.RemoteAddr())
		oldOrder.broadcast(c, fmt.Sprintf("%s Out the Order", c.user))
		log.Printf("%s Out the Order", c.user)
	}
}

func (s *server) del(c *client) {
	// Измерение времени исполнения в Golang
	// https://golang-blog.blogspot.com/2020/04/measure-execution-time-in-golang.html
	start := time.Now()
	//log.Printf("Del")						// 875ns (Пусто)	// 84.002µs (Код измерения)
	//log.Printf("Del Del Del Del Del Del Del Del Del Del ")	// 90.419µs (Код измерения)
	duration := time.Since(start)			// Продолжительность (duration)
	fmt.Println(duration)					// Отформатированная строка, например, "2h3m0.5s" или "4.503μs"
	//fmt.Println(duration.Nanoseconds())	// Nanoseconds как int64 (84002)
}

func (s *server) mem(c *client) {
	/*
	https://ru.stackoverflow.com/questions/1184546/Освобождение-памяти-в-go
	https://coderoad.ru/24863164/Как-анализировать-golang-памяти

	Количество свободной памяти можно узнать, вызвав функцию ReadMemStats. 
	В выходном параметре типа MemStats поле HeapIdle содержит число байтов в чанках памяти, 
	в которых нет объектов Go. 
	То есть это та память, которую, теоретически, можно отдать обратно операционной системе.
	*/

	var ms1, ms2 runtime.MemStats
    big = nil 						// Drop the link to 1 MB (1 куча = 16 MB, 16*66.4 = 1 GB)
    runtime.GC()					// Force GC
    runtime.ReadMemStats(&ms1)
    debug.FreeOSMemory()			// Force memory release (Принудительное освобождение памяти)
    runtime.ReadMemStats(&ms2)
    fmt.Println(" ")
    fmt.Println("Users:                 ", len(c.user))
    fmt.Println("Rooms:                 ", len(s.orders))
    fmt.Println("1MB in bytes:          ", 1<<20)
    fmt.Println("Idle memory before:    ", ms1.HeapIdle)		// Куча Занято До
    fmt.Println("Idle memory after:     ", ms2.HeapIdle)		// Куча Занято После
    fmt.Println("Idle memory delta:     ", int64(ms2.HeapIdle)-int64(ms1.HeapIdle))
    fmt.Println("Released memory before:", ms1.HeapReleased)	// Куча Свободно До
    fmt.Println("Released memory after: ", ms2.HeapReleased)	// Куча Свободно После
    fmt.Println("Released memory delta: ", ms2.HeapReleased - ms1.HeapReleased)

	/*
	1MB in bytes:  			 1048576
	Idle memory before:  	66404352
	Idle memory after:  	66404352
	Idle memory delta:  	       0
	Released memory before: 65241088
	Released memory after:  66404352
	Released memory delta:   1163264

	Idle memory before:     66404352
	Idle memory after:      66437120
	Idle memory delta:         32768
	Released memory before: 65273856
	Released memory after:  66437120
	Released memory delta:   1163264

	Idle memory before:     66437120
	Idle memory after:      66404352
	Idle memory delta:        -32768
	Released memory before: 65273856
	Released memory after:  66371584
	Released memory delta:   1097728
	*/
}

func (s *server) ping(c *client, args []string) {
	// fmt.Println("Ping ", c.user)
	if len(args) < 2 {
		c.msgC("Error, usage: /pimg gps")
		return
	}
	//c.user = args[1]

	args2 := strings.Split(args[1], "|")
	lat := strings.TrimSpace(args2[0])
	lon := strings.TrimSpace(args2[0])
	alt := strings.TrimSpace(args2[0])		// Altitude
	acc := strings.TrimSpace(args2[0])		// Accuracy

	log.Printf("user: %s, gps: %s, %s, %s, %s", c.user, lat, lon, alt, acc)

	c.msgC(fmt.Sprintf("pong"))
}
    
func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":50013")
	if err != nil {log.Fatalf("Start Server Error: %s", err.Error())}
	defer listener.Close()
	log.Printf("TCP Server Started on port: 50013")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection Error: %s", err.Error())
			continue
		}
		go s.newClient(conn)
	}
}
