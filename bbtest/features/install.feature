Feature: Install package

  Scenario: install
    Given package fio-bco is installed
    Then  systemctl contains following active units
      | name         | type    |
      | fio-bco-rest | service |
      | fio-bco      | service |
      | fio-bco      | path    |
