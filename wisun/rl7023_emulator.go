package wisun

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/goburrow/serial"
)

// RL7023Emulator is ROHM RL7023 emulator
type RL7023Emulator struct {
	port     serial.Port
	reader   *bufio.Reader
	writer   *bufio.Writer
	curLine  []byte
	echoback bool
}

// NewRL7023Emulator returns RL7023Emulator instance
func NewRL7023Emulator(addr string) *RL7023Emulator {
	config := serial.Config{
		Address:  addr,
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  30 * time.Second,
	}

	port, err := serial.Open(&config)
	if err != nil {
		log.Fatal("Faild to open serial:", err)
	}

	r := bufio.NewReaderSize(port, 4096)
	w := bufio.NewWriter(port)
	return &RL7023Emulator{port: port, reader: r, writer: w, echoback: true}
}

// Close closees connection
func (e RL7023Emulator) Close() {
	e.port.Close()
}

func (e RL7023Emulator) flush() {
	e.writer.Flush()
}

func (e RL7023Emulator) echoBack() {
	if e.echoback {
		e.writer.Write(e.curLine)
		e.writer.WriteString("\r\n")
		e.flush()
		fmt.Println("=>", string(e.curLine))
	}
}

func (e RL7023Emulator) ok() {
	e.writer.WriteString("OK\r\n")
	e.flush()
	fmt.Println("=>OK")
}

func (e RL7023Emulator) er(num int) {
	code := fmt.Sprintf("ER%02d", num)
	e.writer.WriteString(code + "\r\n")
	e.flush()
	fmt.Println("=>", code)
}

func (e RL7023Emulator) rxUDP(data []byte) {
	dlen := len(data)
	cmd := fmt.Sprintf("ERXUDP FE80:0000:0000:0000:021C:6400:030C:12A4 FF02:0000:0000:0000:0000:0000:0000:0001 0E1A 0E1A 001C6400030C12A4 1 0 %04x ", dlen)
	e.writer.WriteString(cmd)
	e.writer.Write(data)
	e.writer.WriteString("\r\n")
	e.flush()
	fmt.Println(cmd)
}

