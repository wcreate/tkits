package tkits

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	tokenLen    = 1 + 8 + 16 + 8
	TOKEN_SYS   = 0x00 // system general user
	TOKEN_USER  = 0x01 // login user
	TOKEN_ADMIN = 0x01
)

var (
	sysuserid = []byte{
		0x19, 0x82, 0x06, 0x08,
		0x65, 0x78, 0xAB, 0xCD,
	}
)

type SimpleToken struct {
	c      *Crypto
	expire int64 // unit is minute
}

func NewSimpleToken(c *Crypto, expire int64) *SimpleToken {
	return &SimpleToken{c, expire}
}

// generate a token
// flag(1) + uid(8) + ip(16) + currenttime(8)
func (st *SimpleToken) GenToken(clientip string, uid int64, flag byte) (string, error) {
	bs := make([]byte, tokenLen)
	buf := bytes.NewBuffer(bs)

	// flag
	buf.WriteByte(flag)

	if flag == TOKEN_SYS { // system user
		if _, err := buf.Write(sysuserid); err != nil {
			return "", err
		}
	} else {
		// uid
		binary.Write(buf, binary.BigEndian, uid)
	}

	// client ip
	ip := net.ParseIP(clientip)
	buf.Write([]byte(ip))

	// current time
	binary.Write(buf, binary.BigEndian, time.Now().Unix())

	// encrypt
	if t, err := st.c.Encrypt(buf.Bytes()); err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(t), nil
	}
}

func (st *SimpleToken) Validate(token, clientip string, uid int64, vtime bool) bool {
	bs, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	// check the len
	if len(bs) != tokenLen {
		return false
	}

	// remain bytes
	buf := bytes.NewReader(bs[1:])

	if flag := bs[0]; flag == TOKEN_SYS {
		dbs := make([]byte, 16)
		if rszie, _ := buf.Read(dbs); rszie != 16 {
			return false
		}
		if !bytes.Equal(dbs, sysuserid) {
			return false
		}
	} else {
		// read uid
		var iuid int64
		if err := binary.Read(buf, binary.BigEndian, &iuid); err != nil || iuid != uid {
			return false
		}
	}

	// ip
	var ip net.IP
	buf.Read([]byte(ip))
	if !ip.Equal(net.ParseIP(clientip)) {
		return false
	}

	// no need to validate whether the time is expired
	if !vtime {
		return true
	}

	var ttime int64
	if err := binary.Read(buf, binary.BigEndian, &ttime); err != nil {
		return false
	}

	if time.Now().Unix() > st.expire*60 {
		log.Errorf("token %s is expired", token)
		return false
	}

	return true
}
