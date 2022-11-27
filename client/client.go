package client

import (
	"bufio"
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/minmax1996/aoimdb/api/msg_protocol"
	"github.com/vmihailenco/msgpack"
)

type TcpClient struct {
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func NewTcpClient(host string) (*TcpClient, error) {
	connection, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	c := TcpClient{
		connection: connection,
		reader:     bufio.NewReader(connection),
		writer:     bufio.NewWriter(connection),
	}

	return &c, nil
}

//Close closes connection
func (c *TcpClient) Close() error {
	return c.connection.Close()
}

//Send sends command string to establised connection
func (c *TcpClient) Send(name string, s ...string) error {
	_, err := c.writer.WriteString(name + " " + strings.Join(s, " ") + "\n")
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

func sendWithTimeout(fn func() (interface{}, error)) (interface{}, error) {
	var response interface{}
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		response, err = fn()
		ch <- err
	}()
	//listen error channel or context.done
	select {
	case err := <-ch:
		return response, err
	case <-ctx.Done():
		return response, ctx.Err()
	}
}

//AuthWithUserPassPair this command will authenticate your client connection to perform other commands
func (c *TcpClient) AuthWithUserPassPair(user, pass string) error {
	return sendWithTimeout(func() (interface{}, error) {
		if err := c.Send("auth", user, pass); err != nil {
			return nil, err
		}

		//blocked read response
		data, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		var item msg_protocol.MsgPackRootMessage
		err = msgpack.Unmarshal([]byte(data), &item)
		if err != nil {
			return nil, err
		}

		if item.AuthResponse == nil || item.AuthResponse.Message != "authenticated" {
			return nil, errors.New("not authenticated")
		}

		return item.AuthResponse, nil
	})
}

func (c *TcpClient) Get(key string) (*msg_protocol.GetResponse, error) {
	return sendWithTimeout(func() error {

		if err := c.Send("get", key); err != nil {
			return err
		}

		//blocked read response
		data, err := c.reader.ReadString('\n')
		if err != nil {
			ch <- err
			return err
		}

		var item msg_protocol.MsgPackRootMessage
		err = msgpack.Unmarshal([]byte(data), &item)
		if err != nil {
			ch <- err
			return err
		}

		if item.GetResponse == nil {
			ch <- errors.New("not response")
			return errors.New("not response")
		}
		resp = item.GetResponse
		ch <- nil
	})
}

func (c *TcpClient) Set(key string, value string) (*msg_protocol.SetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Send("set", key, value); err != nil {
		return nil, err
	}

	// make error channel
	var resp *msg_protocol.SetResponse
	ch := make(chan error, 1)
	go func() {
		//blocked read response
		data, err := c.reader.ReadString('\n')
		if err != nil {
			ch <- err
			return
		}

		var item msg_protocol.MsgPackRootMessage
		err = msgpack.Unmarshal([]byte(data), &item)
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
