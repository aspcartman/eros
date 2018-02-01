package env

import (
	"fmt"

	"github.com/aspcartman/darkside"
	"github.com/aspcartman/darkside/g"
	"github.com/aspcartman/exceptions"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

var Log logrus.FieldLogger

func init() {
	setupPanicSafenet()
	setupLogging()
	//setupGoroutineLocal()
}

func setupPanicSafenet() {
	raven.SetDSN("https://e4b002434e4a4e89b4bcf2d80291a21e:46a733bd902f41bf940ae676616424ca@sentry.aspc.me/7")

	darkside.SetUnrecoveredPanicHandler(func(p *g.Panic) {
		var title string
		var reason string
		var extra e.Map
		switch arg := p.Arg.(type) {
		case *e.Exception:
			title = arg.Info
			reason = arg.BottommostError().Error()
			extra = arg.Args
		case error:
			title = arg.Error()
		case fmt.Stringer:
			title = arg.String()
		default:
			title = fmt.Sprint(arg)
		}
		ravenCapture(title, reason, extra, 1, true)
		raven.Wait()
		fmt.Println("SENT CRASH PANIC TO SENTRY")
	})

	e.RegisterHook(func(ex *e.Exception) {
		ravenCapture(ex.Info, ex.BottommostError().Error(), ex.Args, 0, false)
	})
}

func ravenCapture(info, value string, extra e.Map, trace int, fatal bool) {
	packet := raven.NewPacket(info, &raven.Exception{
		Type:       info,
		Value:      value,
		Stacktrace: raven.NewStacktrace(trace+1, 3, nil),
	})
	for k, vv := range extra {
		switch v := vv.(type) {
		case string:
			packet.Extra[k] = v
		default:
			packet.Extra[k] = fmt.Sprint(v)
		}
	}
	var tags map[string]string
	if fatal {
		tags = map[string]string{"level": "fatal"}
	}
	raven.Capture(packet, tags)
}

func setupLogging() {
	logger := logrus.StandardLogger()
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true, ForceColors: true}
	e.RegisterHook(func(ex *e.Exception) {
		logger.WithFields(logrus.Fields(ex.Args)).WithError(ex.Cause).Error(ex.Info)
	})
	Log = logger
}
