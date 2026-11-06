// The 4D-guesser page at /4d-guesser. The URL never carries a game id: a game is only
// reachable from the device that started it. Fork-specific: not part of upstream Photoview.
import { ApolloError, useMutation, useQuery } from '@apollo/client'
import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import Layout from '../../components/layout/Layout'
import { MAPBOX_QUERY } from '../../components/layout/MainMenu'
import { NotFoundPage } from '../../components/routes/Routes'
import { mapboxEnabledQuery } from '../../__generated__/mapboxEnabledQuery'
import Loader from '../../primitives/Loader'
import { Button } from '../../primitives/form/Input'
import GameSummary from './GameSummary'
import GuessMap from './GuessMap'
import GuessPhoto from './GuessPhoto'
import MonthYearPicker from './MonthYearPicker'
import RoundReveal from './RoundReveal'
import StartScreen, { StartOptions } from './StartScreen'
import { GuessGameFields } from './__generated__/GuessGameFields'
import { guessGame, guessGameVariables } from './__generated__/guessGame'
import {
  guessGameInfo,
  guessGameInfoVariables,
} from './__generated__/guessGameInfo'
import { guesserEnabled } from './__generated__/guesserEnabled'
import {
  startGuessGame,
  startGuessGameVariables,
} from './__generated__/startGuessGame'
import { submitGuess, submitGuessVariables } from './__generated__/submitGuess'
import { getDeviceKey, getLastGameId, setLastGameId } from './guessDevice'
import {
  GUESSER_ENABLED_QUERY,
  GUESS_GAME_INFO_QUERY,
  GUESS_GAME_QUERY,
  START_GUESS_GAME_MUTATION,
  SUBMIT_GUESS_MUTATION,
} from './guessGameQueries'
import {
  formatScore,
  gameView,
  guessedCount,
  isGameOpen,
  LatLng,
} from './guessGameState'

const errorText = (error?: ApolloError | null) =>
  error?.graphQLErrors[0]?.message ?? error?.message ?? null

const GuessGamePage = () => {
  const { t } = useTranslation()

  const { data: enabledData, loading: enabledLoading } =
    useQuery<guesserEnabled>(GUESSER_ENABLED_QUERY)
  const { data: mapboxData, loading: mapboxLoading } =
    useQuery<mapboxEnabledQuery>(MAPBOX_QUERY)

  const title = t('guess_game.title', '4D-guesser')

  if (enabledLoading || mapboxLoading) {
    return (
      <Layout title={title}>
        <Loader active message={t('general.loading.page', 'Loading page')} />
      </Layout>
    )
  }

  if (!enabledData?.siteInfo.guesserEnabled) {
    return <NotFoundPage t={t} />
  }

  if (!mapboxData?.mapboxToken) {
    return (
      <Layout title={title}>
        <h1 className="text-xl font-semibold mb-2">
          {t('guess_game.no_mapbox_title', 'The map is not set up')}
        </h1>
        <p>
          {t(
            'guess_game.no_mapbox',
            'The 4D-guesser needs a map. Set the MAPBOX_TOKEN environment variable to a Mapbox access token, which can be created for free at mapbox.com.'
          )}
        </p>
      </Layout>
    )
  }

  return (
    <Layout title={title}>
      <GuessGame />
    </Layout>
  )
}

const GuessGame = () => {
  const { t } = useTranslation()
  const [gameId, setGameId] = useState<string | null>(getLastGameId)
  const [playing, setPlaying] = useState(false)
  const [revealRoundId, setRevealRoundId] = useState<string | null>(null)

  const gameQuery = useQuery<guessGame, guessGameVariables>(GUESS_GAME_QUERY, {
    variables: { id: gameId ?? '' },
    skip: gameId === null,
    fetchPolicy: 'cache-and-network',
  })
  const game = gameQuery.data?.guessGame ?? null

  // A game this device no longer reaches (deleted, another account) is forgotten
  useEffect(() => {
    if (gameQuery.error && !gameQuery.data) {
      setLastGameId(null)
      setGameId(null)
      setPlaying(false)
    }
  }, [gameQuery.error, gameQuery.data])

  const [startGame, startState] = useMutation<
    startGuessGame,
    startGuessGameVariables
  >(START_GUESS_GAME_MUTATION)

  const onStart = async (options: StartOptions) => {
    const result = await startGame({
      variables: {
        deviceKey: getDeviceKey(),
        rounds: options.rounds,
        difficulty: options.difficulty,
        filter: { minRating: options.minRating },
        replaceGameId: game && isGameOpen(game) ? game.id : null,
      },
    }).catch(() => null)

    const id = result?.data?.startGuessGame.id
    if (!id) return

    setLastGameId(id)
    setGameId(id)
    setRevealRoundId(null)
    setPlaying(true)
  }

  if (gameId !== null && gameQuery.loading && !game) {
    return (
      <Loader active message={t('general.loading.default', 'Loading...')} />
    )
  }

  if (!playing || !game || game.id !== gameId) {
    return (
      <StartScreen
        lastGame={game}
        onContinue={() => {
          setRevealRoundId(null)
          setPlaying(true)
        }}
        onStart={onStart}
        starting={startState.loading}
        error={errorText(startState.error)}
      />
    )
  }

  return (
    <GamePlay
      game={game}
      revealRoundId={revealRoundId}
      onRevealed={setRevealRoundId}
      onNewGame={() => {
        setRevealRoundId(null)
        setPlaying(false)
      }}
    />
  )
}

