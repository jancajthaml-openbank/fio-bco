from http.server import BaseHTTPRequestHandler
import json


class RequestHandler(BaseHTTPRequestHandler):

  def log_message(self, format, *args):
    pass

  def __set_last_date(self):
    parts = self.path.split('/')

    if len(parts) < 6:
      return self.__respond(404)

    token = parts[4]
    dateFrom = parts[5]

    try:
      dateFrom = datetime.datetime.strptime(dateFrom, '%Y-%d-%m')
    except:
      return self.__respond(400)

    response = self.server.logic.set_last_date(token, dateFrom)

    if not response:
      return self.__respond(404)

    return self.__respond(200)

  def __set_last_id(self):
    parts = self.path.split('/')

    if len(parts) < 6:
      return self.__respond(404)

    token = parts[4]
    idFrom = parts[5]

    response = self.server.logic.set_last_id(token, idFrom)

    if not response:
      return self.__respond(404)

    return self.__respond(200)

  def __get_last_statements(self):
    parts = self.path.split('/')

    if len(parts) < 6 or parts[5] != 'transactions.json':
      return self.__respond(404)

    token = parts[4]

    response = self.server.logic.__get_last_statements(token)

    if not response:
      return self.__respond(404)

    return self.__respond(200, [])

  def do_GET(self):
    if self.path.startswith('/ib_api/rest/set-last-date'):
      return self.__set_last_date()
    elif self.path.startswith('/ib_api/rest/set-last-id'):
      return self.__set_last_id()
    elif self.path.startswith('/ib_api/rest/last'):
      return self.__get_last_statements()
    else:
      return self.__respond(404)

  def __respond(self, status, body=None):
    self.send_response(status)
    self.send_header('Content-type','application/json')
    self.end_headers()
    if body:
      self.wfile.write(json.dumps(body).encode('utf-8'))
