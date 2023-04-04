package wisun

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"

	"github.com/u-one/go-el-controller/transport"
)

// RL7023Client is client for TESSERA RL7023
type RL7023Client struct {
	sendSeq int
	readSeq int
	serial  transport.Serial
	panDesc PanDesc
	joined  bool
}

// NewRL7023Client returns RL7023Client instance
func NewRL7023Client(portaddr string) *RL7023Client {
	fmt.Println("NewRL7023Client: ", portaddr)
	s := transport.NewSerialImpl(portaddr)
	return &RL7023Client{serial: s}
}

// Close closees connection
func (c *RL7023Client) Close() {
	if c.joined {
		c.Term()
	}
	c.serial.Close()
}

// Send sends serial command
func (c *RL7023Client) send(in []byte) error {
	c.sendSeq++
	log.Printf("Send[%d]:%s", c.sendSeq, stringWithBinary(in))
	err := c.serial.Send(in)
	if err != nil {
		return err
	}
	// Echoback
	_, err = c.recv()
	if err != nil {
		return err
	}
	return nil
}

// recv receives serial response by line
func (c *RL7023Client) recv() ([]byte, error) {
	line, err := c.serial.Recv()
	c.readSeq++
	if err != nil {
		return []byte{}, err
	}

	log.Printf("Read[%d]:%s", c.readSeq, stringWithBinary(line))
	line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
	return line, err
}

func (c *RL7023Client) recvOK() error {
	r, err := c.recv()
	if err != nil {
		return err
	}
	if !bytes.Equal(r, []byte("OK")) {
		return fmt.Errorf("command failed [%s]", r)
	}
	return nil
}

// Version is ..
func (c RL7023Client) Version() (string, error) {
	err := c.send([]byte("SKVER\r\n"))
	if err != nil {
		return "", err
	}

	//EVER X.Y.Z
	r, err := c.recv()
	if err != nil {
		return "", err
	}

	if !bytes.HasPrefix(r, []byte("EVER")) {
		return "", fmt.Errorf("unexpected response [%s]", r)
	}

	tokens := bytes.Split(r, []byte{' '})
	if len(tokens) < 2 {
		return "", fmt.Errorf("version string not found")
	}
	ver := string(tokens[1])

	err = c.recvOK()
	if err != nil {
		return ver, err
	}
	return ver, nil
}

// SetBRoutePassword is..
func (c RL7023Client) SetBRoutePassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("b-route password is empty")
	}

	err := c.send([]byte("SKSETPWD C " + password + "\r\n"))
	if err != nil {
		return err
	}

	return c.recvOK()
}

// SetBRouteID  is ..
func (c RL7023Client) SetBRouteID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("b-route ID is empty")
	}

	err := c.send([]byte("SKSETRBID " + id + "\r\n"))
	if err != nil {
		return err
	}

	return c.recvOK()
}

func (c RL7023Client) scan(ctx context.Context, duration int) (bool, error) {

	err := c.send([]byte(fmt.Sprintf("SKSCAN 2 FFFFFFFF %d 0 \r\n", duration)))
	if err != nil {
		return false, err
	}

	err = c.recvOK()
	if err != nil {
		return false, err
	}

	for {
		ch := make(chan error)
		var data []byte
		go func(data *[]byte) {
			res, err := c.recv()
			if err != nil {
				log.Println(err)
				ch <- err
			}

			if bytes.HasPrefix(res, []byte("EVENT 22")) {
				log.Println("found EVENT 22")
				ch <- nil
			}
			if bytes.HasPrefix(res, []byte("EVENT 20")) {
				log.Println("found EVENT 20")
				*data = res
				ch <- nil
			}
		}(&data)

		select {
		case err := <-ch:
			if err == nil {
				return len(data) != 0, nil
			}
		case <-ctx.Done():
			return false, fmt.Errorf("scan timeout: %w", ctx.Err())
		}
	}
}

