package product

import (
	"context"
	"errors"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductTemplateRepo struct {
	*db.Database
	collName string
}

func NewProductTemplateRepo(curDb *db.Database, collName string) *ProductTemplateRepo {
	store := &ProductTemplateRepo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *ProductTemplateRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductTemplateRepo) Create(ctx context.Context, paylaod ProductTemplateModel) error {

	productList := paylaod.ProductList
	productDetail := paylaod.ProductDetail

	objectIds, err := store.SliceOfHexToObjectID(productList.Detail, productList.Category)
	if err != nil {
		return err
	}

	productTemplate := ProductTemplate{
		Name: paylaod.Name,
		List: ProductList{
			Name:       productList.Name,
			Sizes:      productList.Sizes,
			Color:      productList.Color,
			Price:      productList.Price,
			ImgLink:    productList.ImgLink,
			Stock:      productList.Stock,
			Discount:   productList.Discount,
			Detail:     objectIds[0],
			Category:   objectIds[1],
			BatchId:    productList.BatchId,
			Gender:     productList.Gender,
			Collection: productList.Collection,
			Tags:       productList.Tags,
			Timestamps: store.GetCurrentTimestamps(),
		},
		Detail: ProductDetail{
			Name:        productList.Name,
			PreviewImg:  productDetail.PreviewImg,
			Description: productDetail.Description,
			Variations:  productDetail.Variations,
			ImgLink:     productDetail.ImgLink,
			BatchId:     productList.BatchId,
			Timestamps:  store.GetCurrentTimestamps(),
		},
	}

	result, err := s.getColl().InsertOne(ctx, productTemplate)

	if mongo.IsDuplicateKeyError(err) {
		return errutil.ErrDocumentAlreadyExist
	} else if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(bson.ObjectID); ok {
		return nil
	}

	return errutil.ErrDatabase
}

func (s *ProductTemplateRepo) FindByFilter(ctx context.Context, filterProductListModel FilterProductListModel, project []string, inclusive bool) ([]ProductList, error) {

	var productList []ProductList

	query := bson.M{}

	name := filterProductListModel.Name
	page := filterProductListModel.Page
	limit := filterProductListModel.Limit

	if len(name) != 0 {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	findOptions := options.Find().SetSkip(int64(limit * (page - 1))).SetLimit(int64(limit))
	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOptions.SetProjection(projection)
	}

	cursor, err := s.getColl().Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &productList); err != nil {
		return nil, err
	}

	return productList, nil
}

func (s *ProductTemplateRepo) FindById(ctx context.Context, id string, project []string, inclusive bool) (ProductTemplate, error) {

	var productDetail ProductTemplate

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return productDetail, err
	}

	filter := bson.M{"_id": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&productDetail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return productDetail, errutil.NotFound("Product")
		}
		return productDetail, err
	}

	return productDetail, nil

}
