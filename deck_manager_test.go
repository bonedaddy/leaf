package leaf

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeckManager(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "leaf.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	db, err := OpenBoltStore(tmpfile.Name())
	require.NoError(t, err)

	dm, err := NewDeckManager(".", db)
	require.NoError(t, err)

	t.Run("ReviewDecks", func(t *testing.T) {
		decks, err := dm.ReviewDecks(100)
		require.NoError(t, err)
		require.Len(t, decks, 1)

		deck := decks[0]
		assert.Equal(t, "Hiragana", deck.Name)
		assert.Equal(t, 46, deck.CardsReady)
		assert.InDelta(t, time.Since(deck.NextReviewAt), 0, float64(time.Minute))
	})

	t.Run("ReviewSession", func(t *testing.T) {
		session, err := dm.ReviewSession("Hiragana", 20)
		require.NoError(t, err)
		assert.Equal(t, 20, session.Total())

		question := session.Next()
		session.Answer("foo")
		err = db.RangeStats("Hiragana", func(card string, s *Stats) bool {
			if card != question {
				return true
			}

			sm := s.Supermemo.(*Supermemo2Plus)
			assert.InDelta(t, 0.45, sm.Difficulty, 0.01)
			assert.InDelta(t, 0.2, sm.Interval, 0.01)
			return false
		})

		require.NoError(t, err)
	})

	t.Run("DeckStats", func(t *testing.T) {
		stats, err := dm.DeckStats("Hiragana")
		require.NoError(t, err)
		assert.Len(t, stats, 46)

		s := stats[0]
		assert.NotEmpty(t, s.Question)
		sm := s.Supermemo.(*Supermemo2PlusCustom)
		assert.InDelta(t, 0.3, sm.Difficulty, 0.01)
	})
}