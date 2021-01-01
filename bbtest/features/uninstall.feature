Feature: Uninstall package

  Scenario: uninstall
    Given lake is not running
    And   package fio-bco is uninstalled
    Then  systemctl does not contain following active units
      | name         | type    |
      | fio-bco-rest | service |
      | fio-bco      | service |
      | fio-bco      | path    |
