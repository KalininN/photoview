// What this browser remembers about the 4D-guesser: a random device key, used by the
// server to avoid repeating photos, and the id of the last game played here.
// Other people on the same account never see this, and there is no list of games.
// Fork-specific: not part of upstream Photoview.

const DEVICE_KEY = 'guess.deviceKey'
const LAST_GAME_KEY = 'guess.lastGameId'

// Used when localStorage is blocked, so the key at least lasts for this page load.
let memoryDeviceKey: string | null = null

function readStorage(key: string): string | null {
  try {
    return window.localStorage.getItem(key)
  } catch {
    return null
  }
}

function writeStorage(key: string, value: string | null) {
  try {
    if (value === null) {
      window.localStorage.removeItem(key)
    } else {
      window.localStorage.setItem(key, value)
    }
  } catch {
    // Storage is blocked (private window, disabled site data): forget silently.
  }
}

function randomKey(): string {
  const bytes = new Uint8Array(16)
  if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
    crypto.getRandomValues(bytes)
  } else {
    for (let i = 0; i < bytes.length; i++) {
      bytes[i] = Math.floor(Math.random() * 256)
    }
  }
  return Array.from(bytes, b => b.toString(16).padStart(2, '0')).join('')
}

export function getDeviceKey(): string {
  const stored = readStorage(DEVICE_KEY)
  if (stored) return stored

  if (memoryDeviceKey === null) memoryDeviceKey = randomKey()
  writeStorage(DEVICE_KEY, memoryDeviceKey)
  return memoryDeviceKey
}

export function getLastGameId(): string | null {
  return readStorage(LAST_GAME_KEY)
}

export function setLastGameId(id: string | null) {
  writeStorage(LAST_GAME_KEY, id)
}
