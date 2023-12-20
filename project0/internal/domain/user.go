package domain

import "time"

// 领域是抽象，面对服务
type User struct {
	Id       int64
	Email    string
	Password string

	Ctime time.Time

	//Addr Address
}

type UserProfile struct {
	Id           int64
	Gender       string
	NickName     string
	Introduction string
	BirthDate    time.Time
}

//type Address struct {
//	Province string
//	Region   string
//}
