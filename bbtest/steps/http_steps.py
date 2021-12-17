#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from behave import *
import json
import time
from decimal import Decimal
import os
from helpers.http import Request


@given('fio gateway contains following statements')
def setup_fio_mock(context):
  pass


@then('token {tenant}/{token} should exist')
def token_exists(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  assert key in context.tokens, 'token does not exist'

  uri = "https://127.0.0.1/token/{}".format(tenant)

  request = Request(method='GET', url=uri)
  request.add_header('Accept', 'application/json')

  response = request.do()
  if response.status == 504:
    response = request.do()

  assert response.status == 200, str(response.status)

  actual = list()
  for line in response.read().decode('utf-8').split('\n'):
    if line:
      actual.append(json.loads(line)['id'])

  assert context.tokens[key] in actual, 'token {} not found in known tokens {}'.format(context.tokens[key], actual)


@then('token {tenant}/{token} should not exist')
def token_not_exists(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  if key not in context.tokens:
    return

  uri = "https://127.0.0.1/token/{}".format(tenant)

  request = Request(method='GET', url=uri)
  request.add_header('Accept', 'application/json')

  response = request.do()
  if response.status == 504:
    response = request.do()

  assert response.status == 200, str(response.status)

  actual = list()
  for line in response.read().decode('utf-8').split('\n'):
    if line:
      actual.append(json.loads(line)['id'])

  assert context.tokens[key] not in actual, 'token {} found in known tokens {}'.format(context.tokens[key], actual)


@given('token {tenant}/{token} is created')
@when('token {tenant}/{token} is created')
def create_token(context, tenant, token):
  payload = {
    'value': token,
  }

  uri = "https://127.0.0.1/token/{}".format(tenant)

  request = Request(method='POST', url=uri)
  request.add_header('Accept', 'application/json')
  request.add_header('Content-Type', 'application/json')
  request.data = json.dumps(payload)

  response = request.do()
  if response.status == 504:
    response = request.do()

  assert response.status == 200, str(response.status)
  key = '{}/{}'.format(tenant, token)
  context.tokens[key] = response.read().decode('utf-8')


@given('token {tenant}/{token} is deleted')
@when('token {tenant}/{token} is deleted')
def create_token(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  assert key in context.tokens, 'token does not exist'

  uri = "https://127.0.0.1/token/{}/{}".format(tenant, context.tokens[key])

  request = Request(method='DELETE', url=uri)

  response = request.do()
  if response.status == 504:
    response = request.do()

  assert response.status == 200, str(response.status)


@given('token {tenant}/{token} is ordered to synchronize')
@when('token {tenant}/{token} is ordered to synchronize')
def create_token(context, tenant, token):
  key = '{}/{}'.format(tenant, token)
  assert key in context.tokens, 'token does not exist'

  uri = "https://127.0.0.1/token/{}/{}/sync".format(tenant, context.tokens[key])

  request = Request(method='GET', url=uri)

  response = request.do()
  if response.status == 504:
    response = request.do()

  assert response.status == 200, str(response.status)


@when('I request HTTP {uri}')
def perform_http_request(context, uri):
  options = dict()
  if context.table:
    for row in context.table:
      options[row['key']] = row['value']

  request = Request(method=options['method'], url=uri)
  request.add_header('Accept', 'application/json')
  if context.text:
    request.add_header('Content-Type', 'application/json')
    request.data = context.text

  response = request.do()

  context.http_response = {
    'status': str(response.status),
    'body': response.read().decode('utf-8'),
    'content-type': response.info().get_content_type()
  }


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
    except AssertionError as err:
      raise AssertionError('{} with response {}'.format(err, response['body']))

