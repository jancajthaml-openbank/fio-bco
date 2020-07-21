#!/usr/bin/env python
# -*- coding: utf-8 -*-

from http.server import BaseHTTPRequestHandler
import json


class RequestHandler(BaseHTTPRequestHandler):

  def log_message(self, format, *args):
    pass

  def do_GET(self):
    parts = self.path.split('/')

    if len(parts) == 3:
      response = self.server.logic.get_acconts(request, parts[2])
      return self.__respond(200, response)

    if len(parts) == 4:
      response = self.server.logic.get_accont(request, parts[2], parts[3])
      if response:
        return self.__respond(200, response)

    return self.__respond(404)

  def do_POST(self):
    parts = self.path.split('/')

    if len(parts) != 3:
      return self.__respond(404)

    request = json.loads(self.rfile.read(int(self.headers['Content-Length'])).decode('utf-8'))
    request['tenant'] = parts[2]

    if not self.server.logic.create_account(request):
      return self.__respond(409)

    return self.__respond(200)

  def __respond(self, status, body=None):
    self.send_response(status)
    self.send_header('Content-type','application/json')
    self.end_headers()
    if body:
      self.wfile.write(json.dumps(body).encode('utf-8'))
