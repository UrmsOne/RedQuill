package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"redquill-backend/pkg/common"
	"redquill-backend/pkg/models"
)

type UserService struct {
	client *mongo.Client
	dbName string
}

func NewUserService(client *mongo.Client, dbName string) *UserService {
	return &UserService{client: client, dbName: dbName}
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	coll := s.client.Database(s.dbName).Collection("users")
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []models.User
	for cur.Next(ctx) {
		var u models.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, cur.Err()
}

type PagedUsers struct {
	Items      []models.User     `json:"items"`
	Pagination common.Pagination `json:"pagination"`
}

func (s *UserService) ListPaged(ctx context.Context, page, pageSize int64, sortExpr, keyword string) (PagedUsers, error) {
	coll := s.client.Database(s.dbName).Collection("users")

	// Build filter: keyword on name/email
	kwFilter := common.BuildKeywordFilter(keyword, []string{"name", "email"})
	filter := common.MergeFilters(bson.M{}, kwFilter)

	// Options: pagination + sort
	sort := common.BuildSort(sortExpr)
	opts := common.BuildFindOptions(page, pageSize, sort, bson.M{})

	items, total, err := common.FindWithPagination[models.User](ctx, coll, filter, opts)
	if err != nil {
		return PagedUsers{}, err
	}

	totalPages := total / common.NormalizePageSize(pageSize)
	if total%common.NormalizePageSize(pageSize) != 0 {
		totalPages++
	}

	return PagedUsers{
		Items: items,
		Pagination: common.Pagination{
			Page:      common.NormalizePage(page),
			PageSize:  common.NormalizePageSize(pageSize),
			Total:     total,
			TotalPage: totalPages,
		},
	}, nil
}

// PostUsers creates a new user (registration-like behavior)
func (s *UserService) PostUsers(ctx context.Context, name, email, password string) (models.User, error) {
	coll := s.client.Database(s.dbName).Collection("users")
	// unique email constraint at application level (ensure index in DB ideally)
	var existing models.User
	if err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&existing); err == nil {
		return models.User{}, errors.New("email already exists")
	}
	now := time.Now()
	u := models.User{
		Name:     name,
		Email:    email,
		Password: common.HashPassword(password),
		Ctime:    now.Unix(),
		Mtime:    now.Unix(),
	}
	res, err := coll.InsertOne(ctx, u)
	if err != nil {
		return models.User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid.Hex()
	}
	u.Password = ""
	return u, nil
}

// GetUsers returns user by id
func (s *UserService) GetUsers(ctx context.Context, id string) (models.User, error) {
	coll := s.client.Database(s.dbName).Collection("users")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid id")
	}
	var u models.User
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&u); err != nil {
		return models.User{}, err
	}
	u.Password = ""
	return u, nil
}

// PutUsers updates allowed fields; nil pointer means no change
func (s *UserService) PutUsers(ctx context.Context, id string, name *string, email *string, password *string) (models.User, error) {
	coll := s.client.Database(s.dbName).Collection("users")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid id")
	}
	update := bson.M{"mtime": time.Now().Unix()}
	set := bson.M{}
	if name != nil {
		set["name"] = *name
	}
	if email != nil {
		set["email"] = *email
	}
	if password != nil && *password != "" {
		set["password"] = common.HashPassword(*password)
	}
	for k, v := range set {
		update[k] = v
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out models.User
	if err := coll.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&out); err != nil {
		return models.User{}, err
	}
	out.Password = ""
	return out, nil
}

// DeleteUsers removes a user by id
func (s *UserService) DeleteUsers(ctx context.Context, id string) error {
	coll := s.client.Database(s.dbName).Collection("users")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	_, err = coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// Authenticate checks credentials and returns the user when ok
func (s *UserService) Authenticate(ctx context.Context, email, password string) (models.User, error) {
	coll := s.client.Database(s.dbName).Collection("users")
	var u models.User
	if err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
		return models.User{}, errors.New("invalid email or password")
	}
	if !common.ComparePassword(u.Password, password) {
		return models.User{}, errors.New("invalid email or password")
	}
	u.Password = ""
	return u, nil
}
