import { expect, test } from "vitest";
import { getDEK, encrypt, decrypt } from "../../app/utils/crypto";

const key = "e7qvsq0FsUQWidzMHr59RJSi5I92l+bGkjUgfWQvRt0="; // random key
const chatId = "1234";

test("encryption and decryption", async () => {
  const value = "https://voucher.redeem.gov.sg/?lang=en-US";
  const dek = await getDEK(key);

  const encrypted = await encrypt(dek, value, chatId);

  const decrypted = await decrypt(
    dek,
    encrypted.iv,
    encrypted.ciphertext,
    chatId
  );
  expect(decrypted).toEqual(value);
});
