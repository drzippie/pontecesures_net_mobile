package saved

import (
	"testing"
	"time"

	"github.com/drzippie/pontecesures_net_mobile/internal/wp"
)

func post(id int, title string) wp.Post {
	return wp.Post{PostSummary: wp.PostSummary{ID: id, Title: title}, Content: "<p>" + title + "</p>"}
}

func TestPutListRemove(t *testing.T) {
	s := NewStore(t.TempDir())
	if len(s.List()) != 0 {
		t.Fatal("new store not empty")
	}
	_ = s.Put(post(1, "primeira"))
	time.Sleep(5 * time.Millisecond)
	_ = s.Put(post(2, "segunda"))

	list := s.List()
	if len(list) != 2 || list[0].ID != 2 || list[1].ID != 1 {
		t.Fatalf("order = %+v", list)
	}
	// Re-saving keeps the original position (save time).
	time.Sleep(5 * time.Millisecond)
	_ = s.Put(post(1, "primeira actualizada"))
	list = s.List()
	if list[0].ID != 2 || list[1].Title != "primeira actualizada" {
		t.Fatalf("after refresh = %+v", list)
	}
	if p, ok := s.Get(1); !ok || p.Content != "<p>primeira actualizada</p>" {
		t.Fatalf("Get = %+v %v", p, ok)
	}
	if err := s.Remove(1); err != nil {
		t.Fatal(err)
	}
	if err := s.Remove(1); err != nil {
		t.Fatal("second remove should be a no-op")
	}
	if _, ok := s.Get(1); ok || len(s.IDs()) != 1 {
		t.Fatal("not removed")
	}
}
