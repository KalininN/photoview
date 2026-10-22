// Fork-specific: not part of upstream Photoview.
import React from 'react'
import { useTranslation } from 'react-i18next'
import { UrlParams } from '../../hooks/useURLParameters'
import Dropdown from '../../primitives/form/Dropdown'

export const DEFAULT_MIN_RATING = 5

/** Reads and writes the `minRating` URL parameter. */
export const minRatingParams = ({
  getParam,
  setParam,
}: Pick<UrlParams, 'getParam' | 'setParam'>) => {
  const param = getParam('minRating')
  const minRating = param !== null ? parseInt(param, 10) : DEFAULT_MIN_RATING
  const setMinRating = (rating: number) =>
    setParam('minRating', rating.toString())

  return { minRating, setMinRating }
}

export type RatingDropdownProps = {
  minRating?: number
  setMinRating?: (rating: number) => void
}

export const RatingDropdown = ({
  minRating,
  setMinRating,
}: RatingDropdownProps) => {
  const { t } = useTranslation()

  const ratingOptions = [
    {
      value: '5',
      label: t('album_filter.rating_options.5', '5'),
    },
    {
      value: '4',
      label: t('album_filter.rating_options.at_least_4', 'at least 4'),
    },
    {
      value: '3',
      label: t('album_filter.rating_options.at_least_3', 'at least 3'),
    },
  ]

  const handleChange = (value: string) => {
    if (setMinRating) {
      setMinRating(parseInt(value, 10))
    }
  }

  return (
    <fieldset>
      <legend id="filter_group_rating-label" className="inline-block mb-1">
        <span>{t('album_filter.rating', 'Rating')}</span>
      </legend>
      <Dropdown
        aria-labelledby="filter_group_rating-label"
        setSelected={handleChange}
        value={minRating?.toString() || undefined}
        items={ratingOptions}
        placeholder={t(
          'album_filter.rating_options.placeholder',
          'Select rating'
        )}
      />
    </fieldset>
  )
}
