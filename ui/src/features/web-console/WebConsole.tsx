import React, { useState } from 'react'

import ReactConsole from '@webscopeio/react-console'
import { useDispatch } from 'react-redux'
import { addKey } from '../database/databaseSlice';

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
              description: 'Echo',
              fn: (...args: any[]) => {
                return new Promise<any>((resolve, reject) => {
                  setTimeout(() => {
                    resolve(`${args.join(' ')}`)
                    dispatch(addKey(`${args.join(' ')}`))
                  }, 2000)
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
