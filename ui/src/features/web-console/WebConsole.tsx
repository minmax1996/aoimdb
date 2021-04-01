import React, { useState } from 'react'

import ReactConsole from '@webscopeio/react-console'
import { useDispatch } from 'react-redux'
import { addKey, clearKeys } from '../database/databaseSlice';

interface Command {
    description: string,
    fn: (...args: any[]) => Promise<any>
}

export default function WebConsole() {
    const [history, setHistory] = useState<string[]>([])
    const dispatch = useDispatch();
    return (
      <div>
        <ReactConsole
          autoFocus
          welcomeMessage="Welcome"
          lineStyle={{textAlign: "left"}}
          history={history}
          onAddHistoryItem={(item)=>setHistory([...history, item])}
          commands={{
            echo: { 
              fn: (...args: any[]) => {
                return new Promise<any>((resolve, _) => {
                  resolve(`${args.join(' ')}`)
                  dispatch(addKey(`${args.join(' ')}`))
                })
              }
            },
            clearkeys: { 
              fn: () => {
                return new Promise<any>((resolve, _) => {
                  resolve('')
                  dispatch(clearKeys())
                })
              }
            },
            history: {
              description: 'History',
              fn: () => new Promise(resolve => {
                 resolve(`${history.join('\r\n')}`)
              })
            }
          }}
        />
      </div>
    )
  }