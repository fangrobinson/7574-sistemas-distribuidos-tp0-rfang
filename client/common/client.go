package common

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/protocol"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/serialization"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	MaxBatchAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

func (c *Client) Shutdown() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) SendBets() ([]model.Bet, error) {
	startLine := 0
	batchSize := c.config.MaxBatchAmount
	filePath := fmt.Sprintf("/data/agency-%v.csv", c.config.ID)

	allBets := make([]model.Bet, 0)

	for {
		file, err := os.Open(filePath)
		if err != nil {
			log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
			return nil, err
		}
		reader := csv.NewReader(file)

		// Skip already processed lines
		for i := 0; i < startLine; i++ {
			_, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					file.Close()
					return allBets, nil
				}
				log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
				file.Close()
				return nil, err
			}
		}

		lines := make([][]string, 0, batchSize)
		for i := 0; i < batchSize; i++ {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
				file.Close()
				return nil, err
			}
			lines = append(lines, record)
		}
		file.Close()

		err = c.createClientSocket()
		if err != nil {
			log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
			return nil, err
		}

		bets, err := c.ProcessBatch(batchSize, lines)
		c.conn.Close()
		if err != nil {
			log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
			return nil, err
		}

		allBets = append(allBets, bets...)
		startLine += batchSize
	}
}

func (c *Client) SendNoMoreBets() error {
	for {
		err := c.createClientSocket()
		if err != nil {
			log.Criticalf("action: no_more_bests | result: connection_failed")
			return err
		}

		bytes, err := serialization.EncodeNoMoreBets(c.config.ID)
		if err != nil {
			continue
		}
		err = protocol.SendMessage(c.conn, bytes.Bytes())
		if err != nil {
			log.Infof("action: no_more_bets | result: not sent")
			c.conn.Close()
			return err
		}
		log.Infof("action: no_more_bets | result: sent")

		time.Sleep(c.config.LoopPeriod)

		m, err := protocol.ReceiveMessage(c.conn)
		if err != nil {
			log.Infof("action: no_more_bets | result: fail")
			c.conn.Close()
			return err
		}
		if len(m) == 1 && m[0] == serialization.NO_MORE_BETS_ACK {
			log.Infof("action: no_more_bets | result: success")
			break
		}
		c.conn.Close()
	}
	return nil
}

func (c *Client) PollWinner() ([]int, error) {
	for {
		// Create the connection the server in every loop iteration.
		c.createClientSocket()
		bytes, err := serialization.EncodeGetWinners(c.config.ID)
		if err != nil {
			c.conn.Close()
			continue
		}
		err = protocol.SendMessage(c.conn, bytes.Bytes())
		if err != nil {
			return nil, err
		}

		m, err := protocol.ReceiveMessage(c.conn)
		if err != nil {
			return nil, err
		}
		if len(m) == 1 && m[0] == serialization.WAIT {
			// Even though it's the same
			// this case it's exclicitely
			// waiting
			c.conn.Close()
			log.Infof("action: poll_winner | result: success | answer: wait")
			time.Sleep(c.config.LoopPeriod)
			continue
		}

		if len(m) >= 3 && m[0] == serialization.WINNERS {
			c.conn.Close()
			winnersAmountInt := int(binary.BigEndian.Uint16(m[1:3]))
			winnersIDs := make([]int, 0)
			for i := 0; i < winnersAmountInt; i++ {
				idStart := 3 + i*4
				idEnd := idStart + 4
				winnerId := int(binary.BigEndian.Uint32(m[idStart:idEnd]))
				winnersIDs = append(winnersIDs, winnerId)
			}
			log.Infof("action: poll_winner | result: success | answer: %v, %v", winnersAmountInt, len(winnersIDs))
			return winnersIDs, nil
		}

		c.conn.Close()
	}
}

func CheckWinners(bets []model.Bet, winners []int) {
	winners_amount := 0
	if len(winners) > 0 {
		winnersMap := make(map[int]bool)
		for _, id := range winners {
			winnersMap[id] = true
		}
		for _, bet := range bets {
			if winnersMap[bet.ID] {
				winners_amount++
			}
		}
	}
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", winners_amount)
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	bets, err := c.SendBets()
	if err != nil {
		log.Infof("action: send_bets | result: fail | client_id: %v", c.config.ID)
		return
	} else {
		log.Infof("action: send_bets | result: success | client_id: %v", c.config.ID)
	}
	c.conn.Close()
	err = c.SendNoMoreBets()
	if err != nil {
		log.Criticalf("action: error | %v", err)
		return
	}
	winners, err := c.PollWinner()
	if err != nil {
		return
	}
	CheckWinners(bets, winners)
	log.Infof("action: run_agency | result: success | client_id: %v", c.config.ID)
}

func (c *Client) ProcessBatch(maxBatchSize int, lines [][]string) ([]model.Bet, error) {
	bets := make([]model.Bet, 0, maxBatchSize)
	for _, record := range lines {
		id, err := strconv.Atoi(record[2])
		if err != nil {
			log.Fatalf("failed to convert ID: %v", err)
			return nil, err
		}
		number, err := strconv.Atoi(record[4])
		if err != nil {
			log.Fatalf("failed to convert number: %v", err)
			return nil, err
		}

		bets = append(bets, model.Bet{
			Name:      record[0],
			Surname:   record[1],
			ID:        id,
			BirthDate: record[3],
			Number:    number,
		})
	}

	bytes, err := serialization.EncodeMultipleBets(c.config.ID, bets)
	if err != nil {
		return nil, err
	}

	err = protocol.SendMessage(c.conn, bytes.Bytes())
	if err != nil {
		return nil, err
	}

	m, err := protocol.ReceiveMessage(c.conn)
	if err != nil || (len(m) != 1 && m[0] != 4) {
		return nil, err
	}

	log.Infof("action: apuestas_enviadas | result: success | cantidad: %v", len(bets))
	return bets, nil
}
