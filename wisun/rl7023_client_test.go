package wisun

import (
	"context"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/u-one/go-el-controller/transport"
)

type resp_RL7023 struct {
	d string
	e error
}

func mock_RL7023(t *testing.T, m *transport.MockSerial, input string, response []resp_RL7023) {
	t.Helper()

	lastCmd := ""
	respCnt := -1

	m.EXPECT().Send(gomock.Any()).DoAndReturn(func(cmd []byte) error {
		lastCmd = string(cmd)
		respCnt = -1
		return nil
	}).AnyTimes()

	m.EXPECT().Recv().DoAndReturn(func() ([]byte, error) {
		resp := ""
		var err error
		if respCnt == -1 {
			resp = lastCmd
		} else {
			resp = response[respCnt].d
			err = response[respCnt].e
		}
		respCnt++
		return []byte(resp), err
	}).AnyTimes()

}

func Test_RL7023_Close(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := transport.NewMockSerial(ctrl)
	mock_RL7023(t, m, "SKTERM\r\n", []resp_RL7023{
		{"OK\r\n", nil},
	})

	c := &RL7023Client{serial: m}
	c.Term()
}

func Test_RL7023_Version(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		input  string
		output []resp_RL7023
		want   string
		err    error
	}{
		{
			name:  "success",
			input: "SKVER\r\n",
			output: []resp_RL7023{
				{"EVER 1.5.2\r\n", nil},
				{"OK\r\n", nil},
			},
			want: "1.5.2",
			err:  nil,
		},
		{
			name:  "recv EVER fail",
			input: "SKVER\r\n",
			output: []resp_RL7023{
				{"", fmt.Errorf("fail on EVER")},
			},
			want: "",
			err:  fmt.Errorf("fail on EVER"),
		},
		{
			name:  "recv not EVER",
			input: "SKVER\r\n",
			output: []resp_RL7023{
				{"XXXX", nil},
			},
			want: "",
			err:  fmt.Errorf("unexpected response [XXXX]"),
		},
		{
			name:  "No version string",
			input: "SKVER\r\n",
			output: []resp_RL7023{
				{"EVER\r\n", nil},
			},
			want: "",
			err:  fmt.Errorf("version string not found"),
		},
		{
			name:  "result not OK",
			input: "SKVER\r\n",
			output: []resp_RL7023{
				{"EVER 1.5.2\r\n", nil},
				{"FAIL\r\n", nil},
			},
			want: "1.5.2",
			err:  fmt.Errorf("command failed [FAIL]"),
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)
			mock_RL7023(t, m, tc.input, tc.output)

			c := &RL7023Client{serial: m}
			got, err := c.Version()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Diffrent result: -want, +got: \n%s", diff)
			}

			if tc.err != nil && err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
				}
			} else if tc.err != err {
				t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
			}

		})
	}

}

func Test_RL7023_SetBRoutePassword(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		pw     string
		input  string
		output []resp_RL7023
		want   string
		err    error
	}{
		{
			name:  "success",
			pw:    "TESTPWDYYYYY",
			input: "SKSETPWD C TESTPWDYYYYY\r\n",
			output: []resp_RL7023{
				{"OK\r\n", nil},
			},
			err: nil,
		},
		{
			name:  "error",
			pw:    "TESTPWDYYYYY",
			input: "SKSETPWD C TESTPWDYYYYY\r\n",
			output: []resp_RL7023{
				{"ER01\r\n", nil},
			},
			err: fmt.Errorf("command failed [ER01]"),
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)

			mock_RL7023(t, m, tc.input, tc.output)

			c := &RL7023Client{serial: m}
			err := c.SetBRoutePassword(tc.input)

			if tc.err != nil && err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
				}
			} else if tc.err != err {
				t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
			}

		})
	}

}

func Test_RL7023_SetBRouteID(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name   string
		pw     string
		input  string
		output []resp_RL7023
		want   string
		err    error
	}{
		{
			name:  "success",
			pw:    "000000TESTID00000000000000000000",
			input: "SKSETRBID 000000TESTID00000000000000000000\r\n",
			output: []resp_RL7023{
				{"OK\r\n", nil},
			},
			err: nil,
		},
		{
			name:  "error",
			pw:    "000000TESTID00000000000000000000",
			input: "SKSETRBID 000000TESTID00000000000000000000\r\n",
			output: []resp_RL7023{
				{"ER01\r\n", nil},
			},
			err: fmt.Errorf("command failed [ER01]"),
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)
			mock_RL7023(t, m, tc.input, tc.output)
			c := &RL7023Client{serial: m}

			err := c.SetBRouteID(tc.input)

			if tc.err != nil && err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
				}
			} else if tc.err != err {
				t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
			}
		})
	}
}

