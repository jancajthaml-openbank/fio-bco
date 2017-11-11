jest.mock("axios")
jest.mock("../modules/sync.js")
jest.mock("../modules/logger.js")

beforeEach(() => {
  jest.clearAllMocks()
})

test("Tenant.constructor - validates arguments", () => {
  const core = require("../modules/core.js")

  expect(() => { new core.Tenant() }).toThrowError(Error("When creating Tenant you have to provide his name"))
})

test("Tenant.createMissingAccounts - core api returns 500", async () => {
  const axios = require("axios")
  const VError = require("verror")
  const core = require("../modules/core.js")

  const internalServerError = new Error()
  internalServerError.response = {"status": 500}

  const accountMissingError = new Error()
  accountMissingError.response = {"status": 404}

  axios.post = jest.fn()
    .mockImplementationOnce(() => {
      throw internalServerError
    })

  axios.get = jest.fn()
    .mockImplementationOnce(() => {
      throw internalServerError
    })
    .mockImplementationOnce(() => {
      throw accountMissingError
    })

  const testAccounts = [
    {
      "accountNumber": "japa",
      "currency": "JPY",
      "isBalanceCheck": true
    }
  ]

  const tenant = new core.Tenant("test2")
  // Here we test the situation that the core api returns 500 for get account request
  await expect(tenant.createMissingAccounts(testAccounts)).rejects.toEqual(new VError(internalServerError, "Request to core api failed"))
  // Here we test the situation that the core api returns 500 for post account request
  await expect(tenant.createMissingAccounts(testAccounts)).rejects.toEqual(new VError(internalServerError, "Request to core api failed"))
})

test("Tenant.createMissingAccounts - one account from list already exits", async () => {
  const axios = require("axios")
  const core = require("../modules/core.js")

  const accountMissingError = new Error()
  accountMissingError.response = {"status": 404}

  axios.get = jest.fn()
    .mockImplementationOnce(() => {
      throw accountMissingError
    })
    .mockImplementationOnce(() => {
      throw accountMissingError
    })

  const testAccounts = [
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
  ]

  const tenant = new core.Tenant("test")
  await tenant.createMissingAccounts(testAccounts)

  expect(axios.get).toHaveBeenCalledTimes(3)
  expect(axios.post).toHaveBeenCalledTimes(2)
  expect(axios.post.mock.calls[0][0]).toBe("http://127.0.0.1:8080/v1/test/core/account")
  expect(axios.post.mock.calls[0][1]).toEqual({
    "accountNumber": "test1",
    "currency": "USD",
    "isBalanceCheck": false
  })
  expect(axios.post.mock.calls[1][0]).toBe("http://127.0.0.1:8080/v1/test/core/account")
  expect(axios.post.mock.calls[1][1]).toEqual({
    "accountNumber": "test2",
    "currency": "CZK",
    "isBalanceCheck": false
  })
})

test("Tenant.createTransactions - create few transactions", async () => {
  const axios = require("axios")
  const core = require("../modules/core.js")
  const sync = require("../modules/sync.js")

  const testAccountNumber = "test"
  const testToken = "test_token"
  const testTransactions = [
    {
      "blame": "fio-bco",
      "id": "2121115983",
      "transfers": [
        {
          "amount": "0.02",
          "credit": "FIO",
          "currency": "CZK",
          "debit": "CZ7920100000002400222233",
          "id": "1152125621",
          "valueDate": "2016-03-26T23:00:00.000Z"
        },
        {
          "amount": "100",
          "credit": "CZ7920100000002400222233",
          "currency": "CZK",
          "debit": "FIO",
          "id": "1158218819",
          "valueDate": "2016-03-26T23:00:00.000Z"
        }
      ]
    },
    {
      "blame": "fio-bco",
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
    }
  ]

  const tenant = new core.Tenant("test")
  await tenant.createTransactions(testTransactions, testAccountNumber, testToken)

  expect(axios.put).toHaveBeenCalledTimes(2)
  expect(axios.put.mock.calls[0][1]).toEqual(testTransactions[0])
  expect(axios.put.mock.calls[1][1]).toEqual(testTransactions[1])

  expect(sync.setTransactionCheckpoint).toHaveBeenCalledTimes(1)
  expect(sync.setTransactionCheckpoint.mock.calls[0][1]).toBe("test")
  expect(sync.setTransactionCheckpoint.mock.calls[0][2]).toBe(testAccountNumber)
  expect(sync.setTransactionCheckpoint.mock.calls[0][3]).toBe(testToken)
  expect(sync.setTransactionCheckpoint.mock.calls[0][4]).toBe("1158218999")
})

test("Tenant.createTransactions - creating existing transaction in core but with different data", async () => {
  const core = require("../modules/core.js")
  const axios = require("axios")
  const log = require("../modules/logger.js")

  const testTenantName = "test"
  const testToken = "test_token"
  const testAccountNumber = "test_acc_num"
  const testTransactions = [
    {
      "blame": "fio-bco",
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
    }
  ]
  const sameIdsDifferentDataError = new Error()
  sameIdsDifferentDataError.response = {"status": 406}

  axios.put = jest.fn()
    .mockImplementationOnce(() => {
      throw sameIdsDifferentDataError
    })
  log.warn = jest.fn()
    .mockImplementationOnce(() => null)

  const tenant = new core.Tenant(testTenantName)
  await tenant.createTransactions(testTransactions, testAccountNumber, testToken)

  expect(log.warn).toHaveBeenCalledTimes(1)
  expect(log.warn.mock.calls[0][0]).toMatch(new RegExp(`Transaction with ID ${testTransactions[0].id} already exits in core but has different data`))
})

test("Tenant.getCheckpointByToken", async () => {
  const core = require("../modules/core.js")
  const sync = require("../modules/sync.js")

  const testTenantName = "test"
  const testToken = "test_token"

  const tenant = new core.Tenant(testTenantName)

  await tenant.getCheckpointByToken(testToken)

  expect(sync.getTransactionCheckpointByToken).toHaveBeenCalledTimes(1)
  expect(sync.getTransactionCheckpointByToken.mock.calls[0][1]).toBe(testTenantName)
  expect(sync.getTransactionCheckpointByToken.mock.calls[0][2]).toBe(testToken)
})

test("Core.getTransactionCheckpoint", async () => {
  const core = require("../modules/core.js")
  const sync = require("../modules/sync.js")

  const testCheckpoint = 567890
  const tenantName = "test"
  const accountNumber = "AccountNo"

  sync.getTransactionCheckpoint = jest.fn(
    () => testCheckpoint
  )

  const tenant = new core.Tenant(tenantName)

  const checkpoint = await tenant.getCheckpointByAccountNumber(accountNumber)

  expect(sync.getTransactionCheckpoint).toHaveBeenCalledTimes(1)
  expect(sync.getTransactionCheckpoint.mock.calls[0][1]).toBe(tenantName)
  expect(sync.getTransactionCheckpoint.mock.calls[0][2]).toBe(accountNumber)
  expect(checkpoint).toBe(testCheckpoint)
})
