package client

import (
	"bufio"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/lib/goerrors"
	"github.com/dwnGnL/ddos-pow/lib/pow"
	challengeResp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"github.com/goccy/go-json"
	"io"
	"log/slog"
	"net"
)

type Client struct {
	conf *config.Config
}

func (s Client) RequestChallenge() (*pow.HashcashData, error) {
	address := fmt.Sprintf("%s:%d", s.conf.Server.Host, s.conf.Server.Port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Warn("tcp dial", err)
		return nil, err
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	err = sendMsg(challengeResp.Message{
		Header: challengeResp.REQUEST_CHALLENGE,
	}, conn)
	if err != nil {
		return nil, fmt.Errorf("err send request: %w", err)
	}

	msgStr, err := readConnMsg(reader)
	if err != nil {
		return nil, fmt.Errorf("err read msg: %w", err)
	}

	msg, err := challengeResp.ParseMessage(msgStr)
	if err != nil {
		return nil, fmt.Errorf("err parse msg: %w", err)
	}

	var hashcash pow.HashcashData

	err = json.Unmarshal([]byte(msg.Payload), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err parse hashcash: %w", err)
	}

	return &hashcash, nil
}

func (s Client) RequestResource(hashcash pow.HashcashData) (string, error) {
	hashcash, err := hashcash.ComputeHashcash(s.conf.Pow.HashcashMaxIterations)
	if err != nil {
		return "", fmt.Errorf("err compute hashcash: %w", err)
	}

	fmt.Println("hashcash computed:", hashcash)
	// marshal solution to json
	byteData, err := json.Marshal(hashcash)
	if err != nil {
		return "", fmt.Errorf("err marshal hashcash: %w", err)
	}

	address := fmt.Sprintf("%s:%d", s.conf.Server.Host, s.conf.Server.Port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Warn("tcp dial", err)
		return "", err
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	// 3. send challenge solution back to server
	err = sendMsg(challengeResp.Message{
		Header:  challengeResp.REQUEST_RESOURCE,
		Payload: string(byteData),
	}, conn)
	if err != nil {
		return "", fmt.Errorf("err send request: %w", err)
	}

	fmt.Println("challenge sent to server")

	// 4. get result quote from server
	msgStr, err := readConnMsg(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}

	msg, err := challengeResp.ParseMessage(msgStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}

	return msg.Payload, nil
}

func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

func sendMsg(msg challengeResp.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	goerrors.Log().Println("msg = ", msgStr)
	return err
}

func New(conf *config.Config) *Client {
	return &Client{
		conf: conf,
	}
}
