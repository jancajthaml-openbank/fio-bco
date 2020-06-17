import json


class BussinessLogic(object):

  def __init__(self):
    pass

  def set_last_date(self, token, dateFrom):
    return True

  def set_last_id(self, token, idFrom):
    return True

  def get_last_statements(self, token):
    return {
      "accountStatement": {
        "info": {
          "accountId": None,
          "bankId": None,
          "currency": None,
          "iban": None,
          "bic": None,
          "openingBalance": None,
          "closingBalance": None,
          "idFrom": None,
          "idTo": None,
          "idLastDownload": None
        },
        "transactionList": {
          "transaction": []
        }
      }
    }
