package digest

import (
	"context"
	isaacnetwork "github.com/ProtoconNet/mitum2/isaac/network"
	"github.com/ProtoconNet/mitum2/network/quicmemberlist"
	"github.com/ProtoconNet/mitum2/network/quicstream"
	"net/http"
	"time"

	crcydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/network"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/singleflight"
)

var (
	HandlerPathSTOService                  = `/sto/{contract:\w+}`
	HandlerPathSTOHolderPartitions         = `/sto/{contract:\w+}/holder/{address:(?i)` + base.REStringAddressString + `}/partitions`
	HandlerPathSTOHolderPartitionBalance   = `/sto/{contract:\w+}/holder/{address:(?i)` + base.REStringAddressString + `}/partition/{partition:\w+}/balance`
	HandlerPathSTOHolderPartitionOperators = `/sto/{contract:\w+}/holder/{address:(?i)` + base.REStringAddressString + `}/partition/{partition:\w+}/operators`
	HandlerPathSTOPartitionBalance         = `/sto/{contract:\w+}/partition/{partition:\w+}/balance`
	//HandlerPathSTOPartitionControllers     = `/sto/{contract:\w+}/partition/{partition:\w+}/controllers`
	HandlerPathSTOOperatorHolders = `/sto/{contract:\w+}/operator/{address:(?i)` + base.REStringAddressString + `}/holders`
)

func init() {
	if b, err := crcydigest.JSON.Marshal(crcydigest.UnknownProblem); err != nil {
		panic(err)
	} else {
		crcydigest.UnknownProblemJSON = b
	}
}

type Handlers struct {
	*zerolog.Logger
	networkID       base.NetworkID
	encoders        *encoder.Encoders
	encoder         encoder.Encoder
	database        *crcydigest.Database
	cache           crcydigest.Cache
	nodeInfoHandler crcydigest.NodeInfoHandler
	send            func(interface{}) (base.Operation, error)
	client          func() (*isaacnetwork.BaseClient, *quicmemberlist.Memberlist, []quicstream.ConnInfo, error)
	router          *mux.Router
	routes          map[ /* path */ string]*mux.Route
	itemsLimiter    func(string /* request type */) int64
	rg              *singleflight.Group
	expireNotFilled time.Duration
}

func NewHandlers(
	ctx context.Context,
	networkID base.NetworkID,
	encs *encoder.Encoders,
	enc encoder.Encoder,
	st *crcydigest.Database,
	cache crcydigest.Cache,
	router *mux.Router,
	routes map[string]*mux.Route,
) *Handlers {
	var log *logging.Logging
	if err := util.LoadFromContextOK(ctx, launch.LoggingContextKey, &log); err != nil {
		return nil
	}

	return &Handlers{
		Logger:          log.Log(),
		networkID:       networkID,
		encoders:        encs,
		encoder:         enc,
		database:        st,
		cache:           cache,
		router:          router,
		routes:          routes,
		itemsLimiter:    crcydigest.DefaultItemsLimiter,
		rg:              &singleflight.Group{},
		expireNotFilled: time.Second * 3,
	}
}

func (hd *Handlers) Initialize() error {
	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	hd.router.Use(cors)

	hd.setHandlers()

	return nil
}

func (hd *Handlers) SetLimiter(f func(string) int64) *Handlers {
	hd.itemsLimiter = f

	return hd
}

func (hd *Handlers) Cache() crcydigest.Cache {
	return hd.cache
}

func (hd *Handlers) Router() *mux.Router {
	return hd.router
}

func (hd *Handlers) Handler() http.Handler {
	return network.HTTPLogHandler(hd.router, hd.Logger)
}

func (hd *Handlers) setHandlers() {
	_ = hd.setHandler(HandlerPathSTOService, hd.handleSTOService, true).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathSTOHolderPartitions, hd.handleSTOHolderPartitions, true).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathSTOHolderPartitionBalance, hd.handleSTOHolderPartitionBalance, true).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathSTOHolderPartitionOperators, hd.handleSTOHolderPartitionOperators, true).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathSTOPartitionBalance, hd.handleSTOPartitionBalance, true).
		Methods(http.MethodOptions, "GET")
	//_ = hd.setHandler(HandlerPathSTOPartitionControllers, hd.handleSTOPartitionControllers, true).
	//	Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathSTOOperatorHolders, hd.handleSTOOperatorHolders, true).
		Methods(http.MethodOptions, "GET")
}

func (hd *Handlers) setHandler(prefix string, h network.HTTPHandlerFunc, useCache bool) *mux.Route {
	var handler http.Handler
	if !useCache {
		handler = http.HandlerFunc(h)
	} else {
		ch := crcydigest.NewCachedHTTPHandler(hd.cache, h)

		handler = ch
	}

	var name string
	if prefix == "" || prefix == "/" {
		name = "root"
	} else {
		name = prefix
	}

	var route *mux.Route
	if r := hd.router.Get(name); r != nil {
		route = r
	} else {
		route = hd.router.Name(name)
	}

	// if rules, found := hd.rateLimit[prefix]; found {
	// 	handler = process.NewRateLimitMiddleware(
	// 		process.NewRateLimit(rules, limiter.Rate{Limit: -1}), // NOTE by default, unlimited
	// 		hd.rateLimitStore,
	// 	).Middleware(handler)

	// 	hd.Log().Debug().Str("prefix", prefix).Msg("ratelimit middleware attached")
	// }

	route = route.
		Path(prefix).
		Handler(handler)

	hd.routes[prefix] = route

	return route
}

func (hd *Handlers) combineURL(path string, pairs ...string) (string, error) {
	e := util.StringError("failed to combine url")

	if n := len(pairs); n%2 != 0 {
		return "", e.Wrap(errors.Errorf("uneven pairs to combine url"))
	} else if n < 1 {
		u, err := hd.routes[path].URL()
		if err != nil {
			return "", e.Wrap(err)
		}
		return u.String(), nil
	}

	u, err := hd.routes[path].URLPath(pairs...)
	if err != nil {
		return "", e.Wrap(err)
	}
	return u.String(), nil
}
