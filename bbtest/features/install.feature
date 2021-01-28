Feature: Install package

  Scenario: install
    Given package fio-bco is installed
    Then  systemctl contains following active units
      | name            | type    |
      | fio-bco         | service |
      | fio-bco-rest    | service |
      | fio-bco-watcher | path    |
      | fio-bco-watcher | service |
