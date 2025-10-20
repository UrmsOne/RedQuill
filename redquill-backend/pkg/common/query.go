package common

import (
	"context"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DefaultPageSize int64 = 20
	MaxPageSize     int64 = 100
)

type Pagination struct {
	Page      int64 `json:"page"`
	PageSize  int64 `json:"pageSize"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"totalPage"`
}

func NormalizePage(page int64) int64 {
	if page < 1 {
		return 1
	}
	return page
}

func NormalizePageSize(pageSize int64) int64 {
	if pageSize <= 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

func ComputeSkipLimit(page, pageSize int64) (int64, int64) {
	page = NormalizePage(page)
	pageSize = NormalizePageSize(pageSize)
	return (page - 1) * pageSize, pageSize
}

func BuildSort(sortExpr string) bson.D {
	if strings.TrimSpace(sortExpr) == "" {
		return bson.D{}
	}
	parts := strings.Split(sortExpr, ",")
	sort := bson.D{}
	for _, p := range parts {
		field := strings.TrimSpace(p)
		if field == "" {
			continue
		}
		dir := int32(1)
		if strings.HasPrefix(field, "-") {
			dir = -1
			field = strings.TrimPrefix(field, "-")
		}
		sort = append(sort, bson.E{Key: field, Value: dir})
	}
	return sort
}

func BuildKeywordFilter(keyword string, fields []string) bson.M {
	kw := strings.TrimSpace(keyword)
	if kw == "" || len(fields) == 0 {
		return bson.M{}
	}
	escaped := regexp.QuoteMeta(kw)
	pattern := ".*" + escaped + ".*"
	ors := make([]bson.M, 0, len(fields))
	for _, f := range fields {
		field := strings.TrimSpace(f)
		if field == "" {
			continue
		}
		ors = append(ors, bson.M{field: bson.M{"$regex": pattern, "$options": "i"}})
	}
	if len(ors) == 0 {
		return bson.M{}
	}
	return bson.M{"$or": ors}
}

func BuildFindOptions(page, pageSize int64, sort bson.D, projection bson.M) *options.FindOptions {
	skip, limit := ComputeSkipLimit(page, pageSize)
	fo := options.Find().SetSkip(skip).SetLimit(limit)
	if len(sort) > 0 {
		fo.SetSort(sort)
	}
	if len(projection) > 0 {
		fo.SetProjection(projection)
	}
	return fo
}

func ParseCommonQueryParams(values url.Values) (page int64, pageSize int64, sortExpr string, keyword string) {
	page = NormalizePage(parseInt64(values.Get("page"), 1))
	pageSize = NormalizePageSize(parseInt64(values.Get("pageSize"), DefaultPageSize))
	sortExpr = values.Get("sort")
	keyword = values.Get("q")
	return
}

func parseInt64(s string, def int64) int64 {
	if s == "" {
		return def
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	return def
}

func FindWithPagination[T any](ctx context.Context, coll *mongo.Collection, filter interface{}, opts *options.FindOptions) ([]T, int64, error) {
	cur, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var out []T
	for cur.Next(ctx) {
		var item T
		if err := cur.Decode(&item); err != nil {
			return nil, 0, err
		}
		out = append(out, item)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func MergeFilters(filters ...bson.M) bson.M {
	ands := make([]bson.M, 0, len(filters))
	for _, f := range filters {
		if len(f) == 0 {
			continue
		}
		ands = append(ands, f)
	}
	if len(ands) == 0 {
		return bson.M{}
	}
	if len(ands) == 1 {
		return ands[0]
	}
	return bson.M{"$and": ands}
}


