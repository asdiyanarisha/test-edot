package scheduler

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"test-edot/src/app/order"
	"test-edot/src/factory"
	"test-edot/util"
	"time"
)

func RunScheduler(f *factory.Factory) {
	wait := make(chan struct{})

	f.Log.Info("start scheduler")

	go func() {
		s := make(chan os.Signal, 1)

		c := cron.New()

		_, err := c.AddFunc(util.GetEnv("RELEASE_STOCK_ORDER_CRON", ""), order.NewService(f).ReleaseStockOrder)
		if err != nil {
			f.Log.Error("Error failed run disbursement", zap.Error(err))
			return
		}

		c.Start()

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		f.Log.Info("gracefully shutdown")
		c.Stop()
		time.Sleep(time.Second * 5)
		f.Log.Info("shutdown done")

		close(wait)
	}()
	<-wait
}
