package uc

import (
	"../common"
	"encoding/base64"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

func (s *UCServer) RPCCheckPassword(logindata LoginData, reply *bool) error {
	value := common.MD5V(logindata.Password)
	userIndex := &UserIndex{}
	err := s.database.C("user_index").FindId(logindata.Username).One(userIndex)
	user := &UserData{}
	if err != nil {
		err = s.database.C("user").Find(bson.M{
			"_id": logindata.Username,
			"password": value,
		}).One(user)
	} else {
		err = s.database.C("user").Find(bson.M{
			"_id": user.Id,
			"password": value,
		}).One(user)
	}
	if err != nil {
		*reply = false
	} else {
		*reply = true
	}
	return nil
}
func (s *UCServer) RPCSearch(key string, reply *UserDataNoPassword) error {
	userIndex := &UserIndex{}
	err := s.database.C("user_index").FindId(key).One(userIndex)
	user := &UserDataNoPassword{}
	if err != nil {
		err = s.database.C("user").FindId(key).One(user)
		if err != nil {
			return err
		} else {
			*reply = *user
		}
	} else {
		err = s.database.C("user").FindId(userIndex.UserId).One(user)
		if err != nil {
			return err
		} else {
			*reply = *user
		}
	}
	return nil
}
func (s *UCServer) RPCUpdate(user UserDataNoPassword, reply *bool) error {
	var phone bool
	var email bool
	old := &UserDataNoPassword{}
	index := &UserIndex{}
	err := s.database.C("user").FindId(user.Id).One(old)
	if err != nil {
		*reply = false
		return err
	}
	if old.Phone != user.Phone {
		err := s.database.C("user_index").FindId(user.Phone).One(index)
		if err != nil {
			phone = true
		} else {
			*reply = false
			return errors.New("Phone is already exists!")
		}
	}
	if old.EMail != user.EMail {
		err := s.database.C("user_index").FindId(user.EMail).One(index)
		if err != nil {
			email = true
		} else {
			*reply = false
			return errors.New("E-Mail is already exists!")
		}
	}
	if phone == true {
		s.database.C("user_index").RemoveId(old.Phone)
		s.database.C("user_index").Insert(&UserIndex{
			Id: user.Phone,
			Type: "phone",
			UserId: user.Id,
		})
	}
	if email == true {
		s.database.C("user_index").RemoveId(old.EMail)
		s.database.C("user_index").Insert(&UserIndex{
			Id: user.EMail,
			Type: "email",
			UserId: user.Id,
		})
	}
	err = s.database.C("user").UpdateId(user.Id, bson.M{"$set": user})
	if err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}
func (s *UCServer) RPCRegister(user UserData, reply *bool) error {
	index := &UserIndex{}
	err := s.database.C("user_index").FindId(user.Phone).One(index)
	if err == nil {
		*reply = false
		return errors.New("phone is already exists")
	}
	err = s.database.C("user_index").FindId(user.EMail).One(index)
	if err == nil {
		*reply = false
		return errors.New("e-mail is already exists")
	}
	_user := &UserData{}
	err = s.database.C("user").FindId(user.Id).One(_user)
	if err == nil {
		*reply = false
		return errors.New("user id is already exists")
	}
	indexPhone, err1 := _set("phone", user.Id, user.Phone)
	indexEmail, err2 := _set("email", user.Id, user.EMail)
	if err1 == nil && err2 == nil {
		s.database.C("user_index").Insert(indexPhone)
		s.database.C("user_index").Insert(indexEmail)
		password, _ := base64.StdEncoding.DecodeString(user.Password)
		buf, _ := s.rsa.RSADecrypt(password)
		user.Password = string(buf)
		user.Password = common.MD5V(user.Password)
		s.database.C("user").Insert(user)
		*reply = true
	}
	return nil
}
func (s *UCServer) RPCRSA(_ string, reply *string) error {
	*reply = string(s.rsa.Publickey)
	return nil
}