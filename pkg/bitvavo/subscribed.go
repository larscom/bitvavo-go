package bitvavo

import (
	"fmt"

	"github.com/goccy/go-json"
)

var ErrUnexpectedType = func(v any) error { return fmt.Errorf("unexpected type '%s'", v) }

type Unsubscribed = Subscribed

type Subscribed struct {
	// Currently active subscriptions that the broker knows of.
	Subscriptions map[Channel][]string

	// Currently active subscriptions with an interval that the broker knows of.
	SubscriptionsInterval map[Channel]map[Interval][]string
}

func (s *Subscribed) UnmarshalJSON(bytes []byte) error {
	var j map[string]any

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	var (
		all                   = j["subscriptions"].(map[string]any)
		subscriptions         = make(map[Channel][]string)
		subscriptionsInterval = make(map[Channel]map[Interval][]string)
	)

	for key, value := range all {
		channel := *channels.Parse(key)

		switch v := value.(type) {
		// without interval
		case []any:
			subscriptions[channel] = make([]string, len(v))
			for index, market := range v {
				subscriptions[channel][index] = market.(string)
			}
		// with interval
		case map[string]any:
			subscriptionsInterval[channel] = make(map[Interval][]string)
			for i, m := range v {
				interval := *intervals.Parse(i)
				markets := m.([]any)
				subscriptionsInterval[channel][interval] = make([]string, len(markets))
				for index, market := range markets {
					subscriptionsInterval[channel][interval][index] = market.(string)
				}

			}
		default:
			return ErrUnexpectedType(v)
		}
	}

	s.Subscriptions = subscriptions
	s.SubscriptionsInterval = subscriptionsInterval

	return nil
}