func Test_RL7023_scan(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		duration int
		input    string
		response []resp_RL7023
		expect   bool
		err      error
	}{
		{
			name:     "not found",
			duration: 4,
			input:    "SKSCAN 2 FFFFFFFF 4 0 \r\n",
			response: []resp_RL7023{
				{"OK\r\n", nil},
				{"EVENT 22 2001:0DB8:0000:0000:011A:1111:0000:0001 0\r\n", nil},
			},
			expect: false,
		},
		{
			name:     "found",
			duration: 5,
			input:    "SKSCAN 2 FFFFFFFF 5 0 \r\n",
			response: []resp_RL7023{
				{"OK\r\n", nil},
				{"EVENT 20 2001:0DB8:0000:0000:011A:1111:0000:0001 0\r\n", nil},
			},
			expect: true,
		},
		{
			name:     "received error",
			duration: 5,
			input:    "SKSCAN 2 FFFFFFFF 5 0 \r\n",
			response: []resp_RL7023{
				{"ER01\r\n", nil},
			},
			expect: false,
			err:    fmt.Errorf("command failed [ER01]"),
		},
		{
			name:     "error",
			duration: 5,
			input:    "SKSCAN 2 FFFFFFFF 5 0 \r\n",
			response: []resp_RL7023{
				{"", fmt.Errorf("error")},
			},
			expect: false,
			err:    fmt.Errorf("error"),
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)
			mock_RL7023(t, m, tc.input, tc.response)

			c := &RL7023Client{serial: m}
			got, err := c.scan(context.Background(), tc.duration)
			if tc.expect != got {
				t.Errorf("Diffrent result: want:%v, got:%v", tc.expect, got)
			}

			if tc.err != nil && err != nil {
				if tc.err.Error() != err.Error() {
					t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
				}
			} else if tc.err != err {
				t.Errorf("Diffrent result: want:%#v, got:%#v", tc.err, err)
			}

		})
	}
}

func Test_RL7023_receivePanDesc(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := transport.NewMockSerial(ctrl)

	response := []string{
		"EPANDESC\r\n",
		"  Channel:21\r\n",
		"  Channel Page:01\r\n",
		"  Pan ID:0002\r\n",
		"  Addr:001A111100000002\r\n",
		"  LQI:CA\r\n",
		"  Side:0\r\n",
		"  PairID:0112CE67\r\n",
	}

	respCnt := 0
	m.EXPECT().Recv().DoAndReturn(func() ([]byte, error) {
		resp := response[respCnt]
		respCnt++
		return []byte(resp), nil
	}).AnyTimes()

	want := PanDesc{
		Addr:     "001A111100000002",
		IPV6Addr: "",
		Channel:  "21",
		PanID:    "0002",
	}

	c := &RL7023Client{serial: m}
	got, err := c.receivePanDesc()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Diffrent result: -want, +got: \n%s", diff)
	}

	if err != nil {
		t.Errorf("%s", err)
	}
}

func Test_RL7023_Scan(t *testing.T) {
	t.Parallel()
	// TODO: implement
}

func Test_RL7023_LL64(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := transport.NewMockSerial(ctrl)
	mock_RL7023(t, m, "001A111100000002", []resp_RL7023{
		{"2001:0DB8:0000:0000:011A:1111:0000:0002\r\n", nil},
	})

	want := "2001:0DB8:0000:0000:011A:1111:0000:0002"

	c := &RL7023Client{serial: m}
	got, err := c.LL64("001A111100000002")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Diffrent result: -want, +got: \n%s", diff)
	}

	if err != nil {
		t.Errorf("%s", err)
	}
}

func Test_RL7023_SRegS2(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		channel  string
		input    string
		response []resp_RL7023
		err      error
	}{
		{
			name:    "success",
			channel: "21",
			input:   "SKSREG S2 21\r\n",
			response: []resp_RL7023{
				{"OK\r\n", nil},
			},
			err: nil,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)
			mock_RL7023(t, m, tc.input, tc.response)

			c := &RL7023Client{serial: m}
			err := c.SRegS2(tc.channel)
			if tc.err != err {
				t.Errorf("Diffrent result: want:%v, got:%v", tc.err, err)
			}
		})
	}
}

