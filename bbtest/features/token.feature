Feature: Token management

  Scenario: create token
    Given tenant BLACKBOX is onbdoarded

    When token BLACKBOX/testToken is created
    Then token BLACKBOX/testToken should exist

    When fio-bco is restarted
    Then token BLACKBOX/testToken should exist
