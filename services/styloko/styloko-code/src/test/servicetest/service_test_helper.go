package servicetest

import (
	"encoding/json"
	utilhttp "github.com/jabong/floRest/src/common/utils/http"
	gm "github.com/onsi/gomega"
)

func validateHealthCheckResponse(responseBody string) {
	var utilResponse utilhttp.Response
	err := json.Unmarshal([]byte(responseBody), &utilResponse)
	gm.Expect(err).To(gm.BeNil())

	utilResponseData := utilResponse.Data
	if v, ok := utilResponseData.(map[string]interface{}); ok {
		_, serviceNodePresent := v["service"]
		gm.Expect(serviceNodePresent).To(gm.Equal(true))
	}
}
