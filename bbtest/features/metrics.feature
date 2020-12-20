Feature: Metrics test

  Scenario: metrics measures expected stats
    Given tenant M2 is onboarded

    Then metrics reports:
      | key                                   | type  |      tags | value |
      | openbank.bco.fio.token.created        | count | tenant:M2 |     0 |
      | openbank.bco.fio.token.deleted        | count | tenant:M2 |     0 |
      | openbank.bco.fio.transaction.imported | count | tenant:M2 |       |
      | openbank.bco.fio.transfer.imported    | count | tenant:M2 |       |

    When token M2/A is created

    Then metrics reports:
      | key                                   | type  |      tags | value |
      | openbank.bco.fio.token.created        | count | tenant:M2 |     1 |
