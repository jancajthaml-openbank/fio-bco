
test("calculate IBAN - Czech Republic", async () => {
  const iban = require("../modules/iban.js")

  expect(iban.calculateCzech("2222222222", "1111")).toBe("CZ4911110000002222222222")
  expect(iban.calculateCzech("2222222222")).toBe("2222222222")
  expect(iban.calculateCzech()).toBeUndefined()
})
