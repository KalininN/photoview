import React from 'react'
import { MockedProvider } from '@apollo/client/testing'
import { render, screen } from '@testing-library/react'

import * as authentication from '../../helpers/authentication'
import { ADMIN_QUERY } from './Layout'
import { MemoryRouter } from 'react-router-dom'
import MainMenu, { MAPBOX_QUERY } from './MainMenu'
import { GUESSER_ENABLED_QUERY } from '../../Pages/GuessGamePage/guessGameQueries'

vi.mock('../../helpers/authentication.ts')

const authTokenMock = vi.mocked(authentication.authToken)

afterEach(() => {
  authTokenMock.mockClear()
})

test('Layout sidebar component', async () => {
  authTokenMock.mockImplementation(() => 'test-token')

  const mockedGraphql = [
    {
      request: {
        query: ADMIN_QUERY,
      },
      result: {
        data: {
          myUser: {
            admin: true,
          },
        },
      },
    },
    {
      request: {
        query: MAPBOX_QUERY,
      },
      result: {
        data: {
          mapboxToken: true,
        },
      },
    },
  ]

  render(
    <MockedProvider mocks={mockedGraphql} addTypename={false}>
      <MemoryRouter>
        <MainMenu />
      </MemoryRouter>
    </MockedProvider>
  )

  expect(screen.getByText('Timeline')).toBeInTheDocument()
  expect(screen.getByText('Albums')).toBeInTheDocument()

  expect(await screen.findByText('Settings')).toBeInTheDocument()
  expect(await screen.findByText('Places')).toBeInTheDocument()
})

// 4D-guesser, fork-specific
describe('4D-guesser button', () => {
  const renderMenu = (mapboxToken: string | null, guesserEnabled: boolean) => {
    authTokenMock.mockImplementation(() => 'test-token')
    render(
      <MockedProvider
        mocks={[
          {
            request: { query: MAPBOX_QUERY },
            result: { data: { mapboxToken } },
          },
          {
            request: { query: GUESSER_ENABLED_QUERY },
            result: { data: { siteInfo: { guesserEnabled } } },
          },
        ]}
        addTypename={false}
      >
        <MemoryRouter>
          <MainMenu />
        </MemoryRouter>
      </MockedProvider>
    )
  }

  test('shown with a Mapbox token when the game is enabled', async () => {
    renderMenu('token', true)
    expect(await screen.findByText('4D-guesser')).toBeInTheDocument()
  })

  test('hidden when the game is disabled', async () => {
    renderMenu('token', false)
    expect(await screen.findByText('Places')).toBeInTheDocument()
    expect(screen.queryByText('4D-guesser')).toBeNull()
  })

  test('hidden without a Mapbox token', async () => {
    renderMenu(null, true)
    await new Promise(resolve => setTimeout(resolve, 20))
    expect(screen.queryByText('4D-guesser')).toBeNull()
  })
})
