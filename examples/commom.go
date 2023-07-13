package common

import (
	"flag"
	"fmt"
	"os"
	"time"

	esphome "github.com/rockavoldy/esphome_api/pkg/client"
	"google.golang.org/protobuf/proto"
)

const (
	EnvHostAddress   = "ESPHOME_ADDRESS"
	EnvPassword      = "ESPHOME_PASSWORD"
	EnvEncryptionKey = "ESPHOME_ENCRYPTION_KEY"
)

var (
	HostAddressFlag   = flag.String("address", "", "esphome node hostname or IP with port. example: my_esphome.local:6053")
	PasswordFlag      = flag.String("password", "", "esphome node API password")
	EncryptionKeyFlag = flag.String("encryption-key", "", "esphome node API encryption key")
	TimeoutFlag       = flag.Duration("timeout", 10*time.Second, "communication timeout")
)

func GetClient(handlerFunc func(msg proto.Message)) (*esphome.Client, error) {
	flag.Parse()

	// update hostaddress
	if *HostAddressFlag == "" {
		if os.Getenv(EnvHostAddress) != "" {
			*HostAddressFlag = os.Getenv(EnvHostAddress)
		} else {
			*HostAddressFlag = "esphome.local:6053"
		}
	}

	// update password
	if *PasswordFlag == "" {
		*PasswordFlag = os.Getenv(EnvPassword)
	}

	// update encryption key
	if *EncryptionKeyFlag == "" {
		*EncryptionKeyFlag = os.Getenv(EnvEncryptionKey)
	}

	if handlerFunc == nil {
		handlerFunc = handlerFuncImpl
	}

	client, err := esphome.GetClient("mycontroller.org", *HostAddressFlag, *EncryptionKeyFlag, *TimeoutFlag, handlerFunc)
	if err != nil {
		return nil, err
	}

	if err = client.Login(*PasswordFlag); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}

func handlerFuncImpl(msg proto.Message) {
	fmt.Printf("received a message, type: %T, value: [%v]\n", msg, msg)
}
