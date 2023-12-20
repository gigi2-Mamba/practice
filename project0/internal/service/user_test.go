package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// bcrypt 限制密码长度不能超过72字节
// 思考源码真的很痛苦  这样去看源码等死吧
func TestPassowrdEncrypt(t *testing.T) {
	password := []byte("123456#hello")
	encrypyted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	fmt.Println(string(encrypyted))
	err = bcrypt.CompareHashAndPassword(encrypyted, []byte("123456#hello"))
	assert.NoError(t, err)

}
