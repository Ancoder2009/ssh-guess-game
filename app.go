package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"

	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
)

type app struct {
	*ssh.Server
	programs []*tea.Program
}

func (a *app) Send(msg tea.Msg) {
	for _, p := range a.programs {
		go p.Send(msg)
	}
}

func NewApp(host string, port string) *app {
	a := new(app)
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(a.ProgramHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	a.Server = s
	return a
}

func (a *app) ProgramHandler(s ssh.Session) *tea.Program {
	_, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil
	}

	model := NewModel(s.User())
	model.app = a

	p := tea.NewProgram(model, append([]tea.ProgramOption{tea.WithColorProfile(colorprofile.TrueColor)}, bubbletea.MakeOptions(s)...)...)
	a.programs = append(a.programs, p)

	return p
}

func (a *app) Start() {
	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	log.Info("Starting SSH server on " + a.Server.Addr)
	go func() {
		if err := a.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case sig := <-sigCh:
		log.Info("Received shutdown signal", "signal", sig.String())
	case err, ok := <-errCh:
		if !ok {
			return
		}
		if err != nil {
			log.Error("Could not start server", "error", err)
		}
		return
	}

	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := a.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}
