package gql

import (
	"github.com/graphql-go/graphql"
)

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"title": &graphql.Field{
			Type:    graphql.String,
			Resolve: resolverTitle,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"Hello": &graphql.Field{
			Type:    graphql.String,
			Resolve: resolverHello,
		},
		"anyPost": &graphql.Field{
			Type:    postType,
			Resolve: resolverAnyPost,
		},
		"allPosts": &graphql.Field{
			Type:    graphql.NewList(postType),
			Resolve: resolverAllPost,
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})
