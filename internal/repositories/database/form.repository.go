package database

import (
	"Perwatch-case/internal/config"
	"Perwatch-case/internal/helpers"
	"Perwatch-case/internal/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormRepository struct {
	Collection *mongo.Collection
}

type FormRepositoryInterface interface {
	CreateForm(form models.Form) error
	GetForm(userID, formID primitive.ObjectID) (models.Form, error)
	GetForms(userID primitive.ObjectID, page int) ([]models.Form, error)
	DeleteForm(userID, formID primitive.ObjectID) error
	UpdateFormName(userID, formID primitive.ObjectID, formName string) error
	CheckFormNameExists(userID primitive.ObjectID, formName string) (bool, error)

	CreateFormField(userID, formID primitive.ObjectID, field models.Field) error
	GetFormFields(userID, formID primitive.ObjectID) (bson.A, error)
	GetFormField(userID, fieldID, formID primitive.ObjectID) (models.Field, error)
	DeleteFormField(userID, fieldID, formID primitive.ObjectID) error

	CheckFormFieldNameExists(userID, formID primitive.ObjectID, fieldName string) (bool, error)
	CheckFormFieldSortExists(userID, formID primitive.ObjectID, sort int) (bool, error)

	FindByFieldName(userID, formID primitive.ObjectID, fieldName string) ([]models.Field, error)
}

func NewFormRepository() FormRepositoryInterface {
	return &FormRepository{
		Collection: GetCollection(config.GetConfig().GetDatabaseConfig().Collections.Forms),
	}
}

func (f *FormRepository) CreateForm(form models.Form) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	_, err := f.Collection.InsertOne(ctx, form)
	if err != nil {
		return err
	}

	return nil
}

func (f *FormRepository) GetForm(userID, formID primitive.ObjectID) (models.Form, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"is_deleted", false}}

	var form models.Form
	if err := f.Collection.FindOne(ctx, filter).Decode(&form); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return form, errors.New("form not found")
		}

		return models.Form{}, err
	}

	return form, nil
}

func (f *FormRepository) GetForms(userID primitive.ObjectID, page int) ([]models.Form, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	perPage := int64(5)                             // Number of forms per page
	skip := int64(page*int(perPage) - int(perPage)) // Number of forms to skip
	findOptions := options.FindOptions{Limit: &perPage, Skip: &skip}

	var forms []models.Form
	filter := bson.D{{"user_id", userID}, {"is_deleted", false}}

	curr, err := f.Collection.Find(ctx, filter, &findOptions)
	if err != nil {
		return nil, err
	}

	if err := curr.All(ctx, &forms); err != nil {
		return nil, err
	}

	if err := curr.Err(); err != nil {
		return nil, err
	}

	if len(forms) == 0 {
		return nil, errors.New("forms not found")
	}

	return forms, nil
}

func (f *FormRepository) DeleteForm(userID, formID primitive.ObjectID) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}}
	result, err := f.Collection.UpdateOne(ctx, filter, bson.D{{"$set", bson.D{{"is_deleted", true}}}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("form not found")
		}

		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("form not found")
	}

	return nil
}

func (f *FormRepository) CheckFormNameExists(userID primitive.ObjectID, formName string) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"user_id", userID}, {"name", formName}, {"is_deleted", false}}
	projection := bson.D{{"_id", 1}}
	setProjection := options.FindOne().SetProjection(projection)

	var result bson.M
	if err := f.Collection.FindOne(ctx, filter, setProjection).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (f *FormRepository) UpdateFormName(userID, formID primitive.ObjectID, formName string) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"is_deleted", false}}
	update := bson.D{{"$set", bson.D{{"name", formName}}}}

	result, err := f.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("form not found")
		}

		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("form not found")
	}

	return nil
}

func (f *FormRepository) CreateFormField(userID, formID primitive.ObjectID, field models.Field) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"is_deleted", false}}
	result, err := f.Collection.UpdateOne(ctx, filter, bson.D{{"$push", bson.D{{"fields", field}}}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("form not found")
		}

		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("form not found")
	}

	return nil
}

func (f *FormRepository) GetFormFields(userID, formID primitive.ObjectID) (bson.A, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"is_deleted", false}}
	aggregate := bson.A{
		bson.D{{"$match", filter}},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"fields",
						bson.D{
							{"$sortArray",
								bson.D{
									{"input", "$fields"},
									{"sortBy", bson.D{{"sort", 1}}},
								},
							},
						},
					},
				},
			},
		},
	}

	var results []bson.M
	curr, err := f.Collection.Aggregate(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	if err := curr.All(ctx, &results); err != nil {
		return nil, err
	}

	if err := curr.Err(); err != nil {
		return nil, err
	}

	fields, ok := results[0]["fields"].(bson.A)
	if !ok {
		return nil, errors.New("fields not found")
	}

	return fields, nil
}

func (f *FormRepository) GetFormField(userID, fieldID, formID primitive.ObjectID) (models.Field, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"fields._id", fieldID}, {"is_deleted", false}}
	projection := bson.D{{"fields.$", 1}}
	setProjection := options.FindOne().SetProjection(projection)

	var form models.Form
	if err := f.Collection.FindOne(ctx, filter, setProjection).Decode(&form); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Field{}, errors.New("field not found")
		}

		return models.Field{}, err
	}

	if len(form.Fields) == 0 {
		return models.Field{}, errors.New("field not found")
	}

	return form.Fields[0], nil
}

func (f *FormRepository) DeleteFormField(userID, fieldID, formID primitive.ObjectID) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{
		{"_id", formID},
		{"user_id", userID},
		{"is_deleted", false},
	}
	update := bson.D{{
		"$pull",
		bson.D{{"fields", bson.D{{"_id", fieldID}}}},
	}}

	result, err := f.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("field not found")
		}

		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("field not found")
	}

	return nil
}

func (f *FormRepository) CheckFormFieldNameExists(userID, formID primitive.ObjectID, fieldName string) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"_id", formID}},
				bson.D{{"user_id", userID}},
				bson.D{{"is_deleted", false}},
				bson.D{{"fields.name", fieldName}},
			},
		},
	}
	projection := bson.D{{"_id", 1}}
	setProjection := options.FindOne().SetProjection(projection)

	var form models.Form
	if err := f.Collection.FindOne(ctx, filter, setProjection).Decode(&form); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (f *FormRepository) CheckFormFieldSortExists(userID, formID primitive.ObjectID, sort int) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"_id", formID}},
				bson.D{{"user_id", userID}},
				bson.D{{"is_deleted", false}},
				bson.D{{"fields.sort", sort}},
			},
		},
	}
	projection := bson.D{{"_id", 1}}
	setProjection := options.FindOne().SetProjection(projection)

	var result bson.M
	if err := f.Collection.FindOne(ctx, filter, setProjection).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (f *FormRepository) FindByFieldName(userID, formID primitive.ObjectID, fieldName string) ([]models.Field, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{
		{"_id", formID},
		{"user_id", userID},
		{"is_deleted", false},
		{"fields.name", fieldName},
	}
	projection := bson.D{{"fields.$", 1}}
	setProjection := options.FindOne().SetProjection(projection)

	var form models.Form
	if err := f.Collection.FindOne(ctx, filter, setProjection).Decode(&form); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []models.Field{}, errors.New("form not found")
		}

		return []models.Field{}, err
	}

	return form.Fields, nil
}
