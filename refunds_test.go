package dusupay

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

type RefundsTestSuite struct {
	suite.Suite
}

func (suite *RefundsTestSuite) TestRefundsRequestIsValidSuccess() {
	request := RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	assert.Nil(suite.T(), request.isValid())
	assert.NoError(suite.T(), request.isValid())
}

func (suite *RefundsTestSuite) TestRefundsRequestIsValidEmptyInternalReference() {
	request := RefundRequest{
		Amount: 100,
	}
	result := request.isValid()
	assert.Error(suite.T(), result)
	assert.Equal(suite.T(), `parameter "internal_reference" is empty`, result.Error())
}

func TestRefundsTestSuite(t *testing.T) {
	suite.Run(t, new(RefundsTestSuite))
}

type RefundsResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *RefundsResource
}

func (suite *RefundsResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHttpTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &RefundsResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *RefundsResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *RefundsResourceTestSuite) TestCreateSuccess() {
	body, _ := LoadStubResponseData("stubs/refunds/create/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/refund", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusAccepted, result.Code)
	assert.Equal(suite.T(), "accepted", result.Status)
	assert.Equal(suite.T(), "Refund Initiated Successfully", result.Message)
	assert.Equal(suite.T(), int64(65205), result.Data.ID)
	assert.Equal(suite.T(), float64(1054), result.Data.RefundAmount)
	assert.Equal(suite.T(), "UGX", result.Data.RefundCurrency)
	assert.Equal(suite.T(), float64(0), result.Data.TransactionFee)
	assert.Equal(suite.T(), float64(1054), result.Data.TotalDebit)
	assert.Equal(suite.T(), "international_ugx", result.Data.ProviderID)
	assert.Equal(suite.T(), "hAkEROAdhIsHrEnB", result.Data.MerchantReference)
	assert.Equal(suite.T(), "DUSUPAYXYXYXYXYXYXYXYXYX", result.Data.CollectionReference)
	assert.Equal(suite.T(), "RFD-DUSUPAYXYXYXYXYXYXYXYXYX-3486003", result.Data.InternalReference)
	assert.Equal(suite.T(), "PENDING", result.Data.TransactionStatus)
	assert.Equal(suite.T(), "refund", result.Data.TransactionType)
	assert.Equal(suite.T(), "4860610032773134", result.Data.AccountNumber)
	assert.Equal(suite.T(), "Request Initiated", result.Data.Message)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *RefundsResourceTestSuite) TestCreateJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/refund", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.NotEmpty(suite.T(), result)
	//result
	assert.False(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusUnauthorized, result.Code)
	assert.Equal(suite.T(), "error", result.Status)
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", result.Message)
	assert.Empty(suite.T(), result.Data)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
	//error
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}

func (suite *RefundsResourceTestSuite) TestCreateNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.Uri+"/v1/refund", httpmock.NewBytesResponder(http.StatusOK, body))

	request := &RefundRequest{
		Amount:            100,
		InternalReference: "internal_reference",
	}
	result, resp, err := suite.testable.Create(suite.ctx, request)
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
	//response
	defer resp.Body.Close()
	bodyRsp, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}

func (suite *RefundsResourceTestSuite) TestCreateInvalidRequest() {
	req := &RefundRequest{}
	result, rsp, err := suite.testable.Create(suite.ctx, req)
	assert.Nil(suite.T(), result)
	assert.Nil(suite.T(), rsp)
	assert.Error(suite.T(), err)
}

func TestRefundsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(RefundsResourceTestSuite))
}
