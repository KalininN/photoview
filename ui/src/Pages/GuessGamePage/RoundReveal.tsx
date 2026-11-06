// The result of a 4D-guesser round, shown next to the map with both pins.
// Fork-specific: not part of upstream Photoview.
import React from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'
import { Button } from '../../primitives/form/Input'
import { GuessGameFields_rounds } from './__generated__/GuessGameFields'
import { formatDistance, formatMonthYear, formatScore } from './guessGameState'

type RoundRevealProps = {
  round: GuessGameFields_rounds
  hasNext: boolean
  onNext(): void
}

const ScoreLine = ({
  label,
  value,
  score,
}: {
  label: string
  value: string
  score: string
}) => (
  <div className="flex justify-between gap-2">
    <span>
      <span className="text-gray-500 dark:text-gray-400">{label}</span> {value}
    </span>
    <span className="font-semibold tabular-nums">{score}</span>
  </div>
)

const RoundReveal = ({ round, hasNext, onNext }: RoundRevealProps) => {
  const { t, i18n } = useTranslation()
  const result = round.result
  if (result === null) return null

  const locale = i18n.language

  return (
    <div className="flex flex-col gap-2" data-testid="round-reveal">
      <div className="text-2xl font-bold tabular-nums">
        {t('guess_game.round_score', '{{score}} points', {
          score: formatScore(result.score, locale),
        })}
      </div>
      <ScoreLine
        label={t('guess_game.distance', 'Off by')}
        value={formatDistance(result.distanceKm, locale)}
        score={formatScore(result.locationScore, locale)}
      />
      <ScoreLine
        label={t('guess_game.taken', 'Taken')}
        value={formatMonthYear(result.answerYear, result.answerMonth, locale)}
        score={formatScore(result.timeScore, locale)}
      />
      <div className="text-sm text-gray-500 dark:text-gray-400">
        {t('guess_game.your_guess', 'You guessed {{date}}', {
          date: formatMonthYear(result.guessYear, result.guessMonth, locale),
        })}
        {result.monthsOff > 0 &&
          ` · ${t('guess_game.months_off', '{{count}} months off', {
            count: result.monthsOff,
          })}`}
      </div>
      <div className="flex items-center gap-4 mt-2">
        <Button onClick={onNext} variant="positive" autoFocus>
          {hasNext
            ? t('guess_game.next_round', 'Next round')
            : t('guess_game.see_results', 'See results')}
        </Button>
        {round.photo && (
          <Link
            to={`/album/${round.photo.album.id}`}
            target="_blank"
            className="text-sm text-blue-600 dark:text-blue-400 hover:underline"
          >
            {t('guess_game.open_album', 'Open its album')}
          </Link>
        )}
      </div>
    </div>
  )
}

export default RoundReveal
