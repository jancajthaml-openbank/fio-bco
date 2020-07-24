#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import json
import threading


class BussinessLogic(object):

  def __init__(self):
    self.tenants = dict()
    self.mutex = threading.Lock()

  def create_transaction(self, data):
    #for _ in range(20):
      #print('Creating transaction {}'.format(data))

    return True


#    tenant = data.get('tenant', None)
#    if not tenant:
#      return False
#    if tenant in self.tenants and data['name'] in self.tenants[tenant]:
#      return False
#
#    self.mutex.acquire()
#
#    if not tenant in self.tenants:
#      self.tenants[tenant] = dict()
#    self.tenants[tenant][data['name']] = {
#      'currency': data['currency'],
#      'format': data['format'],
#      'isBalanceCheck': data['isBalanceCheck'],
#    }
#
#    self.mutex.release()
#
#    return True
