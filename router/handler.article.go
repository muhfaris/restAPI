package router

import (
	"encoding/json"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/muhfaris/restAPI/api"
	"github.com/muhfaris/restAPI/model"
	"github.com/pkg/errors"
)

func HandlerArticleList(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	//	ctx := r.Context()
	ctx := r.Context()
	response, _ := model.GetAllArticle(ctx, dbPool)

	return response, nil
}

func HandlerArticleDetail(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	params := mux.Vars(r)

	result, err := model.GetDetailArticle(ctx, dbPool, params["id"])
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "article/detail"), "", http.StatusNotFound)
	}

	return result, nil
}

func HandlerArticleUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	params := mux.Vars(r)

	article := model.ModelArticles{}
	id := params["id"]
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		return nil, api.NewError(errors.Wrap(err, "article/create/parse_requet_data"), "", http.StatusNotAcceptable)
	}

	idData := article.ID.Hex()
	if id != "" && id != idData {
		return nil, api.NewError(errors.Wrap(errors.New("ID article is not found"), "article/update/get_id"), "", http.StatusNotFound)
	}

	if err := article.Update(ctx, dbPool); err != nil {
		return nil, api.NewError(errors.Wrap(err, "article/update"), "", http.StatusBadRequest)
	}

	return http.StatusNoContent, nil
}

func HandlerArticleDelete(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	params := mux.Vars(r)

	article := model.ModelArticles{}
	id := params["id"]

	article.ID = bson.ObjectIdHex(id)

	if err := article.Delete(ctx, dbPool); err != nil {
		return nil, api.NewError(errors.Wrap(err, "article/delete"), "", http.StatusBadRequest)
	}

	return http.StatusNoContent, nil
}

func HandlerArticleCreate(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	article := model.ModelArticles{}
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		return nil, api.NewError(errors.Wrap(err, "article/create/parse_requet_data"), "", http.StatusNotAcceptable)
	}

	ctx := r.Context()
	if err := article.Create(ctx, dbPool); err != nil {
		return nil, api.NewError(errors.Wrap(err, "article/create"), "", http.StatusBadRequest)
	}

	return http.StatusCreated, nil
}
