Feature: FIO Gateway import

  Scenario: import from gateway token
    Given fio gateway contains following statements
    """
      {
        "accountStatement": {
          "info": {
            "dateStart": "2016-08-03+0200",
            "idList": null,
            "idLastDownload": null,
            "closingBalance": 2060.52,
            "bic": "FIOBCZPPXXX",
            "yearList": null,
            "idTo": 10000000001,
            "currency": "CZK",
            "openingBalance": 2543.81,
            "iban": "CZ1220100000001234567890",
            "idFrom": 10000000002,
            "bankId": "2010",
            "dateEnd": "2016-08-03+0200",
            "accountId": "1234567890"
          },
          "transactionList": {
            "transaction": [
              {
                "column22": {
                  "value": 17184950004,
                  "name": "ID pohybu",
                  "id": 22
                },
                "column0": {
                  "value": "2018-10-31+0100",
                  "name": "Datum",
                  "id": 0
                },
                "column1": {
                  "value": 1.17,
                  "name": "Objem",
                  "id": 1
                },
                "column14": {
                  "value": "CZK",
                  "name": "Měna",
                  "id": 14
                },
                "column2": null,
                "column10": null,
                "column3": null,
                "column12": null,
                "column4": null,
                "column5": null,
                "column6": null,
                "column7": null,
                "column16": null,
                "column8": {
                  "value": "Připsaný úrok",
                  "name": "Typ",
                  "id": 8
                },
                "column9": null,
                "column18": null,
                "column25": null,
                "column26": null,
                "column17": {
                  "value": 20258187372,
                  "name": "ID pokynu",
                  "id": 17
                }
              },
              {
                "column22": {
                  "value": 17184950005,
                  "name": "ID pohybu",
                  "id": 22
                },
                "column0": {
                  "value": "2018-10-31+0100",
                  "name": "Datum",
                  "id": 0
                },
                "column1": {
                  "value": -0.17,
                  "name": "Objem",
                  "id": 1
                },
                "column14": {
                  "value": "CZK",
                  "name": "Měna",
                  "id": 14
                },
                "column2": null,
                "column10": null,
                "column3": null,
                "column12": null,
                "column4": null,
                "column5": null,
                "column6": null,
                "column7": null,
                "column16": null,
                "column8": {
                  "value": "Odvod daně z úroků",
                  "name": "Typ",
                  "id": 8
                },
                "column9": null,
                "column18": null,
                "column25": null,
                "column26": null,
                "column17": {
                  "value": 20258187372,
                  "name": "ID pokynu",
                  "id": 17
                }
              },
              {
                "column18": null,
                "column26": null,
                "column10": null,
                "column12": null,
                "column14": {
                  "name": "Měna",
                  "value": "CZK",
                  "id": 14
                },
                "column17": {
                  "name": "ID pokynu",
                  "value": 12210748893,
                  "id": 17
                },
                "column16": {
                  "name": "Zpráva pro příjemce",
                  "value": "Nákup: ORDR, PRAGUE, CZ, dne 1.8.2016, částka  130.00 CZK",
                  "id": 16
                },
                "column22": {
                  "name": "ID pohybu",
                  "value": 10000000002,
                  "id": 22
                },
                "column9": {
                  "name": "Provedl",
                  "value": "Javorek, Jan",
                  "id": 9
                },
                "column8": {
                  "name": "Typ",
                  "value": "Platba kartou",
                  "id": 8
                },
                "column25": {
                  "name": "Komentář",
                  "value": "Nákup: ORDR, PRAGUE, CZ, dne 1.8.2016, částka  130.00 CZK",
                  "id": 25
                },
                "column5": {
                  "name": "VS",
                  "value": "5678",
                  "id": 5
                },
                "column4": null,
                "column7": {
                  "name": "Uživatelská identifikace",
                  "value": "Nákup: ORDR, PRAGUE, CZ, dne 1.8.2016, částka  130.00 CZK",
                  "id": 7
                },
                "column6": null,
                "column1": {
                  "name": "Objem",
                  "value": -130.0,
                  "id": 1
                },
                "column0": {
                  "name": "Datum",
                  "value": "2016-08-03+0200",
                  "id": 0
                },
                "column3": null,
                "column2": null
              },
              {
                "column18": null,
                "column26": null,
                "column10": null,
                "column12": null,
                "column14": {
                  "name": "Měna",
                  "value": "CZK",
                  "id": 14
                },
                "column17": {
                  "name": "ID pokynu",
                  "value": 12210832097,
                  "id": 17
                },
                "column16": {
                  "name": "Zpráva pro příjemce",
                  "value": "Nákup: Billa Ul. Konevova, Praha - Vitko, CZ, dne 1.8.2016, částka  353.29 CZK",
                  "id": 16
                },
                "column22": {
                  "name": "ID pohybu",
                  "value": 10000000001,
                  "id": 22
                },
                "column9": {
                  "name": "Provedl",
                  "value": "Javorek, Jan",
                  "id": 9
                },
                "column8": {
                  "name": "Typ",
                  "value": "Platba kartou",
                  "id": 8
                },
                "column25": {
                  "name": "Komentář",
                  "value": "Nákup: Billa Ul. Konevova, Praha - Vitko, CZ, dne 1.8.2016, částka  353.29 CZK",
                  "id": 25
                },
                "column5": {
                  "name": "VS",
                  "value": "1234",
                  "id": 5
                },
                "column4": null,
                "column7": {
                  "name": "Uživatelská identifikace",
                  "value": "Nákup: Billa Ul. Konevova, Praha - Vitko, CZ, dne 1.8.2016, částka  353.29 CZK",
                  "id": 7
                },
                "column6": null,
                "column1": {
                  "name": "Objem",
                  "value": -353.29,
                  "id": 1
                },
                "column0": {
                  "name": "Datum",
                  "value": "2016-08-03+0200",
                  "id": 0
                },
                "column3": null,
                "column2": null
              }
            ]
          }
        }
      }
    """
    Given tenant IMPORT is onbdoarded
    And fio-bco is reconfigured with
    """
      FIO_BCO_SYNC_RATE=1s
      FIO_BCO_HTTP_PORT=443
    """
    And token IMPORT/importToken is created
