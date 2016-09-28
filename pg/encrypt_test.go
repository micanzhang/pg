package pg

import "testing"

func equalSliceByte(a []byte, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestPadAndUnpad(t *testing.T) {
	for i := 0; i < 256; i++ {
		in, _ := randBytes(i)
		padding := pad(in)
		out := unpad(padding)

		if !equalSliceByte(in, out) {
			t.Errorf("expected: %+v, got : %+v", in, out)
		}
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	testCases := []string{
		"",
		"hello, crytography",
		"1qaz@WSX3edc$RFV",
		"242379fafhasdzcvfafasfj42342384023fasfasdnvzxvzxcv&*()43423hcacfdMN?><MCD4324@#$#$#_*cafdajFDFDJSFDS出发地方发的f&(##$)hlahj",
	}

	for _, tc := range testCases {
		key := Key(tc)
		message := []byte(tc)
		ct, err := Encrypt(key, message)
		if err != nil {
			t.Error(err)
			continue
		}

		text, err := Decrypt(key, ct)
		if err != nil {
			t.Error(err)
			continue
		}

		if !equalSliceByte(message, text) {
			t.Errorf("expected: %s, got: %s", string(message), string(text))
		}
	}

}
