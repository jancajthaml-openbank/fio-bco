Feature: Token management

  Scenario: onboard
  	Given tenant BLACKBOX is onboarded

  Scenario: create token
    When token BLACKBOX/testToken is created
    Then token BLACKBOX/testToken should exist

    When restart unit "fio-bco-import@BLACKBOX.service"
    Then token BLACKBOX/testToken should exist
  
  Scenario: delete token
    When token BLACKBOX/testToken is deleted
    Then token BLACKBOX/testToken should not exist

    When restart unit "fio-bco-import@BLACKBOX.service"
    Then token BLACKBOX/testToken should not exist