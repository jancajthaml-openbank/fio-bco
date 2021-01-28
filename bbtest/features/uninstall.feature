Feature: Uninstall package

  Scenario: uninstall
    Given lake is not running
    And   package fio-bco is uninstalled
    Then  systemctl does not contain following active units
      | name            | type    |
      | fio-bco         | service |
      | fio-bco-rest    | service |
      | fio-bco-watcher | path    |
      | fio-bco-watcher | service |

