package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/lib/cache"
	"github.com/dwnGnL/ddos-pow/lib/pow"
	challengeResp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"golang.org/x/exp/rand"
	"log/slog"
	"strconv"
	"time"
)

// Quotes - const array of quotes to respond on client's request
var Quotes = []string{
	"All saints who remember to keep and do these sayings, " +
		"walking in obedience to the commandments, " +
		"shall receive health in their navel and marrow to their bones",

	"And shall find wisdom and great treasures of knowledge, even hidden treasures",

	"And shall run and not be weary, and shall walk and not faint",

	"And I, the Lord, give unto them a promise, " +
		"that the destroying angel shall pass by them, " +
		"as the children of Israel, and not slay them",
}

type Server struct {
	conf  *config.Config
	cache *cache.InMemoryCache
}

func (s Server) Ping() string {
	return "Pong..."
}

func (s Server) ResponseChallenge(clientInfo string) (msg *challengeResp.Message, err error) {
	randValue := rand.Intn(100000)
	err = s.cache.Add(randValue, s.conf.Pow.HashcashDuration)
	if err != nil {
		return nil, fmt.Errorf("err add rand to cache: %w", err)
	}

	hashcash := pow.HashcashData{
		Version:    1,
		ZerosCount: s.conf.Pow.HashcashZerosCount,
		Date:       time.Now().Unix(),
		Resource:   clientInfo,
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
		Counter:    0,
	}

	hashCashMarshaled, err := json.Marshal(hashcash)
	if err != nil {
		return nil, fmt.Errorf("err marshal hashCash: %v", err)
	}

	msg = &challengeResp.Message{
		Header:  challengeResp.RESPONSE_CHALLENGE,
		Payload: string(hashCashMarshaled),
	}

	return
}

func (s Server) ResponseResource(clientInfo string, hashCashSolved string) (msg *challengeResp.Message, err error) {
	fmt.Println("just a check response resource")
	fmt.Printf("client %s requests resource with payload %s\n", clientInfo, hashCashSolved)
	var hashcash pow.HashcashData
	msg = new(challengeResp.Message)

	err = json.Unmarshal([]byte(hashCashSolved), &hashcash)
	if err != nil {
		return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
	}

	slog.Info("hashcash test 4", hashcash)

	if hashcash.Resource != clientInfo {
		msg.Payload = "invalid hashcash resource"
		return msg, nil
	}

	randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
	if err != nil {
		return nil, fmt.Errorf("err decode rand: %w", err)
	}

	randValue, err := strconv.Atoi(string(randValueBytes))
	if err != nil {
		return nil, fmt.Errorf("err decode rand: %w", err)
	}

	exists, err := s.cache.Get(randValue)
	if err != nil {
		return nil, fmt.Errorf("err get rand from cache: %w", err)
	}

	if !exists {
		msg.Payload = "challenge expired or not sent"
		return msg, nil
	}

	if time.Now().Unix()-hashcash.Date > s.conf.Pow.HashcashDuration {
		msg.Payload = "challenge expired"
		return msg, nil
	}

	maxIter := hashcash.Counter
	if maxIter == 0 {
		maxIter = 1
	}

	_, err = hashcash.ComputeHashcash(maxIter)
	if err != nil {
		return nil, fmt.Errorf("invalid hashcash")
	}
	fmt.Printf("client %s succesfully computed hashcash %s\n", clientInfo, hashCashSolved)

	msg = &challengeResp.Message{
		Header:  challengeResp.RESPONSE_RESOURCE,
		Payload: Quotes[rand.Intn(4)],
	}

	return msg, nil
}

func New(conf *config.Config) *Server {
	return &Server{
		conf:  conf,
		cache: cache.InitInMemoryCache(),
	}
}