func (c RL7023Client) receivePanDesc() (PanDesc, error) {
	ed := PanDesc{}
	line, err := c.recv()
	if err == nil && bytes.HasPrefix(line, []byte("EPANDESC")) {
		line, err := c.recv() // Channel
		if err != nil {
			return PanDesc{}, fmt.Errorf("failed to get Channel [%s]", err)
		}
		tokens := bytes.Split(line, []byte("Channel:"))
		ed.Channel = string(bytes.Trim(tokens[1], "\r\n"))
		c.recv()             // Channel Page: XX
		line, err = c.recv() // Pan ID: XXXX
		if err != nil {
			return PanDesc{}, fmt.Errorf("failed to get Pan ID [%s]", err)
		}
		tokens = bytes.Split(line, []byte("Pan ID:"))
		ed.PanID = string(bytes.Trim(tokens[1], "\r\n"))
		line, err = c.recv() // Addr:XXXXXXXXXXXXXXXX
		if err != nil {
			return PanDesc{}, fmt.Errorf("failed to get Addr [%s]", err)
		}
		tokens = bytes.Split(line, []byte("Addr:"))
		ed.Addr = string(bytes.Trim(tokens[1], "\r\n"))
		c.recv() // LQI:CA
		c.recv() // Side:X
		c.recv() // PairID:XXXXXXXX
	}
	return ed, err
}

// Scan is ..
func (c RL7023Client) Scan(ctx context.Context) (PanDesc, error) {
	duration := 4
	for {
		if duration > 8 {
			log.Println("duration limit(8) exceeds")
			break
		}

		found, err := c.scan(ctx, duration)
		if err != nil {
			return PanDesc{}, fmt.Errorf("scan failed: %w", err)
		}
		if found {
			break
		}
		duration = duration + 1
	}

	ed, err := c.receivePanDesc()
	log.Printf("Received EPANDesc:%#v", ed)
	return ed, err
}

// LL64 is .
func (c RL7023Client) LL64(addr string) (string, error) {
	cmd := fmt.Sprintf("SKLL64 %s\r\n", addr)
	c.send([]byte(cmd))
	line, err := c.recv()
	if err != nil {
		return "", err
	}
	ipV6Addr := string(bytes.Trim(line, "\r\n"))
	log.Printf("Translated address:%#v", ipV6Addr)
	return ipV6Addr, nil
}

// SRegS2 is.
func (c RL7023Client) SRegS2(channel string) error {
	cmd := fmt.Sprintf("SKSREG S2 %s\r\n", channel)
	c.send([]byte(cmd))
	c.recv()
	return nil
}

// SRegS3 is ..
func (c RL7023Client) SRegS3(panID string) error {
	cmd := fmt.Sprintf("SKSREG S3 %s\r\n", panID)
	c.send([]byte(cmd))
	c.recv()
	return nil
}

// Join is ..
func (c *RL7023Client) Join(desc PanDesc) (bool, error) {
	cmd := fmt.Sprintf("SKJOIN %s\r\n", desc.IPV6Addr)
	c.send([]byte(cmd))
	c.recv()

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			err := fmt.Errorf("timeout:%w", ctx.Err())
			log.Println(err)
			return false, err
		default:
			res, err := c.recv()
			if err != nil {
				log.Println(err)
				if err.Error() == "serial: timeout" {
					continue
				}
				return false, fmt.Errorf("join failed: %w", err)
			}

			tokens := bytes.Split(res, []byte{' '})
			eventType := string(tokens[0])

			switch eventType {
			case "EVENT":
				if len(tokens) < 2 {
					return false, fmt.Errorf("invalid format [%s]", res)
				}
				num, err := strconv.ParseInt(string(tokens[1]), 16, 8)
				if err != nil {
					return false, fmt.Errorf("invalid EVENT num [%s]", res)
				}
				switch num {
				case 0x24:
					log.Println("Join failed")
					return false, nil
				case 0x25:
					log.Println("Join succeed")
					c.joined = true
					return true, nil
				}
			}
		}
	}

}

