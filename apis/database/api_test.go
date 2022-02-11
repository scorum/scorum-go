package database

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	scorumhttp "github.com/scorum/scorum-go/transport/http"
	"github.com/stretchr/testify/require"
)

const nodeHTTPS = "https://testnet.scorum.work"

func TestGetAccountsCount(t *testing.T) {
	transport := scorumhttp.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	count, err := api.GetAccountsCount(context.Background())
	require.NoError(t, err)
	require.True(t, count > 0)
}

func TestGetConfig(t *testing.T) {
	transport := scorumhttp.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	config, err := api.GetConfig(context.Background())
	require.NoError(t, err)
	require.Equal(t, "SCR", config.ScorumAddressPrefix)
}

func TestLookupAccounts(t *testing.T) {
	transport := scorumhttp.NewTransport(nodeHTTPS)
	api := NewAPI(transport)

	t.Run("from beginning", func(t *testing.T) {
		accounts, err := api.LookupAccounts(context.Background(), "", 1000)
		require.NoError(t, err)
		require.True(t, len(accounts) > 0)
	})

	t.Run("from 'bebe'", func(t *testing.T) {
		accounts, err := api.LookupAccounts(context.Background(), "bebe", 1000)
		t.Log(accounts)
		require.NoError(t, err)
		require.True(t, len(accounts) > 0)
	})

	t.Run("get all cursor", func(t *testing.T) {
		const limit = 5

		var (
			all, add   []string
			lowerBound string
		)

		for {
			accounts, err := api.LookupAccounts(context.Background(), lowerBound, limit)
			require.NoError(t, err)
			if lowerBound == "" {
				add = accounts[:]
			} else {
				add = accounts[1:]
			}

			all = append(all, add...)
			if len(add) == 0 {
				break
			}
			lowerBound = all[len(all)-1]
		}

		count, err := api.GetAccountsCount(context.Background())
		require.NoError(t, err)

		require.Equal(t, count, len(all))
	})

	t.Run("exceeded limit", func(t *testing.T) {
		_, err := api.LookupAccounts(context.Background(), "", 2000)
		require.Error(t, err)
	})
}

func TestAPI_GetAccounts_ReturnNoError_When_ResponseIsNotJSONRPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"comment_count": "9,223,372,036,854,775,807"}`))
		writer.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	transport := scorumhttp.NewTransport(server.URL)
	api := NewAPI(transport)
	_, err := api.GetAccounts(context.Background(), "leo")
	require.NoError(t, err)
}

func TestAPI_GetAccounts_ReturnError_When_CouldNotUnmarshallJSONRPCResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(`{"id":0,"result":[{"comment_count": "9,223,372,036,854,775,807"}]}`))
	}))

	defer server.Close()

	transport := scorumhttp.NewTransport(server.URL)
	api := NewAPI(transport)
	_, err := api.GetAccounts(context.Background(), "leo")

	require.EqualError(t, err, "failed to unmarshal rpc result: [{\"comment_count\": \"9,223,372,036,854,775,807\"}]: json: cannot unmarshal string into Go struct field Account.comment_count of type int32")
}
