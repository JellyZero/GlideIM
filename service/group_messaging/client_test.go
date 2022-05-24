package group_messaging

import (
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/rpc"
	"github.com/glide-im/glideim/service"
	"testing"
)

func TestNewClient(t *testing.T) {

	config, err := service.GetConfig()
	if err != nil {
		panic(err)
	}
	//etcd := config.Etcd.Servers

	client, err := NewClient(&rpc.ClientOptions{
		Name: config.GroupMessaging.Client.Name,
		Addr: config.GroupMessaging.Server.Addr,
		Port: config.GroupMessaging.Server.Port,
		//EtcdServers: etcd,
	})
	if err != nil {
		t.Error(err)
	}
	cm := message.NewChatMessage(1, 1, 1, 1, 1, "123", 1)
	err = client.DispatchMessage(1, message.ActionGroupMessage, &cm)
	if err != nil {
		t.Error(err)
	}

	_ = client.Close()
}
