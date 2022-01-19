#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import threading
from http.server import HTTPServer
from openbank_testkit import SelfSignedCeritifate
import ssl
from .handler import RequestHandler
from .logic import BussinessLogic


class LedgerMock(threading.Thread):

  def __init__(self, context):
    threading.Thread.__init__(self)
    self.context = context
    self.port = 4401

    self.__certificate = SelfSignedCeritifate('ledger-mock')
    self.__certificate.generate()

  def start(self):
    self.httpd = HTTPServer(('127.0.0.1', self.port), RequestHandler)
    self.httpd.socket = ssl.wrap_socket(self.httpd.socket, certfile=self.__certificate.certfile, keyfile=self.__certificate.keyfile, server_side=True)
    self.httpd.logic = BussinessLogic()
    threading.Thread.start(self)

  def run(self):
    self.httpd.serve_forever()

  def stop(self):
    if self.httpd:
      self.httpd.shutdown()
    try:
      self.join()
    except:
      pass
    del self.__certificate
    del self.httpd
