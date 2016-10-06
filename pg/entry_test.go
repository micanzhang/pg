package pg

import "testing"

func TestPassword(t *testing.T) {
	for i := 0; i < 1000; i++ {
		bs, _ := randBytes(i)
		passwd := string(bs)
		ks, _ := randBytes(i)
		key := string(Key(string(ks)))

		ps, err := NewPassword(key, passwd)
		if err != nil {
			t.Error(err)
		}

		passwd1, err := ps.Plaintext(key)
		if err != nil {
			t.Error(err)
		}
		if passwd != passwd1 {
			t.Errorf("expected: %s, got: %s", passwd, passwd1)
		}
	}
}

func TestKey(t *testing.T) {
	tcs := []string{
		"",
		"lol",
		"fas324HJK^&*%^",
		"qeqwertyurewuroweprfasfdhasjkllvczvnm,v213215743243297@#$%$$*%^*&e3<>?K:O&*()324234",
	}

	for _, tc := range tcs {
		if length := len(Key(tc)); length != 64 {
			t.Errorf("expected: %d, got: %d (length of %s)", 64, length, tc)
		}
	}
}