type GamePlayProps = {
  game: GuessGameFields
  revealRoundId: string | null
  onRevealed(roundId: string | null): void
  onNewGame(): void
}

const GamePlay = ({
  game,
  revealRoundId,
  onRevealed,
  onNewGame,
}: GamePlayProps) => {
  const { t, i18n } = useTranslation()
  const view = gameView(game, revealRoundId)

  // Year range and map view over every photo the account could play, whatever the
  // game's rating filter, so the picker never hides the right answer
  const { data: infoData } = useQuery<guessGameInfo, guessGameInfoVariables>(
    GUESS_GAME_INFO_QUERY,
    { variables: { filter: { minRating: 3 } }, fetchPolicy: 'cache-first' }
  )
  const info = infoData?.guessGameInfo
  const bounds =
    info?.southWest && info?.northEast
      ? { southWest: info.southWest, northEast: info.northEast }
      : null
  const thisYear = new Date().getFullYear()
  const minYear = info?.minYear ?? 1950
  const maxYear = info?.maxYear ?? thisYear

  const [pin, setPin] = useState<LatLng | null>(null)
  const [year, setYear] = useState<number | null>(null)
  const [month, setMonth] = useState<number | null>(null)

  const roundId = view.kind === 'summary' ? null : view.round.id
  useEffect(() => {
    setPin(null)
    setYear(null)
    setMonth(null)
  }, [roundId])

  const [submit, submitState] = useMutation<submitGuess, submitGuessVariables>(
    SUBMIT_GUESS_MUTATION,
    { refetchQueries: [GUESS_GAME_QUERY], awaitRefetchQueries: true }
  )

  if (view.kind === 'summary') {
    return <GameSummary game={game} onNewGame={onNewGame} />
  }

  const round = view.round
  const rounds = [...game.rounds].sort((a, b) => a.index - b.index)
  const nextPhoto = rounds[round.index + 1]?.photo ?? null
  const result = view.kind === 'reveal' ? round.result : null

  const canGuess =
    view.kind === 'guess' &&
    pin !== null &&
    year !== null &&
    month !== null &&
    !submitState.loading

  const onGuess = async () => {
    if (!canGuess || pin === null || year === null || month === null) return
    const response = await submit({
      variables: {
        roundId: round.id,
        guess: {
          latitude: pin.latitude,
          longitude: pin.longitude,
          year,
          month,
        },
      },
    }).catch(() => null)
    if (response?.data) onRevealed(round.id)
  }

  const header = t(
    'guess_game.progress',
    'Round {{round}} of {{total}} · {{score}} points',
    {
      round: round.index + 1,
      total: rounds.length,
      score: formatScore(game.totalScore, i18n.language),
    }
  )

  let hint: string | null = null
  if (view.kind === 'guess') {
    if (pin === null) {
      hint = t('guess_game.hint_pin', 'Click the map where the photo was taken')
    } else if (year === null || month === null) {
      hint = t('guess_game.hint_date', 'Pick the year and month')
    }
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex justify-between items-baseline gap-2">
        <h1 className="text-lg font-semibold" data-testid="guess-progress">
          {header}
        </h1>
        <span className="text-sm text-gray-500 dark:text-gray-400">
          {t('guess_game.done_count', '{{done}} of {{total}} guessed', {
            done: guessedCount(game),
            total: rounds.length,
          })}
        </span>
      </div>
      <div className="flex flex-col lg:flex-row gap-3 lg:h-[calc(100vh-170px)] lg:min-h-[480px]">
        <div className="h-[45vh] lg:h-auto lg:flex-[3] min-w-0">
          <GuessPhoto photo={round.photo} preload={nextPhoto} />
        </div>
        <div className="flex flex-col gap-3 lg:flex-[2] min-w-0 lg:max-w-[560px]">
          <div className="h-[45vh] lg:h-auto lg:flex-grow min-h-[260px]">
            <GuessMap
              viewKey={round.id}
              pin={result ? result.guess : pin}
              onPin={view.kind === 'guess' ? setPin : undefined}
              answer={result ? result.answer : null}
              bounds={view.kind === 'guess' ? bounds : null}
            />
          </div>
          {view.kind === 'guess' ? (
            <div className="flex flex-col gap-3 pb-20 lg:pb-0">
              <MonthYearPicker
                minYear={minYear}
                maxYear={maxYear}
                year={year}
                month={month}
                onYear={setYear}
                onMonth={setMonth}
                disabled={submitState.loading}
              />
              <div className="flex items-center gap-3">
                <Button
                  variant="positive"
                  className="text-lg py-1.5"
                  disabled={!canGuess}
                  onClick={onGuess}
                >
                  {t('guess_game.guess', 'Guess')}
                </Button>
                {hint && (
                  <span className="text-sm text-gray-500 dark:text-gray-400">
                    {hint}
                  </span>
                )}
              </div>
              {submitState.error && (
                <div className="flex flex-col gap-2 items-start">
                  <p className="text-red-600">{errorText(submitState.error)}</p>
                  <Button onClick={onNewGame}>
                    {t('guess_game.back_to_start', 'Back to the start')}
                  </Button>
                </div>
              )}
            </div>
          ) : (
            <div className="pb-20 lg:pb-0">
              <RoundReveal
                round={round}
                hasNext={view.hasNext}
                onNext={() => onRevealed(null)}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default GuessGamePage
