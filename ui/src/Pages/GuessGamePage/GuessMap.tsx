// The 4D-guesser map: click or drag to place the guess, and after the guess show the
// answer with a line between the two. Fork-specific: not part of upstream Photoview.
import type mapboxgl from 'mapbox-gl'
import React, { useEffect, useRef, useState } from 'react'
import styled from 'styled-components'
import useMapboxMap from '../../components/mapbox/MapboxMap'
import { LatLng, unwrapLongitude, wrapLongitude } from './guessGameState'

const LINE_SOURCE = 'guess-line'
const GUESS_COLOR = '#3b82f6'
const ANSWER_COLOR = '#16a34a'
const MIN_START_ZOOM = 1.3

const MapWrapper = styled.div<{ pinning: boolean }>`
  width: 100%;
  height: 100%;

  && .mapboxgl-canvas-container.mapboxgl-interactive {
    cursor: ${({ pinning }) => (pinning ? 'crosshair' : 'grab')};
  }
`

export type MapBounds = { southWest: LatLng; northEast: LatLng }

type GuessMapProps = {
  /** The player's pin */
  pin: LatLng | null
  /** Called when the player places or moves the pin; leave out to lock the pin */
  onPin?(pin: LatLng): void
  /** The right answer, shown after the guess */
  answer?: LatLng | null
  /** Where to look while guessing */
  bounds?: MapBounds | null
  /** Changes for every round, so the view resets to `bounds` */
  viewKey: string
}

const emptyLine: GeoJSON.Feature<GeoJSON.LineString> = {
  type: 'Feature',
  properties: {},
  geometry: { type: 'LineString', coordinates: [] },
}

const GuessMap = ({ pin, onPin, answer, bounds, viewKey }: GuessMapProps) => {
  const [loaded, setLoaded] = useState(false)
  const onPinRef = useRef(onPin)
  onPinRef.current = onPin

  const guessMarker = useRef<mapboxgl.Marker | null>(null)
  const answerMarker = useRef<mapboxgl.Marker | null>(null)

  const { mapContainer, mapboxMap, mapboxLibrary } = useMapboxMap({
    configureMapbox: (map, mapbox) => {
      map.addControl(new mapbox.NavigationControl({ showCompass: false }))
      map.dragRotate.disable()
      map.touchZoomRotate.disableRotation()

      map.on('click', event => {
        onPinRef.current?.({
          latitude: event.lngLat.lat,
          longitude: wrapLongitude(event.lngLat.lng),
        })
      })

      map.on('load', () => {
        map.addSource(LINE_SOURCE, { type: 'geojson', data: emptyLine })
        map.addLayer({
          id: LINE_SOURCE,
          type: 'line',
          source: LINE_SOURCE,
          paint: {
            'line-color': '#111827',
            'line-width': 2,
            'line-dasharray': [2, 2],
          },
        })
        setLoaded(true)
      })
    },
    mapboxOptions: {
      projection: { name: 'globe' },
      center: [0, 20],
      zoom: 1,
    },
  })

  const map = loaded ? mapboxMap : null

  // The guess pin
  useEffect(() => {
    if (!map || !mapboxLibrary) return

    if (pin === null) {
      guessMarker.current?.remove()
      guessMarker.current = null
      return
    }

    if (guessMarker.current === null) {
      const marker = new mapboxLibrary.Marker({ color: GUESS_COLOR })
      marker.on('dragend', () => {
        const lngLat = marker.getLngLat()
        onPinRef.current?.({
          latitude: lngLat.lat,
          longitude: wrapLongitude(lngLat.lng),
        })
      })
      guessMarker.current = marker.setLngLat([0, 0]).addTo(map)
    }

    guessMarker.current
      .setLngLat([pin.longitude, pin.latitude])
      .setDraggable(onPin !== undefined)
  }, [map, mapboxLibrary, pin?.latitude, pin?.longitude, onPin !== undefined])

  // The answer and the line to it
  useEffect(() => {
    if (!map || !mapboxLibrary) return
    const lineSource = map.getSource(LINE_SOURCE) as
      | mapboxgl.GeoJSONSource
      | undefined

    if (!answer) {
      answerMarker.current?.remove()
      answerMarker.current = null
      lineSource?.setData(emptyLine)
      return
    }

    const answerLon = pin
      ? unwrapLongitude(pin.longitude, answer.longitude)
      : answer.longitude

    if (answerMarker.current === null) {
      answerMarker.current = new mapboxLibrary.Marker({ color: ANSWER_COLOR })
        .setLngLat([answerLon, answer.latitude])
        .addTo(map)
    }
    answerMarker.current.setLngLat([answerLon, answer.latitude])

    if (pin) {
      lineSource?.setData({
        ...emptyLine,
        geometry: {
          type: 'LineString',
          coordinates: [
            [pin.longitude, pin.latitude],
            [answerLon, answer.latitude],
          ],
        },
      })
      map.fitBounds(
        [
          [
            Math.min(pin.longitude, answerLon),
            Math.min(pin.latitude, answer.latitude),
          ],
          [
            Math.max(pin.longitude, answerLon),
            Math.max(pin.latitude, answer.latitude),
          ],
        ],
        { padding: 60, maxZoom: 14, duration: 800 }
      )
    } else {
      map.flyTo({ center: [answerLon, answer.latitude], zoom: 10 })
    }
  }, [
    map,
    mapboxLibrary,
    answer?.latitude,
    answer?.longitude,
    pin?.latitude,
    pin?.longitude,
  ])

  // A new round starts with the view of the whole pool
  useEffect(() => {
    if (!map || !bounds) return
    const { southWest, northEast } = bounds
    const camera = map.cameraForBounds(
      [
        [southWest.longitude, southWest.latitude],
        [northEast.longitude, northEast.latitude],
      ],
      { padding: 40, maxZoom: 10 }
    )
    if (!camera) return
    // A pool spread over the world would otherwise show a tiny globe
    map.jumpTo({ ...camera, zoom: Math.max(camera.zoom ?? 0, MIN_START_ZOOM) })
  }, [map, viewKey, bounds !== null && bounds !== undefined])

  return (
    <MapWrapper
      className="rounded overflow-hidden"
      data-testid="guess-map"
      pinning={onPin !== undefined}
    >
      {mapContainer}
    </MapWrapper>
  )
}

export default GuessMap