// Send is...
func (c *RL7023Client) Send(data []byte) ([]byte, error) {
	ipv6 := c.panDesc.IPV6Addr
	cmd := []byte(fmt.Sprintf("SKSENDTO 1 %s 0E1A 1 0 %04X ", ipv6, len(data)))
	cmd = append(cmd, data...)
	cmd = append(cmd, []byte("\r\n")...)
	c.send(cmd)

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			return nil, ctx.Err()
		default:
			res, err := c.recv()
			if err != nil {
				log.Println(err)
				if err.Error() == "serial: timeout" {
					continue
				}
				return nil, err
			}

			tokens := bytes.Split(res, []byte{' '})
			eventType := string(tokens[0])

			switch eventType {
			case "EVENT":
				if len(tokens) < 2 {
					log.Printf("invalid format [%s]\n", res)
				}
				num, err := strconv.ParseInt(string(tokens[1]), 16, 8)
				if err != nil {
					log.Printf("invalid EVENT num [%s]\n", res)
				}
				switch num {
				case 0x21:
					log.Println("UDP send succeed")
				default:
					log.Printf("unexpected EVENT %x\n", num)
				}
			case "ERXUDP":
				// ERXUDP <SENDER> <DEST> <RPORT> <LPORT> <SENDERLLA> (<RSSI>) <SECURED> <SIDE> <DATALEN> <DATA>
				if len(tokens) >= 10 {
					dstPort, err := strconv.ParseInt(string(tokens[4]), 16, 16)
					if err != nil {
						return nil, fmt.Errorf("invalid destination port [%s]", res)
					}
					switch dstPort {
					case 3610: // ECHONET Lite
						src := tokens[9]
						dst := make([]byte, hex.DecodedLen(len(src)))
						n, err := hex.Decode(dst, src)
						return dst[:n], err
					case 716: // PANA
						log.Println("PANA data")
					case 19788: // MLE
						log.Println("MLE data")
					}

				}
			}
		}
	}
}

// Connect connects to smart-meter
func (c *RL7023Client) Connect(ctx context.Context, bRouteID, bRoutePW string) error {

	if len(bRouteID) == 0 {
		err := fmt.Errorf("set B-route ID")
		return err
	}
	if len(bRoutePW) == 0 {
		err := fmt.Errorf("set B-route password")
		return err
	}

	err := c.SetBRoutePassword(bRoutePW)
	if err != nil {
		err := fmt.Errorf("SetBRoutePassword failed: %w", err)
		return err
	}

	err = c.SetBRouteID(bRouteID)
	if err != nil {
		err := fmt.Errorf("SetBRouteID failed: %w", err)
		return err
	}

	pd, err := c.Scan(ctx)
	if err != nil {
		err := fmt.Errorf("Scan failed: %w", err)
		return err
	}

	ipv6Addr, err := c.LL64(pd.Addr)
	if err != nil {
		err := fmt.Errorf("LL64 failed: %w", err)
		return err
	}

	pd.IPV6Addr = ipv6Addr
	log.Printf("Translated address:%#v", pd)

	err = c.SRegS2(pd.Channel)
	if err != nil {
		err := fmt.Errorf("SRegS2 failed: %w", err)
		return err
	}

	err = c.SRegS3(pd.PanID)
	if err != nil {
		err := fmt.Errorf("SRegS3 failed: %w", err)
		return err
	}

	// PANA authentication
	joined, err := c.Join(pd)
	if err != nil {
		err := fmt.Errorf("Join failed: %w", err)
		return err
	}

	if !joined {
		return fmt.Errorf("Join failed")
	}

	c.panDesc = pd

	// TODO: return error
	return nil
}

// Term terminates PANA session
func (c RL7023Client) Term() {
	c.send([]byte("SKTERM\r\n"))
	c.recv()
}
