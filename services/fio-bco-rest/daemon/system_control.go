// Copyright (c) 2016-2019, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package daemon

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-rest/config"

	"github.com/coreos/go-systemd/dbus"

	log "github.com/sirupsen/logrus"
)

// SystemControl represents systemctl subroutine
type SystemControl struct {
	Support
	underlying *dbus.Conn
}

// NewSystemControl returns new systemctl fascade
func NewSystemControl(ctx context.Context, cfg config.Configuration) SystemControl {
	conn, err := dbus.New()
	if err != nil {
		panic(fmt.Sprintf("Unable to obtain dbus connection because %+v", err))
	}

	return SystemControl{
		Support:    NewDaemonSupport(ctx),
		underlying: conn,
	}
}

// ListUnits returns list of unit names
func (sys SystemControl) ListUnits(prefix string) ([]string, error) {
	units, err := sys.underlying.ListUnits()
	if err != nil {
		return nil, err
	}

	var result = make([]string, 0)
	for _, unit := range units {
		if unit.LoadState == "not-found" || !strings.HasPrefix(unit.Name, prefix) {
			continue
		}
		result = append(result, strings.TrimSuffix(strings.TrimPrefix(unit.Name, prefix), ".service"))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result, nil
}

// DisableUnit disables unit
func (sys SystemControl) DisableUnit(name string) error {
	log.Debugf("Disabling unit %s", name)

	ch := make(chan string)

	// FIXME
	if _, err := sys.underlying.StopUnit(name, "replace", ch); err != nil {
		return fmt.Errorf("unable to stop unit %s because %+v", name, err)
	}

	select {

	case result := <-ch:
		if result != "done" {
			return fmt.Errorf("unable to stop unit %s", name)
		}
		log.Infof("Stopped unit %s", name)
		log.Infof("Disabling unit %s", name)

		if _, err := sys.underlying.DisableUnitFiles([]string{name}, false); err != nil {
			return fmt.Errorf("unable to disable unit %s because %+v", name, err)
		}

		return nil

	case <-time.After(3 * time.Second):
		return fmt.Errorf("unable to stop unit %s because timeout", name)

	}
}

// EnableUnit enables unit
func (sys SystemControl) EnableUnit(name string) error {
	log.Debugf("Enabling units %s", name)

	if _, _, err := sys.underlying.EnableUnitFiles([]string{name}, false, false); err != nil {
		return fmt.Errorf("unable to enable unit %s because %+v", name, err)
	}

	ch := make(chan string)

	if _, err := sys.underlying.StartUnit(name, "replace", ch); err != nil {
		return fmt.Errorf("unable to start unit %s because %+v", name, err)
	}

	select {

	case result := <-ch:
		if result != "done" {
			return fmt.Errorf("unable to start unit %s", name)
		}
		log.Infof("Started unit %s", name)
		return nil

	case <-time.After(3 * time.Second):
		return fmt.Errorf("unable to start unit %s because timeout", name)

	}

	return nil
}

// WaitReady wait for system to be ready
func (sys SystemControl) WaitReady(deadline time.Duration) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch x := e.(type) {
			case string:
				err = fmt.Errorf(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic")
			}
		}
	}()

	ticker := time.NewTicker(deadline)
	select {
	case <-sys.IsReady:
		ticker.Stop()
		err = nil
		return
	case <-ticker.C:
		err = fmt.Errorf("daemon was not ready within %v seconds", deadline)
		return
	}
}

// Start handles everything needed to start http-server daemon
func (sys SystemControl) Start() {
	defer sys.MarkDone()

	sys.MarkReady()

	select {
	case <-sys.canStart:
		break
	case <-sys.Done():
		return
	}

	log.Info("Start system-control daemon")

	<-sys.exitSignal
}

// Stop shutdowns systemctl fascade
func (sys *SystemControl) Stop() {
	log.Info("Stopping system-control daemon")
	sys.cancel()
	return
}
