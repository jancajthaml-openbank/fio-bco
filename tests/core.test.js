jest.mock("axios");
jest.mock("../modules/sync.js");

beforeEach(() => {
  jest.clearAllMocks();
});

test("create missing acounts, one account from list already exits", async () => {
  let core = require("../modules/core.js");
  let axios = require("axios");
  let accountMissingError = new Error();
  accountMissingError.response = {"status": 404};

  axios.get = jest.fn()
    .mockImplementationOnce(() => {
      throw accountMissingError;
    })
    .mockImplementationOnce(() => {
      throw accountMissingError;
    });

  let testAccounts = [
    {
      "accountNumber": "test1",
      "currency": "USD",
      "isBalanceCheck": false
    },
    {
      "accountNumber": "test2",
      "currency": "CZK",
      "isBalanceCheck": false
    },
    {
      "accountNumber": "test3",
      "currency": "BTC",
      "isBalanceCheck": false
    }
  ];

  let tenant = new core.Tenant("test");
  await tenant.createMissingAccounts(testAccounts);

  expect(axios.get).toHaveBeenCalledTimes(3);
  expect(axios.post).toHaveBeenCalledTimes(2);
  expect(axios.post.mock.calls[0][0]).toBe("http://127.0.0.1:8080/v1/test/core/account/");
  expect(axios.post.mock.calls[0][1]).toEqual({
    "accountNumber": "test1",
    "currency": "USD",
    "isBalanceCheck": false
  });
  expect(axios.post.mock.calls[1][0]).toBe("http://127.0.0.1:8080/v1/test/core/account/");
  expect(axios.post.mock.calls[1][1]).toEqual({
    "accountNumber": "test2",
    "currency": "CZK",
    "isBalanceCheck": false
  });
});

test("Create few transactions", async () => {
  let core = require("../modules/core.js");
  let axios = require("axios");
  let sync = require("../modules/sync.js");

  let testAccountNumber = "test";
  let testTransactions = [{
    "blame": "fio-sync",
    "id": "2121115983",
    "transfers": [{
      "amount": "0.02",
      "credit": "FIO",
      "currency": "CZK",
      "debit": "CZ7920100000002400222233",
      "id": "1152125621",
      "valueDate": "2016-03-26T23:00:00.000Z"
    }, {
      "amount": "100",
      "credit": "CZ7920100000002400222233",
      "currency": "CZK",
      "debit": "FIO",
      "id": "1158218819",
      "valueDate": "2016-03-26T23:00:00.000Z"
    }]
  }, {
    "blame": "fio-sync",
    "id": "2151261787",
    "transfers": [
      {
        "amount": "20",
        "credit": "CZ7920100000002400222233",
        "currency": "CZK",
        "debit": "FIO",
        "id": "1158218999",
        "valueDate": "2016-03-26T23:00:00.000Z"
      }
    ]
  }];

  let tenant = new core.Tenant("test");
  await tenant.createTransactions(testTransactions, testAccountNumber);

  expect(axios.put).toHaveBeenCalledTimes(2);
  expect(axios.put.mock.calls[0][1]).toEqual(testTransactions[0]);
  expect(axios.put.mock.calls[1][1]).toEqual(testTransactions[1]);

  expect(sync.setTransactionCheckpoint).toHaveBeenCalledTimes(1);
  expect(sync.setTransactionCheckpoint.mock.calls[0][1]).toBe("test");
  expect(sync.setTransactionCheckpoint.mock.calls[0][2]).toBe(testAccountNumber);
  expect(sync.setTransactionCheckpoint.mock.calls[0][3]).toBe(1158218999);
});