package tcp_client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/minmax1996/aoimdb/pkg/protocols"
)

type Client struct {
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func NewClient(host string) (*Client, error) {
	connection, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	c := Client{
		connection: connection,
		reader:     bufio.NewReader(connection),
		writer:     bufio.NewWriter(connection),
	}

	return &c, nil
}

// Send sends command string to establised connection
func (c *Client) Close() error {
	return c.connection.Close()
}

// Send sends command string to establised connection
func (c *Client) Send(name string, s ...string) error {
	_, err := c.writer.WriteString(name + " " + strings.Join(s, " ") + "\n")
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Client) AuthWithUserPassPair(user, pass string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Send("auth", user, pass); err != nil {
		return err
	}

	// make error channel
	ch := make(chan error, 1)
	go func() {
		//blocked read response
		data, err := c.reader.ReadString('\n')
		if err != nil {
			ch <- err
			return
		}

		var item protocols.Response
		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			ch <- err
			return
		}

		if item.AuthResponse.Message != "authenticated" {
			ch <- errors.New("not authenticated")
			return
		}

		ch <- nil
	}()

	//listen error channel or context.done
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) Get(key string) (*protocols.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Send("get", key); err != nil {
		return nil, err
	}

	// make error channel
	var resp *protocols.GetResponse
	ch := make(chan error, 1)
	go func() {
		//blocked read response
		data, err := c.reader.ReadString('\n')
		if err != nil {
			ch <- err
			return
		}

		var item protocols.Response
		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			ch <- err
			return
		}

		if item.GetResponse == nil {
			ch <- errors.New("not response")
			return
		}
		resp = item.GetResponse
		ch <- nil
	}()

	//listen error channel or context.done
	select {
	case err := <-ch:
		return resp, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) Set(key string, value string) (*protocols.SetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Send("set", key, value); err != nil {
		return nil, err
	}

	// make error channel
	var resp *protocols.SetResponse
	ch := make(chan error, 1)
	go func() {
		//blocked read response
		data, err := c.reader.ReadString('\n')
		if err != nil {
			ch <- err
			return
		}

		var item protocols.Response
		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			ch <- err
			return
		}

		if item.SetResponse == nil {
			ch <- errors.New("not response")
			return
		}
		resp = item.SetResponse
		ch <- nil
	}()

	//listen error channel or context.done
	select {
	case err := <-ch:
		return resp, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
