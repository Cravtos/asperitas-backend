package postgql

import (
	"github.com/cravtos/asperitas-backend/business/data/posts"
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
	authType      *graphql.Object
	userType      *graphql.Object
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

	postTypeEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "Type",
		Values: graphql.EnumValueConfigMap{
			"LINK": &graphql.EnumValueConfig{
				Value: "link",
			},
			"TEXT": &graphql.EnumValueConfig{
				Value: "text",
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
			post, _ := p.Value.(posts.Info)
			if post.Type == "url" || post.Type == "link" {
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
				Resolve: postsRes,
				Args: graphql.FieldConfigArgument{
					"category": &graphql.ArgumentConfig{
						Type: categoryEnum,
					},
					"user_id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
			},
			"post": &graphql.Field{
				Type:    infoInterface,
				Resolve: postRes,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
			},
		},
	})

	userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type:    graphql.String,
				Resolve: username,
			},
			"user_id": &graphql.Field{
				Type:    graphql.ID,
				Resolve: userID,
			},
		},
	})

	authType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Auth",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type:    graphql.String,
				Resolve: token,
			},
			"user": &graphql.Field{
				Type:    userType,
				Resolve: authUser,
			},
		},
	})

	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create_post": &graphql.Field{
				Type:    infoInterface,
				Resolve: postCreate,
				Args: graphql.FieldConfigArgument{
					"type": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(postTypeEnum),
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"category": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(categoryEnum),
					},
					"payload": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
			},
			"delete_post": &graphql.Field{
				Type:    infoInterface,
				Resolve: postDelete,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
			},
			"create_comment": &graphql.Field{
				Type:    infoInterface,
				Resolve: commentCreate,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
			},
			"delete_comment": &graphql.Field{
				Type:    infoInterface,
				Resolve: commentDelete,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{ //I don`t know why we ask for post_id
						Type: graphql.NewNonNull(graphql.ID),
					},
					"comment_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
			},
			"upvote": &graphql.Field{
				Type:    infoInterface,
				Resolve: upvote,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
			},
			"downvote": &graphql.Field{
				Type:    infoInterface,
				Resolve: downvote,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
			},
			"unvote": &graphql.Field{
				Type:    infoInterface,
				Resolve: unvote,
				Args: graphql.FieldConfigArgument{
					"post_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
			},
			"register": &graphql.Field{
				Type:    authType,
				Resolve: register,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
			},

			"sign_in": &graphql.Field{
				Type:    authType,
				Resolve: signIn,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
			},
		},
	})

	_ = mutationType

	Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
		Types:    []graphql.Type{postTextType, postLinkType, authorType, commentType, voteType, authType, userType},
	})
}

func GetSchema() (graphql.Schema, error) {
	Init()
	return Schema, nil
}
