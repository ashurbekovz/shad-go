package structtags

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	expectedUser1 = User{
		ID:              1,
		Name:            "John",
		Surname:         "Doe",
		Phone:           "88005551234",
		HasSubscription: true,
	}
	userURL1 = fmt.Sprintf(
		"localhost/user?id=%d&name=%s&surname=%s&phone=%s&has_subscription=%t",
		expectedUser1.ID,
		expectedUser1.Name,
		expectedUser1.Surname,
		expectedUser1.Phone,
		expectedUser1.HasSubscription,
	)

	expectedUser2 = User{
		ID:              1,
		Name:            "Mike",
		Surname:         "Crag",
		Phone:           "88005341234",
		HasSubscription: false,
	}
	userURL2 = fmt.Sprintf(
		"localhost/user?id=%d&name=%s&surname=%s&phone=%s&has_subscription=%t",
		expectedUser2.ID,
		expectedUser2.Name,
		expectedUser2.Surname,
		expectedUser2.Phone,
		expectedUser2.HasSubscription,
	)

	expectedGood = Good{
		ID:   45,
		Name: "pizza",
	}
	goodURL = fmt.Sprintf(
		"localhost/good?id=%d&name=%s",
		expectedGood.ID,
		expectedGood.Name,
	)

	expectedOrder = Order{
		ID:      37,
		UserID:  73,
		GoodIds: []int{1, 2, 3},
		Date:    "01.01.1970",
	}
	orderURL = fmt.Sprintf(
		"localhost/order?id=%d&user_id=%d&good_ids=%d&good_ids=%d&good_ids=%d&date=%s",
		expectedOrder.ID,
		expectedOrder.UserID,
		expectedOrder.GoodIds[0],
		expectedOrder.GoodIds[1],
		expectedOrder.GoodIds[2],
		expectedOrder.Date,
	)
)

type User struct {
	ID              int
	Name            string
	Surname         string
	Phone           string
	HasSubscription bool `http:"has_subscription"`
}

type Good struct {
	ID   int
	Name string
}

type Order struct {
	ID      int
	UserID  int   `http:"user_id"`
	GoodIds []int `http:"good_ids"`
	Date    string
}

func TestUnpack_User(t *testing.T) {
	r1, _ := http.NewRequest("GET", userURL1, nil)
	user1 := &User{}
	err := Unpack(r1, user1)
	require.NoError(t, err)
	require.Equal(t, expectedUser1.ID, user1.ID)
	require.Equal(t, expectedUser1.Name, user1.Name)
	require.Equal(t, expectedUser1.Surname, user1.Surname)
	require.Equal(t, expectedUser1.Phone, user1.Phone)
	require.Equal(t, expectedUser1.HasSubscription, user1.HasSubscription)

	r2, _ := http.NewRequest("GET", userURL2, nil)
	user2 := &User{}
	err = Unpack(r2, user2)
	require.NoError(t, err)
	require.Equal(t, expectedUser2.ID, user2.ID)
	require.Equal(t, expectedUser2.Name, user2.Name)
	require.Equal(t, expectedUser2.Surname, user2.Surname)
	require.Equal(t, expectedUser2.Phone, user2.Phone)
	require.Equal(t, expectedUser2.HasSubscription, user2.HasSubscription)
}

func TestUnpack_Good(t *testing.T) {
	r, _ := http.NewRequest("GET", goodURL, nil)
	good := &Good{}
	err := Unpack(r, good)
	require.NoError(t, err)
	require.Equal(t, expectedGood.ID, good.ID)
	require.Equal(t, expectedGood.Name, good.Name)
}

func TestUnpack_Order(t *testing.T) {
	r, _ := http.NewRequest("GET", orderURL, nil)
	order := &Order{}
	err := Unpack(r, order)
	require.NoError(t, err)
	require.Equal(t, expectedOrder.ID, order.ID)
	require.Equal(t, expectedOrder.UserID, order.UserID)
	require.Equal(t, expectedOrder.GoodIds, order.GoodIds)
	require.Equal(t, expectedOrder.Date, order.Date)
}

func TestUnpack_ParseFormError(t *testing.T) {
	r, _ := http.NewRequest("POST", "localhost", nil)
	user := &User{}
	err := Unpack(r, user)
	require.Error(t, err)
}

func TestUnpack_IncorrectBoolData(t *testing.T) {
	url := "localhost/user?id=1&has_subscription=7"
	r, _ := http.NewRequest("GET", url, nil)
	user := &User{}
	err := Unpack(r, user)
	require.Error(t, err)
}

func TestUnpack_IncorrectIntData(t *testing.T) {
	url := "localhost/user?id=abc"
	r, _ := http.NewRequest("GET", url, nil)
	user := &User{}
	err := Unpack(r, user)
	require.Error(t, err)
}

func BenchmarkUnpacker(b *testing.B) {
	userRequest, _ := http.NewRequest("GET", userURL1, nil)
	user := &User{}

	goodRequest, _ := http.NewRequest("GET", goodURL, nil)
	good := &Good{}

	orderRequest, _ := http.NewRequest("GET", orderURL, nil)
	order := &Order{}

	b.Run("user", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			_ = Unpack(userRequest, user)
		}
	})

	b.Run("good", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			_ = Unpack(goodRequest, good)
			_ = Unpack(orderRequest, order)
		}
	})

	b.Run("order", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			_ = Unpack(orderRequest, order)
		}
	})
}
