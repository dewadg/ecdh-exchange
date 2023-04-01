async function generateKeyPair () {
  const key = await window.crypto.subtle.generateKey(
    {
      name: 'ECDH',
      namedCurve: 'P-256'
    },
    false,
    ['deriveKey']
  )

  return [key.privateKey, key.publicKey]
}

async function computeSecret (myPrivateKey, peerPublicKey) {
  const secret = await window.crypto.subtle.deriveKey(
    {
      name: 'ECDH',
      public: peerPublicKey
    },
    myPrivateKey,
    {
      name: 'AES-GCM',
      length: 256
    },
    true,
    ['encrypt', 'decrypt']
  )

  return secret
}

async function importPublicKey (hexEncodedKey = '') {
  const decodedKey = new Uint8Array(hexEncodedKey.match(/.{1,2}/g).map(byte => parseInt(byte, 16))).buffer

  const importedKey = await window.crypto.subtle.importKey(
    'raw',
    decodedKey,
    {
      name: 'ECDH',
      namedCurve: 'P-256'
    },
    false,
    []
  )

  return importedKey
}

async function encodeToHex (key) {
  const byteArray = new Uint8Array(await window.crypto.subtle.exportKey('raw', key))

  return Array.from(byteArray).map(byte => byte.toString(16).padStart(2, '0')).join('')
}
