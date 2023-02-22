// === Go TCP Chat Server ===
// (c) Fork: foxjony, 03.02.2023

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

// ======================================================

var big = make([]byte, 1<<20)		// Allocate 1 MB

type commandID int

const (
	CMD_USER commandID = iota	// 0
	CMD_ORDER	                 // 1
	CMD_LIST			// 2
	CMD_MSG				// 3
	CMD_QUIT			// 4
	CMD_DEL				// 5
	CMD_MEM				// 6
	CMD_PING			// 7
)

type command struct {
	id     commandID		// 0-7
	client *client
	args   []string
}

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

// ======================================================

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
			case "/quit":  	c.commands <- command{id: CMD_QUIT,  client: c,}
			case "/del":  	c.commands <- command{id: CMD_DEL,   client: c,}
			case "/mem":  	c.commands <- command{id: CMD_MEM,   client: c,}
			case "/ping":  	c.commands <- command{id: CMD_PING,  client: c, args: args}
			default:	c.err(fmt.Errorf("Unknown Command %s", cmd))
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
	// Вимір часу виконання в Golang
	// https://golang-blog.blogspot.com/2020/04/measure-execution-time-in-golang.html
	start := time.Now()
	//log.Printf("Del")						// 875ns (Пусто)	// 84.002µs (Код виміру)
	//log.Printf("Del Del Del Del Del Del Del Del Del Del ")	// 90.419µs (Код виміру)
	duration := time.Since(start)		// Тривалість (duration)
	fmt.Println(duration)			// Відформатована строка, наприклад, "2h3m0.5s" або "4.503μs"
	//fmt.Println(duration.Nanoseconds())	// Nanoseconds як int64 (84002)
}

func (s *server) mem(c *client) {
	/*
	https://ru.stackoverflow.com/questions/1184546/Освобождение-памяти-в-go
	https://coderoad.ru/24863164/Как-анализировать-golang-памяти

	Кількість вільної пам'яті можна дізнатися, викликавши функцію ReadMemStats.
	У вихідному параметрі типу MemStats поле HeapIdle містить число байтів у чанках пам'яті,
	у яких немає об'єктів Go.
	Тобто це та пам'ять, яку, теоретично, можна віддати назад операційній системі.
	*/

	var ms1, ms2 runtime.MemStats
	big = nil 				// Drop the link to 1 MB (1 куча = 16 MB, 16*66.4 = 1 GB)
	runtime.GC()				// Force GC
	runtime.ReadMemStats(&ms1)
	debug.FreeOSMemory()			// Force memory release (Примусове вивільнення пам'яті)
	runtime.ReadMemStats(&ms2)
	fmt.Println(" ")
	fmt.Println("Users:                 ", len(c.user))
	fmt.Println("Rooms:                 ", len(s.orders))
	fmt.Println("1MB in bytes:          ", 1<<20)
	fmt.Println("Idle memory before:    ", ms1.HeapIdle)			// Куча Зайнято До
	fmt.Println("Idle memory after:     ", ms2.HeapIdle)			// Куча Зайнято Після
	fmt.Println("Idle memory delta:     ", int64(ms2.HeapIdle)-int64(ms1.HeapIdle))
	fmt.Println("Released memory before:", ms1.HeapReleased)		// Куча Свобідно До
	fmt.Println("Released memory after: ", ms2.HeapReleased)		// Куча Свобідно Після
	fmt.Println("Released memory delta: ", ms2.HeapReleased - ms1.HeapReleased)

	/*
	1MB in bytes:  		1048576
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
	if len(args) < 2 {
		c.msgC("Error, usage: /pimg gps")
		return
	}

	gps := strings.Split(args[1], "|")
	if len(gps) == 5 {
		if len(gps[0]) > 1 {
			//lat := strings.TrimSpace(args_gps[0])		// Lat
			//lon := strings.TrimSpace(args_gps[1])		// Lon
			//alt := strings.TrimSpace(args_gps[2])		// Altitude
			//acc := strings.TrimSpace(args_gps[3])		// Accuracy
			//spd := strings.TrimSpace(args_gps[4])		// Speed
			//log.Printf("user: %s, gps: %s, %s, %s, %s, %s", c.user, lat, lon, alt, acc, spd)
			log.Printf("iduser: %s, gps: %s, %s, %s, %s, %s", c.user, gps[0], gps[1], gps[2], gps[3], gps[4])
		}
	}

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
