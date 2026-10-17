//go:build integration

package integration

func checkRewards(x *fixture) {
	c, d := x.coordinator("rewards-owner"), x.coordinator("rewards-partner")
	a, b := x.login("buyer"), x.login("outsider")
	x.cooperate(c, a)
	x.cooperate(d, a)
	eventID := x.event(c, "Начисление", 2)
	participationID := x.participate(c, a, eventID)
	x.finishTime(eventID)
	x.call(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         eventID,
		"request_id": "reward-funds",
		"results": []any{
			map[string]any{
				"participation_id":  participationID,
				"confirmed_minutes": "180",
			},
		},
	})
	response := x.call(d, "reward.v1.RewardService/CreateReward", map[string]any{
		"title":          "Сувенир",
		"price_coins":    "20",
		"quantity_limit": 2,
	})
	id := text(x.t, response["id"])
	x.fail(b, "reward.v1.RewardService/GetReward", map[string]any{"id": id}, "permission_denied")
	x.fail(c, "reward.v1.RewardService/UpdateReward", map[string]any{
		"id":    id,
		"title": "Чужое",
	}, "permission_denied")
	x.call(a, "reward.v1.RewardService/GetReward", map[string]any{"id": id})
	x.fail(a, "reward.v1.RewardService/ListRewards", map[string]any{"scope": 2}, "permission_denied")
	catalog := x.call(a, "reward.v1.RewardService/ListRewards", map[string]any{
		"scope":          1,
		"coordinator_id": d.id,
	})
	require(x.t, len(list(x.t, catalog, "rewards")) == 1, "partner reward catalog failed")
	x.fail(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            id,
		"expected_price_coins": "19",
		"request_id":           "buy",
	}, "aborted")
	bought := object(x.t, x.call(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            id,
		"expected_price_coins": "20",
		"request_id":           "buy",
	}), "reward_order")
	orderID := text(x.t, bought["id"])
	balance := object(x.t, x.call(a, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(x.t, integer(x.t, balance, "coin_balance") == 10, "cross-coordinator MVP balance is not 10")
	x.call(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            id,
		"expected_price_coins": "20",
		"request_id":           "buy",
	})
	x.fail(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            id,
		"expected_price_coins": "21",
		"request_id":           "buy",
	}, "already_exists")
	x.fail(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            id,
		"expected_price_coins": "20",
		"request_id":           "overdraft",
	}, "failed_precondition")
	x.fail(b, "reward_order.v1.RewardOrderService/GetRewardOrder", map[string]any{"id": orderID}, "permission_denied")
	x.call(d, "reward.v1.RewardService/UpdateReward", map[string]any{
		"id":          id,
		"price_coins": "25",
		"title":       "Новое название",
	})
	original := object(x.t, x.call(a, "reward_order.v1.RewardOrderService/GetRewardOrder", map[string]any{
		"id": orderID,
	}), "reward_order")
	require(
		x.t,
		original["reward_title"] == "Сувенир" && integer(x.t, original, "price_coins") == 20,
		"order snapshots changed",
	)
	x.fail(d, "reward.v1.RewardService/UpdateReward", map[string]any{
		"id":             id,
		"quantity_limit": 0,
	}, "failed_precondition")
	x.call(a, "reward_order.v1.RewardOrderService/CancelRewardOrder", map[string]any{"id": orderID})
	x.call(a, "reward_order.v1.RewardOrderService/CancelRewardOrder", map[string]any{"id": orderID})
	balance = object(x.t, x.call(a, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(x.t, integer(x.t, balance, "coin_balance") == 30, "refund used current price or duplicated credit")
	issued := object(x.t, x.call(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            id,
		"expected_price_coins": "25",
		"request_id":           "issued",
	}), "reward_order")
	issuedID := text(x.t, issued["id"])
	x.call(d, "reward_order.v1.RewardOrderService/IssueRewardOrder", map[string]any{"id": issuedID})
	x.call(d, "reward_order.v1.RewardOrderService/IssueRewardOrder", map[string]any{"id": issuedID})
	x.fail(a, "reward_order.v1.RewardOrderService/CancelRewardOrder", map[string]any{
		"id": issuedID,
	}, "failed_precondition")
	inexpensive := text(x.t, x.call(d, "reward.v1.RewardService/CreateReward", map[string]any{
		"title":          "Значок",
		"price_coins":    "5",
		"quantity_limit": 1,
	})["id"])
	pending := object(x.t, x.call(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            inexpensive,
		"expected_price_coins": "5",
		"request_id":           "pending",
	}), "reward_order")
	pendingID := text(x.t, pending["id"])

	x.call(d, "reward.v1.RewardService/DeleteReward", map[string]any{"id": inexpensive})
	x.call(d, "reward.v1.RewardService/DeleteReward", map[string]any{"id": inexpensive})
	x.fail(a, "reward.v1.RewardService/GetReward", map[string]any{"id": inexpensive}, "not_found")
	canceled := object(x.t, x.call(a, "reward_order.v1.RewardOrderService/GetRewardOrder", map[string]any{
		"id": pendingID,
	}), "reward_order")
	require(
		x.t,
		canceled["cancellation_reason"] == "REWARD_ORDER_CANCELLATION_REASON_REWARD_DELETED",
		"delete did not cancel pending order",
	)
	x.call(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            inexpensive,
		"expected_price_coins": "5",
		"request_id":           "pending",
	})
	x.call(d, "reward.v1.RewardService/DeleteReward", map[string]any{"id": id})
	x.fail(d, "reward.v1.RewardService/UpdateReward", map[string]any{"id": id}, "failed_precondition")

	terminationReward := text(x.t, x.call(d, "reward.v1.RewardService/CreateReward", map[string]any{
		"title":          "Заказ перед выходом",
		"price_coins":    "1",
		"quantity_limit": 1,
	})["id"])
	terminationOrder := object(x.t, x.call(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            terminationReward,
		"expected_price_coins": "1",
		"request_id":           "termination",
	}), "reward_order")
	x.call(a, "cooperation.v1.CooperationService/TerminateCooperation", map[string]any{"coordinator_id": d.id})

	terminatedOrder := object(
		x.t,
		x.call(a, "reward_order.v1.RewardOrderService/GetRewardOrder", map[string]any{"id": terminationOrder["id"]}),
		"reward_order",
	)
	require(
		x.t,
		terminatedOrder["cancellation_reason"] == "REWARD_ORDER_CANCELLATION_REASON_COOPERATION_TERMINATED",
		"cooperation termination did not cancel pending order",
	)
	x.call(a, "reward_order.v1.RewardOrderService/CreateRewardOrder", map[string]any{
		"reward_id":            id,
		"expected_price_coins": "25",
		"request_id":           "issued",
	})
	history := x.call(a, "reward_order.v1.RewardOrderService/ListRewardOrders", map[string]any{
		"scope":  1,
		"status": 2,
	})
	require(x.t, len(list(x.t, history, "reward_orders")) == 1, "issued order disappeared")
	incoming := x.call(d, "reward_order.v1.RewardOrderService/ListRewardOrders", map[string]any{
		"scope":        2,
		"volunteer_id": a.id,
	})
	require(x.t, len(list(x.t, incoming, "reward_orders")) == 4, "incoming scope or historical orders lost")
	operations := x.call(a, "volunteer.v1.VolunteerService/ListVolunteerOperations", nil)
	require(x.t, len(list(x.t, operations, "operations")) == 8, "wallet ledger lost operations")
	refunds := x.call(a, "volunteer.v1.VolunteerService/ListVolunteerOperations", map[string]any{"type": 3})
	require(x.t, len(list(x.t, refunds, "operations")) == 3, "wallet operation filter failed")
	balance = object(x.t, x.call(a, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	require(
		x.t,
		integer(x.t, balance, "coin_balance") == 5 && integer(x.t, balance, "confirmed_minutes") == 180,
		"reward flow changed confirmed time or balance",
	)
}
