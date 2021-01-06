package post

import "github.com/cravtos/asperitas-backend/business/auth"

func upvotePercentage(votes []Vote) int {
	var positive float32

	for _, vote := range votes {
		if vote.Vote == 1 {
			positive++
		}
	}

	return int(positive / float32(len(votes)) * 100)
}

func toInfo(post PostDB, claims auth.Claims) Info {
	var info Info
	if post.Type == "url" {
		info = InfoLink{
			ID:          post.ID,
			Score:       post.Score,
			Views:       post.Views,
			Title:       post.Title,
			Payload:     post.Payload,
			Category:    post.Category,
			DateCreated: post.DateCreated,
			Author: Author{
				Username: claims.User.Username,
				ID:       claims.User.ID,
			},
			Votes: []Vote{
				{User: claims.User.ID, Vote: 1},
			},
			Comments:         []Comment{},
			UpvotePercentage: 100,
		}
	} else {
		info = InfoText{
			ID:          post.ID,
			Score:       post.Score,
			Views:       post.Views,
			Title:       post.Title,
			Payload:     post.Payload,
			Category:    post.Category,
			DateCreated: post.DateCreated,
			Author: Author{
				Username: claims.User.Username,
				ID:       claims.User.ID,
			},
			Votes: []Vote{
				{User: claims.User.ID, Vote: 1},
			},
			Comments:         []Comment{},
			UpvotePercentage: 100,
		}
	}
	return info
}

// todo: function to make InfoLink or InfoText from all fields and payload
