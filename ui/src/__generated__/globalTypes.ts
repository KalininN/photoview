/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

/**
 * How fast scores fall with distance. AUTO derives it from the spread of the photos
 */
export enum GuessDifficulty {
  AUTO = "AUTO",
  CITY = "CITY",
  REGION = "REGION",
  WORLD = "WORLD",
}

/**
 * Supported language translations of the user interface
 */
export enum LanguageTranslation {
  Basque = "Basque",
  Danish = "Danish",
  Dutch = "Dutch",
  English = "English",
  French = "French",
  German = "German",
  Italian = "Italian",
  Japanese = "Japanese",
  Polish = "Polish",
  Portuguese = "Portuguese",
  Russian = "Russian",
  SimplifiedChinese = "SimplifiedChinese",
  Spanish = "Spanish",
  Swedish = "Swedish",
  TraditionalChineseHK = "TraditionalChineseHK",
  TraditionalChineseTW = "TraditionalChineseTW",
  Turkish = "Turkish",
  Ukrainian = "Ukrainian",
}

export enum MediaType {
  Photo = "Photo",
  Video = "Video",
}

/**
 * Specified the type a particular notification is of
 */
export enum NotificationType {
  Close = "Close",
  Message = "Message",
  Progress = "Progress",
}

/**
 * Used to specify which order to sort items in
 */
export enum OrderDirection {
  ASC = "ASC",
  DESC = "DESC",
}

/**
 * Which photos a 4D-guesser game picks from
 */
export interface GuessGameFilter {
  minRating?: number | null;
  albumIds?: string[] | null;
}

/**
 * A guess for one round
 */
export interface GuessInput {
  latitude: number;
  longitude: number;
  year: number;
  month: number;
}

//==============================================================
// END Enums and Input Objects
//==============================================================
