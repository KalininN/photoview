// Year slider and month buttons for a 4D-guesser guess.
// Fork-specific: not part of upstream Photoview.
import classNames from 'classnames'
import React from 'react'
import { useTranslation } from 'react-i18next'
import { shortMonthNames } from './guessGameState'

type MonthYearPickerProps = {
  minYear: number
  maxYear: number
  year: number | null
  month: number | null
  onYear(year: number): void
  onMonth(month: number): void
  disabled?: boolean
}

const stepButton =
  'w-8 h-8 rounded border border-gray-200 bg-gray-50 hover:bg-gray-100 disabled:opacity-40 dark:bg-dark-input-bg dark:border-dark-input-border'

const MonthYearPicker = ({
  minYear,
  maxYear,
  year,
  month,
  onYear,
  onMonth,
  disabled,
}: MonthYearPickerProps) => {
  const { t, i18n } = useTranslation()

  // A slider needs room to move, even when every photo is from one year
  const low = Math.min(minYear, maxYear - 1)
  const high = Math.max(maxYear, minYear + 1)
  const sliderValue = year ?? Math.round((low + high) / 2)
  const clampYear = (value: number) => Math.min(high, Math.max(low, value))

  const months = shortMonthNames(i18n.language)

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <button
          type="button"
          className={stepButton}
          disabled={disabled || (year !== null && year <= low)}
          onClick={() => onYear(clampYear((year ?? sliderValue) - 1))}
          aria-label={t('guess_game.year_earlier', 'Earlier year')}
        >
          −
        </button>
        <input
          type="range"
          className="flex-grow"
          min={low}
          max={high}
          step={1}
          value={sliderValue}
          disabled={disabled}
          aria-label={t('guess_game.year', 'Year')}
          onChange={event => onYear(Number(event.target.value))}
          // A click or key press on the thumb without moving it still picks that year
          onPointerUp={event => onYear(Number(event.currentTarget.value))}
          onKeyUp={event => onYear(Number(event.currentTarget.value))}
        />
        <button
          type="button"
          className={stepButton}
          disabled={disabled || (year !== null && year >= high)}
          onClick={() => onYear(clampYear((year ?? sliderValue) + 1))}
          aria-label={t('guess_game.year_later', 'Later year')}
        >
          +
        </button>
        <span
          className={classNames(
            'w-16 text-center text-lg font-semibold tabular-nums',
            year === null && 'text-gray-400'
          )}
          data-testid="guess-year"
        >
          {year ?? '----'}
        </span>
      </div>
      <div
        className="grid grid-cols-6 gap-1"
        role="group"
        aria-label={t('guess_game.month', 'Month')}
      >
        {months.map((name, index) => (
          <button
            key={index}
            type="button"
            disabled={disabled}
            aria-pressed={month === index + 1}
            onClick={() => onMonth(index + 1)}
            className={classNames(
              'py-1 text-sm rounded border',
              month === index + 1
                ? 'bg-blue-500 border-blue-600 text-white'
                : 'bg-gray-50 border-gray-200 hover:bg-gray-100 dark:bg-dark-input-bg dark:border-dark-input-border dark:text-dark-input-text'
            )}
          >
            {name}
          </button>
        ))}
      </div>
    </div>
  )
}

export default MonthYearPicker
