package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// 新建教室文档数据
func CreateClassroomDoc(instance *ClassroomModel) error {
	collection := DB.Self.Database(DBName).Collection(ClassroomCol)
	_, err := collection.InsertOne(context.TODO(), instance)

	return err
}

// 批量新建教室文档数据
func CreateMultipleClassroomDocs(instances []*ClassroomModel) error {
	var docs []interface{}
	for _, instance := range instances {
		docs = append(docs, instance)
	}

	_, err := DB.Self.Database(DBName).Collection(ClassroomCol).InsertMany(context.TODO(), docs)
	return err
}

// 更新教室文档
func UpdateClassroom(instance *ClassroomModel) error {
	collection := DB.Self.Database(DBName).Collection(ClassroomCol)
	_, err := collection.ReplaceOne(
		context.TODO(),
		bson.M{"week": instance.Week, "weekday": instance.Day, "building": instance.Building},
		instance,
	)

	return err
}

// 获取文档数据
func GetClassroomDoc(week, day int, building string) (*ClassroomModel, error) {
	var classroom ClassroomModel

	err := DB.Self.Database(DBName).Collection(ClassroomCol).
		FindOne(context.TODO(), bson.M{"week": week, "day": day, "building": building}).
		Decode(&classroom)

	if err != nil {
		return nil, err
	}

	return &classroom, nil
}