// Start starts to emulate
func (e *RL7023Emulator) Start() {
	fmt.Println("Start")
	for {
		line, _, err := e.reader.ReadLine()
		if err != nil {
			fmt.Println("[Error]", err)
		}
		e.curLine = line
		fmt.Println("<=", string(line))
		if bytes.HasPrefix(line, []byte("SKVER")) {
			e.echoBack()
			e.writer.WriteString("EVER 1.0.0\r\n")
			e.flush()
			e.ok()
		} else if bytes.HasPrefix(line, []byte("SKSETPWD")) {
			e.echoBack()
			e.ok()
		} else if bytes.HasPrefix(line, []byte("SKSETRBID")) {
			e.echoBack()
			e.ok()
		} else if bytes.HasPrefix(line, []byte("SKSCAN")) {
			e.echoBack()
			e.ok()

			time.Sleep(3 * time.Second)

			cmd := string(line)
			params := strings.Split(cmd, " ")
			if params[3] == "4" {
				e.writer.WriteString("EVENT 22 FE80:0000:0000:0000:021D:1290:1234:5678 0\r\n")
				e.flush()
			} else if params[3] == "5" {
				e.writer.WriteString("EVENT 20 FE80:0000:0000:0000:021D:1290:1234:5678 0\r\n")
				e.flush()

				e.writer.WriteString("EPANDESC\r\n")
				e.writer.WriteString(" Channel:21\r\n")
				e.writer.WriteString(" Channel Page:09\r\n")
				e.writer.WriteString(" Pan ID:8888\r\n")
				e.writer.WriteString(" Addr:12345678ABCDEF01\r\n")
				e.writer.WriteString(" LQI:E1\r\n")
				e.writer.WriteString(" Side:0\r\n")
				e.writer.WriteString(" PairID:AABBCCDD\r\n")
				e.flush()
			}

		} else if bytes.HasPrefix(line, []byte("SKLL64")) {
			e.echoBack()

			e.writer.WriteString("FE80:0000:0000:0000:021D:1290:1234:ABCD\r\n")
			e.flush()

		} else if bytes.HasPrefix(line, []byte("SKSREG")) {
			e.echoBack()
			e.ok()
		} else if bytes.HasPrefix(line, []byte("SKJOIN")) {
			e.echoBack()
			e.ok()

			e.writer.WriteString("EVENT 22 FE80:0000:0000:0000:021D:1290:1234:5678 0\r\n")
			e.flush()

			e.writer.WriteString("EVENT 21 FE80:0000:0000:0000:021D:1290:1234:5678 0 00\r\n")
			e.flush()
			e.writer.WriteString("ERXUDP FE80:0000:0000:0000:021D:1290:1234:ABCD FE80:0000:0000:0000:021D:1290:1234:5678 02CC 02CC 12345678ABCDEF01 0 0 0028 00000028C00000021C2FF4B92384415800060000000400000000000500030000000400000000000C\r\n")
			e.flush()

			e.writer.WriteString("EVENT 21 FE80:0000:0000:0000:021D:1290:1234:5678 0 00\r\n")
			e.flush()
			e.writer.WriteString("ERXUDP FE80:0000:0000:0000:021D:1290:1234:ABCD FE80:0000:0000:0000:021D:1290:1234:5678 02CC 02CC 12345678ABCDEF01 0 0 0068 00000068800000021C2FF4B92384415900050000001000004B5AD190998EE92EFD8A95FC29389FA20002000000380000014800382F00983B5C1B723326E7BE2B4C07E8090CF1534D3030303030303939303231313030303030303030303030303031313243453637\r\n")
			e.flush()

			e.writer.WriteString("EVENT 21 FE80:0000:0000:0000:021D:1290:1234:5678 0 00\r\n")
			e.flush()
			e.writer.WriteString("ERXUDP FE80:0000:0000:0000:021D:1290:1234:ABCD FE80:0000:0000:0000:021D:1290:1234:5678 02CC 02CC 12345678ABCDEF01 0 0 0054 00000054800000021C2FF4B92384415A00020000003B00000149003B2F80983B5C1B723326E7BE2B4C07E8090CF11988788F485E6994ED46F0F5361E9AB700000000739716C580AD496217B7688AE46EFFAF7E00\r\n")
			e.flush()

			e.writer.WriteString("EVENT 21 FE80:0000:0000:0000:021D:1290:1234:5678 0 00\r\n")
			e.flush()
			e.writer.WriteString("ERXUDP FE80:0000:0000:0000:021D:1290:1234:ABCD FE80:0000:0000:0000:021D:1290:1234:5678 02CC 02CC 12345678ABCDEF01 0 0 0058 00000058A00000021C2FF4B92384415B0007000000040000000000000002000000040000034900040004000000040000000012100008000000040000000151800001000000100000108E7957B9B26C04345D7025C24A2472\r\n")
			e.flush()

			e.writer.WriteString("EVENT 21 FE80:0000:0000:0000:021D:1290:1234:5678 0 00\r\n")
			e.flush()
			e.writer.WriteString("EVENT 25 FE80:0000:0000:0000:021D:1290:1234:5678 0\r\n")
			e.flush()

		} else if bytes.HasPrefix(line, []byte("SKSENDTO")) {
			e.echoBack()

			e.writer.WriteString("EVENT 21 FE80:0000:0000:0000:021D:1290:1234:5678 0 00\r\n")
			e.flush()

			e.ok()
			e.rxUDP([]byte{0x10, 0x81, 0x00, 0x01, 0x02, 0x88, 0x01, 0x05, 0xff, 0x01, 'r', 0x01, 0xe7, 0x04, 0x00, 0x00, 0x01, 0xf8})
		} else if bytes.HasPrefix(line, []byte("SKTERM")) {
			e.echoBack()
			e.ok()
		}

	}
}