func Test_RL7023_SRegS3(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		panID    string
		input    string
		response []resp_RL7023
		err      error
	}{
		{
			name:  "success",
			panID: "0002",
			input: "SKSREG S3 0002\r\n",
			response: []resp_RL7023{
				{"OK\r\n", nil},
			},
			err: nil,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)
			mock_RL7023(t, m, tc.input, tc.response)

			c := &RL7023Client{serial: m}
			err := c.SRegS3(tc.panID)
			if tc.err != err {
				t.Errorf("Diffrent result: want:%v, got:%v", tc.err, err)
			}
		})
	}
}

func Test_RL7023_Join(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		panDesc  PanDesc
		input    string
		response []resp_RL7023
		want     bool
		err      error
	}{
		{
			name:    "success",
			panDesc: PanDesc{},
			input:   "SKJOIN 2001:0DB8:0000:0000:011A:1111:0000:0002\r\n",
			response: []resp_RL7023{
				{"OK\r\n", nil},
				{"EVENT 21 2001:0DB8:0000:0000:011A:1111:0000:0002 0 00\r\n", nil},
				{"ERXUDP 2001:0DB8:0000:0000:011A:1111:0000:0002 2001:0DB8:0000:0000:011A:1111:0000:0001 02CC 02CC 001C6400030C12A4 0 0 0054 00000054800000021C2FF4B1D1A295AA00020000003B0000015B003B2F808A7F842BD41B4C0902258E791A8FF605914DA6E50C90EDE8ACC05E035979BA4A00000000C639CE0D16CC74349C8F6DF9FF9D6EA1D200\r\n", nil},
				{"EVENT 25 2001:0DB8:0000:0000:011A:1111:0000:0002 0\r\n", nil},
				{"ERXUDP 2001:0DB8:0000:0000:011A:1111:0000:0002 2001:0DB8:0000:0000:011A:1111:0000:0001 02CC 02CC 001A111100000002 0 0 0028 108100000EF0010EF0017301D50401028801\r\n", nil},
				{"\r\n", nil},
				{"\r\n", nil},
				{"\r\n", nil},
				{"\r\n", nil},
			},
			want: true,
			err:  nil,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)
			mock_RL7023(t, m, tc.input, tc.response)

			c := &RL7023Client{serial: m}
			got, err := c.Join(tc.panDesc)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Diffrent result: -want, +got: \n%s", diff)
			}

			if tc.err != err {
				t.Errorf("Diffrent error: want:%v, got:%v", tc.err, err)
			}
		})
	}
}

func Test_RL7023_Send(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		ipv6addr string
		data     []byte
		input    string
		response []resp_RL7023
		want     []byte
		err      error
	}{
		{
			name:  "success",
			data:  []byte{'X', 'X', 'X', 'X'},
			input: "SKSENDTO 1 2001:0DB8:0000:0000:011A:1111:0000:0002 0E1A 1 0 000E \r\n",
			response: []resp_RL7023{
				{"EVENT 21 2001:0DB8:0000:0000:011A:1111:0000:0002 0 00\r\n", nil},
				{"OK\r\n", nil},
				{"ERXUDP FE80:0000:0000:0000:021C:6400:030C:12A4 FE80:0000:0000:0000:021D:1291:0000:0574 0E1A 0E1A 001C6400030C12A4 1 0 0012 1081000102880105FF017201E704000001F8\r\n", nil},
			},
			want: []byte{0x10, 0x81, 0x00, 0x01, 0x02, 0x88, 0x01, 0x05, 0xff, 0x01, 'r', 0x01, 0xe7, 0x04, 0x00, 0x00, 0x01, 0xf8},
			err:  nil,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := transport.NewMockSerial(ctrl)
			mock_RL7023(t, m, tc.input, tc.response)

			c := &RL7023Client{serial: m, panDesc: PanDesc{IPV6Addr: "2001:0DB8:0000:0000:011A:1111:0000:0002"}}
			got, err := c.Send(tc.data)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Diffrent result: -want, +got: \n%s", diff)
			}

			if tc.err != err {
				t.Errorf("Diffrent error: want:%v, got:%v", tc.err, err)
			}
		})
	}
}
