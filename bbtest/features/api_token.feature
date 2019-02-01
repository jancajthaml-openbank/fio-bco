Feature: API test

  Scenario: Token API - get tokens when application is from scratch
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      SYNC_RATE=1s
      HTTP_PORT=443
    """

    When I request curl GET https://localhost/token/API
    Then curl responds with 200
    """
      []
    """

  Scenario: Token API - create non existant token
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      SYNC_RATE=1s
      HTTP_PORT=443
    """

    When I request curl POST https://localhost/token/API/A
    Then curl responds with 200
    """
      {}
    """

  Scenario: Token API - get tokens
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      SYNC_RATE=1s
      HTTP_PORT=443
    """

    When I request curl POST https://localhost/token/API/B
    Then curl responds with 200
    """
      {}
    """

    When I request curl GET https://localhost/token/API
    Then curl responds with 200
    """
      [
        "A",
        "B"
      ]
    """

  Scenario: Token API - delete existant token
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      SYNC_RATE=1s
      HTTP_PORT=443
    """

    When I request curl POST https://localhost/token/API/C
    Then curl responds with 200
    """
      {}
    """

    When I request curl DELETE https://localhost/token/API/C
    Then curl responds with 200
    """
      {}
    """
