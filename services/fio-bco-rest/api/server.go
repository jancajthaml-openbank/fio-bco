// Copyright (c) 2016-2020, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-rest/actor"
	"github.com/jancajthaml-openbank/fio-bco-rest/system"

	localfs "github.com/jancajthaml-openbank/local-fs"
	"github.com/labstack/echo/v4"
)

const READ_TIMEOUT = 5 * time.Second
const WRITE_TIMEOUT = 5 * time.Second

// Server is a fascade for http-server following handler api of Gin and
// lifecycle api of http
type Server struct {
	underlying *http.Server
	listener   *net.TCPListener
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// NewServer returns new secure server instance
func NewServer(port int, certPath string, keyPath string, rootStorage string, actorSystem *actor.System, systemControl system.Control, diskMonitor system.CapacityCheck, memoryMonitor system.CapacityCheck) *Server {
	storage, err := localfs.NewPlaintextStorage(rootStorage)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}

	router := echo.New()

	certificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Error().Msgf("Invalid cert %s and key %s", certPath, keyPath)
		return nil
	}

	router.GET("/health", HealtCheck(memoryMonitor, diskMonitor))
	router.HEAD("/health", HealtCheckPing(memoryMonitor, diskMonitor))

	router.GET("/tenant", ListTenants(systemControl))
	router.POST("/tenant/:tenant", CreateTenant(systemControl))
	router.DELETE("/tenant/:tenant", DeleteTenant(systemControl))

	router.DELETE("/token/:tenant/:id", DeleteToken(actorSystem))
	router.POST("/token/:tenant", CreateToken(actorSystem))
	router.GET("/token/:tenant", GetTokens(storage))

	return &Server{
		underlying: &http.Server{
			Addr:         fmt.Sprintf("127.0.0.1:%d", port),
			ReadTimeout:  READ_TIMEOUT,
			WriteTimeout: WRITE_TIMEOUT,
			Handler:      router,
			TLSConfig: &tls.Config{
				MinVersion:               tls.VersionTLS12,
				MaxVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
				InsecureSkipVerify:       false,
				CurvePreferences: []tls.CurveID{
					tls.CurveP521,
					tls.CurveP384,
					tls.CurveP256,
				},
				CipherSuites: CipherSuites,
				Certificates: []tls.Certificate{
					certificate,
				},
			},
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		},
		listener: nil,
	}
}

func (server *Server) Setup() error {
	if server == nil {
		return fmt.Errorf("nil pointer")
	}
	ln, err := net.Listen("tcp", server.underlying.Addr)
	if err != nil {
		return err
	}
	server.listener = ln.(*net.TCPListener)
	return nil
}

func (server *Server) Done() <-chan interface{} {
	done := make(chan interface{})
	close(done)
	return done
}

func (server *Server) Cancel() {
	if server == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), WRITE_TIMEOUT)
	defer cancel()
	server.underlying.Shutdown(ctx)
}

func (server *Server) Work() {
	if server == nil {
		return
	}
	log.Info().Msgf("Server listening on %s", server.underlying.Addr)
	tlsListener := tls.NewListener(tcpKeepAliveListener{server.listener}, server.underlying.TLSConfig)
	err := server.underlying.Serve(tlsListener)
	if err != nil && err != http.ErrServerClosed {
		log.Error().Msg(err.Error())
	}
}
