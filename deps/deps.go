// Only purpose for this package is to provide a way to expose the `Dependencies`
// struct without introducing import cycle problems.

package deps

import (
	"github.com/cactus/go-statsd-client/statsd"
)

type Dependencies struct {
	StatsD statsd.Statter
}
