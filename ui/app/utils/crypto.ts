function base64ToArrayBuffer(base64: string) {
  // Decode the Base64 string into a binary string
  const binaryString = atob(base64);

  // Get the length of the binary string
  const length = binaryString.length;

  // Create a Uint8Array to hold the byte data
  const bytes = new Uint8Array(length);

  // Populate the Uint8Array with the character codes of the binary string
  for (let i = 0; i < length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }

  // Return the underlying ArrayBuffer of the Uint8Array
  return bytes.buffer;
}

function arrayBufferToBase64(buffer: ArrayBuffer | Uint8Array) {
  // Create a Uint*Array to hold the buffer
  const bytes = new Uint8Array(buffer);

  // Get the length of the buffer
  const len = bytes.byteLength;

  // Populate binary string with the character codes of the buffer
  let binary = "";
  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i] as number);
  }

  // Convert to string
  return btoa(binary);
}

async function fetchEncryptionKey(
  keyVersion: number,
  chatId: number,
  initDataRaw: string
): Promise<string> {
  const data = await $fetch<string>("/api/encryption/key", {
    method: "GET",
    params: {
      keyVersion: keyVersion,
      chatId: chatId,
    },
    headers: {
      Authorization: `tma ${initDataRaw}`,
    },
  });
  return data as string;
}

const getDEK = async (base64: string): Promise<CryptoKey> => {
  const raw = base64ToArrayBuffer(base64);
  return await crypto.subtle.importKey(
    "raw",
    raw,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"]
  );
};

const encrypt = async (
  key: CryptoKey,
  plaintext: string,
  addData: string
): Promise<{ iv: string; ciphertext: string }> => {
  const encoder = new TextEncoder();
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const ct = await crypto.subtle.encrypt(
    {
      name: "AES-GCM",
      iv: iv,
      additionalData: encoder.encode(addData),
    },
    key,
    encoder.encode(plaintext)
  );
  return { iv: arrayBufferToBase64(iv), ciphertext: arrayBufferToBase64(ct) };
};

const decrypt = async (
  key: CryptoKey,
  iv: string,
  ciphertext: string,
  addData: string
): Promise<string> => {
  try {
    const plaintext = await crypto.subtle.decrypt(
      {
        name: "AES-GCM",
        iv: base64ToArrayBuffer(iv),
        additionalData: new TextEncoder().encode(addData),
      },
      key,
      base64ToArrayBuffer(ciphertext)
    );
    return new TextDecoder().decode(plaintext);
  } catch (err: any) {
    console.error(err);
    throw new Error("decryption failed");
  }
};

export { fetchEncryptionKey, getDEK, encrypt, decrypt };
