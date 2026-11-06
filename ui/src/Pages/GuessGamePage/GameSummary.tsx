// The end of a 4D-guesser game: every round and the total.
// Fork-specific: not part of upstream Photoview.
import React from 'react'
import { useTranslation } from 'react-i18next'
import { ProtectedImage } from '../../components/photoGallery/ProtectedMedia'
import { Button } from '../../primitives/form/Input'
import { GuessGameFields } from './__generated__/GuessGameFields'
import { formatDistance, formatMonthYear, formatScore } from './guessGameState'

type GameSummaryProps = {
  game: GuessGameFields
  onNewGame(): void
}

const GameSummary = ({ game, onNewGame }: GameSummaryProps) => {
  const { t, i18n } = useTranslation()
  const locale = i18n.language
  const rounds = [...game.rounds].sort((a, b) => a.index - b.index)

  return (
    <div className="max-w-2xl mx-auto flex flex-col gap-4">
      <h1 className="text-2xl font-bold">
        {t('guess_game.total_score', '{{score}} of {{max}} points', {
          score: formatScore(game.totalScore, locale),
          max: formatScore(game.maxScore, locale),
        })}
      </h1>
      <table className="w-full text-sm">
        <thead>
          <tr className="text-left text-gray-500 dark:text-gray-400">
            <th className="py-1 font-normal">
              {t('guess_game.round', 'Round')}
            </th>
            <th className="py-1 font-normal" />
            <th className="py-1 font-normal">
              {t('guess_game.distance', 'Off by')}
            </th>
            <th className="py-1 font-normal">
              {t('guess_game.taken', 'Taken')}
            </th>
            <th className="py-1 font-normal text-right">
              {t('guess_game.score', 'Score')}
            </th>
          </tr>
        </thead>
        <tbody>
          {rounds.map(round => (
            <tr
              key={round.id}
              className="border-t border-gray-200 dark:border-dark-border"
            >
              <td className="py-2">{round.index + 1}</td>
              <td className="py-2">
                {round.photo?.thumbnail && (
                  <ProtectedImage
                    src={round.photo.thumbnail.url}
                    alt={t('guess_game.round_photo', 'Photo of round {{n}}', {
                      n: round.index + 1,
                    })}
                    className="h-12 w-16 object-cover rounded"
                  />
                )}
              </td>
              <td className="py-2">
                {round.result &&
                  formatDistance(round.result.distanceKm, locale)}
              </td>
              <td className="py-2">
                {round.result &&
                  formatMonthYear(
                    round.result.answerYear,
                    round.result.answerMonth,
                    locale
                  )}
              </td>
              <td className="py-2 text-right tabular-nums">
                {round.result && formatScore(round.result.score, locale)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div>
        <Button onClick={onNewGame} variant="positive">
          {t('guess_game.new_game', 'New game')}
        </Button>
      </div>
    </div>
  )
}

export default GameSummary
