#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import zmq
import threading
import time


class ZMQHelper(threading.Thread):

  def __init__(self, context):
    threading.Thread.__init__(self)
    self.__cancel = threading.Event()
    self.__mutex = threading.Lock()
    self.backlog = []
    self.context = context
    self.working = True

  def clear(self):
    self.working = True

  def silence(self):
    self.working = False

  def start(self):
    ctx = zmq.Context.instance()

    self.__pull_url = 'tcp://127.0.0.1:5562'
    self.__pub_url = 'tcp://127.0.0.1:5561'

    self.__pub = ctx.socket(zmq.PUB)
    self.__pub.bind(self.__pub_url)

    self.__pull = ctx.socket(zmq.PULL)
    self.__pull.bind(self.__pull_url)
    self.__pull.set_hwm(100)

    threading.Thread.start(self)

  def run(self):
    while not self.__cancel.is_set():
      try:
        data = self.__pull.recv(zmq.NOBLOCK)
        if not (data and self.working):
          continue
        self.backlog.append(data)
        self.__pub.send(data)
      except Exception as ex:
        if ex.errno != 11:
          return

  def send(self, data):
    self.__pub.send(data.encode())

  def ack(self, data):
    self.__mutex.acquire()
    self.backlog = [item for item in self.backlog if item != data]
    self.__mutex.release()

  def stop(self):
    if self.__cancel.is_set():
      return
    self.__cancel.set()
    try:
      self.join()
    except:
      pass
    self.__pub.close()
    self.__pull.close()
