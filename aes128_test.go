package tkits

import "testing"

func TestKey(t *testing.T) {
	factor := "ywlSRb80TaCQ4b7b"
	t.Logf("%s", KermitStr(factor))
}

func TestCrypto(t *testing.T) {
	usr, _ := scrypto.EncryptStr("oto")
	t.Logf("%s", usr)
	if d, _ := scrypto.DecryptStr(usr); d != "oto" {
		t.Fail()
	}
	pwd, _ := scrypto.EncryptStr("123456")
	t.Logf("%s", pwd)
	if d, _ := scrypto.DecryptStr(pwd); d != "123456" {
		t.Fail()
	}
}
