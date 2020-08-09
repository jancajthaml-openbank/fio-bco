#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import tempfile
import string
import random


class SelfSignedCeritifate(object):

  def __init__(self, name):
    self.__name = name
    self.__ca_key = tempfile.NamedTemporaryFile()
    self.__ca_crt = tempfile.NamedTemporaryFile()
    self.__key = tempfile.NamedTemporaryFile()
    self.__csr = tempfile.NamedTemporaryFile()
    self.__crt = tempfile.NamedTemporaryFile()
    self.__ext = tempfile.NamedTemporaryFile()

    fd = open(self.__ext.name, 'w')
    fd.write('[v3_ca]\n')
    fd.write('subjectAltName = DNS:localhost,IP:127.0.0.1\n')
    fd.write('[v3_req]\n')
    fd.write('extendedKeyUsage=serverAuth\n')
    fd.write('subjectAltName = DNS:localhost,IP:127.0.0.1')

  @property
  def keyfile(self):
    return self.__key.name

  @property
  def certfile(self):
    return self.__crt.name

  def generate(self):
    pwd = ''.join(random.SystemRandom().choice(string.ascii_lowercase) for i in range(10))
    # ca
    os.system('openssl genrsa -des3 -passout pass:{} -out "{}" 2048 > /dev/null 2>&1'.format(pwd, self.__ca_key.name))
    os.system('openssl req -x509 -new -nodes -passin pass:{} -key "{}" -sha256 -days 1825 -out "{}" -subj /CN={}.ca > /dev/null 2>&1'.format(pwd, self.__ca_key.name, self.__ca_crt.name, self.__name))
    # key
    os.system('openssl genrsa -out "{}" 2048 > /dev/null 2>&1'.format(self.__key.name))
    # csr
    os.system('openssl req -new -sha256 -key "{}" -subj /CN=localhost -out "{}" > /dev/null 2>&1'.format(self.__key.name, self.__csr.name))
    # crt
    os.system('openssl x509 -req -passin pass:{} -extfile "{}" -extensions v3_req -extensions v3_ca -in "{}" -CA "{}" -CAkey "{}" -CAcreateserial -out "{}" -days 1 -sha256 > /dev/null 2>&1'.format(pwd, self.__ext.name, self.__csr.name, self.__ca_crt.name, self.__ca_key.name, self.__crt.name))
    del pwd
    # trust ca
    os.system('cp {} /usr/local/share/ca-certificates/{}.crt'.format(self.__ca_crt.name, os.path.basename(self.__ca_crt.name)))
    os.system('update-ca-certificates > /dev/null 2>&1')

  def cleanup(self):
    self.__ca_key.close()
    self.__ca_crt.close()
    self.__key.close()
    self.__csr.close()
    self.__crt.close()
    self.__ext.close()
