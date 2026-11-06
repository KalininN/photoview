// 4D-guesser GraphQL documents. Fork-specific: not part of upstream Photoview.
import { gql } from '@apollo/client'

export const GUESSER_ENABLED_QUERY = gql`
  query guesserEnabled {
    siteInfo {
      guesserEnabled
    }
  }
`

export const GUESS_GAME_INFO_QUERY = gql`
  query guessGameInfo($filter: GuessGameFilter) {
    guessGameInfo(filter: $filter) {
      eligibleCount
      minYear
      maxYear
      southWest {
        latitude
        longitude
      }
      northEast {
        latitude
        longitude
      }
    }
  }
`

// The photo's title, path and EXIF are never requested, so the UI cannot show them.
const GUESS_GAME_FIELDS = gql`
  fragment GuessGameFields on GuessGame {
    id
    createdAt
    finishedAt
    abandoned
    difficulty
    totalScore
    maxScore
    rounds {
      id
      index
      photo {
        id
        blurhash
        album {
          id
        }
        thumbnail {
          url
          width
          height
        }
        highRes {
          url
          width
          height
        }
      }
      result {
        guess {
          latitude
          longitude
        }
        guessYear
        guessMonth
        answer {
          latitude
          longitude
        }
        answerYear
        answerMonth
        distanceKm
        monthsOff
        locationScore
        timeScore
        score
      }
    }
  }
`

export const GUESS_GAME_QUERY = gql`
  ${GUESS_GAME_FIELDS}
  query guessGame($id: String!) {
    guessGame(id: $id) {
      ...GuessGameFields
    }
  }
`

export const START_GUESS_GAME_MUTATION = gql`
  ${GUESS_GAME_FIELDS}
  mutation startGuessGame(
    $deviceKey: String!
    $rounds: Int
    $difficulty: GuessDifficulty
    $filter: GuessGameFilter
    $replaceGameId: String
  ) {
    startGuessGame(
      deviceKey: $deviceKey
      rounds: $rounds
      difficulty: $difficulty
      filter: $filter
      replaceGameId: $replaceGameId
    ) {
      ...GuessGameFields
    }
  }
`

export const SUBMIT_GUESS_MUTATION = gql`
  mutation submitGuess($roundId: String!, $guess: GuessInput!) {
    submitGuess(roundId: $roundId, guess: $guess) {
      score
    }
  }
`
