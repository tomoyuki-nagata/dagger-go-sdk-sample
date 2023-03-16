package repository

import (
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://zipcloud.ibsnet.co.jp/api/search?zipcode=100-0001",
		httpmock.NewStringResponder(200, "mocked"),
	)
	res, err := Get()
	assert.Nil(t, err)
	assert.Equal(t, "mocked", res)
	// assert.Equal(t, "これは失敗するアサーション", res)
}
