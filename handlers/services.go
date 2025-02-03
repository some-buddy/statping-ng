package handlers

import (
	"github.com/gorilla/mux"
	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/types/errors"
	"github.com/statping-ng/statping-ng/types/failures"
	"github.com/statping-ng/statping-ng/types/hits"
	"github.com/statping-ng/statping-ng/types/services"
	"github.com/statping-ng/statping-ng/utils"
	"net/http"
)

type serviceOrder struct {
	Id    int64 `json:"service"`
	Order int   `json:"order"`
}

func findService(r *http.Request) (*services.Service, error) {
	vars := mux.Vars(r)
	id := utils.ToInt(vars["id"])
	servicer, err := services.Find(id)
	if err != nil {
		return nil, err
	}
	if !servicer.Public.Bool && !IsReadAuthenticated(r) {
		return nil, errors.NotAuthenticated
	}
	return servicer, nil
}

func reorderServiceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var newOrder []*serviceOrder
	if err := DecodeJSON(r, &newOrder); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	for _, s := range newOrder {
		service, err := services.Find(s.Id)
		if err != nil {
			sendErrorJson(err, w, r)
			return
		}
		service.Order = s.Order
		service.Update()
	}
	returnJson(newOrder, w, r)
}

func apiServiceHandler(r *http.Request) interface{} {
	srv, err := findService(r)
	if err != nil {
		return err
	}
	srv = srv.UpdateStats()
	return *srv
}

func apiCreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service *services.Service
	if err := DecodeJSON(r, &service); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	if err := service.Create(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	go services.ServiceCheckQueue(service, true)

	sendJsonAction(service, "create", w, r)
}

type servicePatchReq struct {
	Online  bool   `json:"online"`
	Issue   string `json:"issue,omitempty"`
	Latency int64  `json:"latency,omitempty"`
}

func apiServicePatchHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	var req servicePatchReq
	if err := DecodeJSON(r, &req); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	service.Online = req.Online
	service.Latency = req.Latency

	issueDefault := "Service was triggered to be offline"
	if req.Issue != "" {
		issueDefault = req.Issue
	}

	if !req.Online {
		services.RecordFailure(service, issueDefault, "trigger")
	} else {
		services.RecordSuccess(service)
	}

	if err := service.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	sendJsonAction(service, "update", w, r)
}

func apiServiceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := DecodeJSON(r, &service); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := service.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	go service.CheckService(true)
	sendJsonAction(service, "update", w, r)
}

func apiServiceDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupQuery, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	objs, err := groupQuery.GraphData(database.ByAverage("latency", 1000))
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	returnJson(objs, w, r)
}

func apiServiceFailureDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupQuery, err := database.ParseQueries(r, service.AllFailures())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	objs, err := groupQuery.GraphData(database.ByCount)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	returnJson(objs, w, r)
}

func apiServicePingDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupQuery, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	objs, err := groupQuery.GraphData(database.ByAverage("ping_time", 1000))
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	returnJson(objs, w, r)
}

func apiServiceTimeDataHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupHits, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	groupFailures, err := database.ParseQueries(r, service.AllFailures())
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	var allFailures []*failures.Failure
	var allHits []*hits.Hit

	if err := groupHits.Find(&allHits); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	if err := groupFailures.Find(&allFailures); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	uptimeData, err := service.UptimeData(allHits, allFailures)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	returnJson(uptimeData, w, r)
}

func apiServiceHitsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := service.AllHits().DeleteAll(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(service, "delete", w, r)
}

func apiServiceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	err = service.Delete()
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	sendJsonAction(service, "delete", w, r)
}

func apiAllServicesHandler(r *http.Request) interface{} {
	var srvs []services.Service
	for _, v := range services.AllInOrder() {
		if !v.Public.Bool && !IsUser(r) {
			continue
		}
		srvs = append(srvs, v)
	}
	return srvs
}

func servicesDeleteFailuresHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := service.AllFailures().DeleteAll(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(service, "delete_failures", w, r)
}

func apiServiceFailuresHandler(r *http.Request) interface{} {
	service, err := findService(r)
	if err != nil {
		return err
	}
	var fails []*failures.Failure
	query, err := database.ParseQueries(r, service.AllFailures())
	if err != nil {
		return err
	}
	query.Find(&fails)
	return fails
}

// apiServiceHitsHandler handles the request to retrieve hits for a service.
func apiServiceHitsHandler(r *http.Request) interface{} {
	service, err := findService(r)
	if err != nil {
		return err
	}
	var hts []*hits.Hit
	query, err := database.ParseQueries(r, service.AllHits())
	if err != nil {
		return err
	}
	query.Find(&hts)
	return hts
}

// apiServiceOutageHandler handles the GET request to retrieve the outage status of a service.
func apiServiceOutageHandler(r *http.Request) {
	service, err := findService(r)
	if err != nil {
		return err
	}

	outage := map[string]interface{}{
		"enable_outage": service.EnableIntermediate,
		"minor_outage_name":   service.StatusMinorOutageName,
		"minor_outage_color":  service.StatusMinorOutageColor,
		"major_outage_name":   service.StatusMajorOutageName,
		"major_outage_color":  service.StatusMajorOutageColor,
	}

	returnJson(outage, w, r)
}

// apiServiceUpdateOutageHandler handles the PUT request to update the outage status of a service.
func apiServiceUpdateOutageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := utils.ToInt(vars["id"])

	service, err := services.Find(serviceID)
	if err != nil {
		sendErrorJson(errors.New("service not found"), w, r)
		return
	}

	var outageStatus map[string]interface{}
	if err := DecodeJSON(r, &outageStatus); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	if enableIntermediate, ok := outageStatus["enable_intermediate"].(bool); ok {
		service.EnableIntermediate = enableIntermediate
	}
	if minorOutageName, ok := outageStatus["minor_outage_name"].(string); ok {
		service.StatusMinorOutageName = minorOutageName
	}
	if minorOutageColor, ok := outageStatus["minor_outage_color"].(string); ok {
		service.StatusMinorOutageColor = minorOutageColor
	}
	if majorOutageName, ok := outageStatus["major_outage_name"].(string); ok {
		service.StatusMajorOutageName = majorOutageName
	}
	if majorOutageColor, ok := outageStatus["major_outage_color"].(string); ok {
		service.StatusMajorOutageColor = majorOutageColor
	}

	if err := service.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	sendJsonAction(service, "update", w, r)
}
