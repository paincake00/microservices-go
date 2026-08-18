package entity

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}
