from behave import *
import ssl
import urllib.request
import json
import time
from decimal import Decimal
import os


@given('fio gateway contains following statements')
def setup_fio_mock(context):
  pass


@then('token {tenant}/{token} should exist')
def token_exists(context, tenant, token):
  uri = "https://127.0.0.1/token/{}".format(tenant)

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE

  request = urllib.request.Request(method='GET', url=uri)
  request.add_header('Accept', 'application/json')

  response = urllib.request.urlopen(request, timeout=10, context=ctx)

  assert response.status == 200


@then('token {tenant}/{token} should not exist')
def token_not_exists(context, tenant, token):
  uri = "https://127.0.0.1/token/{}".format(tenant)

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE

  request = urllib.request.Request(method='GET', url=uri)
  request.add_header('Accept', 'application/json')

  assert response.status == 200


@given('token {tenant}/{token} is created')
@when('token {tenant}/{token} is created')
def create_token(context, tenant, token):
  payload = {
    'value': token,
  }

  uri = "https://127.0.0.1/token/{}".format(tenant)

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE

  request = urllib.request.Request(method='POST', url=uri)
  request.add_header('Accept', 'application/json')
  request.add_header('Content-Type', 'application/json')
  request.data = json.dumps(payload).encode('utf-8')

  response = urllib.request.urlopen(request, timeout=10, context=ctx)

  assert response.status == 200

  response = response.read().decode('utf-8')
  response = json.loads(response)


@when('I request HTTP {uri}')
def perform_http_request(context, uri):
  options = dict()
  if context.table:
    for row in context.table:
      options[row['key']] = row['value']

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE

  request = urllib.request.Request(method=options['method'], url=uri)
  request.add_header('Accept', 'application/json')
  if context.text:
    request.add_header('Content-Type', 'application/json')
    request.data = context.text.encode('utf-8')
  context.http_response = dict()

  try:
    response = urllib.request.urlopen(request, timeout=10, context=ctx)
    context.http_response['status'] = str(response.status)
    context.http_response['body'] = response.read().decode('utf-8')
  except urllib.error.HTTPError as err:
    context.http_response['status'] = str(err.code)
    context.http_response['body'] = err.read().decode('utf-8')


@then('HTTP response is')
def check_http_response(context):
  options = dict()
  if context.table:
    for row in context.table:
      options[row['key']] = row['value']

  assert context.http_response
  response = context.http_response
  del context.http_response

  if 'status' in options:
    assert response['status'] == options['status'], 'expected status {} actual {}'.format(options['status'], response['status'])

  if context.text:
    def diff(path, a, b):
      if type(a) == list:
        assert type(b) == list, 'types differ at {} expected: {} actual: {}'.format(path, list, type(b))
        for idx, item in enumerate(a):
          assert item in b, 'value {} was not found at {}[{}]'.format(item, path, idx)
          diff('{}[{}]'.format(path, idx), item, b[b.index(item)])
      elif type(b) == dict:
        assert type(b) == dict, 'types differ at {} expected: {} actual: {}'.format(path, dict, type(b))
        for k, v in a.items():
          assert k in b
          diff('{}.{}'.format(path, k), v, b[k])
      else:
        assert type(a) == type(b), 'types differ at {} expected: {} actual: {}'.format(path, type(a), type(b))
        assert a == b, 'values differ at {} expected: {} actual: {}'.format(path, a, b)

    diff('', json.loads(context.text), json.loads(response['body']))
