import React from 'react'
import { fireEvent, render, screen } from '@testing-library/react'
import MonthYearPicker from './MonthYearPicker'

const setup = (year: number | null, month: number | null) => {
  const onYear = vi.fn()
  const onMonth = vi.fn()
  render(
    <MonthYearPicker
      minYear={2005}
      maxYear={2025}
      year={year}
      month={month}
      onYear={onYear}
      onMonth={onMonth}
    />
  )
  return { onYear, onMonth }
}

test('nothing is picked at first', () => {
  setup(null, null)
  expect(screen.getByTestId('guess-year')).toHaveTextContent('----')
  expect(screen.queryAllByRole('button', { pressed: true })).toHaveLength(0)
})

test('the slider and the buttons pick a year within the bounds', () => {
  const { onYear } = setup(2025, null)
  const slider = screen.getByRole('slider', { name: 'Year' })
  expect(slider).toHaveAttribute('min', '2005')
  expect(slider).toHaveAttribute('max', '2025')

  fireEvent.change(slider, { target: { value: '2010' } })
  expect(onYear).toHaveBeenLastCalledWith(2010)

  expect(screen.getByRole('button', { name: 'Later year' })).toBeDisabled()
  fireEvent.click(screen.getByRole('button', { name: 'Earlier year' }))
  expect(onYear).toHaveBeenLastCalledWith(2024)
})

test('months are buttons, January is 1', () => {
  const { onMonth } = setup(2020, 3)
  expect(screen.getByRole('button', { name: 'Mar' })).toHaveAttribute(
    'aria-pressed',
    'true'
  )
  fireEvent.click(screen.getByRole('button', { name: 'Jan' }))
  expect(onMonth).toHaveBeenCalledWith(1)
  fireEvent.click(screen.getByRole('button', { name: 'Dec' }))
  expect(onMonth).toHaveBeenCalledWith(12)
})

test('a pool from one year still gets a usable slider', () => {
  render(
    <MonthYearPicker
      minYear={2020}
      maxYear={2020}
      year={null}
      month={null}
      onYear={vi.fn()}
      onMonth={vi.fn()}
    />
  )
  const slider = screen.getByRole('slider', { name: 'Year' })
  expect(slider).toHaveAttribute('min', '2019')
  expect(slider).toHaveAttribute('max', '2021')
})
