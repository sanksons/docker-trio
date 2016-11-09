package rating

type UploadRatingHeathCheck struct {
}

func (n UploadRatingHeathCheck) GetName() string {
	return "Seller get api health check successfull"
}

func (n UploadRatingHeathCheck) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"status": "success",
	}
}
