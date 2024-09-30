package repository

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bfbarry/CollabSource/back-end/model"
	"github.com/bfbarry/CollabSource/back-end/mongoClient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	mongoClient *mongo.Database
}

var mongoRepository *Repository

func GetMongoRepository() *Repository {
	return mongoRepository
}

func init() {
	// defer mongoClient.CloseMongoClient()
	mongoRepository = &Repository{mongoClient: mongoClient.GetMongoDb()}
}

func (self *Repository) getCollection(coll string) *mongo.Collection {
	return self.mongoClient.Collection(coll)
}

func (self *Repository) Insert(coll string, obj model.Model) (primitive.ObjectID, error) {

	res, mongoerr := self.getCollection(coll).InsertOne(context.TODO(), obj)
	if mongoerr != nil {
		fmt.Printf("Error inserting object e message: %s", mongoerr)
		return primitive.NilObjectID, mongoerr
	}
	switch val := res.InsertedID.(type) { //"type switch"
	case primitive.ObjectID:
		return val, nil
	default:
		return primitive.NilObjectID, fmt.Errorf("Expected an ObjectID, received %T instead", val)
	}
}

func (self *Repository) FindByID(coll string, id primitive.ObjectID, obj model.Model) error {
	// var op errors.Op = "repository.FindByID"

	filter := bson.M{"_id": id}
	err := self.getCollection(coll).FindOne(context.TODO(), filter).Decode(obj)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return err
		default:
			return err
		}
	}
	return nil
}

func (self *Repository) FindOne(coll string, uniqueField bson.M, obj model.Model) error {
	// TODO maybe constraints on name to ensure it's one field in the model
	err := self.getCollection(coll).FindOne(context.TODO(), uniqueField).Decode(obj)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return err
		default:
			return err
		}
	}
	return nil
}

func (self *Repository) Update(coll string, id primitive.ObjectID, obj model.Model) (int64, error) {
	// var op errors.Op = "repository.Update"

	result, err := self.getCollection(coll).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": obj})
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (self *Repository) Delete(coll string, id primitive.ObjectID) (int64, error) {
	//var op errors.Op = "repository.Delete"

	// var del_err error
	// switch deleteMode {
	// case SoftDelete:
	// 	_, del_err = self.getCollection(coll).UpdateOne(context.TODO(), bson.M{"_id": objId}, bson.M{"$set": bson.M{"deleted": true}})
	// case HardDelete:
	result, err := self.getCollection(coll).DeleteOne(context.TODO(), bson.M{"_id": id})
	// }
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (self *Repository) FindManyByPage(coll string, results interface{}, pageNum int, pageSize int, filter bson.M) (bool, error) {

	findOptions := options.Find()
	skip := (pageNum - 1) * pageSize
	findOptions.SetLimit(int64(pageSize+1))
	findOptions.SetSkip(int64(skip))
	cursor, findErr := self.getCollection(coll).Find(context.TODO(), filter, findOptions)
	if findErr != nil {
		return false, findErr
	}

    // Get the reflect Value of results (which should be a pointer to a slice)
    sliceValue := reflect.ValueOf(results)
    if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
        return false, fmt.Errorf("results argument must be a pointer to a slice")
    }

    sliceElemType := sliceValue.Elem().Type().Elem()
	count := 0
    for cursor.Next(context.TODO()) {
		count++
		if count <= pageSize {
			result := reflect.New(sliceElemType).Elem()
	
			if err := cursor.Decode(result.Addr().Interface()); err != nil {
				fmt.Printf(err.Error())
				return false, err
			}
	
			sliceValue.Elem().Set(reflect.Append(sliceValue.Elem(), result))
		}
    }

	var hasNext bool
	if count == pageSize+1 {
		hasNext = true
	} else {
		hasNext = false
	}
	return hasNext, nil
}

func (self *Repository) FindManyByJunction(coll string, fromKey string, fromKeyVal primitive.ObjectID, toKey string, toColl string, 
											pageNum int, pageSize int, results interface{}) (bool, error) {
	skip := (pageNum - 1) * pageSize
	newRoot := fmt.Sprintf("$%s", toColl)
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{fromKey, fromKeyVal}}}},
		{{"$lookup", bson.D{
			{"from", toColl},              
			{"localField", toKey},      
			{"foreignField", "_id"},           
			{"as", toColl},                
		}}},
		{{"$unwind", newRoot}},
		{{"$replaceRoot", bson.D{{"newRoot", newRoot}}}},
		{{"$skip", skip}},                    
		{{"$limit", pageSize+1}},                
	}

	cursor, err := self.getCollection(coll).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return false, err
	}
	defer cursor.Close(context.TODO())

    // Get the reflect Value of results (which should be a pointer to a slice)
    sliceValue := reflect.ValueOf(results)
    if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
        return false, fmt.Errorf("results argument must be a pointer to a slice")
    }

    sliceElemType := sliceValue.Elem().Type().Elem()
	count := 0
    for cursor.Next(context.TODO()) {
		count++
		if count <= pageSize {
			result := reflect.New(sliceElemType).Elem()
	
			if err := cursor.Decode(result.Addr().Interface()); err != nil {
				fmt.Printf(err.Error())
				return false, err
			}
	
			sliceValue.Elem().Set(reflect.Append(sliceValue.Elem(), result))
		}
    }

	var hasNext bool
	if count == pageSize+1 {
		hasNext = true
	} else {
		hasNext = false
	}
	return hasNext, nil
}
