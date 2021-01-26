package postgql

func parseCatAndUser(category interface{}, userID interface{}) (string, string) {
	if category == nil && userID == nil {
		return "all", ""
	}
	if category == nil {
		if user, ok := userID.(string); ok {
			return "all", user
		}
		return "all", ""
	}
	if userID == nil {
		if cat, ok := category.(string); ok {
			return cat, ""
		}
		return "all", ""
	}
	if cat, ok := category.(string); ok {
		if user, ok := userID.(string); ok {
			return cat, user
		}
		return cat, ""
	}
	if user, ok := userID.(string); ok {
		return "all", user
	}
	return "all", ""

}
