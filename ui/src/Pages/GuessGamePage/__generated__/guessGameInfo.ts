/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

import { GuessGameFilter } from './../../../__generated__/globalTypes'

// ====================================================
// GraphQL query operation: guessGameInfo
// ====================================================

export interface guessGameInfo_guessGameInfo_southWest {
  __typename: 'Coordinates'
  /**
   * GPS latitude in degrees
   */
  latitude: number
  /**
   * GPS longitude in degrees
   */
  longitude: number
}

export interface guessGameInfo_guessGameInfo_northEast {
  __typename: 'Coordinates'
  /**
   * GPS latitude in degrees
   */
  latitude: number
  /**
   * GPS longitude in degrees
   */
  longitude: number
}

export interface guessGameInfo_guessGameInfo {
  __typename: 'GuessGameInfo'
  /**
   * How many photos have GPS coordinates and an EXIF date
   */
  eligibleCount: number
  /**
   * Earliest year of the photos, for the year picker
   */
  minYear: number | null
  /**
   * Latest year of the photos, for the year picker
   */
  maxYear: number | null
  /**
   * South-west corner of the photos' bounding box, for the initial map view
   */
  southWest: guessGameInfo_guessGameInfo_southWest | null
  /**
   * North-east corner of the photos' bounding box, for the initial map view
   */
  northEast: guessGameInfo_guessGameInfo_northEast | null
}

export interface guessGameInfo {
  /**
   * Describe the photos a game with this filter would use
   */
  guessGameInfo: guessGameInfo_guessGameInfo
}

export interface guessGameInfoVariables {
  filter?: GuessGameFilter | null
}
