// Pure helpers for the 4D-guesser page: which screen to show for a game, and formatting.
// Fork-specific: not part of upstream Photoview.
import {
  GuessGameFields,
  GuessGameFields_rounds,
} from './__generated__/GuessGameFields'

export type LatLng = { latitude: number; longitude: number }

export type GameView =
  | { kind: 'guess'; round: GuessGameFields_rounds }
  | {
      kind: 'reveal'
      round: GuessGameFields_rounds
      hasNext: boolean
    }
  | { kind: 'summary' }

function sortedRounds(game: GuessGameFields) {
  return [...game.rounds].sort((a, b) => a.index - b.index)
}

/** The first round without a guess, or null when every round is guessed. */
export function nextRound(
  game: GuessGameFields
): GuessGameFields_rounds | null {
  return sortedRounds(game).find(round => round.result === null) ?? null
}

export function guessedCount(game: GuessGameFields): number {
  return game.rounds.filter(round => round.result !== null).length
}

/** A game this device can still play. */
export function isGameOpen(game: GuessGameFields): boolean {
  return !game.abandoned && game.finishedAt === null && nextRound(game) !== null
}

/**
 * The screen for a game. revealRoundId is the round the player just guessed: its result
 * stays on screen until they move on. After a reload it is gone, and play continues at
 * the next unguessed round.
 */
export function gameView(
  game: GuessGameFields,
  revealRoundId: string | null
): GameView {
  const next = nextRound(game)

  if (revealRoundId !== null) {
    const revealed = game.rounds.find(round => round.id === revealRoundId)
    if (revealed && revealed.result !== null) {
      return { kind: 'reveal', round: revealed, hasNext: next !== null }
    }
  }

  if (next !== null) return { kind: 'guess', round: next }
  return { kind: 'summary' }
}

/** The shorter way around the globe from `from` to `to`, as a longitude next to from's. */
export function unwrapLongitude(from: number, to: number): number {
  let lon = to
  while (lon - from > 180) lon -= 360
  while (lon - from < -180) lon += 360
  return lon
}

/** Longitude folded into [-180, 180], for clicks on a map that was panned around the globe. */
export function wrapLongitude(lon: number): number {
  const wrapped = ((((lon + 180) % 360) + 360) % 360) - 180
  return wrapped === -180 && lon > 0 ? 180 : wrapped
}

export function formatDistance(km: number, locale?: string): string {
  if (km < 1) {
    return `${Math.round(km * 1000).toLocaleString(locale)} m`
  }
  if (km < 10) {
    return `${km.toLocaleString(locale, { maximumFractionDigits: 1 })} km`
  }
  return `${Math.round(km).toLocaleString(locale)} km`
}

export function formatScore(score: number, locale?: string): string {
  return score.toLocaleString(locale)
}

function safeLocale(locale?: string): string | undefined {
  try {
    if (locale) new Intl.DateTimeFormat(locale)
    return locale
  } catch {
    return undefined
  }
}

/** "March 2019" in the given locale. */
export function formatMonthYear(
  year: number,
  month: number,
  locale?: string
): string {
  return new Intl.DateTimeFormat(safeLocale(locale), {
    year: 'numeric',
    month: 'long',
    timeZone: 'UTC',
  }).format(new Date(Date.UTC(year, month - 1, 15)))
}

/** Short month names, January first. */
export function shortMonthNames(locale?: string): string[] {
  const format = new Intl.DateTimeFormat(safeLocale(locale), {
    month: 'short',
    timeZone: 'UTC',
  })
  return Array.from({ length: 12 }, (_, i) =>
    format.format(new Date(Date.UTC(2001, i, 15)))
  )
}
