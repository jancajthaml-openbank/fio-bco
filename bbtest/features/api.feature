Feature: REST

  Scenario: Tenant API
    Given fio-bco is running
    And   tenant API is onbdoarded

    When I request curl GET https://127.0.0.1:4002/tenant
    Then curl responds with 200
    """
      []
    """

    When I request curl POST https://127.0.0.1:4002/tenant/APITESTA
    Then curl responds with 200
    """
      {}
    """

    When I request curl POST https://127.0.0.1:4002/tenant/APITESTB
    Then curl responds with 200
    """
      {}
    """

    When I request curl GET https://127.0.0.1:4002/tenant
    Then curl responds with 200
    """
      [
        "APITESTB"
      ]
    """

    When I request curl POST https://127.0.0.1:4002/tenant/APITESTC
    Then curl responds with 200
    """
      {}
    """

    When I request curl DELETE https://127.0.0.1:4002/tenant/APITESTC
    Then curl responds with 200
    """
      {}
    """

  Scenario: Token API
    Given fio-bco is running
    And   tenant API is onbdoarded

    When I request curl GET https://127.0.0.1:4002/token/API
    Then curl responds with 200
    """
      []
    """

    When I request curl POST https://127.0.0.1:4002/token/API
    """
      {
        "value": "A"
      }
    """
    Then curl responds with 200

    When I request curl POST https://127.0.0.1:4002/token/API
    """
      {
        "value": "B"
      }
    """
    Then curl responds with 200

    When I request curl GET https://127.0.0.1:4002/token/API
    Then curl responds with 200
