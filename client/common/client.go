package common

import (
	"encoding/csv"
	"fmt"
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
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	startLine := 0
	batchSize := c.config.MaxBatchAmount
	filePath := fmt.Sprintf("/data/agency-%v.csv", c.config.ID)

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		file, err := os.Open(filePath)
		if err != nil {
			log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
		}

		reader := csv.NewReader(file)

		for i := 0; i < startLine; i++ {
			_, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
			}
		}

		lines := make([][]string, 0, batchSize)
		for i := 0; i < batchSize; i++ {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
			}
			lines = append(lines, record)
		}
		file.Close()

		err = c.ProcessBatch(batchSize, lines)
		if err != nil {
			log.Criticalf("action: apuestas_enviadas | result: fail | cantidad: 0")
			return
		}
		startLine += batchSize

		c.conn.Close()

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) ProcessBatch(maxBatchSize int, lines [][]string) error {
	bets := make([]model.Bet, 0, maxBatchSize)
	for _, record := range lines {
		id, err := strconv.Atoi(record[2])
		if err != nil {
			log.Fatalf("failed to convert ID: %v", err)
			return err
		}
		number, err := strconv.Atoi(record[4])
		if err != nil {
			log.Fatalf("failed to convert number: %v", err)
			return err
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
		return err
	}
	err = protocol.SendMessage(c.conn, bytes.Bytes())
	if err != nil {
		return err
	}
	m, err := protocol.ReceiveMessage(c.conn)
	if err != nil || (len(m) != 1 && m[0] != 4) {
		return err
	}
	log.Infof("action: apuestas_enviadas | result: success | cantidad: %v", len(bets))
	return nil
}
