package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "samplebank/db/mock"
	db "samplebank/db/sqlc"
	"samplebank/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt64(1, 1000),
		Username: util.RandomUsername(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}
func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	controller := gomock.NewController(t)

	store := mockdb.NewMockStore(controller)
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/acounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

}
