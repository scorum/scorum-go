package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scorum/scorum-go/rpc"
)

const nodeHTTPS = "https://testnet.scorum.work"

func TestGetAccountsCount(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	count, err := api.GetAccountsCount(context.Background())
	require.NoError(t, err)
	require.True(t, count > 0)
}

func TestGetConfig(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
	api := NewAPI(transport)

	config, err := api.GetConfig(context.Background())
	require.NoError(t, err)
	require.Equal(t, "SCR", config.ScorumAddressPrefix)
}

func TestLookupAccounts(t *testing.T) {
	transport := rpc.NewHTTPTransport(nodeHTTPS)
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
