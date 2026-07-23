package user

import "testing"

func TestPasswordHashUsesIndependentRandomSalt(t *testing.T) {
	firstHash, firstSalt, err := passwordHash("example-password")
	if err != nil {
		t.Fatalf("第一次生成密码摘要: %v", err)
	}
	secondHash, secondSalt, err := passwordHash("example-password")
	if err != nil {
		t.Fatalf("第二次生成密码摘要: %v", err)
	}
	if firstSalt == "" || secondSalt == "" || firstSalt == secondSalt {
		t.Fatalf("每个账号应使用独立随机盐: first=%q second=%q", firstSalt, secondSalt)
	}
	if firstHash == secondHash {
		t.Fatal("不同随机盐不应生成相同摘要")
	}
	ok, err := verifyPassword("example-password", firstSalt, firstHash)
	if err != nil || !ok {
		t.Fatalf("验证正确密码: ok=%t err=%v", ok, err)
	}
	ok, err = verifyPassword("wrong-password", firstSalt, firstHash)
	if err != nil || ok {
		t.Fatalf("验证错误密码: ok=%t err=%v", ok, err)
	}
}
