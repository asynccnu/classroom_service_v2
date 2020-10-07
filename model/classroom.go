package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UnavailableClassrooms struct {
	Week     [2]int
	Weekday  int
	Time     [2]int
	WeekType string // 单双周
	Place    string
}

func InsertClassroomsInDB(instance *Classrooms) error {
	collection := DB.Self.Database(DBName).Collection(ClassroomCol)
	_, err := collection.InsertOne(context.TODO(), *instance)

	return err
}

func UpdateAvailableClassroomInDB(newClassrooms *Classrooms) error {
	collection := DB.Self.Database(DBName).Collection(ClassroomCol)
	_, err := collection.ReplaceOne(context.TODO(),
		bson.M{"week": newClassrooms.Week, "weekday": newClassrooms.Weekday, "building": newClassrooms.Building},
		*newClassrooms)

	return err
}

func GetClassroomsFromDB(week int, weekday int, building string) (*Classrooms, error) {
	classroom := Classrooms{}

	err := DB.Self.Database(DBName).Collection(ClassroomCol).
		FindOne(context.TODO(), bson.M{"week": week, "weekday": weekday, "building": building}).
		Decode(&classroom)

	if err != nil {
		return nil, err
	}

	return &classroom, nil
}
