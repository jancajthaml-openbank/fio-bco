#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from behave import *
import ssl
import urllib.request
import socket
import http
import json
import time
from decimal import Decimal
import os


@given('fio gateway contains following statements')
def setup_fio_mock(context):
  pass


@then('token {tenant}/{token} should exist')
def token_exists(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  assert key in context.tokens, 'token does not exist'

  uri = "https://127.0.0.1/token/{}".format(tenant)

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE

  request = urllib.request.Request(method='GET', url=uri)
  request.add_header('Accept', 'application/json')

  try:
    response = urllib.request.urlopen(request, timeout=10, context=ctx)
    assert response.status == 200, str(response.status)

    actual = list()
    for line in response.read().decode('utf-8').split('\n'):
      if not line:
        continue
      actual.append(json.loads(line)['id'])

    assert context.tokens[key] in actual, 'token {} not found in known tokens {}'.format(context.tokens[key], actual)

  except (http.client.RemoteDisconnected, socket.timeout):
    raise AssertionError('timeout')


@then('token {tenant}/{token} should not exist')
def token_not_exists(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  if key not in context.tokens:
    return

  uri = "https://127.0.0.1/token/{}".format(tenant)

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE


  request = urllib.request.Request(method='GET', url=uri)
  request.add_header('Accept', 'application/json')

  try:
    response = urllib.request.urlopen(request, timeout=10, context=ctx)
    assert response.status == 200, str(response.status)

    actual = list()
    for line in response.read().decode('utf-8').split('\n'):
      if not line:
        continue
      actual.append(json.loads(line)['id'])

    assert context.tokens[key] not in actual, 'token {} found in known tokens {}'.format(context.tokens[key], actual)

  except (http.client.RemoteDisconnected, socket.timeout):
    raise AssertionError('timeout')


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

  try:
    response = urllib.request.urlopen(request, timeout=10, context=ctx)
    assert response.status == 200
  except (http.client.RemoteDisconnected, socket.timeout):
    raise AssertionError('timeout')

  key = '{}/{}'.format(tenant, token)

  context.tokens[key] = response.read().decode('utf-8')


@given('token {tenant}/{token} is deleted')
@when('token {tenant}/{token} is deleted')
def create_token(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  assert key in context.tokens, 'token does not exist'

  uri = "https://127.0.0.1/token/{}/{}".format(tenant, context.tokens[key])

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE

  request = urllib.request.Request(method='DELETE', url=uri)

  try:
    response = urllib.request.urlopen(request, timeout=10, context=ctx)
    assert response.status == 200, str(response.status)
  except (http.client.RemoteDisconnected, socket.timeout):
    raise AssertionError('timeout')


@given('token {tenant}/{token} is ordered to synchronize')
@when('token {tenant}/{token} is ordered to synchronize')
def create_token(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  assert key in context.tokens, 'token does not exist'

  uri = "https://127.0.0.1/token/{}/{}/sync".format(tenant, context.tokens[key])

  ctx = ssl.create_default_context()
  ctx.check_hostname = False
  ctx.verify_mode = ssl.CERT_NONE

  request = urllib.request.Request(method='GET', url=uri)

  try:
    response = urllib.request.urlopen(request, timeout=10, context=ctx)
    assert response.status == 200, str(response.status)
  except (http.client.RemoteDisconnected, socket.timeout):
    raise AssertionError('timeout')


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
    context.http_response['content-type'] = response.info().get_content_type()
  except (http.client.RemoteDisconnected, socket.timeout):
    context.http_response['status'] = '504'
    context.http_response['body'] = ""
    context.http_response['content-type'] = 'text-plain'
  except urllib.error.HTTPError as err:
    context.http_response['status'] = str(err.code)
    context.http_response['body'] = err.read().decode('utf-8')
    context.http_response['content-type'] = 'text-plain'


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

    actual = None

    if response['content-type'].startswith('text/plain'):
      actual = list()
      for line in response['body'].split('\n'):
        if not line:
          continue
        if line.startswith('{'):
          actual.append(json.loads(line))
        else:
          actual.append(line)
    elif response['content-type'].startswith('application/json'):
      actual = json.loads(response['body'])
    else:
      actual = response['body']

    try:
      expected = json.loads(context.text)
      diff('', expected, actual)
    except AssertionError as ex:
      raise AssertionError('{} with response {}'.format(ex, response['body']))

