package hlt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Connection struct {
	width, height int
	PlayerTag     int
	reader        *bufio.Reader
	writer        io.Writer
}

func (c *Connection) sendString(input string) {
	fmt.Println(input)
}

func (c *Connection) getString() string {
	retstr, _ := c.reader.ReadString('\n')
	retstr = strings.TrimSpace(retstr)
	return retstr
}

func (c *Connection) getInt() int {
	i, err := strconv.Atoi(c.getString())
	if err != nil {
		log.Println("Errored on initial tag: ", err)
	}
	return i
}

func NewConnection(botName string) Connection {
	conn := Connection{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
	}
	conn.PlayerTag = conn.getInt()
	sizeInfo := strings.Split(conn.getString(), " ")
	width, _ := strconv.Atoi(sizeInfo[0])
	height, _ := strconv.Atoi(sizeInfo[1])
	conn.width = width
	conn.height = height
	conn.sendString(botName)
	return conn
}

func (c *Connection) UpdateMap(turn int) GameMap {
	log.Println("--- NEW TURN ---", turn)
	gameString := c.getString()

	gameMap := GameMap{
		MyId:         c.PlayerTag,
		Turn:         turn,
		Width:        c.width,
		Height:       c.height,
		Planets:      []int{},
		Players:      [4]Player{},
		EnemyShips:   []int{},
		MyShips:      []int{},
		PlanetLookup: make(map[int]*Planet),
		ShipLookup:   make(map[int]*Ship),
	}
	gameMap.ParseGameString(gameString)
	log.Println("    Parsed map")
	return gameMap
}

func (c *Connection) SubmitCommands(commandQueue []string) {
	commandString := strings.Join(commandQueue, " ")
	log.Println("Final string :", commandString)
	c.sendString(commandString)
}
