package advert

type AdvertData struct {
	Id       string  `json:"id" bson:"_id"`
	Name     string  `json:"name" bson:"name"`
	Contents string  `json:"contents" bson:"contents"`
	Link     string  `json:"link" bson:"link"`
	Image    string  `json:"image" bson:"image"`
	UserId   string  `json:"user_id" bson:"user_id"`
	Price    float32 `json:"price" bson:"price"`
	Level    string  `json:"level" bson:"level"`
	Type     string  `json:"type" bson:"type"`
}
