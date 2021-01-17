package postgql

import (
	"github.com/graphql-go/graphql"
)

var (
	Schema graphql.Schema

	infoInterface *graphql.Interface
	postLinkType  *graphql.Object
	postTextType  *graphql.Object
	authorType    *graphql.Object
	voteType      *graphql.Object
	commentType   *graphql.Object
)

func Init() {
	categoryEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "Category",
		Values: graphql.EnumValueConfigMap{
			"ALL": &graphql.EnumValueConfig{
				Value: "all",
			},
			"MUSIC": &graphql.EnumValueConfig{
				Value: "music",
			},
			"FUNNY": &graphql.EnumValueConfig{
				Value: "funny",
			},
			"VIDEOS": &graphql.EnumValueConfig{
				Value: "videos",
			},
			"PROGRAMMING": &graphql.EnumValueConfig{
				Value: "programming",
			},
			"NEWS": &graphql.EnumValueConfig{
				Value: "news",
			},
			"FASHION": &graphql.EnumValueConfig{
				Value: "fashion",
			},
		},
	})
	authorType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type:    graphql.String,
				Resolve: authorUsername,
			},
			"id": &graphql.Field{
				Type:    graphql.ID,
				Resolve: authorID,
			},
		},
	})
	voteType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Vote",
		Fields: graphql.Fields{
			"vote": &graphql.Field{
				Type:    graphql.String,
				Resolve: voteVote,
			},
			"user_id": &graphql.Field{
				Type:    graphql.ID,
				Resolve: voteUserID,
			},
			"user": &graphql.Field{
				Type:    authorType,
				Resolve: voteUser,
			},
		},
	})
	commentType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"body": &graphql.Field{
				Type:    graphql.String,
				Resolve: commentBody,
			},
			"id": &graphql.Field{
				Type:    graphql.ID,
				Resolve: commentID,
			},
			"date_created": &graphql.Field{
				Type:    graphql.DateTime,
				Resolve: commentDateCreated,
			},
			"author": &graphql.Field{
				Type:    authorType,
				Resolve: commentAuthor,
			},
		},
	})
	infoInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Info",
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			post, _ := p.Value.(postDB)
			if post.Type == "url" {
				return postLinkType
			}
			return postTextType
		},
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"score": &graphql.Field{
				Type: graphql.Int,
			},
			"views": &graphql.Field{
				Type: graphql.Int,
			},
			"category": &graphql.Field{
				Type: graphql.String,
			},
			"date_created": &graphql.Field{
				Type: graphql.DateTime,
			},
			"upvote_percentage": &graphql.Field{
				Type: graphql.Int,
			},
			"author": &graphql.Field{
				Type: authorType,
			},
			"votes": &graphql.Field{
				Type: graphql.NewList(voteType),
			},
			"comments": &graphql.Field{
				Type: graphql.NewList(commentType),
			},
		},
	})

	postLinkType = graphql.NewObject(graphql.ObjectConfig{
		Name: "PostLink",
		Interfaces: []*graphql.Interface{
			infoInterface,
		},
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type:    graphql.String,
				Resolve: postTitle,
			},
			"url": &graphql.Field{
				Type:    graphql.String,
				Resolve: postURL,
			},
			"type": &graphql.Field{
				Type:    graphql.String,
				Resolve: postType,
			},
			"id": &graphql.Field{
				Type:    graphql.ID,
				Resolve: postID,
			},
			"score": &graphql.Field{
				Type:    graphql.Int,
				Resolve: postScore,
			},
			"views": &graphql.Field{
				Type:    graphql.Int,
				Resolve: postViews,
			},
			"category": &graphql.Field{
				Type:    graphql.String,
				Resolve: postCategory,
			},
			"date_created": &graphql.Field{
				Type:    graphql.DateTime,
				Resolve: postDateCreated,
			},
			"upvote_percentage": &graphql.Field{
				Type:    graphql.Int,
				Resolve: postUpvotePercentage,
			},
			"author": &graphql.Field{
				Type:    authorType,
				Resolve: postAuthor,
			},
			"votes": &graphql.Field{
				Type:    graphql.NewList(voteType),
				Resolve: postVotes,
			},
			"comments": &graphql.Field{
				Type:    graphql.NewList(commentType),
				Resolve: postComments,
			},
		},
	})

	postTextType = graphql.NewObject(graphql.ObjectConfig{
		Name: "PostText",
		Interfaces: []*graphql.Interface{
			infoInterface,
		},
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type:    graphql.String,
				Resolve: postTitle,
			},
			"text": &graphql.Field{
				Type:    graphql.String,
				Resolve: postText,
			},
			"type": &graphql.Field{
				Type:    graphql.String,
				Resolve: postType,
			},
			"id": &graphql.Field{
				Type:    graphql.ID,
				Resolve: postID,
			},
			"score": &graphql.Field{
				Type:    graphql.Int,
				Resolve: postScore,
			},
			"views": &graphql.Field{
				Type:    graphql.Int,
				Resolve: postViews,
			},
			"category": &graphql.Field{
				Type:    graphql.String,
				Resolve: postCategory,
			},
			"date_created": &graphql.Field{
				Type:    graphql.DateTime,
				Resolve: postDateCreated,
			},
			"upvote_percentage": &graphql.Field{
				Type:    graphql.Int,
				Resolve: postUpvotePercentage,
			},
			"author": &graphql.Field{
				Type:    authorType,
				Resolve: postAuthor,
			},
			"votes": &graphql.Field{
				Type:    graphql.NewList(voteType),
				Resolve: postVotes,
			},
			"comments": &graphql.Field{
				Type:    graphql.NewList(commentType),
				Resolve: postComments,
			},
		},
	})

	authorType.AddFieldConfig("posts", &graphql.Field{
		Type:    graphql.NewList(infoInterface),
		Resolve: authorPosts,
		Args: graphql.FieldConfigArgument{
			"category": &graphql.ArgumentConfig{
				Type: categoryEnum,
			},
		},
	})

	commentType.AddFieldConfig("post", &graphql.Field{
		Type:    infoInterface,
		Resolve: commentPost,
	})

	var queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"Hello": &graphql.Field{
				Type:    graphql.String,
				Resolve: Hello,
			},
			"anyPost": &graphql.Field{
				Type:    infoInterface,
				Resolve: anyPost,
			},
			"posts": &graphql.Field{
				Type:    graphql.NewList(infoInterface),
				Resolve: posts,
				Args: graphql.FieldConfigArgument{
					"category": &graphql.ArgumentConfig{
						Type: categoryEnum,
					},
					"user_id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
			},
		},
	})

	Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		Types: []graphql.Type{postTextType, postLinkType, authorType, commentType, voteType},
	})
}

func GetSchema() (graphql.Schema, error) {
	Init()
	return Schema, nil
}
