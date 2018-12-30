Feature: API test

  Scenario: Token API - get tokens when application is from scratch
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      FIO_BCO_SYNC_RATE=1s
      FIO_BCO_HTTP_PORT=443
    """

    When I request curl GET https://localhost/tokens/API
    Then curl responds with 200
    """
      []
    """

  Scenario: Token API - create non existant token
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      FIO_BCO_SYNC_RATE=1s
      FIO_BCO_HTTP_PORT=443
    """

    When I request curl POST https://localhost/token/API
    """
      {
        "value": "A"
      }
    """
    Then curl responds with 200
    """
      {}
    """

  Scenario: Token API - get tokens
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      FIO_BCO_SYNC_RATE=1s
      FIO_BCO_HTTP_PORT=443
    """

    When I request curl POST https://localhost/token/API
    """
      {
        "value": "B"
      }
    """
    Then curl responds with 200
    """
      {}
    """

    When I request curl GET https://localhost/tokens/API
    Then curl responds with 200
    """
      [
        {
          "value": "A"
        },
        {
          "value": "B"
        }
      ]
    """

  Scenario: Token API - delete existant token
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      FIO_BCO_SYNC_RATE=1s
      FIO_BCO_HTTP_PORT=443
    """

    When I request curl POST https://localhost/token/API
    """
      {
        "value": "C"
      }
    """
    Then curl responds with 200
    """
      {}
    """

    When I request curl DELETE https://localhost/token/API/C
    Then curl responds with 200
    """
      {}
    """
