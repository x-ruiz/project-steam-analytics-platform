import React from 'react'
import { ResponsivePie } from '@nivo/pie'

// const data = [
//   {
//     "id": "scala",
//     "label": "scala",
//     "value": 0.2,
//     "color": "hsl(204, 70%, 50%)"
//   },
//   {
//     "id": "javascript",
//     "label": "javascript",
//     "value": 0.4,
//     "color": "hsl(146, 70%, 50%)"
//   },
//   {
//     "id": "go",
//     "label": "go",
//     "value": 0.1,
//     "color": "hsl(116, 70%, 50%)"
//   },
//   {
//     "id": "css",
//     "label": "css",
//     "value": 0.2,
//     "color": "hsl(46, 70%, 50%)"
//   },
//   {
//     "id": "lisp",
//     "label": "lisp",
//     "value": 0.1,
//     "color": "hsl(98, 70%, 50%)"
//   }
// ]

export function PlaytimeDistributionChart({ lifetimePlaytimeData: lifetimePlaytimeData}) {
  var hue = Math.floor(Math.random() * 360)
  var data = []
  for (var i in lifetimePlaytimeData) {
    var dataObj = {}
    var obj = lifetimePlaytimeData[i]

    dataObj["id"] = obj.name
    dataObj["label"] = obj.name
    dataObj["value"] = Math.floor(obj.playtime_percentage * 100) 
    dataObj["color"] = `hsl(${hue}, 70%, 50%)`
    data.push(dataObj)
  }

  console.log("DATA", data)
return <ResponsivePie
      data={data}
      margin={{ top: 40, right: 80, bottom: 80, left: 80 }}
      innerRadius={0.5}
      padAngle={0.7}
      cornerRadius={3}
      activeOuterRadiusOffset={8}
      borderWidth={1}
      borderColor={{
          from: 'color',
          modifiers: [
              [
                  'darker',
                  0.2
              ]
          ]
      }}
      arcLinkLabelsSkipAngle={10}
      arcLinkLabelsTextColor="#ffffff"
      arcLinkLabelsThickness={2}
      arcLinkLabelsColor={{ from: 'color' }}
      arcLabel={d => `${d.value}%`}
      arcLabelsSkipAngle={10}
      arcLabelsTextColor={{
          from: 'color',
          modifiers: [
              [
                  'darker',
                  3
              ]
          ]
      }}
      defs={[
          {
              id: 'dots',
              type: 'patternDots',
              background: 'inherit',
              color: 'rgb(255, 255, 255)',
              size: 4,
              padding: 1,
              stagger: true
          },
          {
              id: 'lines',
              type: 'patternLines',
              background: 'inherit',
              color: 'rgba(255, 255, 255, 0.3)',
              rotation: -45,
              lineWidth: 6,
              spacing: 10
          }
      ]}
      fill={[
          {
              match: {
                  id: 'ruby'
              },
              id: 'dots'
          },
          {
              match: {
                  id: 'c'
              },
              id: 'dots'
          },
          {
              match: {
                  id: 'go'
              },
              id: 'dots'
          },
          {
              match: {
                  id: 'python'
              },
              id: 'dots'
          },
          {
              match: {
                  id: 'scala'
              },
              id: 'lines'
          },
          {
              match: {
                  id: 'lisp'
              },
              id: 'lines'
          },
          {
              match: {
                  id: 'elixir'
              },
              id: 'lines'
          },
          {
              match: {
                  id: 'javascript'
              },
              id: 'lines'
          }
      ]}
      legends={[
          {
              anchor: 'bottom',
              direction: 'row',
              justify: false,
              translateX: 0,
              translateY: 46,
              itemsSpacing: 5,
              itemWidth: 150,
              itemHeight: 18,
              itemTextColor: '#999',
              itemDirection: 'top-to-bottom',
              itemOpacity: 1,
              symbolSize: 18,
              symbolShape: 'circle',
              effects: [
                  {
                      on: 'hover',
                      style: {
                          itemTextColor: '#000'
                      }
                  }
              ]
          }
      ]}
  />
}
