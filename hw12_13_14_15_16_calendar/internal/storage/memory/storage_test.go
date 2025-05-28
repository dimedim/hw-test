package memorystorage

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Skip()
}

var (
	UserID  = "user_id1"
	TimeNow = time.Date(2025, 5, 28, 16, 12, 23, 0, time.UTC)
)

var events = []*models.Event{
	{
		ID:       "0",
		Title:    "dawdwad",
		UserID:   UserID,
		StartsAt: TimeNow,
	},
	{
		ID:       "1",
		Title:    "title",
		UserID:   UserID,
		StartsAt: TimeNow,
	},
	{
		ID:       "2",
		Title:    "week",
		UserID:   UserID,
		StartsAt: TimeNow.AddDate(0, 0, 2),
	},
	{
		ID:       "3",
		Title:    "month",
		UserID:   UserID,
		StartsAt: TimeNow.AddDate(0, 0, 10),
	},
	{
		ID:       "4",
		Title:    "title",
		UserID:   UserID,
		StartsAt: TimeNow,
	},
}

func TestBasic(t *testing.T) {
	storage := New()
	ctx := context.Background()

	err := storage.CreateEvent(ctx, events[0])
	require.NoError(t, err)

	res, err := storage.ListEventsByDay(ctx, UserID, TimeNow.UTC())
	require.NoError(t, err)
	require.Equal(t, events[0], res[0])

	res, err = storage.ListEventsByMonth(ctx, UserID, TimeNow.UTC())
	require.NoError(t, err)
	require.Equal(t, events[0], res[0])

	res, err = storage.ListEventsByWeek(ctx, UserID, TimeNow.UTC())
	require.NoError(t, err)
	require.Equal(t, events[0], res[0])

	err = storage.DeleteEvent(ctx, events[0].ID)
	require.NoError(t, err)

	for _, v := range events {
		err = storage.CreateEvent(ctx, v)
		require.NoError(t, err)
	}

	res, err = storage.ListEventsByWeek(ctx, UserID, events[2].StartsAt)
	require.NoError(t, err)
	require.Equal(t, events[2], res[0])

	res, err = storage.ListEventsByMonth(ctx, UserID, events[3].StartsAt)
	require.NoError(t, err)
	require.Equal(t, events[3], res[0])

	closeErr := storage.Close()
	require.NoError(t, closeErr)
}

func TestErrorCases(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stor := New()

	err := stor.CreateEvent(ctx, events[0])
	require.Error(t, err)

	err = stor.UpdateEvent(ctx, "1", events[0])
	require.Error(t, err)

	err = stor.DeleteEvent(ctx, "1")
	require.Error(t, err)

	e, err := stor.ListEventsByDay(ctx, UserID, TimeNow)
	require.Error(t, err)
	require.Nil(t, e)

	e, err = stor.ListEventsByWeek(ctx, UserID, TimeNow)
	require.Error(t, err)
	require.Nil(t, e)

	e, err = stor.ListEventsByMonth(ctx, UserID, TimeNow)
	require.Error(t, err)
	require.Nil(t, e)
}

func TestGetEventListByAnyDate_Empty(t *testing.T) {
	storage := New()
	res, err := storage.GetEventListByAnyDate("user",
		time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, res, 0)
}

func TestGetEventListByAnyDate_FilterAndSort(t *testing.T) {
	start := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	e1 := &models.Event{ID: "e1", UserID: "u1", StartsAt: start.Add(-time.Hour)}
	e2 := &models.Event{ID: "e2", UserID: "u1", StartsAt: start}
	e3 := &models.Event{ID: "e3", UserID: "u2", StartsAt: start.Add(time.Hour)}
	e4 := &models.Event{ID: "e4", UserID: "u1", StartsAt: end.Add(-time.Hour)}
	e5 := &models.Event{ID: "e5", UserID: "u1", StartsAt: end}

	storage := New()
	ctx := context.Background()
	for _, ev := range []*models.Event{e1, e2, e3, e4, e5} {
		require.NoError(t, storage.CreateEvent(ctx, ev))
	}

	res, err := storage.GetEventListByAnyDate("u1", start, end)
	require.NoError(t, err)

	require.Len(t, res, 2)
	require.Equal(t, "e2", res[0].ID)
	require.Equal(t, "e4", res[1].ID)
}

func TestUpdateEvent_NotExists(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*4)
	defer canc()
	storage := New()
	err := storage.UpdateEvent(ctx, "invalid", &models.Event{})
	require.Equal(t, models.ErrEventNotExists, err)
}

func TestUpdateEvent_CancelledContext(t *testing.T) {
	storage := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := storage.UpdateEvent(ctx, "any", &models.Event{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "update event:")
}

func TestUpdateEvent_Success(t *testing.T) {
	storage := New()
	ctx := context.Background()

	e := &models.Event{ID: "e1", UserID: "u1", StartsAt: time.Now().UTC()}
	err := storage.CreateEvent(ctx, e)
	require.NoError(t, err)

	updated := &models.Event{ID: "e1", UserID: "u1", Title: "updated title", StartsAt: e.StartsAt}
	err = storage.UpdateEvent(ctx, "e1", updated)
	require.NoError(t, err)

	got := storage.DB["e1"]
	require.Equal(t, "updated title", got.Title)
	require.False(t, got.UpdatedAt.IsZero())
	require.WithinDuration(t, time.Now().UTC(), got.UpdatedAt, time.Second)
}

func TestConcurrency1(t *testing.T) {
	stor := New()

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	const n = 1000
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := stor.CreateEvent(ctx, &models.Event{ID: strconv.Itoa(i), UserID: UserID, StartsAt: TimeNow})
			require.NoError(t, err)
		}()
	}

	wg.Wait()
	res, err := stor.ListEventsByDay(ctx, UserID, TimeNow)
	require.NoError(t, err)
	require.Equal(t, n, len(res))
}

func TestConcurrency2(t *testing.T) {
	storage := New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	wg := sync.WaitGroup{}
	const n = 1000
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("id_%d", i)
			e := &models.Event{
				ID:       id,
				UserID:   "user_id",
				StartsAt: time.Date(2025, 5, 28, 0, 0, i, 0, time.UTC),
			}
			require.NoError(t, storage.CreateEvent(ctx, e))
			e.Title = "updated"
			require.NoError(t, storage.UpdateEvent(ctx, id, e))
		}(i)
	}
	wg.Wait()

	events, err := storage.ListEventsByMonth(ctx, "user_id", time.Date(2025, 5, 28, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, events, n)
}
