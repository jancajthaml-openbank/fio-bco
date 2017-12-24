jest.mock("axios")

beforeEach(() => {
  jest.clearAllMocks()
})

test("FIO Api crash is caught", async () => {
  const axios = require("axios")
  const VError = require("verror")
  const fio = require("../modules/fio.js")

  const internalServerError = new VError(new Error(), "Request to FIO api failed")

  axios.get
    .mockImplementationOnce(() => {
      throw new Error()
    })

  await expect(fio.getFioAccountStatement("s4cret", null, false)).rejects.toEqual(internalServerError)
})

test("Extract unique core accounts from fio account statement", () => {
  const fio = require("../modules/fio.js")
  const sampleFioStatement = require("./test-fio-statement.json")

  expect(fio.extractUniqueCoreAccounts(sampleFioStatement))
    .toEqual(expect.arrayContaining([
      {
        "accountNumber": "CZ7920100000002400222233",
        "currency": "CZK",
        "isBalanceCheck": false
      },
      {
        "accountNumber": "FIO",
        "currency": "CZK",
        "isBalanceCheck": false
      },
      {
        "accountNumber": "Counterpart",
        "currency": "CZK",
        "isBalanceCheck": false
      }
    ]))
})

test("Extract core account statement from fio account statement", () => {
  const fio = require("../modules/fio.js")
  const sampleFioStatement = require("./test-fio-statement.json")
  const sampleCoreStatement = require("./test-core-statement.json")

  expect(fio.toCoreAccountStatement(sampleFioStatement))
    .toEqual(sampleCoreStatement)
})

test("Retrieve fio statement data", async () => {
  const fio = require("../modules/fio.js")
  const axios = require("axios")

  axios.get
    .mockReturnValueOnce(null)
    .mockReturnValueOnce({
      "data": {
        "accountStatement": {
          "info": {
            "iban": "test"
          }
        }
      }
    })

  const result = await fio.getFioAccountStatement("s4cret", null, false)

  expect(result.accountStatement.info.iban).toBe("test")
  expect(axios.get.mock.calls[1][0])
})

test("Set position to the beginning", async () => {
  const axios = require("axios")
  const fio = require("../modules/fio.js")

  axios.get
    .mockReturnValueOnce(null)
    .mockReturnValueOnce({
      "data": {
        "accountStatement": {
          "info": {
            "iban": "test"
          }
        }
      }
    })

  const result = await fio.getFioAccountStatement("s4cret", null, false)

  expect(axios.get.mock.calls[0][0]).toBe("https://www.fio.cz/ib_api/rest/set-last-date/s4cret/2012-07-27/")
})

test("Set position to the specific transaction", async () => {
  const axios = require("axios")
  const fio = require("../modules/fio.js")

  axios.get
    .mockReturnValueOnce(null)
    .mockReturnValueOnce({
      "data": {
        "accountStatement": {
          "info": {
            "iban": "test"
          }
        }
      }
    })

  const result = await fio.getFioAccountStatement("s4cret", "12345", false)

  expect(axios.get.mock.calls[0][0]).toBe("https://www.fio.cz/ib_api/rest/set-last-id/s4cret/12345/")
})

test("Test exception on FIO timeout", async () => {
  const axios = require("axios")
  const fio = require("../modules/fio.js")

  const conflictError = new Error()
  conflictError.response = { "status": 409 }

  axios.get
    .mockImplementationOnce(() => null)
    .mockImplementationOnce(() => {
      throw conflictError
    })

  let error
  try {
    await fio.getFioAccountStatement("s4cret", null, false)
  } catch (e) {
    error = e
  }
  expect(error.jse_cause.response.status).toBe(409)
})

test("Test wait on FIO timeout", async () => {
  const fio = require("../modules/fio.js")
  const axios = require("axios")

  const conflictError = new Error()
  conflictError.response = { "status": 409 }

  global.setTimeout = jest.fn((cb, timeout) => {
    cb()
  })

  axios.get
    .mockImplementationOnce(() => null)
    .mockImplementationOnce(() => {
      throw conflictError
    })
    .mockImplementationOnce(() => {
      return {
        "data": {
          "accountStatement": {
            "info": {
              "iban": "test"
            }
          }
        }
      }
    })

  const result = await fio.getFioAccountStatement("s4cret", null, true)
  expect(result.accountStatement.info.iban).toBe("test")
  expect(global.setTimeout.mock.calls[0][1]).toBe(20 * 1000)
})

test("Rethrow unexpected error", async () => {
  const fio = require("../modules/fio.js")
  const axios = require("axios")
  const mockedError = new Error()
  mockedError.response = {"status": 111}

  axios.get
    .mockImplementationOnce(() => null)
    .mockImplementationOnce(() => {
      throw mockedError
    })

  let error
  try {
    await fio.getFioAccountStatement("s4cret", null, true)
  } catch (e) {
    error = e
  }
  expect(error.jse_cause.response.status).toBe(111)
})
