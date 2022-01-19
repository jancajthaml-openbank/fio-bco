#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from openbank_testkit import Shell, Platform, Package


class UnitHelper(object):

  @staticmethod
  def default_config():
    return {
      "LOG_LEVEL": "DEBUG",
      "FIO_GATEWAY": "https://127.0.0.1:4000",
      "SYNC_RATE": "1h",
      "MEMORY_THRESHOLD": 0,
      "STORAGE_THRESHOLD": 0,
      "VAULT_GATEWAY": "https://127.0.0.1:4400",
      "LEDGER_GATEWAY": "https://127.0.0.1:4401",
      "LAKE_HOSTNAME": "127.0.0.1",
      "HTTP_PORT": 443,
      "SERVER_KEY": "/etc/fio-bco/secrets/domain.local.key",
      "SERVER_CERT": "/etc/fio-bco/secrets/domain.local.crt",
      "ENCRYPTION_KEY": "/etc/fio-bco/secrets/fs_encryption.key",
      "STATSD_ENDPOINT": "127.0.0.1:8125",
      "STORAGE": "/data"
    }

  def __init__(self, context):
    self.store = dict()
    self.units = list()
    self.context = context

  def download(self):
    version = os.environ.get('VERSION', '')
    meta = os.environ.get('META', '')

    if version.startswith('v'):
      version = version[1:]

    assert version, 'VERSION not provided'
    assert meta, 'META not provided'

    package = Package('fio-bco')

    cwd = os.path.realpath('{}/../..'.format(os.path.dirname(__file__)))

    assert package.download(version, meta, '{}/packaging/bin'.format(cwd)), 'unable to download package fio-bco'

    self.binary = '{}/packaging/bin/fio-bco_{}_{}.deb'.format(cwd, version, Platform.arch)

  def configure(self, params = None):
    options = dict()
    options.update(UnitHelper.default_config())
    if params:
      options.update(params)

    os.makedirs('/etc/fio-bco/conf.d', exist_ok=True)
    with open('/etc/fio-bco/conf.d/init.conf', 'w') as fd:
      fd.write(str(os.linesep).join("FIO_BCO_{!s}={!s}".format(k, v) for (k, v) in options.items()))

  def collect_logs(self):
    cwd = os.path.realpath('{}/../..'.format(os.path.dirname(__file__)))

    os.makedirs('{}/reports/blackbox-tests/logs'.format(cwd), exist_ok=True)

    (code, result, error) = Shell.run(['journalctl', '-o', 'cat', '--no-pager'])
    if code == 'OK':
      with open('{}/reports/blackbox-tests/logs/journal.log'.format(cwd), 'w') as fd:
        fd.write(result)

    for unit in set(self.__get_systemd_units() + self.units):
      (code, result, error) = Shell.run(['journalctl', '-o', 'cat', '-u', unit, '--no-pager'])
      if code != 'OK' or not result:
        continue
      with open('{}/reports/blackbox-tests/logs/{}.log'.format(cwd, unit), 'w') as fd:
        fd.write(result)

  def teardown(self):
    self.collect_logs()
    for unit in self.__get_systemd_units():
      Shell.run(['systemctl', 'stop', unit])
    self.collect_logs()

  def __get_systemd_units(self):
    (code, result, error) = Shell.run(['systemctl', 'list-units', '--all', '--no-legend'])
    result = [item.replace('*', '').strip().split(' ')[0].strip() for item in result.split(os.linesep)]
    result = [item for item in result if "fio-bco" in item and not item.endswith('unit.slice')]
    return result
