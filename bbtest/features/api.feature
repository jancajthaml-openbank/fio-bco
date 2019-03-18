Feature: REST

  Scenario: Tenant API
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      LOG_LEVEL=DEBUG
      HTTP_PORT=443
    """

    When I request curl GET https://localhost/tenant
    Then curl responds with 200
    """
      []
    """

    When I request curl POST https://localhost/tenant/APITESTA
    Then curl responds with 200
    """
      {}
    """

    When I request curl POST https://localhost/tenant/APITESTB
    Then curl responds with 200
    """
      {}
    """

    When I request curl GET https://localhost/tenant
    Then curl responds with 200
    """
      [
        "APITESTB"
      ]
    """

    When I request curl POST https://localhost/tenant/APITESTC
    Then curl responds with 200
    """
      {}
    """

    When I request curl DELETE https://localhost/tenant/APITESTC
    Then curl responds with 200
    """
      {}
    """

  Scenario: Token API
    Given tenant API is onbdoarded
    And fio-bco is reconfigured with
    """
      LOG_LEVEL=DEBUG
      HTTP_PORT=443
    """

    When I request curl GET https://localhost/token/API
    Then curl responds with 200
    """
      []
    """

    When I request curl POST https://localhost/token/API
    """
      {
        "value": "A"
      }
    """
    Then curl responds with 200

    When I request curl POST https://localhost/token/API
    """
      {
        "value": "B"
      }
    """
    Then curl responds with 200

    When I request curl GET https://localhost/token/API
    Then curl responds with 200
