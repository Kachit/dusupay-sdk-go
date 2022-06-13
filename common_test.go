package dusupay

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

type CommonTestSuite struct {
	suite.Suite
}

func (suite *CommonTestSuite) TestTransformStructToMapSuccess() {
	req := &CollectionRequest{}
	req.Amount = 100
	req.ProviderId = "foo"
	req.Currency = CurrencyCodeEUR
	req.Method = TransactionMethodCard
	result, err := transformStructToMap(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), req.Amount, result["amount"])
	assert.Equal(suite.T(), req.ProviderId, result["provider_id"])
	assert.Equal(suite.T(), string(req.Method), result["method"])
	assert.Equal(suite.T(), string(req.Currency), result["currency"])
}

func (suite *CommonTestSuite) TestTransformStructToMapErrorUnmarshal() {
	result, err := transformStructToMap("foo")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *CommonTestSuite) TestTransformStructToMapErrorMarshal() {
	result, err := transformStructToMap(math.Inf(1))
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func TestCommonTestSuite(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}
