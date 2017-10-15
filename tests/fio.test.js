jest.mock("axios");

beforeEach(() => {
  jest.clearAllMocks();
});

test("Extract unique core accounts from fio account statement", () => {
  let fio = require("../modules/fio.js");
  let sampleFioStatement = require("./test-fio-statement.json");

  expect(fio.extractUniqueCoreAccounts(sampleFioStatement))
    .toEqual(expect.arrayContaining([{
      "accountNumber": "CZ7920100000002400222233",
      "currency": "CZK",
      "isBalanceCheck": false
    },{
      "accountNumber": "FIO",
      "currency": "CZK",
      "isBalanceCheck": false
    },{
      "accountNumber": "Counterpart",
      "currency": "CZK",
      "isBalanceCheck": false
    }]));
});

test("Extract core account statement from fio account statement", () => {
  let fio = require("../modules/fio.js");
  let sampleFioStatement = require("./test-fio-statement.json");
  let sampleCoreStatement = require("./test-core-statement.json");

  expect(fio.extractCoreAccountStatement(sampleFioStatement))
    .toEqual(sampleCoreStatement);
});

test("Retrieve fio statement data", async () => {
  let fio = require("../modules/fio.js");
  let axios = require("axios");

  axios.get
    .mockReturnValueOnce(null)
    .mockReturnValueOnce({"data": "test"});

  let result = await fio.getFioAccountStatement("s4cret", null, false);

  expect(result).toBe("test");
  expect(axios.get.mock.calls[1][0]);
});

test("Set position to the beginning", async () => {
  let fio = require("../modules/fio.js");
  let axios = require("axios");

  axios.get
    .mockReturnValueOnce(null)
    .mockReturnValueOnce({"data": "test"});

  let result = await fio.getFioAccountStatement("s4cret", null, false);

  expect(axios.get.mock.calls[0][0]).toBe("https://www.fio.cz/ib_api/rest/set-last-date/s4cret/1900-01-01/");
});

test("Set position to the specific transaction", async () => {
  let fio = require("../modules/fio.js");
  let axios = require("axios");

  axios.get
    .mockReturnValueOnce(null)
    .mockReturnValueOnce({"data": "test"});

  let result = await fio.getFioAccountStatement("s4cret", "12345", false);

  expect(axios.get.mock.calls[0][0]).toBe("https://www.fio.cz/ib_api/rest/set-last-id/s4cret/12345/");
});

test("Test exception on FIO timeout", async () => {
  let fio = require("../modules/fio.js");
  let axios = require("axios");
  let mockedError = new Error();
  mockedError.response = { "status": 409 };

  axios.get
    .mockImplementationOnce(() => null)
    .mockImplementationOnce(() => {
      throw mockedError;
    });

  let error;
  try {
    await fio.getFioAccountStatement("s4cret", null, false);
  } catch (e) {
    error = e;
  }
  expect(error).toBe(mockedError);
});

test("Test wait on FIO timeout", async () => {
  let fio = require("../modules/fio.js");
  let axios = require("axios");

  global.setTimeout = jest.fn((cb, timeout) => {
    cb();
  });

  axios.get
    .mockImplementationOnce(() => null)
    .mockImplementationOnce(() => {
      let mockedError = new Error();
      mockedError.response = { "status": 409 };
      throw mockedError;
    })
    .mockImplementationOnce(() => {
      return {"data": "test"};
    });

  let result = await fio.getFioAccountStatement("s4cret", null, true);
  expect(result).toBe("test");
  expect(global.setTimeout.mock.calls[0][1]).toBe(20 * 1000);
});

test("Rethrow unexpected error", async () => {
  let fio = require("../modules/fio.js");
  let axios = require("axios");
  let mockedError = new Error();
  mockedError.response = {"status": 111};

  axios.get
    .mockImplementationOnce(() => null)
    .mockImplementationOnce(() => {
      throw mockedError;
    });

  let error;
  try {
    let result = await fio.getFioAccountStatement("s4cret", null, true);
  } catch (e) {
    error = e;
  }
  expect(error).toBe(mockedError);
});

