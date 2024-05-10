package mcstatus

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Client is used to make calls to the configured Minecraft server
type Client struct {
	Config
}

// Config is used to configure a Client. Address should be a
// fully-qualified domain name
type Config struct {
	Addr    string
	Timeout time.Duration
}

type status struct {
	Version           string
	Motd              string
	ActivePlayerCount int
	MaxPlayerCount    int
}

// New returns a Client configured to perform operations against the
// configured server
func New(config Config) (*Client, error) {
	if config.Addr == "" {
		return nil, fmt.Errorf("Config must include server address")
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}

	// test to ensure a connection can be established
	client := &Client{Config: config}
	_, err := client.status()
	if err != nil {
		return nil, err
	}
	return &Client{Config: config}, nil
}

// MaxPlayerCount return the number of players allowed to sign in on the
// configured Minecraft server.
func (c *Client) MaxPlayerCount() (int, error) {
	status, err := c.status()
	if err != nil {
		return 0, err
	}
	return status.MaxPlayerCount, nil
}

// GetActivePlayerCount returns the number of players currently signed in on
// the configured Minecraft server
//func (c *Client) ActivePlayerCount() (int, error) {
//	status, err := c.status()
//	if err != nil {
//		return 0, err
//	}
//	return status.ActivePlayerCount, nil
//}

// GetMotd returns the current "message of the day" for the configured
// Minecraft server
func (c *Client) Motd() (string, error) {
	status, err := c.status()
	if err != nil {
		return "", err
	}
	return status.Motd, nil
}

// GetVersion returns the current version of Minecraft running on the server
// for this client.
func (c *Client) Version() (string, error) {
	status, err := c.status()
	if err != nil {
		return "", err
	}
	return status.Version, nil
}

func (c *Client) status() (*status, error) {
	conn, err := net.DialTimeout("tcp", c.Config.Addr, c.Config.Timeout)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte("\xFE\x01"))
	if err != nil {
		return nil, err
	}

	rawData := make([]byte, 512)
	_, err = conn.Read(rawData)
	if err != nil {
		return nil, err
	}

	if rawData == nil || len(rawData) == 0 {
		return nil, fmt.Errorf("no data returned from server")
	}

	data := strings.Split(string(rawData[:]), "\x00\x00\x00")

	fmt.Println(data)

	if len(data) < 6 {
		return nil, fmt.Errorf("invalid data returned: %s", data)
	}

	apc, err := strconv.Atoi(data[4])
	if err != nil {
		return nil, err
	}

	//	mpc, err := strconv.Atoi(data[5])
	//	if err != nil {
	//		return nil, err
	//	}

	status := &status{
		Version:           data[2],
		Motd:              data[3],
		ActivePlayerCount: apc,
		//		MaxPlayerCount:    mpc,
	}

	return status, nil
}
