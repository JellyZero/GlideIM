package client

import (
	"go_im/im/conn"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
	"strconv"
	"sync"
	"sync/atomic"
)

type DefaultClientManager struct {
	clients     *clients
	clientCount int64
}

func NewDefaultManager() *DefaultClientManager {
	ret := new(DefaultClientManager)
	ret.clients = newClients()
	return ret
}

// ClientConnected 当一个用户连接建立后, 由该方法创建 IClient 实例 Client 并管理该连接, 返回该由连接创建客户端的标识 id
// 返回的标识 id 是一个临时 id, 后续连接认证后会改变
func (c *DefaultClientManager) ClientConnected(conn conn.Connection) int64 {
	statistics.SConnEnter()

	// 获取一个临时 uid 标识这个连接
	connUid := uid.GenTemp()
	ret := newClient(conn)
	ret.SetID(connUid, 0)
	c.clients.add(connUid, 0, ret)
	atomic.AddInt64(&c.clientCount, 1)
	// 开始处理连接的消息
	ret.Run()
	return connUid
}

func (c *DefaultClientManager) AddClient(uid int64, cs IClient) {
	c.clients.add(uid, 0, cs)
	atomic.AddInt64(&c.clientCount, 1)
}

// ClientSignIn 客户端登录, id 为连接时使用的临时标识, uid 为z用户标识, device 用于区分不同设备
func (c *DefaultClientManager) ClientSignIn(id, uid_ int64, device int64) error {
	logger.D("client sign in origin-id=%d, uid=%d", id, uid_)
	tempDs := c.clients.get(id)
	if tempDs == nil || tempDs.size() == 0 {
		// 该客户端不存在
		logger.W("attempt to sign in a nonexistent client, id=%d", id)
		return nil
	}
	client := tempDs.get(0)
	logged := c.clients.get(uid_)
	if logged != nil && logged.size() > 0 {
		// 多设备登录
		existing := logged.get(device)
		if existing != nil {
			logger.D("multi device login mutex, uid=%d, device=%d", uid_, device)
			existing.SetID(uid.GenTemp(), 0)
			// "Your account is logged in on another device"
			existing.EnqueueMessage(message.NewMessage(0, message.ActionNotifyKickOut, "Your account is logged in on another device"))
			existing.Exit()
			logged.remove(device)
			atomic.AddInt64(&c.clientCount, 1)
		}
		if logged.size() > 0 {
			msg := "multi device login, device=" + strconv.FormatInt(device, 10)
			EnqueueMessage(uid_, message.NewMessage(0, message.ActionNotifyAccountLogin, msg))
		}
		logged.put(device, client)
	} else {
		// 单设备登录
		c.clients.add(uid_, device, client)
	}
	client.SetID(uid_, device)
	// 删除临时 id
	c.clients.delete(id, 0)
	return nil
}

func (c *DefaultClientManager) ClientLogout(uid_ int64, device int64) error {
	cl := c.clients.get(uid_)
	if cl == nil || cl.size() == 0 {
		logger.E("uid is not sign in, uid=%d", uid_)
		return nil
	}
	logDevice := cl.get(device)
	if logDevice == nil {
		logger.E("device not exist")
		return nil
	}
	logger.I("client logout, uid=%d, device=%d", uid_, device)
	logDevice.SetID(uid.GenTemp(), 0)
	logDevice.Exit()
	cl.remove(device)
	atomic.AddInt64(&c.clientCount, -1)
	statistics.SConnExit()
	return nil
}

func (c *DefaultClientManager) EnqueueMessage(uid int64, device int64, msg *message.Message) error {
	ds := c.clients.get(uid)
	if ds == nil || ds.size() == 0 {
		// offline
		return nil
	}
	if device != 0 {
		d := ds.get(device)
		if d == nil {
			return nil
		}
		d.EnqueueMessage(msg)
	}
	ds.foreach(func(deviceId int64, c IClient) {
		if device != 0 && deviceId != device {
			return
		}
		if c.Closed() {
			// TODO 2021-10-27 client is offline, store
		} else {
			c.EnqueueMessage(msg)
		}
	})
	return nil
}

func (c *DefaultClientManager) isOnline(uid int64) bool {
	ds := c.clients.get(uid)
	if ds == nil {
		return false
	}
	return ds.size() > 0
}

func (c *DefaultClientManager) isDeviceOnline(uid, device int64) bool {
	ds := c.clients.get(uid)
	if ds == nil {
		return false
	}
	return ds.get(device) != nil
}

func (c *DefaultClientManager) allClient() []int64 {
	var ret []int64
	for k := range c.clients.clients {
		if k > 0 {
			ret = append(ret, k)
		}
	}
	return ret
}

//////////////////////////////////////////////////////////////////////////////

type devices struct {
	ds map[int64]IClient
}

func (d *devices) put(device int64, cli IClient) {
	d.ds[device] = cli
}

func (d *devices) get(device int64) IClient {
	return d.ds[device]
}

func (d *devices) remove(device int64) {
	delete(d.ds, device)
}

func (d *devices) foreach(f func(device int64, c IClient)) {
	for k, v := range d.ds {
		f(k, v)
	}
}
func (d *devices) size() int {
	return len(d.ds)
}

type clients struct {
	m       sync.RWMutex
	clients map[int64]*devices
}

func newClients() *clients {
	ret := new(clients)
	ret.m = sync.RWMutex{}
	ret.clients = make(map[int64]*devices)
	return ret
}

func (g *clients) size() int {
	g.m.RLock()
	defer g.m.RUnlock()
	return len(g.clients)
}

func (g *clients) get(uid int64) *devices {
	g.m.RLock()
	defer g.m.RUnlock()
	cl, ok := g.clients[uid]
	if ok && cl.size() != 0 {
		return cl
	}
	return nil
}

func (g *clients) contains(uid int64) bool {
	g.m.RLock()
	defer g.m.RUnlock()
	_, ok := g.clients[uid]
	return ok
}

func (g *clients) add(uid int64, device int64, c IClient) {
	g.m.Lock()
	defer g.m.Unlock()
	cs, ok := g.clients[uid]
	if ok {
		cs.put(device, c)
	} else {
		d := &devices{map[int64]IClient{}}
		d.put(device, c)
		g.clients[uid] = d
	}
}

func (g *clients) delete(uid int64, device int64) {
	g.m.Lock()
	defer g.m.Unlock()
	d, ok := g.clients[uid]
	if ok {
		d.remove(device)
		if d.size() == 0 {
			delete(g.clients, uid)
		}
	}
}
