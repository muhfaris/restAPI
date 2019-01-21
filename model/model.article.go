package model

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

const collection = "articles"

type (
	modelLinks struct {
		Self    string `json:"self,omitempty"`
		Next    string `json:"next,omitempty"`
		Last    string `json:"last,omitempty"`
		Related string `json:"related,omitempty"`
	}

	modelAttributes struct {
		Title string `json:"title,omitempty"`
	}

	modelAuthor struct {
		Links modelLinks      `json:"links,omitempty"`
		Data  modelAuthorData `json:"data,omitempty"`
	}

	modelAuthorData struct {
		Type string `json:"type,omitempty"`
		ID   string `json:"id,omitempty"`
	}

	modelRelationshipsComment struct {
		Links modelLinks `json:"links,omitempty"`
		Data  []struct {
			Type string `json:"type,omitempty"`
			ID   string `json:"id,omitempty"`
		} `json:"data,omitempty"`
	}

	modelIncludedAttributes struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Twitter   string `json:"twitter,omitempty"`
	}

	ModelArticles struct {
		ID    bson.ObjectId `json:"_id" bson:"_id,omitempty"`
		Links modelLinks    `json:"links,omitempty"`
		Data  []struct {
			Type          string          `json:"type"`
			ID            string          `json:"id"`
			Attributes    modelAttributes `json:"attributes"`
			Relationships struct {
				Author   modelAuthor               `json:"author,omitempty"`
				Comments modelRelationshipsComment `json:"comments,omitempty"`
			} `json:"relationships,omitempty"`

			Links modelLinks `json:"links,omitempty"`
		} `json:"data,omitempty"`

		Included []struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				FirstName string `json:"firstName,omitempty" bson:"firstName"`
				LastName  string `json:"lastName,omitempty" bson:"lastName"`
				Twitter   string `json:"twitter,omitempty"`
				Facebook  string `json:"facebook,omitempty"`
			} `json:"attributes,omitempty"`

			Links modelLinks `json:"links,omitempty"`

			Relationships struct {
				Author modelAuthor `json:"author,omitempty"`
			} `json:"relationships,omitempty"`
		} `json:"included"`
	}
)

func GetAllArticle(ctx context.Context, db *mgo.Database) ([]ModelArticles, error) {
	c := db.C(collection)

	var result []ModelArticles
	err := c.Find(nil).All(&result)
	if err != nil {
		return nil, errors.Wrap(err, "article/list")
	}
	return result, nil
}

func GetDetailArticle(ctx context.Context, db *mgo.Database, id string) (*ModelArticles, error) {
	c := db.C(collection)

	fmt.Println(id)
	var result ModelArticles
	err := c.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return &result, errors.Wrap(err, "article/detail")
	}
	return &result, nil
}

func (m *ModelArticles) Create(ctx context.Context, db *mgo.Database) error {
	c := db.C(collection)
	m.ID = bson.NewObjectId()

	if err := c.Insert(m); err != nil {
		return errors.Wrap(err, "article/insert")
	}

	return nil
}

func (m *ModelArticles) Update(ctx context.Context, db *mgo.Database) error {
	c := db.C(collection)

	if err := c.UpdateId(m.ID, &m); err != nil {
		return errors.Wrap(err, "article/insert")
	}

	return nil
}
