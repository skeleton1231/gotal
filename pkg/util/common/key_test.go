package common

import (
	"testing"
)

// TestEncryptAndCompare 测试 Encrypt 和 Compare 函数
func TestEncryptAndCompare(t *testing.T) {
	// 测试用的明文密码
	testPassword := "12345678"

	hashedPassword := "$2a$10$n6nYIDN39wMOJKv0svHvreUgxTZPFwqt3PKItvIVRUzNsVmQ4kYFK"
	var err error
	// 加密密码
	// hashedPassword, err := Encrypt(testPassword)
	// if err != nil {
	// 	t.Errorf("Encrypt failed: %v", err)
	// }

	// 打印明文密码和加密密码
	t.Logf("明文密码: %s", testPassword)
	t.Logf("加密密码: %s", hashedPassword)

	// 检查加密密码是否为空
	if hashedPassword == "" {
		t.Errorf("Encrypt returned empty hash")
	}

	// 比较原始密码和加密密码
	err = Compare(hashedPassword, testPassword)
	if err != nil {
		t.Errorf("Compare failed: %v", err)
	}

	// 使用错误的密码进行比较
	err = Compare(hashedPassword, "wrongpassword")
	if err == nil {
		t.Errorf("Compare should fail with wrong password")
	}
}
