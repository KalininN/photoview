// The 4D-guesser start screen: continue this device's last game or start a new one.
// Fork-specific: not part of upstream Photoview.
import { useQuery } from '@apollo/client'
import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { GuessDifficulty } from '../../__generated__/globalTypes'
import Dropdown from '../../primitives/form/Dropdown'
import { Button } from '../../primitives/form/Input'
import { GuessGameFields } from './__generated__/GuessGameFields'
import {
  guessGameInfo,
  guessGameInfoVariables,
} from './__generated__/guessGameInfo'
import { GUESS_GAME_INFO_QUERY } from './guessGameQueries'
import { formatScore, guessedCount, isGameOpen } from './guessGameState'

export const DEFAULT_MIN_RATING = 4
export const DEFAULT_ROUNDS = 5

export type StartOptions = {
  rounds: number
  difficulty: GuessDifficulty
  minRating: number
}

type StartScreenProps = {
  lastGame: GuessGameFields | null
  onContinue(): void
  onStart(options: StartOptions): void
  starting: boolean
  error?: string | null
}

const Field = ({
  label,
  children,
}: {
  label: string
  children: React.ReactNode
}) => (
  <label className="flex flex-col gap-1 text-sm">
    <span className="text-xs uppercase font-semibold">{label}</span>
    {children}
  </label>
)

const StartScreen = ({
  lastGame,
  onContinue,
  onStart,
  starting,
  error,
}: StartScreenProps) => {
  const { t, i18n } = useTranslation()
  const locale = i18n.language

  const [rounds, setRounds] = useState(DEFAULT_ROUNDS)
  const [difficulty, setDifficulty] = useState(GuessDifficulty.AUTO)
  const [minRating, setMinRating] = useState(DEFAULT_MIN_RATING)

  const { data: infoData, loading: infoLoading } = useQuery<
    guessGameInfo,
    guessGameInfoVariables
  >(GUESS_GAME_INFO_QUERY, {
    variables: { filter: { minRating } },
    fetchPolicy: 'cache-and-network',
  })
  const eligible = infoData?.guessGameInfo.eligibleCount ?? null
  const gameRounds = eligible === null ? rounds : Math.min(rounds, eligible)

  const canContinue = lastGame !== null && isGameOpen(lastGame)
  const lastFinished =
    lastGame !== null && !lastGame.abandoned && lastGame.finishedAt !== null

  let poolText: string | null = null
  if (eligible === 0) {
    poolText = t(
      'guess_game.pool_empty',
      'No photos can be used: a photo needs GPS coordinates and a date in its EXIF data, and a high enough rating.'
    )
  } else if (eligible !== null && eligible < rounds) {
    poolText = t(
      'guess_game.pool_small',
      'Only {{count}} photos can be used, so the game has {{count}} rounds.',
      { count: eligible }
    )
  } else if (eligible !== null) {
    poolText = t(
      'guess_game.pool_size',
      '{{count}} photos with a place and a date',
      { count: eligible }
    )
  }

  return (
    <div className="max-w-xl mx-auto flex flex-col gap-6 py-4">
      <div>
        <h1 className="text-2xl font-bold mb-2">
          {t('guess_game.title', '4D-guesser')}
        </h1>
        <p className="text-gray-600 dark:text-gray-300">
          {t(
            'guess_game.intro',
            'Guess where and when each photo was taken: drop a pin on the map and pick the year and month. Each round scores up to 5,000 points for the place and 5,000 for the time.'
          )}
        </p>
      </div>

      {canContinue && lastGame && (
        <div className="flex flex-col gap-2">
          <Button
            variant="positive"
            className="self-start text-lg py-2"
            onClick={onContinue}
            autoFocus
          >
            {t(
              'guess_game.continue',
              'Continue (round {{round}} of {{total}})',
              {
                round: guessedCount(lastGame) + 1,
                total: lastGame.rounds.length,
              }
            )}
          </Button>
        </div>
      )}

      {lastFinished && lastGame && (
        <p className="text-gray-600 dark:text-gray-300">
          {t(
            'guess_game.last_score',
            'Your last game here: {{score}} of {{max}} points.',
            {
              score: formatScore(lastGame.totalScore, locale),
              max: formatScore(lastGame.maxScore, locale),
            }
          )}
        </p>
      )}

      <div className="flex flex-col gap-4 border-t border-gray-200 dark:border-dark-border pt-4">
        <h2 className="text-lg font-semibold">
          {t('guess_game.new_game', 'New game')}
        </h2>
        <div className="flex flex-wrap gap-4">
          <Field label={t('guess_game.rounds', 'Rounds')}>
            <Dropdown
              items={[3, 5, 10].map(n => ({
                value: String(n),
                label: String(n),
              }))}
              selected={String(rounds)}
              setSelected={value => setRounds(Number(value))}
            />
          </Field>
          <Field label={t('guess_game.difficulty', 'Difficulty')}>
            <Dropdown
              items={[
                {
                  value: GuessDifficulty.AUTO,
                  label: t('guess_game.difficulty_auto', 'Auto'),
                },
                {
                  value: GuessDifficulty.CITY,
                  label: t('guess_game.difficulty_city', 'City'),
                },
                {
                  value: GuessDifficulty.REGION,
                  label: t('guess_game.difficulty_region', 'Region'),
                },
                {
                  value: GuessDifficulty.WORLD,
                  label: t('guess_game.difficulty_world', 'World'),
                },
              ]}
              selected={difficulty}
              setSelected={value => setDifficulty(value as GuessDifficulty)}
            />
          </Field>
          <Field label={t('guess_game.min_rating', 'Photos rated')}>
            <Dropdown
              items={[
                { value: '5', label: '★★★★★' },
                { value: '4', label: t('guess_game.rating_4', '★★★★ and up') },
                { value: '3', label: t('guess_game.rating_all', 'All') },
              ]}
              selected={String(minRating)}
              setSelected={value => setMinRating(Number(value))}
            />
          </Field>
        </div>
        <p
          className="text-sm text-gray-600 dark:text-gray-300 min-h-[1.25rem]"
          data-testid="guess-pool"
        >
          {infoLoading && eligible === null
            ? t('general.loading.default', 'Loading...')
            : poolText}
        </p>
        <div>
          <Button
            variant={canContinue ? 'default' : 'positive'}
            className="text-lg py-2"
            disabled={starting || !eligible}
            onClick={() =>
              onStart({ rounds: gameRounds, difficulty, minRating })
            }
          >
            {canContinue
              ? t('guess_game.start_over', 'Start a new game instead')
              : t('guess_game.start', 'Start')}
          </Button>
        </div>
        {error && <p className="text-red-600">{error}</p>}
      </div>
    </div>
  )
}

export default StartScreen
