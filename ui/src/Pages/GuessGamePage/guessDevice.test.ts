import { getDeviceKey, getLastGameId, setLastGameId } from './guessDevice'

afterEach(() => {
  vi.restoreAllMocks()
  localStorage.clear()
})

test('the device key is random and kept', () => {
  const key = getDeviceKey()
  expect(key).toMatch(/^[0-9a-f]{32}$/)
  expect(getDeviceKey()).toBe(key)
  expect(localStorage.getItem('guess.deviceKey')).toBe(key)
})

test('the last game id is remembered and forgotten', () => {
  expect(getLastGameId()).toBeNull()
  setLastGameId('abc')
  expect(getLastGameId()).toBe('abc')
  setLastGameId(null)
  expect(getLastGameId()).toBeNull()
})

test('blocked storage does not break the game', () => {
  vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
    throw new Error('blocked')
  })
  vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
    throw new Error('blocked')
  })

  const key = getDeviceKey()
  expect(key).toMatch(/^[0-9a-f]{32}$/)
  expect(getDeviceKey()).toBe(key)
  expect(getLastGameId()).toBeNull()
  expect(() => setLastGameId('abc')).not.toThrow()
})
