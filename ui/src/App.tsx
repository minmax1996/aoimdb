import React from 'react';
import { Counter } from './features/counter/Counter';
import { Database } from './features/database/Database';
// @ts-ignore  
import WebConsole from './features/web-console/WebConsole'
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <WebConsole />
        {/* <Counter /> */}
        <Database />
      </header>
    </div>
  );
}

export default App;
