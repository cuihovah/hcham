package advert

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func (s *Server) RPCGetTextRecommends(prop map[string]int, reply *[]AdvertData) error {
	if len(prop) == 0 {
		pipe := s.database.C("advert").Pipe([]bson.M{
			bson.M{"$match": bson.M{"contents": bson.M{"$ne": ""}}},
			bson.M{"$sample": bson.M{"size": 10}},
		})
		err := pipe.All(reply)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else {
		adverts := s.recommend(prop, 5)
		*reply = adverts
	}
	return nil
}

func (s *Server) RPCGetIMaxRecommends(prop int, reply *[]AdvertData) error {
	pipe := s.database.C("advert").Pipe([]bson.M{
		bson.M{"$match": bson.M{"level": "imax"}},
		bson.M{"$sample": bson.M{"size": prop}},
	})
	err := pipe.All(reply)
	return err
}



