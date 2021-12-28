package http

import (
	"net/http"
	"strings"

	"github.com/ah-its-andy/ecosystem/config"
	"github.com/ah-its-andy/ecosystem/i18n"
	"github.com/ah-its-andy/ecosystem/logging"
	"github.com/gin-gonic/gin"
)

type Application interface {
	UseRouteMap(f func(*RouteMap))
	UseInterceptors(interceptors ...Interceptor)
	Run()
}

var _ Application = (*webApplication)(nil)

type webApplication struct {
	engine       *gin.Engine
	routes       *RouteMap
	interceptors []Interceptor
	logger       logging.Logger
	appConfig    config.ConfigSection

	finalize Interceptor
}

func NewWebApplication(logger logging.Logger, appConfig config.ConfigSection) Application {
	app := &webApplication{}
	app.engine = gin.New()
	app.routes = &RouteMap{
		areaItems: make([]*areaItem, 0),
	}
	app.interceptors = make([]Interceptor, 0)
	app.logger = logger
	app.appConfig = appConfig
	app.finalize = &finalizeInterceptor{}
	return app
}

func (app *webApplication) UseRouteMap(f func(*RouteMap)) {
	f(app.routes)
	for _, area := range app.routes.areaItems {
		group := app.engine.Group(area.uri)
		for _, r := range area.routeItems {
			if strings.Compare(strings.ToLower(r.method), "get") == 0 {
				group.GET(r.uri, app.wrapGinHandler(r))
			} else if strings.Compare(strings.ToLower(r.method), "post") == 0 {
				group.POST(r.uri, app.wrapGinHandler(r))
			} else {
				panic(NewI18NError(i18n.ERR_UNSUPPORTED_METHOD, i18n.GetMessage(i18n.ERR_UNSUPPORTED_METHOD)))
			}
		}
	}
}

func (app *webApplication) UseInterceptors(interceptors ...Interceptor) {
	app.interceptors = append(app.interceptors, interceptors...)
}

func (app *webApplication) Run() {
	err := app.engine.Run(app.appConfig.MustGetString("service.addr"))
	if err != nil {
		panic(err)
	}
}

func (app *webApplication) wrapGinHandler(routerItem *routeItem) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Errorf("panic: %v\n", err)
			}
		}()
		for _, interceptor := range app.interceptors {
			err := interceptor.OnExecuting(routerItem.uri, routerItem.method, c)
			if err != nil {
				app.logger.Error(err)
				c.JSON(http.StatusOK, &ActionResult{
					Code:    i18n.ERR_SERVICE_UNAVAILABLE_498,
					Msg:     i18n.GetMessage(i18n.ERR_SERVICE_UNAVAILABLE_498),
					TipType: TIPTYPE_ERROR,
				})
			}
		}
		result, execErr := routerItem.f(c)
		for _, interceptor := range app.interceptors {
			err := interceptor.OnExecuted(routerItem.uri, routerItem.method, result, execErr, c)
			if err != nil {
				app.logger.Error(err)
				c.JSON(http.StatusOK, &ActionResult{
					Code:    i18n.ERR_SERVICE_UNAVAILABLE_499,
					Msg:     i18n.GetMessage(i18n.ERR_SERVICE_UNAVAILABLE_499),
					TipType: TIPTYPE_ERROR,
				})
			}
		}
		app.finalize.OnExecuted(routerItem.uri, routerItem.method, result, execErr, c)
	}
}

type ActionResult struct {
	Code    string
	Msg     string
	TipType string
	Data    interface{}
}

var _ Interceptor = (*finalizeInterceptor)(nil)

type finalizeInterceptor struct {
	logger logging.Logger
}

func (*finalizeInterceptor) OnExecuting(uri string, method string, ctx *gin.Context) error {
	return nil
}
func (interceptor *finalizeInterceptor) OnExecuted(uri string, method string, result interface{}, err error, ctx *gin.Context) error {
	if r, ok := result.(*ActionResult); ok {
		ctx.JSON(http.StatusOK, r)
	} else if err != nil {
		if i18nErr, ok := err.(*I18NError); ok {
			ctx.JSON(http.StatusOK, &ActionResult{
				Code:    i18nErr.Error(),
				Msg:     i18nErr.Message(),
				TipType: TIPTYPE_ERROR,
			})
		} else {
			interceptor.logger.Error(err)
			ctx.JSON(http.StatusOK, &ActionResult{
				Code:    i18n.ERR_SERVICE_UNAVAILABLE_500,
				Msg:     i18n.GetMessage(i18n.ERR_SERVICE_UNAVAILABLE_500),
				TipType: TIPTYPE_ERROR,
			})
		}
	}
	return nil
}
