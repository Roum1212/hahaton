//go:build integration

package integration

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
)

func checkGeoImages(x *fixture) {
	const coordinatorDescription = "Фото профиля"

	c, a := x.coordinator("geo-owner"), x.login("geo-volunteer")
	x.cooperate(c, a)
	missing := x.fail(actor{}, "geo.v1.GeoService/GetCity", map[string]any{"id": "unknown-city"}, "not_found")
	info := resourceInfo(x.t, missing)
	require(x.t, info.GetResourceType() == "geo.v1.City", "wrong missing city resource")
	require(x.t, info.GetResourceName() == "id:unknown-city", "missing city id not included")
	cities := x.call(actor{}, "geo.v1.GeoService/ListCities", map[string]any{"query": "  обн  "})
	require(x.t, len(list(x.t, cities, "cities")) == 1, "city search did not normalize query")
	cities = x.call(actor{}, "geo.v1.GeoService/ListCities", map[string]any{
		"limit":  1000,
		"offset": "1",
	})
	require(x.t, len(list(x.t, cities, "cities")) == 1, "city-specific pagination limit was lost")

	for _, paging := range []map[string]any{
		{"limit": -1},
		{"limit": 1001},
		{"offset": "-1"},
	} {
		x.fail(actor{}, "geo.v1.GeoService/ListCities", paging, "invalid_argument")
	}

	city := object(x.t, x.call(actor{}, "geo.v1.GeoService/GetCity", map[string]any{"id": "city"}), "city")
	require(x.t, city["has_boundary"] == true, "persisted boundary flag missing")
	noBoundary := object(x.t, x.call(actor{}, "geo.v1.GeoService/GetCity", map[string]any{"id": "no-boundary"}), "city")
	require(x.t, noBoundary["center"] == nil && noBoundary["has_boundary"] != true, "missing geometry was fabricated")
	id := x.event(c, "Геопоиск", 1)
	area := map[string]any{"city_id": "city"}
	catalog := x.call(a, "event.v1.EventService/ListEvents", map[string]any{
		"scope":       1,
		"search_area": area,
	})
	require(x.t, len(list(x.t, catalog, "events")) == 1, "city polygon failed")
	x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id": id,
		"point": map[string]any{
			"latitude":  55.85,
			"longitude": 36.85,
		},
	})
	hole := x.call(a, "event.v1.EventService/ListEvents", map[string]any{
		"scope":       1,
		"search_area": area,
	})
	require(x.t, len(list(x.t, hole, "events")) == 0, "polygon hole was not excluded")
	x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id": id,
		"point": map[string]any{
			"latitude":  55,
			"longitude": 36.5,
		},
	})
	boundary := x.call(a, "event.v1.EventService/ListEvents", map[string]any{
		"scope":       1,
		"search_area": area,
	})
	require(x.t, len(list(x.t, boundary, "events")) == 1, "polygon boundary excluded")

	radius := map[string]any{"radius": map[string]any{
		"center": map[string]any{
			"latitude":  55,
			"longitude": 36.5,
		},
		"radius_meters": 100,
	}}
	circle := x.call(a, "event.v1.EventService/ListEvents", map[string]any{
		"scope":       1,
		"search_area": radius,
	})
	require(x.t, len(list(x.t, circle, "events")) == 1, "radius search failed")
	x.fail(a, "event.v1.EventService/ListEvents", map[string]any{
		"scope": 1,
		"search_area": map[string]any{
			"city_id": "no-boundary",
		},
	}, "failed_precondition")

	query := map[string]any{
		"query": "  Обнинск  ",
		"from_bound": map[string]any{
			"value": "city",
		},
		"to_bound": map[string]any{
			"value": "settlement",
		},
	}
	suggestions := x.call(a, "geo.v1.GeoService/SuggestAddresses", query)
	data := object(x.t, dictionary(x.t, list(x.t, suggestions, "suggestions")[0]), "data")
	require(x.t, object(x.t, data, "point")["latitude"] != nil, "dadata coordinates missing")
	x.dadataStatus.Store(http.StatusTooManyRequests)
	x.fail(a, "geo.v1.GeoService/SuggestAddresses", query, "resource_exhausted")
	x.dadataStatus.Store(http.StatusInternalServerError)
	x.fail(a, "geo.v1.GeoService/SuggestAddresses", query, "unavailable")
	x.dadataStatus.Store(http.StatusOK)

	for _, body := range []string{`{`, `{}`, `{"suggestions": {}}`} {
		x.dadataBody = body
		x.fail(a, "geo.v1.GeoService/SuggestAddresses", query, "unavailable")
	}

	x.dadataBody = `{
  "suggestions": [
    {
      "value": "Обнинск",
      "unrestricted_value": "Обнинск",
      "data": {
        "geo_lat": null,
        "geo_lon": null,
        "qc_geo": null
      }
    }
  ]
}`
	suggestions = x.call(a, "geo.v1.GeoService/SuggestAddresses", query)
	data = object(x.t, dictionary(x.t, list(x.t, suggestions, "suggestions")[0]), "data")
	require(x.t, data["point"] == nil && data["qc_geo"] == nil, "missing dadata coordinates replaced by zeros")

	var buffer bytes.Buffer
	must(x.t, png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 2, 2))))
	imageResult := x.call(c, "image.v1.ImageService/UploadImage", map[string]any{
		"data": base64.StdEncoding.EncodeToString(buffer.Bytes()),
	})
	path := text(x.t, imageResult["image_path"])
	photoEvent := x.event(c, "С фотографией", 1, map[string]any{"image_path": path})
	event := object(x.t, x.call(c, "event.v1.EventService/GetEvent", map[string]any{"id": photoEvent}), "event")
	require(x.t, event["image_url"] == imageResult["image_url"], "event create lost typed image path")
	x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{"id": photoEvent})
	event = object(x.t, x.call(c, "event.v1.EventService/GetEvent", map[string]any{"id": photoEvent}), "event")
	require(x.t, event["image_url"] == imageResult["image_url"], "omitted event image was cleared")
	x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":         photoEvent,
		"image_path": "",
	})
	event = object(x.t, x.call(c, "event.v1.EventService/GetEvent", map[string]any{"id": photoEvent}), "event")
	require(x.t, event["image_url"] == nil || event["image_url"] == "", "event image clearing failed")

	ownImage := x.call(a, "image.v1.ImageService/UploadImage", map[string]any{
		"data": base64.StdEncoding.EncodeToString(buffer.Bytes()),
	})
	created := x.request(actor{}, "admin.v1.AdminService/CreateCoordinator", map[string]any{
		"volunteer_id":                         a.id,
		"avatar_path":                          ownImage["image_path"],
		"cooperation_request_cooldown_seconds": "60",
	}, true)
	require(x.t, created.Code == http.StatusOK, created.Body.String())
	ownCoordinator := object(x.t, x.call(a, "coordinator.v1.CoordinatorService/GetCoordinator", nil), "coordinator")
	require(x.t, ownCoordinator["avatar_url"] == ownImage["image_url"], "coordinator create lost typed avatar")

	x.call(c, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"avatar_path": path})
	x.call(c, "coordinator.v1.CoordinatorService/UpdateCoordinator", map[string]any{"avatar_path": path})
	x.call(
		c,
		"coordinator.v1.CoordinatorService/UpdateCoordinator",
		map[string]any{"description": coordinatorDescription},
	)
	x.call(c, "event.v1.EventService/UpdateEvent", map[string]any{
		"id":         id,
		"image_path": path,
	})

	rewardID := text(x.t, x.call(c, "reward.v1.RewardService/CreateReward", map[string]any{
		"title":          "Фото",
		"price_coins":    "1",
		"quantity_limit": 1,
		"image_path":     path,
	})["id"])
	for _, scenario := range []struct {
		method string
		body   map[string]any
	}{
		{
			"volunteer.v1.VolunteerService/UpdateVolunteer",
			map[string]any{"avatar_path": "invalid"},
		},
		{
			"coordinator.v1.CoordinatorService/UpdateCoordinator",
			map[string]any{"avatar_path": "invalid"},
		},
		{
			"event.v1.EventService/UpdateEvent",
			map[string]any{
				"id":         id,
				"image_path": "invalid",
			},
		},
		{
			"reward.v1.RewardService/UpdateReward",
			map[string]any{
				"id":         rewardID,
				"image_path": "invalid",
			},
		},
		{
			"event.v1.EventService/CreateEvent",
			map[string]any{
				"title":      "Фото",
				"image_path": "invalid",
			},
		},
		{
			"reward.v1.RewardService/CreateReward",
			map[string]any{"image_path": "invalid"},
		},
	} {
		x.fail(c, scenario.method, scenario.body, "invalid_argument")
	}

	missingPath := c.id + "/00000000-0000-7000-8000-000000000001.png"
	x.fail(
		c,
		"coordinator.v1.CoordinatorService/UpdateCoordinator",
		map[string]any{"avatar_path": missingPath},
		"not_found",
	)

	x.call(c, "reward.v1.RewardService/UpdateReward", map[string]any{
		"id":         rewardID,
		"image_path": "",
	})
	x.fail(a, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"avatar_path": path}, "permission_denied")
	x.fail(c, "image.v1.ImageService/UploadImage", map[string]any{"data": base64.StdEncoding.EncodeToString([]byte(
		"not an image",
	))}, "invalid_argument")

	file := httptest.NewRecorder()
	x.handler.ServeHTTP(file, httptest.NewRequest(
		http.MethodGet,
		"https://demo.test"+text(x.t, imageResult["image_url"]),
		nil,
	))
	require(x.t, file.Code == http.StatusOK && bytes.Equal(file.Body.Bytes(), buffer.Bytes()), "stored file differs")
	x.call(c, "volunteer.v1.VolunteerService/UpdateVolunteer", map[string]any{"avatar_path": ""})
	profile := object(x.t, x.call(c, "volunteer.v1.VolunteerService/GetVolunteer", nil), "volunteer")
	coordinator := object(x.t, x.call(c, "coordinator.v1.CoordinatorService/GetCoordinator", nil), "coordinator")
	require(x.t, profile["avatar_url"] == nil || profile["avatar_url"] == "", "avatar clearing failed")
	require(x.t, coordinator["avatar_url"] == imageResult["image_url"], "coordinator avatar not independent")
	x.call(c, "coordinator.v1.CoordinatorService/UpdateCoordinator", map[string]any{"avatar_path": ""})
	coordinator = object(x.t, x.call(c, "coordinator.v1.CoordinatorService/GetCoordinator", nil), "coordinator")
	require(x.t, coordinator["description"] == coordinatorDescription, "omitted description was changed")
	require(x.t, coordinator["avatar_url"] == nil || coordinator["avatar_url"] == "", "coordinator avatar not cleared")
}
