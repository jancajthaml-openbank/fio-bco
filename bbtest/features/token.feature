Feature: Token management

  Scenario: create token
    Given tenant BLACKBOX is onboarded

    When token BLACKBOX/testToken is created
    Then token BLACKBOX/testToken should exist

    When restart unit "fio-bco-import@BLACKBOX.service"
    Then token BLACKBOX/testToken should exist
