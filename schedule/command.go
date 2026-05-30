package schedule

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/iVampireSP/beacon/foundation"
	"github.com/iVampireSP/beacon/lock"
	"github.com/iVampireSP/beacon/logger"
	jobqueue "github.com/iVampireSP/beacon/queue"
	"github.com/iVampireSP/beacon/tracing"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

// schedulerCommand holds the scheduler command's dependencies.
type schedulerCommand struct {
	app      *foundation.Application
	cron     *cron.Cron
	locker   *lock.Locker
	mq       *jobqueue.Queue
	cronjobs []CronJob
}

// NewSchedulerCommand builds the scheduler command. app is injected by the
// container; cron, lock, queue and cronjobs are resolved lazily in
// PersistentPreRunE so unrelated commands don't spin up the scheduler.
func NewSchedulerCommand(app *foundation.Application) *cobra.Command {
	return (&schedulerCommand{app: app}).Command()
}

// Command constructs the scheduler cobra command.
func (s *schedulerCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "scheduler", Short: "Start cron queue scheduler",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if err := s.app.Invoke(func(c *cron.Cron, l *lock.Locker, mq *jobqueue.Queue) {
				s.cron = c
				s.locker = l
				s.mq = mq
			}); err != nil {
				return err
			}
			// Collect cronjobs contributed by any provider implementing
			// CronProvider, instead of one aggregated dig binding.
			for _, cp := range foundation.ProvidersImplementing[CronProvider](s.app) {
				s.cronjobs = append(s.cronjobs, cp.CronJobs()...)
			}
			return nil
		},
		RunE: func(c *cobra.Command, _ []string) error { return s.Handle(c) },
	}
	cmd.Flags().Bool("once", false, "Run all jobs once and exit")
	cmd.Flags().StringP("queue", "j", "", "Run a specific job by name and exit")
	cmd.Flags().BoolP("list", "l", false, "List all registered cronjobs")
	return cmd
}

// Handle runs the scheduler service.
func (s *schedulerCommand) Handle(cmd *cobra.Command) error {
	runOnce, _ := cmd.Flags().GetBool("once")
	runJob, _ := cmd.Flags().GetString("queue")
	listJobs, _ := cmd.Flags().GetBool("list")

	ctx := cmd.Context()

	tp, err := tracing.GetService("foundation-scheduler")
	if err != nil {
		return err
	}
	defer tracing.ShutdownWithTimeout(tp)

	if s.cron == nil {
		return errors.New("scheduler bootstrap failed: cron instance is nil")
	}

	mutex := NewRedisMutex(s.locker)
	sched := NewScheduler(s.cron, mutex, s.mq, cmd.Root())
	sched.RegisterAll(s.cronjobs)

	if listJobs {
		listAllEvents(sched)
		return nil
	}

	if runJob != "" {
		logger.Info("scheduler: running single", "job", runJob)
		if err := sched.RunEvent(ctx, runJob); err != nil {
			logger.Error("scheduler", "job", runJob, "failed", err)
			return err
		}
		logger.Info("scheduler completed", "job", runJob)
		return nil
	}

	if runOnce {
		logger.Info("scheduler: running jobs once", "all", len(s.cronjobs))
		if err := sched.RunAllEvents(ctx); err != nil {
			logger.Error("scheduler: run-all", "failed", err)
			return err
		}
		logger.Info("scheduler: all jobs completed")
		return nil
	}

	if err := sched.Start(ctx); err != nil {
		return err
	}
	defer sched.Stop()

	s.cron.Start()
	defer s.cron.Stop()

	logger.Info("scheduler: started cronjobs, waiting for signals", "with", len(s.cronjobs))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("scheduler: received signal, stopping", "signal", sig)
	return nil
}

func listAllEvents(sched *Scheduler) {
	events := sched.ListEvents()
	if len(events) == 0 {
		fmt.Println("No cronjobs registered.")
		return
	}

	nameW, exprW, descW := len("NAME"), len("SCHEDULE"), len("DESCRIPTION")
	for _, e := range events {
		if len(e.Name) > nameW {
			nameW = len(e.Name)
		}
		if len(e.Expression) > exprW {
			exprW = len(e.Expression)
		}
		if len(e.Description) > descW {
			descW = len(e.Description)
		}
	}

	rowFmt := fmt.Sprintf("  %%-%ds  %%-%ds  %%-%ds  %%s\n", nameW, exprW, descW)
	fmt.Printf(rowFmt, "NAME", "SCHEDULE", "DESCRIPTION", "FLAGS")
	fmt.Printf("  %s  %s  %s  %s\n",
		strings.Repeat("-", nameW),
		strings.Repeat("-", exprW),
		strings.Repeat("-", descW),
		strings.Repeat("-", 5),
	)
	for _, e := range events {
		var flags []string
		if e.OnOneServer {
			flags = append(flags, "one-server")
		}
		if e.WithoutOverlapping {
			flags = append(flags, "no-overlap")
		}
		fmt.Printf(rowFmt, e.Name, e.Expression, e.Description, strings.Join(flags, ", "))
	}
	fmt.Printf("\n  %d cronjob(s) total\n", len(events))
}
