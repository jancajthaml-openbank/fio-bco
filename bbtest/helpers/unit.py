#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
from openbank_testkit import Shell, Platform, Package
from systemd import journal


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

  def __fetch_logs(self, unit=None):
    reader = journal.Reader()
    reader.this_boot()
    reader.log_level(journal.LOG_DEBUG)
    if unit:
      reader.add_match(_SYSTEMD_UNIT=unit)
    for entry in reader:
      yield entry['MESSAGE']

  def collect_logs(self):
    cwd = os.path.realpath('{}/../..'.format(os.path.dirname(__file__)))

    logs_dir = '{}/reports/blackbox-tests/logs'.format(cwd)
    os.makedirs(logs_dir, exist_ok=True)

    with open('{}/journal.log'.format(logs_dir), 'w') as fd:
      for line in self.__fetch_logs():
        fd.write(line)
        fd.write(os.linesep)
    
    for unit in set(self.__get_systemd_units() + self.units):
      with open('{}/{}.log'.format(logs_dir, unit), 'w') as fd:
        for line in self.__fetch_logs(unit):
          fd.write(line)
          fd.write(os.linesep)

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
