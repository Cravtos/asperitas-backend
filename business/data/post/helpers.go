package post

func UpvotePercentage(votes []Vote) int {
	var positive float32

	for _, vote := range votes {
		if vote.Vote == 1 {
			positive++
		}
	}

	return int(positive / float32(len(votes)) * 100)
}

// todo: function to make InfoLink or InfoText from all fields and payload