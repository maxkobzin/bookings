package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/maxkobzin/bookings/internal/driver"
	"github.com/maxkobzin/bookings/internal/models"
)

//type postData struct {
//	key   string
//	value string
//}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatucCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	//	{"post-search-avail", "/search-availability", "POST", []postData{
	//		{key: "start", value: "2020-01-01"},
	//		{key: "end", value: "2020-01-02"},
	//	}, http.StatusOK},
	//	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
	//		{key: "start", value: "2020-01-01"},
	//		{key: "end", value: "2020-01-02"},
	//	}, http.StatusOK},
	//	{"mr post", "/make-reservation", "POST", []postData{
	//		{key: "first_name", value: "John"},
	//		{key: "last_name", value: "Smith"},
	//		{key: "email", value: "me@here.com"},
	//		{key: "phone", value: "555-555-5555"},
	//	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatucCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatucCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
		StartDate: time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2050, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	postData := url.Values{}
	postData.Add("first_name", "John")
	postData.Add("last_name", "Smith")
	postData.Add("email", "john@smith.com")
	postData.Add("phone", "555-555-5555")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d,  wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid data
	postData = url.Values{}
	postData.Add("first_name", "J")
	postData.Add("last_name", "Smith")
	postData.Add("email", "john@smith.com")
	postData.Add("phone", "555-555-5555")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to insert reservation into database
	postData = url.Values{}
	postData.Add("first_name", "John")
	postData.Add("last_name", "Smith")
	postData.Add("email", "john@smith.com")
	postData.Add("phone", "555-555-5555")
	reservation.RoomID = 2

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure to insert restriction into database
	reservation.RoomID = 1000

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting restriction: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// test rooms are available
	postData := url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostAvailability handler returned wrong response code if rooms are available: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test rooms are not available
	postData = url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "2049-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code if rooms are not available: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid data
	postData = url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "2050-01-01")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code if data invalid: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for missing post body: got %d,  wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date
	postData = url.Values{}
	postData.Add("start", "invalid")
	postData.Add("end", "2050-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for invalid start date: got %d,  wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid end date
	postData = url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "invalid")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned wrong response code for invalid end date: got %d,  wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// first case - rooms are not available
	postData := url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "2050-01-02")
	postData.Add("room_id", "1")

	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	// make handler handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// get responce recorder
	rr := httptest.NewRecorder()

	// make request fo our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}

	if j.OK {
		t.Error("Got not availability when available was expected in AvailabilityJSON")
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
