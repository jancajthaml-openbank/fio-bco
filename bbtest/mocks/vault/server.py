#!/usr/bin/env python
# -*- coding: utf-8 -*-

import threading
from http.server import HTTPServer
import ssl
import os
import tempfile

from .handler import RequestHandler
from .logic import BussinessLogic


class VaultMock(threading.Thread):

  def __init__(self, context):
    threading.Thread.__init__(self)
    self.context = context
    self.port = 4400
    self.__keyfile = tempfile.NamedTemporaryFile()
    self.__certfile = tempfile.NamedTemporaryFile()
    os.system('openssl req -x509 -nodes -newkey rsa:2048 -keyout "{}" -out "{}" -days 1 -subj "/C=CZ/ST=Czechia/L=Prague/O=OpenBanking/OU=IT/CN=localhost/emailAddress=jan.cajthaml@gmail.com" > /dev/null 2>&1'.format(self.__keyfile.name, self.__certfile.name))
    os.system('cp {} /usr/local/share/ca-certificates/ > /dev/null 2>&1'.format(self.__certfile.name))
    os.system('update-ca-certificates > /dev/null 2>&1')

  def start(self):
    self.httpd = HTTPServer(('127.0.0.1', self.port), RequestHandler)
    self.httpd.socket = ssl.wrap_socket(self.httpd.socket, certfile=self.__certfile.name, keyfile=self.__keyfile.name, server_side=True)
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
    self.__keyfile.close()
    self.__certfile.close()
    del self.httpd
