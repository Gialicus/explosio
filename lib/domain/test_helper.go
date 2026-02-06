package domain

// BuildActivityForTest crea un'attività minimale per test (senza DSL).
// Usato da engine e altri package nei test.
func BuildActivityForTest(id, name string, duration int, subs []*Activity) *Activity {
	a := &Activity{
		ID:            id,
		Name:          name,
		Description:   name,
		Duration:      duration,
		MinDuration:   duration,
		SubActivities: subs,
	}
	for _, s := range subs {
		s.Next = append(s.Next, a.ID)
	}
	return a
}
