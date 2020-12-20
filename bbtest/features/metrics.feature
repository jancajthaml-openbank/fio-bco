Feature: Metrics test

  Scenario: metrics can remembers previous values after reboot
    Given tenant M2 is onboarded

    Then metrics reports:
      | key                                      | type  | value |
      | openbank.bco.fio.M2.token.created        | count |     0 |
      | openbank.bco.fio.M2.token.deleted        | count |     0 |
      | openbank.bco.fio.M2.transaction.imported | count |     0 |
      | openbank.bco.fio.M2.transfer.imported    | count |     0 |

    When token M2/A is created

    Then metrics reports:
      | key                                      | type  | value |
      | openbank.bco.fio.M2.token.created        | count |     1 |
