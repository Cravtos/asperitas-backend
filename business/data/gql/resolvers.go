package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func resolverHello(p graphql.ResolveParams) (interface{}, error) {
	return "World", nil
}

func resolverTitle(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(postDB)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Title, nil
}

func resolverAnyPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(Key).(Access)
	if !ok {
		return nil, errors.New("claims missing from context")
	}
	posts, err := a.selectAllPosts(p.Context)
	if err != nil {
		return nil, err
	}
	return posts[0], nil
}

func resolverAllPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(Key).(Access)
	if !ok {
		return nil, errors.New("claims missing from context")
	}
	posts, err := a.selectAllPosts(p.Context)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
