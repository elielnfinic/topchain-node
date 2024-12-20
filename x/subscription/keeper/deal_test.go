package keeper_test

import (
	"testing"

	"topchain/x/subscription/types"

	"github.com/stretchr/testify/require"
)

func TestGetSetDeal(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// Create a deal
	deal := types.Deal{Id: "12345", CroId: "alicecro", Requester: Alice, Status: types.Deal_SCHEDULED, AvailableAmount: 1000, TotalAmount: 100, StartBlock: 10, EndBlock: 20}
	k.SetDeal(ctx, deal)

	dealResponse, found := k.GetDeal(ctx, deal.Id)
	require.True(t, found)
	require.EqualValues(t, deal, dealResponse)
}

func TestDealActive(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// create a subscription
	sub := types.Subscription{Id: "123", DealId: "12345", Provider: "provider1", StartBlock: 10, EndBlock: 15}
	k.SetSubscription(ctx, sub)
	// Create a deal
	deal := types.Deal{Id: "12345", CroId: "alicecro", Requester: Alice, Status: types.Deal_UNDEFINED, AvailableAmount: 1000, TotalAmount: 100, StartBlock: 10, EndBlock: 20, SubscriptionIds: []string{"123"}}
	k.SetDeal(ctx, deal)

	// The deal must be inactive at block number 0
	isActive := k.IsDealActive(ctx, deal)
	require.False(t, isActive)

	// Jump to block 12
	ctx = MockBlockHeight(ctx, am, 12)

	// The deal must be active at block number 12
	isActive = k.IsDealActive(ctx, deal)
	require.True(t, isActive)

	// Jump to block 18
	ctx = MockBlockHeight(ctx, am, 18)
	// The deal must be inactive at block number 18
	isActive = k.IsDealActive(ctx, deal)
	require.False(t, isActive)

}

func TestGetAllActiveProviders(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// create a subscription
	sub := types.Subscription{Id: "123", DealId: "12345", Provider: "provider1", StartBlock: 10, EndBlock: 15}
	k.SetSubscription(ctx, sub)
	// Create a deal
	deal := types.Deal{Id: "12345", CroId: "alicecro", Requester: Alice, Status: types.Deal_UNDEFINED, AvailableAmount: 1000, TotalAmount: 100, StartBlock: 10, EndBlock: 20, SubscriptionIds: []string{"123"}}
	k.SetDeal(ctx, deal)

	activeSubs := k.GetAllActiveSubscriptions(ctx, deal)
	// there shouldn't be any active subs at block 0
	require.True(t, len(activeSubs) == 0)

	// Jump to block 12
	ctx = MockBlockHeight(ctx, am, 12)
	activeSubs = k.GetAllActiveSubscriptions(ctx, deal)
	// there should be an active subs at block 12
	require.True(t, len(activeSubs) == 1)
	_, ok := activeSubs[sub.Id]
	require.True(t, ok)

	// Jump to block 18
	ctx = MockBlockHeight(ctx, am, 18)
	activeSubs = k.GetAllActiveSubscriptions(ctx, deal)
	// there shouldn't be an active subs at block 18
	require.True(t, len(activeSubs) == 0)

}

func TestIsDealUnavailable(t *testing.T) {
	k, ms, ctx, _ := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	require.False(t, k.IsDealUnavailable(types.Deal_ACTIVE))
	require.False(t, k.IsDealUnavailable(types.Deal_SCHEDULED))
	require.False(t, k.IsDealUnavailable(types.Deal_INITIALIZED))
	require.True(t, k.IsDealUnavailable(types.Deal_CANCELLED))
	require.True(t, k.IsDealUnavailable(types.Deal_EXPIRED))

}

func TestDealHasProvider(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)

	// create a subscription
	sub := types.Subscription{Id: "123", DealId: "12345", Provider: "provider1", StartBlock: 10, EndBlock: 15}
	k.SetSubscription(ctx, sub)
	// Create a deal
	deal := types.Deal{Id: "12345", CroId: "alicecro", Requester: Alice, Status: types.Deal_UNDEFINED, AvailableAmount: 1000, TotalAmount: 100, StartBlock: 10, EndBlock: 20, SubscriptionIds: []string{"123"}}
	k.SetDeal(ctx, deal)

	hasProvider := k.DealHasProvider(ctx, deal, "provider1")
	require.True(t, hasProvider)

	hasProvider = k.DealHasProvider(ctx, deal, "provider2")
	require.False(t, hasProvider)

	// Jump to block 18
	ctx = MockBlockHeight(ctx, am, 18)
	hasProvider = k.DealHasProvider(ctx, deal, "provider1")
	require.False(t, hasProvider)
}
