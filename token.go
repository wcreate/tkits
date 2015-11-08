package tkits

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"time"
)

type SimpleToken struct {
	c      *Crypto
	expire float64 // unit is minute
}

func NewSimpleToken(c *Crypto, expire float64) *SimpleToken {
	return &SimpleToken{c, expire}
}

func (st *SimpleToken) GenToken(clientip, uid string) (string, error) {
	// uid(8) + ip(16) + currenttime(15)
	bs := make([]byte, 8+16+15)
	buf := bytes.NewBuffer(bs)

	// uid
	if iuid, err := strconv.ParseInt(uid, 10, 0); err != nil {
		return "", err
	} else {
		binary.Write(buf, binary.BigEndian, iuid)
	}

	// client ip
	ip := net.ParseIP(clientip)
	buf.Write([]byte(ip))

	// current time
	if ct, err := time.Now().MarshalBinary(); err != nil {
		return "", err
	} else {
		buf.Write([]byte(ct))
	}

	// encrypt
	if t, err := st.c.Encrypt(buf.Bytes()); err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(t), nil
	}
}

func (st *SimpleToken) Validate(token string, clientip, uid string, vtime bool) bool {
	bs, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	// check the len
	if len(bs) != 8+16+15 {
		return false
	}

	buf := bytes.NewReader(bs)
	// read uid
	var iuid int64
	if err := binary.Read(buf, binary.BigEndian, &iuid); err != nil {
		return false
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

	var ct time.Time
	if err := ct.UnmarshalBinary(bs[8+16:]); err != nil {
		return false
	}

	if time.Now().Sub(ct).Minutes() > st.expire {
		log.Printf("token %s is expired", token)
		return false
	}

	return true
}

func writeInt8(dest []byte, offset int, num int64) {
	for i := 0; i < 8; i++ {
		dest[offset+i] = byte(num >> uint((7-i)*8) & 0xFF)
	}

}
